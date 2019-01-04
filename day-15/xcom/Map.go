package xcom

import (
	"github.com/DomBlack/advent-of-code-2018/lib/vectors"
	"log"
	"sort"
	"strings"
)

// The Map
type Map struct {
	Cells         map[vectors.Vec2]*Cell // The cells of this map
	Units         Units                  // All the units on this map
	width, height int                    // The width and height of the map
}

// Creates a new map
func NewMap(inputMap string, elfAttackPower int) (res *Map) {
	// Create an empty map structure
	res = &Map{
		make(map[vectors.Vec2]*Cell),
		make(Units, 0),
		0, 0,
	}

	// Parse the map string
	x, y := 0, 0

	for _, r := range inputMap {
		pos := vectors.NewVec2(x, y)

		switch r {
		case '\n':
			x = 0
			y++
			continue
		case '#':
			res.Cells[pos] = &Cell{true, nil}
		case '.':
			res.Cells[pos] = &Cell{false, nil}
		case 'E':
			unit := NewElf(pos, elfAttackPower)
			res.Units = append(res.Units, unit)
			res.Cells[pos] = &Cell{false, unit}
		case 'G':
			unit := NewGoblin(pos)
			res.Units = append(res.Units, unit)
			res.Cells[pos] = &Cell{false, unit}
		default:
			continue
		}

		x++
	}

	res.width = x
	res.height = y + 1

	return
}

// Process a single round of combat
func (m *Map) Round() (combatOver bool) {
	// Units take their turns in the reading order of their starting position;
	// top-to-bottom, left-to-right
	sort.Sort(m.Units)

	for _, unit := range m.Units {
		if unit.IsDead() {
			continue
		}

		// "Each unit begins its turn by identifying all possible targets"
		possibleTargets := unit.FindTargets(m)

		// "If no targets remain, combat ends"
		if len(possibleTargets) == 0 {
			combatOver = true
			return
		}

		// "If the unit is already in range of a target, it does not move"
		if unit.GetAdjacentTarget(m) == nil {
			// "Otherwise, since it is not in range of a target, it moves."

			inRangeCells := possibleTargets.GetEmptyAdjacentCells(m)

			// "if the unit cannot reach (find an open path to) any of the squares that are in range, it ends its turn"
			if len(inRangeCells) == 0 {
				continue
			}

			floodMap := m.NewFloodMap(unit)

			// "If multiple squares are in range and tied for being reachable in the fewest steps, the square which is first in reading order is chosen."
			targetCell := floodMap.GetNearestReachableFrom(inRangeCells)

			// "if the unit cannot reach (find an open path to) any of the squares that are in range, it ends its turn"
			if targetCell == nil {
				continue
			}

			// "The unit then takes a single step toward the chosen square along the shortest path to that square"
			nextPosition := floodMap.FindNextStepTowards(targetCell)

			// Do the actual move
			m.Cells[unit.Position].Unit = nil
			unit.Position = nextPosition
			m.Cells[unit.Position].Unit = unit
		}

		// "After moving (or if the unit began its turn in range of a target), the unit attacks."
		if adjacentTarget := unit.GetAdjacentTarget(m); adjacentTarget != nil {
			if wasKilled := unit.Attack(adjacentTarget); wasKilled {
				m.Cells[adjacentTarget.Position].Unit = nil
			}
		}
	}

	return
}

func (m *Map) RunCombatSim() int {
	numRounds := 0

	for !m.Round() {
		numRounds++
	}

	hitPointsRemaining := 0
	for _, unit := range m.Units {
		if !unit.IsDead() {
			hitPointsRemaining += unit.Health
		}
	}

	return numRounds * hitPointsRemaining
}

func (m Map) DrawMap(floodMap *FloodMap) string {
	var str strings.Builder

	for y := 0; y < m.height; y++ {
		unitsOnRow := make([]*Unit, 0)

		// Display the row
		for x := 0; x < m.width; x++ {
			pos := vectors.NewVec2(x, y)
			cell, found := m.Cells[pos]

			if !found {
				log.Fatalf("Unable to find cell %v", pos)
			}

			if cell.Unit != nil {
				unitsOnRow = append(unitsOnRow, cell.Unit)
			}

			if floodMap != nil && !cell.IsWall {
				// If we have a flood map print that
				cost, found := (*floodMap)[pos]
				if found {
					str.WriteRune(rune('0' + cost))
				} else {
					str.WriteString(cell.String())
				}
			} else {
				str.WriteString(cell.String())
			}
		}

		// Now the row is drawn, display unit health for the row
		str.WriteString("   ")

		for index, unit := range unitsOnRow {
			if index > 0 {
				str.WriteString(", ")
			}

			str.WriteString(unit.String())
		}

		str.WriteRune('\n')
	}

	return str.String()
}

// Convert the map to a string
func (m Map) String() string {
	return m.DrawMap(nil)
}
