package main

import (
	"container/ring"
	"fmt"
)

func main() {
	fmt.Println("Part 1: ", part1(465, 71940))
	fmt.Println("Part 2: ", part1(465, 71940 * 100))
}

func part1(numPlayers int, lastMarble int) int {
	scores := make([]int, numPlayers)
	highScore := 0
	currentPlayer := 0

	circle := ring.New(1)
	circle.Value = 0

	for currentMarble := 1; currentMarble <= lastMarble; currentMarble++ {
		if currentMarble % 23 == 0 {
			// If the current marble is a multiple of 23, then they get the score of that instead
			scores[currentPlayer] += currentMarble

			// Plus they remove the marble 7 positions counter-clockwise and get that score too
			circle = circle.Move(-8)

			scores[currentPlayer] += circle.Unlink(1).Value.(int)
			circle = circle.Next()

			// Update high score
			if scores[currentPlayer] > highScore {
				highScore = scores[currentPlayer]
			}
		} else {
			// Place the new marble 1-2 marbles clockwise from the current position
			circle = circle.Move(2)
			entry := ring.New(1)
			entry.Value = currentMarble

			// Append the new marble at the new position
			circle = entry.Link(circle)
		}

		// Increment the Player
		currentPlayer++
		if currentPlayer >= numPlayers {
			currentPlayer = 0
		}
	}

	return highScore
}
