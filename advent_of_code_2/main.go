package main

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
	// "fmt"
	"log"
	"os"
)

type CubeHandfulls struct {
	Red   int
	Green int
	Blue  int
}


func processGameId(s *string) (int, error) {
	var id int

	gameIdSlice := strings.Fields(*s)

	if len(gameIdSlice) != 2 {
		return id, fmt.Errorf("No id found in game id: %s", *s)
	}

	id, err := strconv.Atoi(gameIdSlice[1])

	if err != nil {
		return id, err
	}

	return id, nil
}

func processCubeHandfulls(s *string) (CubeHandfulls, error) {
	var cubeHandfulls CubeHandfulls

	cubeHandfullsSlice := strings.Split(*s, ",")

	for _, colorValue := range cubeHandfullsSlice {
		cubeHandfullsSplit := strings.Fields(colorValue)

		if len(cubeHandfullsSplit) != 2 {
			return cubeHandfulls, fmt.Errorf("Couldnt parse cube handfull from string: %s", colorValue)
		}

		value, err := strconv.Atoi(cubeHandfullsSplit[0])
		if err != nil {
			return cubeHandfulls, err
		}

		color := cubeHandfullsSplit[1]

		switch color {
		case "red":
			cubeHandfulls.Red += value
		case "green":
			cubeHandfulls.Green += value
		case "blue":
			cubeHandfulls.Blue += value
		default:
			return cubeHandfulls, fmt.Errorf("Unknown color: %s", color)
		}

	}

	return cubeHandfulls, nil
}

func processLine(s *string) (int, []CubeHandfulls, error) {
	var id int
	var cubeHandfullss []CubeHandfulls

	gameIdGameSlice := strings.Split(*s, ":")

	if len(gameIdGameSlice) != 2 {
		return id, cubeHandfullss, fmt.Errorf("No `:` found in line, or more than one `:` in line: %s", *s)
	}

	id, err := processGameId(&gameIdGameSlice[0])
	if err != nil {
		return id, cubeHandfullss, err
	}

	cubeHandfullsStrings := strings.Split(gameIdGameSlice[1], ";")
	for _, cubeHandfullsString := range cubeHandfullsStrings {
		cubeHandfulls, err := processCubeHandfulls(&cubeHandfullsString)
		if err != nil {
			return id, cubeHandfullss, err
		}

		cubeHandfullss = append(cubeHandfullss, cubeHandfulls)
	}

	return id, cubeHandfullss, nil

}

func main() {
	file, err := os.Open("input.txt")

	if err != nil {
		log.Fatal(err)
	}

	fileScanner := bufio.NewScanner(file)

	var result int

	for fileScanner.Scan() {
		line := fileScanner.Text()

		_, cubeHandfulls, err := processLine(&line)

		if err != nil {
			log.Fatal(err)
		}

		var maxBlue int = 0
		var maxRed int = 0
		var maxGreen int = 0

		for _, cubeHandfull := range cubeHandfulls {
			maxBlue = max(maxBlue, cubeHandfull.Blue)
			maxRed = max(maxRed, cubeHandfull.Red)
			maxGreen = max(maxGreen, cubeHandfull.Green)
		}

		result += (maxBlue * maxRed * maxGreen)
	}

	fmt.Println(result)

}
