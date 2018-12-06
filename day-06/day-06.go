package main

import (
	"fmt"
	"github.com/DomBlack/advent-of-code-2018/lib"
)

func main() {
	strInput := lib.InputAsStrings("day-06")
	input := PointsFromStrings(strInput)

	fmt.Println("Part 1:", part1(input))
}

func part1(input []Point) int {
	topLeft, bottomRight := getBoundingBox(input)

	fmt.Println("Bounding Box", topLeft, bottomRight)
	return 0
}

// Returns the bounding box
func getBoundingBox(points []Point) (Point, Point) {
	topLeft := Point{1000, 1000}
	bottomRight := Point{0, 0}

	for _, point := range points {
		topLeft = topLeft.Min(point)
		bottomRight = bottomRight.Max(point)
	}

	return topLeft, bottomRight
}
