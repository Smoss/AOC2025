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

func check_string(id_string string, str_len int) bool {
	median := str_len / 2

	for i := 1; i <= median; i++ {
		check_substr := id_string[:i]
		check_substr_len := utf8.RuneCountInString(check_substr)
		if str_len%check_substr_len != 0 {
			continue
		}

		is_repeated := true
		for j := check_substr_len; j < str_len; j += check_substr_len {
			if id_string[j:j+check_substr_len] != check_substr {
				is_repeated = false
				break
			}
		}

		if is_repeated {
			return true
		}
	}
	return false
}

func calc_range(pair string, invalid_id_total *atomic.Int64, invalid_id_total_p_2 *atomic.Int64) {
	split_pair := strings.Split(strings.Trim(pair, "\r\n"), "-")
	range_start, err := strconv.Atoi(split_pair[0])

	if err != nil {
		fmt.Println("Bad start")
		fmt.Printf("%s", split_pair[0])
		fmt.Println()
		return
	}

	range_end, err := strconv.Atoi(split_pair[1])
	if err != nil {
		fmt.Println("Bad end")
		fmt.Printf("%s", split_pair[1])
		fmt.Println()
		return
	}

	for i := range_start; i <= range_end; i++ {
		string_id := strconv.Itoa(i)
		str_len := utf8.RuneCountInString(string_id)
		stored_i := int64(i)

		if check_string(string_id, str_len) {
			invalid_id_total_p_2.Add(stored_i)
			// fmt.Println(stored_i, invalid_id_total_p_2.Load())
		}
		if str_len%2 != 0 {
			continue
		}

		first_half := string_id[:str_len/2]
		second_half := string_id[str_len/2:]
		if first_half == second_half {
			invalid_id_total.Add(stored_i)
			// fmt.Println(stored_i, invalid_id_total.Load())
		}
	}
}

func main() {
	data, _ := os.ReadFile("input.txt")
	id_ranges := strings.Split(string(data), "\n")
	var invalid_id_total atomic.Int64
	var invalid_id_total_p_2 atomic.Int64
	invalid_id_total.Store(0)
	invalid_id_total_p_2.Store(0)
	var wait_group sync.WaitGroup
	for _, pair := range id_ranges {
		wait_group.Go(
			func() {
				calc_range(pair, &invalid_id_total, &invalid_id_total_p_2)
			})
	}
	wait_group.Wait()
	fmt.Printf("Part 1: %d, Part 2: %d", invalid_id_total.Load(), invalid_id_total_p_2.Load())
	fmt.Println()
}
