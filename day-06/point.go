package main

import (
	"fmt"
	"github.com/DomBlack/advent-of-code-2018/lib"
	"log"
)

// A single point
type Point struct{ x, y int }

// Calculates the manhattan distance between two points
func (p1 Point) Distance(p2 Point) int {
	return lib.Abs(p1.x-p2.x) + lib.Abs(p1.y-p2.y)
}

// Returns a point presenting the largest x & y of these two points
func (p1 Point) Max(p2 Point) Point {
	return Point{
		lib.Max(p1.x, p2.x),
		lib.Max(p1.y, p2.y),
	}
}

// Returns a point presenting the smallest x & y of these two points
func (p1 Point) Min(p2 Point) Point {
	return Point{
		lib.Min(p1.x, p2.x),
		lib.Min(p1.y, p2.y),
	}
}

// Creates a single point
func NewPoint(coords string) (res Point) {
	num, err := fmt.Sscanf(coords, "%d, %d", &res.x, &res.y)

	if err != nil {
		log.Fatal(err)
	}

	if num != 2 {
		log.Fatalf("Expected 2, got %d", num)
	}

	return
}

// Create a slice of points from a slice of strings
func PointsFromStrings(coords []string) []Point {
	points := make([]Point, len(coords))

	for i, coord := range coords {
		points[i] = NewPoint(coord)
	}

	return points
}
