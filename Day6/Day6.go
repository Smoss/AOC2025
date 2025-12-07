package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"unicode/utf8"
)

type numberResult struct {
	rowIndex int
	numbers  []string
}

func newNumberResult(index int, numbers []string) *numberResult {
	return &numberResult{
		rowIndex: index,
		numbers:  numbers,
	}
}

func product(numbers []int) int {
	currProd := 1
	for _, val := range numbers {
		currProd *= val
	}
	return currProd
}
func sum(numbers []int) int {
	currSum := 0
	for _, val := range numbers {
		currSum += val
	}
	return currSum
}

func parseLeftAligned(numberLine []string, currIndex int) (string, int) {
	seenSpace := false
	number := ""
	for currIndex < len(numberLine) {
		val := numberLine[currIndex]
		currIndex++
		if val != " " {
			if !seenSpace {
				number += val
			} else {
				number = number[:utf8.RuneCountInString(number)-1]
				break
			}
		} else {
			number += "0"
			seenSpace = true
		}
	}
	return number, currIndex - 1
}

func parseRightAligned(numberLine []string, currIndex int) (string, int) {
	seenNumber := false
	number := ""
	for currIndex < len(numberLine) {
		val := numberLine[currIndex]
		currIndex++
		if val == " " {
			if !seenNumber {
				number += "0"
			} else {
				break
			}
		} else {
			number += val
			seenNumber = true
		}
	}
	return number, currIndex
}

func parseLine(numberLine string, rowIndex int, columnWidths []int) *numberResult {
	numbers := make([]string, 0)
	currIndex := 0
	for _, val := range columnWidths {
		newNumber := numberLine[currIndex : currIndex+val]
		numbers = append(numbers, strings.ReplaceAll(newNumber, " ", "0"))
		currIndex += val + 1
	}

	return newNumberResult(rowIndex, numbers)
}

func parseCephNumbers(cephColumn []string, columnWidth int, action string) int {
	numbers := make([]int, 0)
	for i := 0; i < columnWidth; i++ {
		number := make([]string, 0)
		for _, line := range cephColumn {
			number = append(number, strings.Split(line, "")[i])
		}
		// slices.Reverse(number)
		result, _ := strconv.Atoi(strings.TrimRight(strings.Join(number, ""), "0"))
		numbers = append(numbers, result)
	}

	fmt.Println(numbers, cephColumn, columnWidth, action)
	if action == "*" {
		return product(numbers)
	} else {
		return sum(numbers)
	}
}

func main() {
	data, _ := os.ReadFile("input.txt")
	raw_split := strings.Split(strings.ReplaceAll(string(data), "\r", ""), "\n")
	totalRows := len(raw_split)
	allNumbers := raw_split[:totalRows-1]
	r1, _ := regexp.Compile(" +")
	r2, _ := regexp.Compile("( +$)|(^ )")
	actions := strings.Split(r2.ReplaceAllString(r1.ReplaceAllString(raw_split[totalRows-1], " "), ""), " ")
	totalColumns := len(actions)
	columns := make([][]int, totalColumns)
	for _, line := range allNumbers {
		splitLine := strings.Split(r2.ReplaceAllString(r1.ReplaceAllString(line, " "), ""), " ")
		for j, val := range splitLine {
			intVal, _ := strconv.Atoi(val)
			columns[j] = append(columns[j], intVal)
		}
	}
	var wg sync.WaitGroup
	paddedNumbers := make([][]string, totalRows-1)
	numberStation := make(chan *numberResult)
	actionLine := raw_split[totalRows-1]
	columnWidths := make([]int, 0)
	curr_width := 1
	for _, val := range strings.Split(actionLine, "")[1:] {
		if val != " " {
			columnWidths = append(columnWidths, curr_width-1)
			curr_width = 1
		} else {
			curr_width++
		}
	}
	columnWidths = append(columnWidths, curr_width)

	var total_sum atomic.Int64
	total_sum.Store(0)

	wg.Go(
		func() {
			for i := 0; i < totalRows-1; i++ {
				result := <-numberStation
				paddedNumbers[result.rowIndex] = result.numbers
			}
		})
	for i, column := range columns {
		wg.Go(
			func() {
				val := 0
				if actions[i] == "+" {
					val = sum(column)
				} else {
					val = product(column)
				}
				total_sum.Add(int64(val))
			})
	}

	for i, line := range allNumbers {
		wg.Go(
			func() {

				numberStation <- parseLine(line, i, columnWidths)
			})
	}
	wg.Wait()
	// fmt.Println(paddedNumbers)
	cephColumns := make([][]string, totalColumns)
	for _, row := range paddedNumbers {
		for j, paddedNumber := range row {
			cephColumns[j] = append(cephColumns[j], paddedNumber)
		}
	}
	// fmt.Println(cephColumns)
	var cephSum atomic.Int64
	cephSum.Store(0)
	for i, columnWidth := range columnWidths {
		wg.Go(
			func() {
				result := parseCephNumbers(cephColumns[i], columnWidth, actions[i])
				cephSum.Add(
					int64(result),
				)
			},
		)
	}
	wg.Wait()

	fmt.Printf("Part 1: %d, Part 2: %d", total_sum.Load(), cephSum.Load())
}
