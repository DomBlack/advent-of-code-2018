package particle

import (
	"fmt"
	"github.com/DomBlack/advent-of-code-2018/lib/vectors"
	"log"
	"strings"
)

type Particle struct {
	Position, Velocity vectors.Vec2
}

// Step time
func (particle *Particle) Step() {
	particle.Position = particle.Position.Add(particle.Velocity)
}

// Creates a new line from a string formatted like this:
// "position=<-6, 10> velocity=< 2, -2>"
func New(line string) (res Particle) {
	res.Position = vectors.Vec2{}
	res.Velocity = vectors.Vec2{}

	// Parse the string (removing whitespace first due to variable lengths)
	num, err := fmt.Sscanf(
		strings.Replace(line, " ", "", -1),
		"position=<%d,%d>velocity=<%d,%d>",
		&res.Position.X, &res.Position.Y, &res.Velocity.X, &res.Velocity.Y,
	)

	if err != nil {
		log.Fatal(err)
	}

	if num != 4 {
		log.Fatal("Expected 4 inputs, got ", num)
	}

	return
}
