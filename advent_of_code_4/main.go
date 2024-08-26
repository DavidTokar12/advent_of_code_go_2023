package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type Card struct {
	score  int
	copies int
}

func processLine(s *string) ([]int, []int, error) {
	var your_numbers []int
	var winning_numbers []int

	cardsSlice := strings.Split(*s, ":")
	if len(cardsSlice) != 2 {
		return your_numbers, winning_numbers, fmt.Errorf("No `:` found in line, or more than one `:` in line: %s", *s)
	}

	numbersSlice := strings.Split(cardsSlice[1], "|")
	if len(numbersSlice) != 2 {
		return your_numbers, winning_numbers, fmt.Errorf("No `|` found in line, or more than one `|` in line: %s", *s)
	}

	your_numbers_string := strings.Fields(numbersSlice[0])
	winning_numbers_string := strings.Fields(numbersSlice[1])

	for _, num_string := range your_numbers_string {
		value, err := strconv.Atoi(num_string)

		if err != nil {
			return your_numbers, winning_numbers, err
		}

		your_numbers = append(your_numbers, value)
	}

	for _, num_string := range winning_numbers_string {
		value, err := strconv.Atoi(num_string)

		if err != nil {
			return your_numbers, winning_numbers, err
		}

		winning_numbers = append(winning_numbers, value)
	}

	return your_numbers, winning_numbers, nil

}

func getLineScore(your_numbers []int, winning_numbers []int) int {
	var result int = 0

	var hashMap [100]bool

	for _, num := range winning_numbers {
		hashMap[num] = true
	}

	for _, num := range your_numbers {

		if hashMap[num] {

			if result == 0 {
				result = 1
			} else {
				result = result * 2
			}

		}

	}

	return result
}

func main() {

	file, err := os.Open("input.txt")

	if err != nil {
		log.Fatal(err)
	}

	fileScanner := bufio.NewScanner(file)

	var cards []Card

	for fileScanner.Scan() {
		line := fileScanner.Text()
		your_numbers, winning_numbers, err := processLine(&line)
		if err != nil {
			log.Fatal(err)
		}

		card_result := getLineScore(your_numbers, winning_numbers)

		cards = append(cards, Card{card_result, 1})

	}

	for i := 0; i < len(cards); i++ {
		var number_of_matches int

		switch s := cards[i].score; s {
		case 0:
			number_of_matches = 0
		case 1:
			number_of_matches = 1
		default:
			number_of_matches = int(math.Log2(float64(s)) + 1)

		}
		for j := (i + 1); j < min(len(cards), (i+1)+number_of_matches); j++ {
			cards[j].copies += cards[i].copies
		}
	}

	var result = 0

	for _, card := range cards {
		result += card.copies
	}

	fmt.Println(cards)
	fmt.Println(result)

}
