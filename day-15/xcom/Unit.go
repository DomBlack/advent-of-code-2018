package xcom

import (
	"github.com/DomBlack/advent-of-code-2018/lib/vectors"
	"strconv"
	"strings"
)

const DefaultHealth = 200

// All adjacent cells in "reading order"
var AdjacentCells = [4]vectors.Vec2{
	{0, -1},
	{-1, 0},
	{1, 0},
	{0, 1},
}

type Unit struct {
	IsElf       bool // Is this unit an elf? If not then it's a goblin
	Health      int
	AttackPower int
	Position    vectors.Vec2
}

// Creates a new goblin at the given position
func NewGoblin(pos vectors.Vec2) *Unit {
	return &Unit{
		false,
		DefaultHealth,
		3,
		pos,
	}
}

// Creates a new elf at the given position
func NewElf(pos vectors.Vec2, attackPower int) *Unit {
	return &Unit{
		true,
		DefaultHealth,
		attackPower,
		pos,
	}
}

// Finds all enemy units on the map
func (u *Unit) FindTargets(m *Map) Units {
	targets := make(Units, 0)

	for _, possibleTarget := range m.Units {
		if possibleTarget.IsElf != u.IsElf && !possibleTarget.IsDead() {
			targets = append(targets, possibleTarget)
		}
	}

	return targets
}

// Finds an adjacent target or nil, with the lowest health (in reading order if tied)
func (u *Unit) GetAdjacentTarget(m *Map) *Unit {
	var adjacentUnit *Unit

	for _, dir := range AdjacentCells {
		adjacentPos := u.Position.Add(dir)

		adjacentCell, found := m.Cells[adjacentPos]

		// Is there an enemy in the adjacent cell?
		if found &&
			adjacentCell.Unit != nil &&
			adjacentCell.Unit.IsElf != u.IsElf &&
			!adjacentCell.Unit.IsDead() {
			if adjacentUnit == nil || adjacentCell.Unit.Health < adjacentUnit.Health {
				// No possible target yet, so this unit automatically becomes it
				// Or the just found enemies health is lower than the one already found
				adjacentUnit = adjacentCell.Unit
			}
		}
	}

	return adjacentUnit
}

// Finds all empty adjacent cells
func (u *Unit) GetAdjacentEmptyCells(m *Map) []vectors.Vec2 {
	adjacentCells := make([]vectors.Vec2, 0)

	for _, dir := range AdjacentCells {
		adjacentPos := u.Position.Add(dir)

		adjacentCell, found := m.Cells[adjacentPos]

		if found && adjacentCell.IsEmpty() {
			adjacentCells = append(adjacentCells, adjacentPos)
		}
	}

	return adjacentCells
}

// This unit attacks the "other"
func (u *Unit) Attack(other *Unit) (wasKilled bool) {
	other.Health -= u.AttackPower
	wasKilled = other.IsDead()
	return
}

// Is this unit dead?
func (u *Unit) IsDead() bool {
	return u.Health <= 0
}

// Writes out unit information to the string
func (u Unit) String() string {
	var str strings.Builder

	if u.IsElf {
		str.WriteRune('E')
	} else {
		str.WriteRune('G')
	}

	str.WriteRune('(')
	str.WriteString(strconv.Itoa(u.Health))
	str.WriteRune(')')

	return str.String()
}

// A slice of units (for the reading order sort; top to bottom, then left to right)
type Units []*Unit

func (u Units) Len() int {
	return len(u)
}

func (u Units) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}

func (u Units) Less(i, j int) bool {
	return u[i].Position.IsReadingOrderLess(u[j].Position)
}

// All empty cells adjacent to these units
func (u Units) GetEmptyAdjacentCells(m *Map) []vectors.Vec2 {
	emptyCells := make([]vectors.Vec2, 0)

	for _, unit := range u {
		for _, cell := range unit.GetAdjacentEmptyCells(m) {
			emptyCells = append(emptyCells, cell)
		}
	}

	return emptyCells
}
