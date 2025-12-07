package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"unicode/utf8"
)

func calcMaxJoltage(line string, total_joltage *atomic.Int64) {
	line = strings.Trim(line, "\r\n")
	max_joltage := 0
	for i, char_val_1 := range line {
		for _, char_val_2 := range line[i+1:] {
			if curr_joltage, _ := strconv.Atoi(string(char_val_1) + string(char_val_2)); curr_joltage > max_joltage {
				max_joltage = curr_joltage
			}
		}
	}
	total_joltage.Add(int64(max_joltage))
}

func recurJoltage(curr_joltage []rune, curr_line string, depth int, curr_max *int) {
	for i, char_val := range curr_line {
		curr_max_string := strconv.Itoa(*curr_max)
		curr_max_clipped, _ := strconv.Atoi(curr_max_string[:depth])
		new_joltage := string(curr_joltage) + string(char_val)
		new_joltage_int, _ := strconv.Atoi(new_joltage)
		// fmt.Println(new_joltage_int, *curr_max)
		// if depth == 11 {
		// 	fmt.Println(new_joltage_int, char_val < curr_max_runes[depth-1], string(curr_max_runes[depth-1]), string(char_val))
		// }
		if depth == 12 && new_joltage_int > *curr_max {
			*curr_max = new_joltage_int
		} else if new_joltage_int < curr_max_clipped {
			continue
		} else if i+1 < utf8.RuneCountInString(curr_line) && depth < 12 {
			recurJoltage([]rune(new_joltage), curr_line[i+1:], depth+1, curr_max)
		}
	}
}

func calcMaxJoltage2(line string, total_joltage *atomic.Int64) {
	line = strings.Trim(line, "\r\n")
	max_joltage, _ := strconv.Atoi(line[:12])
	recurJoltage(
		[]rune(""),
		line,
		1,
		&max_joltage,
	)

	fmt.Println(total_joltage.Load(), max_joltage)
	total_joltage.Add(int64(max_joltage))
}

func main() {
	data, _ := os.ReadFile("input.txt")
	id_ranges := strings.Split(string(data), "\n")

	var total_joltage atomic.Int64
	var total_joltage_2 atomic.Int64
	total_joltage.Store(0)
	total_joltage_2.Store(0)

	var wait_group sync.WaitGroup
	for _, line := range id_ranges {
		wait_group.Go(
			func() {
				calcMaxJoltage(line, &total_joltage)
			})
		wait_group.Go(
			func() {
				calcMaxJoltage2(line, &total_joltage_2)
			})
	}
	wait_group.Wait()
	fmt.Printf("Part 1: %d, Part 2: %d", total_joltage.Load(), total_joltage_2.Load())
	fmt.Println()
}
