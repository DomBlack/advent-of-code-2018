package main

import (
	"fmt"
	"github.com/DomBlack/advent-of-code-2018/lib"
	"log"
	"regexp"
	"strconv"
)

func main() {
	// Read input and convert to pieces
	inputStr := lib.InputAsStrings("day-03")
	var input = make([]Piece, len(inputStr))

	for index, value := range inputStr {
		input[index] = NewPiece(value)
	}

	part1, part2 := day3(input)

	fmt.Println("Part 1:", part1)
	fmt.Println("Part 2:", part2)
}

func day3(pieces []Piece) (int, int) {
	// Find the size of the grid
	height := 0
	width := 0

	for _, piece := range pieces {
		if piece.Bottom() > height {
			height = piece.Bottom()
		}

		if piece.Right() > width {
			width = piece.Right()
		}
	}

	// Now make the grid
	grid := make([][][]int, width)
	for i := range grid {
		grid[i] = make([][]int, height)
	}

	// Fill the grid and look for overlaps
	overlaps := 0

	piecesOverlapped := make([]bool, len(pieces)) // pieces overlapped

	for _, piece := range pieces {
		for _, p := range piece.Positions() {
			if grid[p.x][p.y] == nil {
				grid[p.x][p.y] = make([]int, 0)
			}

			// Count if this square has been overlapped
			if len(grid[p.x][p.y]) == 1 {
				overlaps++
			}

			// Mark any overlapped pieces
			for _, overlappedID := range grid[p.x][p.y] {
				piecesOverlapped[overlappedID - 1] = true
				piecesOverlapped[piece.id - 1] = true
			}

			grid[p.x][p.y] = append(grid[p.x][p.y], piece.id)
		}
	}


	// Find the piece which was not overlapped
	nonOverlappedPiece := 0
	for id, overlapped := range piecesOverlapped {
		if !overlapped {
			nonOverlappedPiece = id + 1
			break
		}
	}

	// Now set all pieces
	return overlaps, nonOverlappedPiece
}

type Point struct {
	x int
	y int
}

// Represents a piece of the fabric
type Piece struct {
	id       int
	position Point
	size     Point
}

// The right edge of a piece
func (p *Piece) Right() int {
	return p.position.x + p.size.x
}

// The bottom edge of a piece
func (p *Piece) Bottom() int {
	return p.position.y + p.size.y
}

// All the square inches on the original fabric which are within this piece
func (p *Piece) Positions() []Point {
	positions := make([]Point, p.size.x*p.size.y)

	i := 0

	for x := 0; x < p.size.x; x++ {
		for y := 0; y < p.size.y; y++ {
			positions[i] = Point{p.position.x + x, p.position.y + y}
			i++
		}
	}

	return positions
}

// Parse a piece
func NewPiece(input string) Piece {
	r, err := regexp.Compile("^#(\\d+) @ (\\d+),(\\d+): (\\d+)x(\\d+)$")
	if err != nil {
		log.Fatal(err)
	}

	parts := r.FindStringSubmatch(input)
	if parts == nil {
		log.Fatal("Unable to parse: ", input)
	}

	return Piece{
		toInt(parts[1]),
		Point{toInt(parts[2]), toInt(parts[3])},
		Point{toInt(parts[4]), toInt(parts[5])},
	}
}

func toInt(what string) int {
	num, err := strconv.Atoi(what)
	if err != nil {
		log.Fatal("Unable to convert to int: ", what)
	}

	return num
}
