package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
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

func main() {
	data, _ := os.ReadFile("input.txt")
	id_ranges := strings.Split(string(data), "\n")
	invalid_id_total := 0
	invalid_id_total_p_2 := 0

	for _, pair := range id_ranges {
		split_pair := strings.Split(strings.Trim(pair, "\r\n"), "-")
		range_start, err := strconv.Atoi(split_pair[0])

		if err != nil {
			fmt.Println("Bad start")
			fmt.Printf("%s", split_pair[0])
			fmt.Println()
			break
		}

		range_end, err := strconv.Atoi(split_pair[1])
		if err != nil {
			fmt.Println("Bad end")
			fmt.Printf("%s", split_pair[1])
			fmt.Println()
			break
		}

		for i := range_start; i <= range_end; i++ {
			string_id := strconv.Itoa(i)
			str_len := utf8.RuneCountInString(string_id)

			if check_string(string_id, str_len) {
				invalid_id_total_p_2 += i
			}
			if str_len%2 != 0 {
				continue
			}

			first_half := string_id[:str_len/2]
			second_half := string_id[str_len/2:]
			if first_half == second_half {
				invalid_id_total += i
			}
		}
	}
	fmt.Printf("Part 1: %d, Part 2: %d", invalid_id_total, invalid_id_total_p_2)
	fmt.Println()
}
