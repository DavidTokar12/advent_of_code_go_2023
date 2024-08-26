package main

import (
	"advent_of_code_3/circular_buffer"
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"unicode"
)

type Coordinate struct {
	x int
	y int
}

var m = make(map[Coordinate][]int)

func findNumericIndexes(s *[]rune) chan [2]int {
	ch := make(chan [2]int)

	go func() {
		defer close(ch)

		var begin_digit_index int = -1

		for index, char := range *s {
			if unicode.IsDigit(char) {
				if begin_digit_index == -1 {
					begin_digit_index = index
				}
			} else {
				if begin_digit_index != -1 {
					end_digit_index := index
					ch <- [2]int{begin_digit_index, end_digit_index}
					begin_digit_index = -1
				}
			}
		}

		if begin_digit_index != -1 {
			ch <- [2]int{begin_digit_index, len(*s)}
		}

	}()

	return ch
}

func sumTwoLines(line_to_sum *[]rune, line_to_check *[]rune, line_number int, upper bool) int {
	var result int

	for indexes := range findNumericIndexes(line_to_sum) {
		begin_number, end_number := indexes[0], indexes[1]

		if isAdjacentToSymbolTwoLine(line_to_sum, line_to_check, begin_number, end_number, line_number, upper) {

			num, err := strconv.Atoi(string((*line_to_sum)[begin_number:end_number]))

			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(string((*line_to_sum)[begin_number:end_number]))

			result += num
		}

	}

	return result
}

func isAdjacentToSymbolTwoLine(line_to_sum *[]rune, line_to_check *[]rune, begin_number int, end_number int, line_number int, upper bool) bool {
	begin := max(0, begin_number-1)
	end := min(len(*line_to_sum)-1, end_number)

	var result bool = false

	num, err := strconv.Atoi(string((*line_to_sum)[begin_number:end_number]))
	if err != nil {
		log.Fatal(err)
	}

	if begin_number != 0 && (*line_to_sum)[begin] != '.' {
		result = true

		if (*line_to_sum)[begin] == '*' {
			coord := Coordinate{line_number, begin}
			m[coord] = append(m[coord], num)
		}

	}

	if end_number != len(*line_to_sum) && (*line_to_sum)[end] != '.' {
		result = true

		if (*line_to_sum)[end] == '*' {
			coord := Coordinate{line_number, end}
			m[coord] = append(m[coord], num)
		}
	}

	for i := begin; i <= end; i++ {
		if (*line_to_check)[i] != '.' {
			result = true

			if (*line_to_check)[i] == '*' {

				var ln int

				if upper {
					ln = line_number + 1
				} else {
					ln = line_number - 1
				}

				coord := Coordinate{ln, i}
				m[coord] = append(m[coord], num)
			}

		}
	}

	return result

}

func isAdjacentToSymbol(first_line *[]rune, second_line *[]rune, third_line *[]rune, begin_number int, end_number int, line_number int) bool {
	var result bool = false

	num, err := strconv.Atoi(string((*second_line)[begin_number:end_number]))
	if err != nil {
		log.Fatal(err)
	}

	begin := max(0, begin_number-1)
	end := min(len(*second_line)-1, end_number)

	if begin_number != 0 && (*second_line)[begin] != '.' {

		if (*second_line)[begin] == '*' {
			coord := Coordinate{line_number, begin}
			m[coord] = append(m[coord], num)
		}

		result = true
	}

	if end_number != len(*second_line) && (*second_line)[end] != '.' {

		if (*second_line)[end] == '*' {
			coord := Coordinate{line_number, end_number}
			m[coord] = append(m[coord], num)
		}

		result = true
	}

	for i := begin; i <= end; i++ {
		if (*first_line)[i] != '.' || (*third_line)[i] != '.' {
			result = true

			if (*first_line)[i] == '*' {
				coord := Coordinate{line_number - 1, i}
				m[coord] = append(m[coord], num)
			}

			if (*third_line)[i] == '*' {
				coord := Coordinate{line_number + 1, i}
				m[coord] = append(m[coord], num)
			}

		}

	}

	return result
}

func sumLine(first_line *[]rune, second_line *[]rune, third_line *[]rune, line_number int) int {
	var result int = 0

	for indexes := range findNumericIndexes(second_line) {
		begin_number, end_number := indexes[0], indexes[1]

		if isAdjacentToSymbol(first_line, second_line, third_line, begin_number, end_number, line_number) {

			num, err := strconv.Atoi(string((*second_line)[begin_number:end_number]))

			if err != nil {
				log.Fatal(err)
			}

			result += num
		}

	}

	return result
}

func main() {

	buf := circular_buffer.NewCircularTripleBuffer()

	file, err := os.Open("input.txt")

	if err != nil {
		log.Fatal(err)
	}

	fileScanner := bufio.NewScanner(file)

	var result int = 0
	var line_number int = -1

	for fileScanner.Scan() {
		line := fileScanner.Text()
		runes := []rune(line)
		line_number++

		if buf.GetSize() == 2 {
			rune_thripplet := buf.GetAll()
			result += sumTwoLines(rune_thripplet[0], rune_thripplet[1], line_number-2, true)
		}

		buf.Append(&runes)
		if buf.IsFull() {
			rune_thripplet := buf.GetAll()
			result += sumLine(rune_thripplet[0], rune_thripplet[1], rune_thripplet[2], line_number-1)
		}
	}

	rune_thripplet := buf.GetAll()
	result += sumTwoLines(rune_thripplet[2], rune_thripplet[1], line_number, false)

	var result_2 int = 0

	for _, value := range m {
		if len(value) >= 2 {

			product := 1
			for _, num := range value {
				product *= num
			}
			result_2 += product

		}
	}

	fmt.Println(m)
	fmt.Println(result_2)
	fmt.Println(result)

}
