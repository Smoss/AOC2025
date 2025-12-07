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

func check_rune(i int, j int, rollsChan chan<- *readOp, total *atomic.Int64) bool {
	coord := make_coord(i, j)
	resultChan := make(chan *readResult)
	rollsChan <- &readOp{
		key:        coord,
		resultChan: resultChan,
	}
	result := <-resultChan
	if val, prs := result.val, result.prs; !prs || val != '@' {
		return false
	}
	nearby_rolls := 0
	for x := -1; x < 2; x++ {
		for y := -1; y < 2; y++ {
			if x != 0 || y != 0 {
				checkResultChan := make(chan *readResult)

				check_coord := make_coord(i+x, y+j)
				rollsChan <- &readOp{
					key:        check_coord,
					resultChan: checkResultChan,
				}
				checkResult := <-checkResultChan
				check_val, check_prs := checkResult.val, checkResult.prs
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

type readResult struct {
	val rune
	prs bool
}

type readOp struct {
	key        string
	resultChan chan *readResult
}

type writeOp struct {
	key         string
	val         rune
	confirmChan chan *struct{}
}

func listen(readChan <-chan *readOp, writeChan <-chan *writeOp, rolls *map[string]rune) {
	for {
		select {
		case read := <-readChan:
			val, prs := (*rolls)[read.key]
			read.resultChan <- &readResult{
				val: val,
				prs: prs,
			}
		case write := <-writeChan:
			(*rolls)[write.key] = write.val
			write.confirmChan <- &struct{}{}
		}
	}
}

func main() {
	data, _ := os.ReadFile("input.txt")
	raw_rolls := strings.Split(string(data), "\n")

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
	readChan := make(chan *readOp)
	writeChan := make(chan *writeOp)
	go listen(
		readChan,
		writeChan,
		&rolls,
	)

	var wait_group sync.WaitGroup
	for i, line := range raw_rolls {
		for j := range line {
			wait_group.Go(
				func() {
					check_rune(i, j, readChan, &total_rolls)
				})

		}
	}
	wait_group.Wait()
	less_rolls := true
	for less_rolls {
		less_rolls = false
		for i, line := range raw_rolls {
			for j := range line {
				wait_group.Go(
					func() {
						can_remove := check_rune(i, j, readChan, &total_rolls_2)
						coord := make_coord(i, j)
						if can_remove {
							less_rolls = true
							confirmChan := make(chan *struct{})
							writeChan <- &writeOp{
								key:         coord,
								val:         '.',
								confirmChan: confirmChan,
							}
							<-confirmChan
						}
					})
			}
		}
		wait_group.Wait()
	}

	fmt.Printf("Part 1: %d, Part 2: %d", total_rolls.Load(), total_rolls_2.Load())
}
