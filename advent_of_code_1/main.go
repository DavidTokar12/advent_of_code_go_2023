package main

import (
	"bufio"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"unicode"
)

type TextDigit struct {
	text         []rune
	currentIndex int
	digit        rune
	toReachIndex int
}

func createTextDigit(s string, val rune) TextDigit {
	return TextDigit{[]rune(s), 0, rune(val), len([]rune(s))}
}

var TextDigits = []TextDigit{
	createTextDigit("zero", '0'),
	createTextDigit("one", '1'),
	createTextDigit("two", '2'),
	createTextDigit("three", '3'),
	createTextDigit("four", '4'),
	createTextDigit("five", '5'),
	createTextDigit("six", '6'),
	createTextDigit("seven", '7'),
	createTextDigit("eight", '8'),
	createTextDigit("nine", '9'),
}

var TextDigitsBackward = []TextDigit{
	createTextDigit("orez", '0'),
	createTextDigit("eno", '1'),
	createTextDigit("owt", '2'),
	createTextDigit("eerht", '3'),
	createTextDigit("ruof", '4'),
	createTextDigit("evif", '5'),
	createTextDigit("xis", '6'),
	createTextDigit("neves", '7'),
	createTextDigit("thgie", '8'),
	createTextDigit("enin", '9'),
}

func resetTextDigits() {
	for i := range TextDigits {
		digit := &TextDigits[i]
		digit.currentIndex = 0
	}

	for i := range TextDigitsBackward {
		digit := &TextDigitsBackward[i]
		digit.currentIndex = 0
	}
}

func getFirstDigit(s *string) (rune, error) {

	for _, runeValue := range *s {

		if unicode.IsDigit(runeValue) {
			return runeValue, nil
		}

		for i := range TextDigits {
			digit := &TextDigits[i]

			if digit.text[digit.currentIndex] == runeValue {
				digit.currentIndex += 1
				if digit.currentIndex == digit.toReachIndex {
					return digit.digit, nil
				}
			} else {
				digit.currentIndex = 0
			}

		}

	}

	return rune(0), fmt.Errorf("No intiger character found in string: %s", *s)
}

func getLastDigit(s *string) (rune, error) {

	input := []rune(*s)

	for i := len(input) - 1; i >= 0; i-- {

		runeValue := input[i]
		if unicode.IsDigit(runeValue) {
			return runeValue, nil
		}

		for j := range TextDigitsBackward {
			digit := &TextDigitsBackward[j]

			if digit.text[digit.currentIndex] == runeValue {
				digit.currentIndex++
				if digit.currentIndex == digit.toReachIndex {
					return digit.digit, nil
				}
			} else {
				digit.currentIndex = 0
			}
		}
	}

	return rune(0), fmt.Errorf("no digit found in string: %s", *s)
}

func main() {
	readFile, err := os.Open("input.txt")
	if err != nil {
		log.Fatal("Error opening file:", err)
	}
	defer readFile.Close()

	result := new(big.Int)
	fileScanner := bufio.NewScanner(readFile)
	lineCount := 0

	for fileScanner.Scan() {
		lineCount++
		line := fileScanner.Text()

		resetTextDigits()
		firstDigit, err1 := getFirstDigit(&line)
		lastDigit, err2 := getLastDigit(&line)

		if err1 != nil || err2 != nil {
			log.Fatalf("Error processing line %d: %v, %v", lineCount, err1, err2)
		}

		num := string(firstDigit) + string(lastDigit)
		if len(num) != 2 {
			log.Printf("Warning: Line %d produced invalid number: %s", lineCount, num)
		}

		intNum, err3 := strconv.Atoi(num)
		if err3 != nil {
			log.Fatalf("Error converting number on line %d: %v", lineCount, err3)
		}

		if intNum < 0 {
			log.Printf("Warning: Negative number created on line %d: %d", lineCount, intNum)
		}

		result.Add(result, big.NewInt(int64(intNum)))
		fmt.Printf("Line %d: %s -> %s = %d, Running total: %s\n", lineCount, line, num, intNum, result.String())
	}

	if err := fileScanner.Err(); err != nil {
		log.Fatal("Error reading file:", err)
	}

	fmt.Printf("Processed %d lines\n", lineCount)
	fmt.Printf("Final Result: %s\n", result.String())
}
