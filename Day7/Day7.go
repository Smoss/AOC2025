package main

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

func main() {
	data, _ := os.ReadFile("input.txt")
	raw_split := strings.Split(strings.ReplaceAll(string(data), "\r", ""), "\n")
	splitterWidth := utf8.RuneCountInString(raw_split[0])
	currIndices := make(map[int]int)
	startIndex := strings.IndexRune(raw_split[0], 'S')
	currIndices[startIndex] = 1
	splitCount := 0
	// fmt.Println(currIndices)
	for _, line := range raw_split[1:] {
		newIndices := make(map[int]int)
		for k := range currIndices {
			if k >= 0 && k < splitterWidth && []rune(line)[k] == '^' {
				if _, prs1 := newIndices[k-1]; !prs1 {
					newIndices[k-1] = 0
				}
				newIndices[k-1] += currIndices[k]
				if _, prs2 := newIndices[k+1]; !prs2 {
					newIndices[k+1] = 0
				}
				newIndices[k+1] += currIndices[k]
				splitCount++
			} else {
				if _, prs := newIndices[k]; !prs {
					newIndices[k] = 0
				}
				newIndices[k] += currIndices[k]
			}
		}
		currIndices = newIndices
		// fmt.Println(currIndices)
		// testLine := strings.Split(line, "")
		// for k, v := range currIndices {
		// 	testLine[k] = strconv.Itoa(v)
		// }
		// fmt.Println(strings.Join(testLine, " "))
	}
	timelines := 0
	for _, v := range currIndices {
		timelines += v
	}
	fmt.Printf("Part 1: %d, Part 2: %d", splitCount, timelines)
	fmt.Println()

}
