package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

type Mapping struct {
	dest_range   [2]int
	source_range [2]int
}

type MappingCategory struct {
	name     string
	mappings []Mapping
}

func stringSliceToIntSlice(ss *[]string) ([]int, error) {
	var result []int

	for _, char := range *ss {
		var num, err = strconv.Atoi(char)
		if err != nil {
			return result, err
		}
		result = append(result, num)
	}

	return result, nil
}

func getFirstRune(s *string) rune {
	var first rune
	for _, c := range *s {
		first = c
		break
	}
	return first
}

func processInput() ([][2]int, []MappingCategory, error) {
	var seed_ranges [][2]int
	var mappingCategories []MappingCategory

	file, err := os.Open("input.txt")

	if err != nil {
		return seed_ranges, mappingCategories, err
	}

	fileScanner := bufio.NewScanner(file)

	// read first line
	fileScanner.Scan()
	line := fileScanner.Text()

	split := strings.Split(line, ":")
	seeds_strings := strings.Fields(split[1])

	seeds, err := stringSliceToIntSlice(&seeds_strings)

	if err != nil {
		return seed_ranges, mappingCategories, err
	}

	for i := 0; i < len(seeds); i += 2 {
		seed_ranges = append(seed_ranges, [2]int{seeds[i], seeds[i] + seeds[i+1]})
	}

	for fileScanner.Scan() {
		line := fileScanner.Text()

		if len(line) == 0 {
			continue
		}

		if unicode.IsDigit(getFirstRune(&line)) {
			mapping_strings := strings.Fields(line)
			mapping_nums, err := stringSliceToIntSlice(&mapping_strings)

			if err != nil {
				return seed_ranges, mappingCategories, err
			}

			dest_start := mapping_nums[0]
			source_start := mapping_nums[1]
			length := mapping_nums[2]

			mapping := Mapping{[2]int{dest_start, dest_start + length}, [2]int{source_start, source_start + length}}

			mapping_index := len(mappingCategories) - 1
			mappingCategories[mapping_index].mappings = append(mappingCategories[mapping_index].mappings, mapping)

		} else {
			split := strings.Fields(line)
			mappingCategories = append(mappingCategories, MappingCategory{split[0], []Mapping{}})
		}

	}

	for idx := range mappingCategories {
		sort.Slice(mappingCategories[idx].mappings, func(i, j int) bool {
			return mappingCategories[idx].mappings[i].source_range[0] < mappingCategories[idx].mappings[j].source_range[0]
		})
	}

	return seed_ranges, mappingCategories, nil
}

func transFormInterval(interval [2]int, source [2]int, dest [2]int) ([2]int, error) {

	if !(interval[0] >= source[0] && interval[1] <= source[1]) {
		return [2]int{}, fmt.Errorf("Bad interval transformation: Interval %v Source: %v Destination: %v", interval, source, dest)
	}

	intervalLength := interval[1] - interval[0]

	newLeft := dest[0] + (interval[0] - source[0])
	newRight := newLeft + intervalLength

	return [2]int{newLeft, newRight}, nil
}

type IntervallMapping struct {
	interval    [2]int
	source      [2]int
	destination [2]int
}

func applyMapping(seeds *[][2]int, mc MappingCategory) [][2]int {
	var result [][2]int

	for _, seed_range := range *seeds {
		var seed_interval_mappings []IntervallMapping

		for _, mapping_range := range mc.mappings {
			// too right
			if mapping_range.source_range[0] > seed_range[1] {
				break
			}

			// too left
			if mapping_range.source_range[1] < seed_range[0] {
				continue
			}

			// Calculate the overlapping part of the seed range with the mapping
			partial_intervall_source := [2]int{
				max(seed_range[0], mapping_range.source_range[0]),
				min(seed_range[1], mapping_range.source_range[1]),
			}

			sim := IntervallMapping{
				interval:    partial_intervall_source,
				source:      mapping_range.source_range,
				destination: mapping_range.dest_range,
			}

			seed_interval_mappings = append(seed_interval_mappings, sim)
		}

		if len(seed_interval_mappings) == 0 {
			// No mappings found, return the seed range as is
			result = append(result, seed_range)
			continue
		}

		// Calculate untransformed intervals
		last_end := seed_range[0]

		for _, seed_interval_mapping := range seed_interval_mappings {
			// Check for gaps between intervals
			if last_end < seed_interval_mapping.interval[0] {
				// There's a gap between the last mapped interval and this one
				// Add the unmapped interval as is
				result = append(result, [2]int{last_end, seed_interval_mapping.interval[0] - 1})
			}

			// Transform the current interval
			mapped_interval, err := transFormInterval(seed_interval_mapping.interval, seed_interval_mapping.source, seed_interval_mapping.destination)
			if err != nil {
				log.Fatal(err)
			}

			result = append(result, mapped_interval)

			// Update the last_end to the end of the current interval
			last_end = seed_interval_mapping.interval[1] + 1
		}

		// If there's still a portion left at the end, add it as is
		if last_end <= seed_range[1] {
			result = append(result, [2]int{last_end, seed_range[1]})
		}

		fmt.Println("Seed ", seed_range, mc, seed_interval_mappings, result)
	}

	return result
}

func main() {

	var seeds, mappings, err = processInput()

	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(applyMapping(&seeds, mappings[0]))

	var result [][2]int = seeds

	for _, mapping := range mappings {
		result = applyMapping(&result, mapping)
	}

	var min_result int = result[0][0]

	for _, res := range result {
		min_result = min(min_result, res[0])
	}

	fmt.Println(min_result)
}
