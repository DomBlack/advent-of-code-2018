package main

import (
	"fmt"
	"github.com/DomBlack/advent-of-code-2018/lib"
)

func main() {
	input := lib.InputAsStrings("day-02")

	fmt.Println("Part 1:", part1(input))
	fmt.Println("Part 2:", part2(input))
}

func part1(boxes []string) int {
	numWithPairs := 0
	numWithTriplets := 0

	for _, id := range boxes {
		pairs, triplets := countPairsAndTriplets(id)

		if pairs > 0 {
			numWithPairs++
		}

		if triplets > 0 {
			numWithTriplets++
		}
	}

	return numWithPairs * numWithTriplets
}

// In the boxes array, there will only be two boxes which have ID that are almost
// identical apart from one character. This returns the identical characters
func part2(boxes []string) string {
	for i, boxA := range boxes {
		for j := i + 1; j < len(boxes); j++ {
			boxB := boxes[j]

			same, matches := isSame(boxA, boxB)

			if same {
				return matches
			}
		}
	}

	return ""
}

func countPairsAndTriplets(input string) (int, int) {
	// First count all the number of times a character appears
	counts := make(map[int32]int)

	for _, char := range input {
		count, _ := counts[char]

		count++

		counts[char] = count
	}

	// Then count how many appear exactly twice or three times
	pairs := 0
	triplets := 0

	for _, count := range counts {
		switch count {
		case 2:
			pairs++
		case 3:
			triplets++
		}
	}

	return pairs, triplets
}

func isSame(a string, b string) (bool, string) {
	if len(a) != len(b) {
		return false, ""
	}

	diffences := 0
	sameChars := make([]rune, 0)

	runeA := []rune(a)
	runeB := []rune(b)

	for i, charA := range runeA {
		charB := runeB[i]

		if charA != charB {
			diffences++

			if diffences > 1 {
				return false, ""
			}
		} else {
			sameChars = append(sameChars, charA)
		}
	}

	if diffences == 1 {
		return true, string(sameChars)
	} else {
		return false, ""
	}
}
