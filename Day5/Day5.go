package main

import (
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync/atomic"
)

type idRange struct {
	bottom int
	top    int
}

func newIdRange(bottom int, top int) *idRange {
	return &idRange{
		bottom: bottom,
		top:    top,
	}
}

func (idRangeToCheck *idRange) checkIdRange(val int) bool {
	return val >= idRangeToCheck.bottom && val <= idRangeToCheck.top
}
func (idRangeToCheck *idRange) totalLength() int {
	return idRangeToCheck.top - idRangeToCheck.bottom + 1
}

func checkRangeOverlap(left idRange, right idRange) bool {
	return left.bottom <= right.top && right.bottom <= left.top
}

func main() {
	data, _ := os.ReadFile("input.txt")
	raw_split := strings.Split(strings.ReplaceAll(string(data), "\r", ""), "\n\n")
	fresh_ranges := strings.Split(raw_split[0], "\n")
	ids_to_check := strings.Split(raw_split[1], "\n")

	idRanges := make([]*idRange, 0)

	for _, idRange := range fresh_ranges {
		splitRange := strings.Split(idRange, "-")
		bottom, err := strconv.Atoi(splitRange[0])
		if err != nil {
			fmt.Println(splitRange)
			panic("Bad Val1")
		}
		top, err := strconv.Atoi(splitRange[1])
		if err != nil {
			fmt.Println(splitRange)
			panic("Bad Val 2")
		}
		idRanges = append(idRanges, newIdRange(
			bottom,
			top,
		))
	}

	var total_fresh atomic.Int64
	total_fresh.Store(0)

	for _, id_to_check := range ids_to_check {
		id_to_check_int, err := strconv.Atoi(id_to_check)
		if err != nil {
			fmt.Println(id_to_check)
			panic("Bad Id To check")
		}
		for _, idRangeCheck := range idRanges {
			if idRangeCheck.checkIdRange(id_to_check_int) {
				total_fresh.Add(1)
				break
			}
		}
	}
	fmt.Println("Checking full range size")
	reducedCount := true
	for reducedCount {
		reducedCount = false
		newRanges := make([]*idRange, 0)
		pairs := make([]string, 0)
		rangeIndicesToRemove := make([]int, 0)
		for i, idRangeCheck0 := range idRanges {
			for j, ididRangeCheck1 := range idRanges[i+1:] {
				if checkRangeOverlap(*idRangeCheck0, *ididRangeCheck1) {
					reducedCount = true
					if !slices.Contains(rangeIndicesToRemove, i) {
						rangeIndicesToRemove = append(rangeIndicesToRemove, i)

					}
					if !slices.Contains(rangeIndicesToRemove, j+i+1) {

						rangeIndicesToRemove = append(rangeIndicesToRemove, j+i+1)
					}
					newBottom := int(math.Min(float64(idRangeCheck0.bottom), float64(ididRangeCheck1.bottom)))
					newTop := int(math.Max(float64(idRangeCheck0.top), float64(ididRangeCheck1.top)))
					newRange := fmt.Sprintf("%d,%d", newBottom, newTop)
					if !slices.Contains(pairs, newRange) {
						newRanges = append(newRanges,
							newIdRange(
								int(math.Min(float64(idRangeCheck0.bottom), float64(ididRangeCheck1.bottom))),
								int(math.Max(float64(idRangeCheck0.top), float64(ididRangeCheck1.top))),
							))
						pairs = append(pairs, newRange)
					}
				}
			}
		}
		slices.Sort(rangeIndicesToRemove)
		slices.Reverse(rangeIndicesToRemove)
		fmt.Printf("Consolidating down by %d", len(rangeIndicesToRemove))
		fmt.Println()
		for _, indexToRemove := range rangeIndicesToRemove {
			idRanges = append(idRanges[:indexToRemove], idRanges[indexToRemove+1:]...)
		}
		idRanges = append(idRanges, newRanges...)
	}
	total_range_len := 0
	for _, idRangeCheck := range idRanges {
		total_range_len += idRangeCheck.totalLength()
	}
	fmt.Printf("Part 1: %d, Part 2: %d", total_fresh.Load(), total_range_len)
}
