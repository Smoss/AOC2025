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

		switch dir {
		case "L":
			for quantity > 0 {
				curr_point--
				quantity -= 1
				if curr_point == 0 {
					true_crossings++
				}
				if curr_point < 0 {
					curr_point = 99
				}
			}
		default:
			for quantity > 0 {
				curr_point++
				quantity -= 1
				if curr_point >= 100 {
					curr_point = 0
					true_crossings++
				}
			}
		}
		for curr_point < 0 {
			curr_point = curr_point + 100
		}
		for curr_point >= 100 {
			curr_point = curr_point - 100
		}

		if curr_point == 0 {
			crossings += 1
		}
		fmt.Printf("Curr dial: %d, %d, %s, %d", quantity, curr_point, twist[:1], true_crossings)
		fmt.Println()
	}
	fmt.Printf("%d, %d", crossings, true_crossings)
	fmt.Println()
}
