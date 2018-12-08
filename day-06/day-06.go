package main

import (
	"fmt"
	"github.com/DomBlack/advent-of-code-2018/lib"
)

func main() {
	strInput := lib.InputAsStrings("day-06")
	input := PointsFromStrings(strInput)

	fmt.Println("Part 1:", part1(input))
	fmt.Println("Part 2:", part2(10000, input))
}

func part1(input []Point) int {
	topLeft, bottomRight := getBoundingBox(input)

	infinite := make(map[Point]bool)
	size     := make(map[Point]int)

	// Loop over the grid
	for x := topLeft.x; x <= bottomRight.x; x++ {
		for y := topLeft.y; y <= bottomRight.y; y++ {
			coord := Point { x, y }
			minCount := 0
			minDistance := 9999
			closetPoint := input[0]

			for _, point := range input {
				distance := coord.Distance(point)

				if distance < minDistance {
					closetPoint = point
					minDistance = distance
					minCount = 1
				} else if distance == minDistance {
					minCount++
				}
			}

			if minCount == 1 {
				size[closetPoint]++

				if x == topLeft.x || x == bottomRight.x || y == topLeft.y || y == bottomRight.y {
					infinite[closetPoint] = true
				}
			}
		}
	}

	maxCount := 0
	for point, count := range size {
		if count > maxCount {
			_, found := infinite[point]

			if !found {
				maxCount = count
			}
		}
	}

	return maxCount
}

func part2(maxDistance int, input []Point) int {
	topLeft, bottomRight := getBoundingBox(input)

	regionSize := 0

	// Loop over the grid
	for x := topLeft.x; x <= bottomRight.x; x++ {
		for y := topLeft.y; y <= bottomRight.y; y++ {
			coord := Point { x, y }
			totalDistance := 0

			for _, point := range input {
				totalDistance += coord.Distance(point)
			}

			if totalDistance < maxDistance {
				regionSize++
			}
		}
	}

	return regionSize
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
