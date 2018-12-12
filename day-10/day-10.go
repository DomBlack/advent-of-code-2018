package main

import (
	"fmt"
	"github.com/DomBlack/advent-of-code-2018/day-10/particle"
	"github.com/DomBlack/advent-of-code-2018/lib"
	"log"
	"os"
)

func main() {
	inputStr := lib.InputAsStrings("day-10")
	input := make([]*particle.Particle, len(inputStr))

	for i, line := range inputStr {
		p := particle.New(line)
		input[i] = &p
	}

	part1(input)
}

func part1(input []*particle.Particle) {
	const maxBounds = 80
	time := 0

	for {
		_, _, width, height := getBounds(input)

		// Only display the current state if we are small enough
		if width < maxBounds && height < maxBounds {
			// Display the current state
			display(input)

			// check if the user wants to stop
			fmt.Println("Time is", time)
			fmt.Println("Press q to quit, or anything else to increment time...")
			b := make([]byte, 1)
			_, err := os.Stdin.Read(b)
			if err != nil {
				log.Fatal(err)
			}

			if b[0] == 'q' {
				break
			}
		}

		// Increment time in the particles
		for _, p := range input {
			p.Step()
		}
		time++
	}
}

func getBounds(input []*particle.Particle) (int, int, int, int) {
	// Find the bounds
	topLeft := input[0].Position
	bottomRight := input[0].Position

	for _, p := range input {
		topLeft = topLeft.Min(p.Position)
		bottomRight = bottomRight.Max(p.Position)
	}

	x := -topLeft.X
	y := -topLeft.Y
	width := bottomRight.X + x
	height := bottomRight.Y + y

	return x, y, width + 1, height + 1
}

func display(particles []*particle.Particle) {
	fmt.Println("\033[2J")

	offsetX, offsetY, width, height := getBounds(particles)

	// Init the "display"
	count := width * height
	points := make([]string, count)
	for i := 0; i < count; i++ {
		points[i] = "."
	}

	// Position the particles
	for _, p := range particles {
		x := p.Position.X + offsetX
		y := p.Position.Y + offsetY
		i := (y * width) + x

		points[i] = "#"
	}

	// Print the display
	for i := 0; i < count; i++ {
		fmt.Print(points[i])

		if i % width == width - 1 {
			fmt.Println()
		}
	}
}
