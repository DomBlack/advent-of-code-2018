package xcom

import (
	"github.com/DomBlack/advent-of-code-2018/lib/vectors"
	"log"
	"sort"
)

type FloodMap map[vectors.Vec2]int

func (m *Map) NewFloodMap(starting *Unit) FloodMap {
	type FloodCell struct {
		position vectors.Vec2
		cost     int
	}

	toVisit := []FloodCell{{starting.Position, 0}}
	visited := make(map[vectors.Vec2]bool)

	floodMap := make(map[vectors.Vec2]int)

	// While we still have cells to visit
	for len(toVisit) > 0 {
		cell := toVisit[0]
		toVisit = toVisit[1:]

		floodMap[cell.position] = cell.cost

		// Check all adjacent cells
		for _, dir := range AdjacentCells {
			adjacentPos := cell.position.Add(dir)

			_, found := visited[adjacentPos]
			visited[adjacentPos] = true

			// If not visited already
			if !found {

				// If the adjacent cell is empty
				adjacentCell, found := m.Cells[adjacentPos]
				if found && adjacentCell.IsEmpty() {
					// Add the cell to the list of which we need to visit
					toVisit = append(toVisit, FloodCell{adjacentPos, cell.cost + 1})
				}
			}
		}
	}

	// If here all reachable cells are
	return floodMap
}

func (m FloodMap) GetNearestReachableFrom(inRangeCells []vectors.Vec2) *vectors.Vec2 {
	// Sort the slice for reading order, so the first we loop over is the min
	sort.Slice(inRangeCells, func(i, j int) bool {
		return inRangeCells[i].IsReadingOrderLess(inRangeCells[j])
	})

	// Now find the reachable position which is the lowest cost to get to
	var nearest *vectors.Vec2
	var nearestCost = 0
	for _, pos := range inRangeCells {
		cost, found := m[pos]

		if found && (nearest == nil || cost < nearestCost) {
			nearestCost = cost
			nPos := pos
			nearest = &nPos
		}
	}

	return nearest
}

func (m FloodMap) FindNextStepTowards(target *vectors.Vec2) vectors.Vec2 {
	// "If multiple steps would put the unit equally closer to its destination, "
	// "the unit chooses the step which is first in reading order. "

	type Node struct { position vectors.Vec2; cost int }

	cost, found := m[*target]
	if !found {
		log.Fatalf("Initial position %v was not in flood map", target)
	}

	if cost == 1 {
		return *target
	}

	toVisit := []Node{ { *target, cost} }
	visited := make(map[vectors.Vec2]bool)
	var possibleStep *vectors.Vec2

	for len(toVisit) > 0 {
		node := toVisit[0]
		toVisit = toVisit[1:]
		visited[node.position] = true

		for _, dir := range AdjacentCells {
			adjacentPos := node.position.Add(dir)

			if _, found := visited[adjacentPos]; found {
				continue
			}

			cost, found := m[adjacentPos]

			if cost == 1 && (possibleStep == nil || adjacentPos.IsReadingOrderLess(*possibleStep)) {
				possibleStep = &adjacentPos
			} else if found && cost < node.cost && cost > 0 {
				toVisit = append(toVisit, Node { adjacentPos, cost })
			}
		}
	}

	if possibleStep == nil {
		log.Fatalf("Unable to compute next step from %v", target)
	}

	return *possibleStep
}
