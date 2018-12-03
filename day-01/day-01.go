package main

import (
	"../lib"
	"fmt"
)

func main() {
	input := lib.InputAsIntegers("day-01")

	fmt.Println("Part 1:", part1(input))
	fmt.Println("Part 2:", part2(input))
}

// Works out the sum of all the lines in the given file
func part1(input []int) int {
	result := 0

	for _, i := range input {
		result += i
	}

	return result
}

// Loops constantly over the scanner summing the lines one by one
// until it finds the first time the sum is the same as it's had before
func part2(input []int) int {
	result := 0

	set := make(map[int]bool)
	set[0] = true

	for {
		for _, i := range input {
			result += i

			// Check if we've see this before
			_, found := set[result]
			if found {
				return result
			}

			// otherwise mark it as found
			set[result] = true
		}
	}
}

