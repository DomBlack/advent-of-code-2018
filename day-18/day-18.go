package main

import (
	"fmt"
	"github.com/DomBlack/advent-of-code-2018/lib"
	"github.com/DomBlack/advent-of-code-2018/lib/algos"
	"strings"
)

func main() {
	input := lib.InputAsString("day-18")

	fmt.Println("Part 1", part1(input))
	fmt.Println("Part 2", part2(input))
}

func part1(input string) int {
	c := NewCollectionArea(input)

	for i := 0; i < 10; i++ {
		c.Tick()
	}

	return c.TotalResourceValue()
}

func part2(input string) int {
	const target = 1000000000

	c := NewCollectionArea(input)

	cycleLength, cycleStart := algos.FloydCycleDetection(c)

	loopCount := ((target - cycleStart) % cycleLength) + cycleStart
	for i := 0; i < loopCount; i++ {
		c.Tick()
	}

	return c.TotalResourceValue()
}

type AcreType int

const (
	OpenGround AcreType = iota
	Trees
	Lumberyard
)

type CollectionArea struct {
	acres [][]AcreType
}

func NewCollectionArea(str string) (res *CollectionArea) {
	str = strings.TrimSpace(str)
	width := strings.Index(str, "\n")
	if width == 0 {
		panic("No line ending found")
	}

	res = &CollectionArea{
		make([][]AcreType, 0),
	}

	res.acres = append(res.acres, make([]AcreType, 0))
	var y, x int

	for _, char := range str {
		switch char {
		case '\r':
			// no-op
		case '\n':
			res.acres = append(res.acres, make([]AcreType, 0))
			x = -1
			y++
		case '.':
			res.acres[y] = append(res.acres[y], OpenGround)
		case '|':
			res.acres[y] = append(res.acres[y], Trees)
		case '#':
			res.acres[y] = append(res.acres[y], Lumberyard)
		default:
			panic(fmt.Sprintf("Unknown char `%v` at row %d col %d", char, y, x))
		}

		x++
	}

	return
}

func (c CollectionArea) String() string {
	var str strings.Builder

	for _, row := range c.acres {
		str.WriteRune('\n')

		for _, acre := range row {
			switch acre {
			case OpenGround:
				str.WriteRune('.')
			case Trees:
				str.WriteRune('|')
			case Lumberyard:
				str.WriteRune('#')
			default:
				panic("Unknown acre type")
			}
		}
	}

	return str.String()[1:]
}

func (c *CollectionArea) Copy() (res *CollectionArea) {
	res = &CollectionArea{make([][]AcreType, len(c.acres))}

	for y, row := range c.acres {
		res.acres[y] = make([]AcreType, len(row))
		copy(res.acres[y], row)
	}

	return
}

func (c *CollectionArea) CopyTickable() algos.Tickable {
	return c.Copy()
}

func (c *CollectionArea) Tick() {
	previous := c.Copy()

	for y, row := range c.acres {
		for x, cellType := range row {
			trees, lumberyards := previous.adjacentCounts(y, x)

			switch cellType {
			case OpenGround:
				if trees >= 3 {
					cellType = Trees
				}
			case Trees:
				if lumberyards >= 3 {
					cellType = Lumberyard
				}
			case Lumberyard:
				if trees == 0 || lumberyards == 0 {
					cellType = OpenGround
				}
			}

			c.acres[y][x] = cellType
		}
	}
}

func (c *CollectionArea) adjacentCounts(y, x int) (trees, lumberyards int) {
	count := func(y, x int) {
		if y >= 0 && y < len(c.acres) {
			if x >= 0 && x < len(c.acres[y]) {
				switch c.acres[y][x] {
				case Trees:
					trees++
				case Lumberyard:
					lumberyards++
				}
			}
		}
	}

	count(y-1, x-1)
	count(y-1, x)
	count(y-1, x+1)
	count(y, x-1)
	count(y, x+1)
	count(y+1, x-1)
	count(y+1, x)
	count(y+1, x+1)

	return
}

func (c CollectionArea) TotalResourceValue() int {
	trees, lumberyards := 0, 0

	for _, row := range c.acres {
		for _, cell := range row {
			switch cell {
			case Trees:
				trees++
			case Lumberyard:
				lumberyards++
			}
		}
	}

	return trees * lumberyards
}
