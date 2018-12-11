package main

import (
	"fmt"
	"github.com/DomBlack/advent-of-code-2018/lib"
)

func main() {
	input := lib.InputAsIntegers("day-08")

	fmt.Println("Part 1:", part1(input))
	fmt.Println("Part 2:", part2(input))
}

func part1(input []int) int {
	return part1Node(&input)
}

func part1Node(input *[]int) int {
	sum := 0

	numChildren := (*input)[0]
	numMeta := (*input)[1]

	*input = (*input)[2:]

	// Get the meta from the children
	for i := 0; i < numChildren; i++ {
		sum += part1Node(input)
	}

	// Sum the meta up
	for i := 0; i < numMeta; i++ {
		sum += (*input)[i]
	}
	*input = (*input)[numMeta:]

	return sum
}

func part2(input []int) int {
	return part2Node(&input)
}

func part2Node(input *[]int) int {
	sum := 0

	numChildren := (*input)[0]
	numMeta := (*input)[1]

	*input = (*input)[2:]

	childValue := make([]int, numChildren)

	// Get the meta from the children
	for i := 0; i < numChildren; i++ {
		childValue[i] = part2Node(input)
	}

	// If no children, then this node's sum is the sum of all it's meta
	if numChildren == 0 {
		for i := 0; i < numMeta; i++ {
			sum += (*input)[i]
		}
	} else {
		// Else this nodes value is the sum of the children referenced by the meta
		for i := 0; i < numMeta; i++ {
			childNum := (*input)[i]

			if childNum <= numChildren && childNum > 0 {
				sum += childValue[childNum - 1]
			}
		}
	}
	*input = (*input)[numMeta:]

	return sum
}
