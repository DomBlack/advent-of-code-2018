package main

import (
	"fmt"
	"github.com/DomBlack/advent-of-code-2018/lib"
	"log"
)


func main() {
	strInput := lib.InputAsStrings("day-12")
	initialState, rules := parseInput(strInput)
	fmt.Println("Part 1:", CalculatePotSum(initialState, rules, 20))
	fmt.Println("Part 2:", CalculatePotSum(initialState, rules, 50000000000))
}

// The Maximun number of rules we can have (5 bits = 32)
const MaxRuleIndexMask = 31

// A fixed array representing the rules
type Rules [MaxRuleIndexMask + 1]bool

// Parse the input into a current state and the rules for generational change
func parseInput(input []string) (state []bool, rules Rules) {
	// Read the initial state
	var initialState string
	num, err := fmt.Sscanf(input[0], "initial state: %s", &initialState)
	if num != 1 || err != nil {
		log.Fatal("Unable to read initial state")
	}

	state = make([]bool, len(initialState))
	for i, char := range initialState {
		if char == '#' {
			state[i] = true
		} else if char != '.' {
			log.Fatal("Unknown initial state char: ", char)
		}
	}

	// Read the rules
	rules = Rules{}
	input = input[2:]

	for _, line := range input {
		var maskStr, result string
		num, err := fmt.Sscanf(line, "%s => %s", &maskStr, &result)
		if num != 2 || err != nil {
			log.Fatal("Unable to read rule: ", line)
		}

		// Convert the rule string into it's bitwise index
		ruleIndex := 0
		for _, char := range maskStr {
			if char == '#' {
				ruleIndex += 1
			}

			ruleIndex = ruleIndex << 1
		}
		ruleIndex = ruleIndex >> 1

		// Set the rule within our rules map
		if result == "#" {
			rules[ruleIndex] = true
		} else if result == "." {
			rules[ruleIndex] = false
		} else {
			log.Fatal("Unknown result: ", result)
		}
	}

	return
}

func CalculatePotSum(state []bool, rules Rules, generations int) int {
	leftIndex := 0
	lastSum := 0
	lastDelta := 0
	stableDeltaCount := 0

	for gen := 0; gen < generations; gen++ {
		// Check if pots before index 0 are going to sprot plants
		secondLeft := rules[boolToInt(state[0])]
		nearPots := (boolToInt(state[0]) << 1) + boolToInt(state[1])
		firstLeft := rules[nearPots]

		loopOffset := 0

		if secondLeft {
			loopOffset = 2
			state = append([]bool {secondLeft, firstLeft}, state...)
		} else if firstLeft {
			loopOffset = 1
			state = append([]bool {firstLeft}, state...)
		}
		leftIndex -= loopOffset

		for i := loopOffset; i < len(state)-2; i++ {
			nearPots = ((nearPots << 1) + boolToInt(state[i+2])) & MaxRuleIndexMask
			state[i] = rules[nearPots]
		}

		// We're adding zeros for the last 2 right edges
		for i := len(state) - 2; i < len(state); i++ {
			nearPots = nearPots << 1 & MaxRuleIndexMask

			state[i] = rules[nearPots]
		}


		// Check if the final two pots cause more pots to the right to sprot plants
		firstRight := rules[nearPots << 1 & MaxRuleIndexMask]
		secondRight := rules[nearPots << 2 & MaxRuleIndexMask]
		if secondRight {
			state = append(state, firstRight, secondRight)
		} else if firstRight {
			state = append(state, firstRight)
		}

		// See if we can find a pattern and quit early
		sum := sumState(leftIndex, state)
		delta := sum - lastSum

		// If the delta between generations is not changing and hasn't changed for a while
		if lastDelta == delta {
			stableDeltaCount++

			if stableDeltaCount > 5 {
				// Then we can assume it's stable for ever and work out the answer
				// at the requested generation
				return sum + (delta * (generations - gen - 1))
			}
		} else {
			stableDeltaCount = 0
		}

		lastSum = sum
		lastDelta = delta
	}

	// Calculate the sum of all the pots with plants now
	return sumState(leftIndex, state)
}

func boolToInt(b bool) int {
	if b {
		return 1
	} else {
		return 0
	}
}

// Prints a generation to screen
func displayGeneration(state []bool) {
	for _, pot := range state {
		if pot {
			fmt.Print("#")
		} else {
			fmt.Print(".")
		}
	}
	fmt.Println()
}

func sumState(leftIndex int, state []bool) int {
	sum := 0
	for index, pot := range state {
		if pot {
			sum += index + leftIndex
		}
	}

	return sum
}
