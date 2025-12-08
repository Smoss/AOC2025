package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

func make_coord(i int, j int) string {
	return strconv.Itoa(i) + "," + strconv.Itoa(j)
}

func check_rune(i int, j int, rolls map[string]rune, total *atomic.Int64) bool {
	coord := make_coord(i, j)
	if val, prs := rolls[coord]; !prs || val != '@' {
		return false
	}
	nearby_rolls := 0
	for x := -1; x < 2; x++ {
		for y := -1; y < 2; y++ {
			if x != 0 || y != 0 {
				check_coord := make_coord(i+x, y+j)
				check_val, check_prs := rolls[check_coord]
				if check_prs && check_val == '@' {
					nearby_rolls++
				}
			}
		}
	}
	if nearby_rolls < 4 {
		total.Add(1)
		return true
	}
	return false
}

type writeOp struct {
	key         string
	confirmChan chan *struct{}
}

func listen(resetChan <-chan chan *struct{}, writeChan <-chan *writeOp, rollsToRemove *[]string) {
	for {
		select {
		case write := <-writeChan:
			*rollsToRemove = append(*rollsToRemove, write.key)
			write.confirmChan <- &struct{}{}
		case responseChan := <-resetChan:
			*rollsToRemove = make([]string, 0)
			responseChan <- &struct{}{}
		}
	}
}

func main() {
	data, _ := os.ReadFile("input.txt")
	raw_rolls := strings.Split(strings.ReplaceAll(string(data), "\r", ""), "\n")

	rolls := make(map[string]rune)
	var total_rolls atomic.Int64
	var total_rolls_2 atomic.Int64
	total_rolls.Store(0)
	total_rolls_2.Store(0)

	for i, line := range raw_rolls {
		for j, val := range line {
			rolls[make_coord(i, j)] = val
		}
	}

	var wait_group sync.WaitGroup
	for i, line := range raw_rolls {
		for j := range line {
			wait_group.Go(
				func() {
					check_rune(i, j, rolls, &total_rolls)
				})

		}
	}
	wait_group.Wait()
	writeChan := make(chan *writeOp)
	resetChan := make(chan chan *struct{})
	less_rolls := true
	rollsToRemove := make([]string, 0)
	go listen(
		resetChan,
		writeChan,
		&rollsToRemove,
	)
	for less_rolls {
		less_rolls = false
		for i, line := range raw_rolls {
			for j := range line {
				wait_group.Go(
					func() {
						can_remove := check_rune(i, j, rolls, &total_rolls_2)
						coord := make_coord(i, j)
						if can_remove {
							confirmChan := make(chan *struct{})
							writeChan <- &writeOp{
								key:         coord,
								confirmChan: confirmChan,
							}
							<-confirmChan
						}
					})
			}
		}
		wait_group.Wait()
		for _, coord := range rollsToRemove {
			rolls[coord] = '.'
		}
		less_rolls = len(rollsToRemove) > 0
		responseChan := make(chan *struct{})
		resetChan <- responseChan
		<-responseChan
	}

	fmt.Printf("Part 1: %d, Part 2: %d", total_rolls.Load(), total_rolls_2.Load())
}
