package main

import (
	"fmt"
	"github.com/DomBlack/advent-of-code-2018/lib"
	"strings"
)

func main() {
	input := lib.InputAsString("day-05")

	fmt.Println("Part 1:", part1(input))
	fmt.Println("Part 2:", part2(input))
}

// Loop through the string and remove any matching neighbours which are different cases of each other and return the
// length of the final string
// e.g. "aBcCba" becomes "aBba" and then "aa" which has length 2
func part1(input string) int {
	for i := 0; i < len(input)-1; i++ {
		first := string(input[i])
		second := string(input[i+1])

		// If they are a match by opposite polarity
		if strings.ToLower(first) == strings.ToLower(second) && first != second {
			// remove the matching pair
			input = input[:i] + input[i+2:]
			i -= 2

			if i < 1 {
				i = -1
			}
		}
	}

	return len(input)
}

// Removes one pair of characters such as "a/A" such that part1 would give the
// shortest answer.
func part2(input string) int {
	options := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}

	shortest := len(input)

	for _, char := range options {
		if strings.Contains(input, char) {
			shorterInput := strings.Replace(
				strings.Replace(input, char, "", -1),
				strings.ToUpper(char),
				"",
				-1,
			)

			result := part1(shorterInput)

			if result < shortest {
				shortest = result
			}
		}
	}

	return shortest
}
