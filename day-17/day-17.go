package main

import (
	"fmt"
	"github.com/DomBlack/advent-of-code-2018/lib"
	"github.com/DomBlack/advent-of-code-2018/lib/collections"
	"github.com/DomBlack/advent-of-code-2018/lib/vectors"
	"strings"
)

func main() {
	input := lib.InputAsStrings("day-17")
	r := NewReservoir(input)
	touchedByWater, settledWater := r.RunSimulation()

	fmt.Println("Part 1", touchedByWater)
	fmt.Println("Part 2", settledWater)
}

type CellState int

const (
	WaterSource CellState = iota
	ClayWall
	FlowingWater
	SettledWater
)

type Reservoir struct {
	min, max     vectors.Vec2               // The min and max x/y values of this reservoir
	cells        map[vectors.Vec2]CellState // The cells in the reservoir
	flowingWater collections.Vec2Stack      // Cells with flowing water
}

func NewReservoir(input []string) (r *Reservoir) {
	r = &Reservoir{
		vectors.NewVec2(500, 5000),
		vectors.NewVec2(500, 0),
		make(map[vectors.Vec2]CellState),
		collections.NewVec2Stack(),
	}

	for _, line := range input {
		var axisA, axisB string
		var cordA, cordBStart, cordBEnd int
		num, err := fmt.Sscanf(line, "%1s=%d, %1s=%d..%d", &axisA, &cordA, &axisB, &cordBStart, &cordBEnd)

		if err != nil {
			panic(err)
		}

		if num != 5 {
			panic("expected 5 fields")
		}

		for b := cordBStart; b <= cordBEnd; b++ {
			var vec vectors.Vec2
			if axisA == "x" {
				vec = vectors.NewVec2(cordA, b)
			} else {
				vec = vectors.NewVec2(b, cordA)
			}

			r.cells[vec] = ClayWall
			r.min = r.min.Min(vec)
			r.max = r.max.Max(vec)
		}
	}

	// Water can flow around the edges of here
	r.min.X--
	r.max.X++

	// We ignore anything with a Y smaller than our minY, so we'll place the source there
	waterSource := vectors.NewVec2(500, r.min.Y-1)
	r.min.Y--
	r.cells[waterSource] = WaterSource
	r.flowingWater.Push(waterSource)

	return
}

var Down = vectors.NewVec2(0, 1)
var Left = vectors.NewVec2(-1, 0)
var Right = vectors.NewVec2(1, 0)

func (r *Reservoir) RunSimulation() (touchedByWater int, settledWater int) {
	for !r.flowingWater.IsEmpty() {
		r.FlowUntilSettledCreated()
	}

	for pos, cell := range r.cells {
		if pos.Y > r.max.Y {
			continue
		}

		if cell == FlowingWater {
			touchedByWater++
		} else if cell == SettledWater {
			settledWater++
			touchedByWater++
		}
	}

	return
}

// Keeps flowing the water until some settled water is created
func (r *Reservoir) FlowUntilSettledCreated() {
	// While there is still flowing water left
	for !r.flowingWater.IsEmpty() {
		source := r.flowingWater.Pop()

		below := source.Add(Down)

		if r.CanFlowTo(below) && below.Y <= r.max.Y {
			r.flowingWater.Push(source) // Re-add this as once the flowing water has settled this might change
			r.AddFlowingWaterTo(below)
		} else if r.IsWallOrSettled(below) {
			floodCells := collections.NewVec2Stack()

			// Perform the flood
			floodCells.Push(source)
			willSettleLeft := r.PreformFlood(&floodCells, source, Left)
			willSettleRight := r.PreformFlood(&floodCells, source, Right)

			if willSettleLeft && willSettleRight {
				for !floodCells.IsEmpty() {
					cell := floodCells.Pop()
					r.cells[cell] = SettledWater
				}
				return
			} else {
				for !floodCells.IsEmpty() {
					cell := floodCells.Pop()
					r.cells[cell] = FlowingWater

					// Only add this flowing water to the system if it can can flow downwards
					below := cell.Add(Down)
					if r.CanFlowTo(below) {
						r.flowingWater.Push(cell)
					}
				}
			}
		}
	}
}

// Checks in the given direction if we will settle or flow
func (r *Reservoir) PreformFlood(floodCells *collections.Vec2Stack, pos, direction vectors.Vec2) (willSettle bool) {
	willSettle = true

	for {
		pos = pos.Add(direction)

		if r.IsWallOrSettled(pos) {
			break
		}

		floodCells.Push(pos)
		below := pos.Add(Down)

		if !r.IsWallOrSettled(below) || pos.X >= r.max.X || pos.X <= r.min.X {
			willSettle = false
			break
		}
	}

	return
}

func (r *Reservoir) AddFlowingWaterTo(pos vectors.Vec2) {
	r.cells[pos] = FlowingWater
	r.flowingWater.Push(pos)
}

func (r *Reservoir) CanFlowTo(pos vectors.Vec2) bool {
	_, found := r.cells[pos]
	return !found
}

func (r *Reservoir) IsWallOrSettled(pos vectors.Vec2) bool {
	cell, found := r.cells[pos]

	return found && (cell == ClayWall || cell == SettledWater)
}

// Converts the reservoir into a displayable form
func (r Reservoir) String() string {
	var str strings.Builder

	for y := r.min.Y; y <= r.max.Y; y++ {
		if y > r.min.Y {
			str.WriteRune('\n')
		}

		for x := r.min.X; x <= r.max.X; x++ {
			cell, found := r.cells[vectors.NewVec2(x, y)]

			if !found {
				str.WriteRune('.')
			} else {
				switch cell {
				case WaterSource:
					str.WriteRune('+')
				case ClayWall:
					str.WriteRune('#')
				case FlowingWater:
					str.WriteRune('|')
				case SettledWater:
					str.WriteRune('~')
				}
			}
		}
	}

	return str.String()
}
