package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	data, _ := os.ReadFile("input.txt")
	twists := strings.Split(string(data), "\n")
	curr_point := 50
	crossings := 0
	true_crossings := 0
	for _, twist := range twists {
		dir := twist[0:1]
		quantity, _ := strconv.Atoi(strings.Trim(twist[1:], " \r\n"))
		orig := curr_point
		switch dir {
		case "L":
			curr_point -= quantity
		default:
			curr_point += quantity
		}
		true_crossings += (quantity / 100)
		for curr_point < 0 {
			curr_point = curr_point + 100
		}
		for curr_point >= 100 {
			curr_point = curr_point - 100
		}
		if dir == "L" && curr_point > orig && orig != 0 {
			true_crossings += 1
		} else if dir == "R" && curr_point < orig && curr_point != 0 {
			true_crossings += 1
		}
		if curr_point == 0 {
			crossings += 1
			true_crossings++
		}
	}
	fmt.Printf("%d, %d", crossings, true_crossings)
	fmt.Println()
}
