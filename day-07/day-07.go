package main

import (
	"fmt"
	"github.com/DomBlack/advent-of-code-2018/day-06/graph"
	"github.com/DomBlack/advent-of-code-2018/lib"
	"log"
	"sort"
)

func main() {
	strInput := lib.InputAsStrings("day-07")
	input := GraphFromStrings(strInput)

	fmt.Println("Part 1:", part1(input))

	_, clock := part2(input, 60, 5)
	fmt.Println("Part 1:", clock)
}

func GraphFromStrings(input []string) []*graph.Node {
	// All known nodes
	nodes := make(map[string]*graph.Node)

	getNode := func(id string) *graph.Node {
		node, found := nodes[id]
		if !found {
			newNode := graph.NewNode(id)
			node = &newNode
			nodes[id] = node
		}

		return node
	}

	// Read the input
	for _, line := range input {
		var firstID, secondID string

		num, err := fmt.Sscanf(line, "Step %s must be finished before step %s can begin.", &firstID, &secondID)

		if err != nil || num != 2 {
			log.Fatal("Unable to parse ", line)
		}

		first := getNode(firstID)
		second := getNode(secondID)

		first.AddEdgeTo(second)
	}

	// Find the starting node
	entryNodes := make([]*graph.Node, 0)
	for _, node := range nodes {
		if len(node.Inbound) == 0 {
			entryNodes = append(entryNodes, node)
		}
	}

	return entryNodes
}

func sortToVisitSlice(toVisit []*graph.Node) {
	sort.Slice(toVisit, func(i, j int) bool {
		return toVisit[i].ID < toVisit[j].ID
	})
}

func part1(entry []*graph.Node) string {
	order := ""
	toVisit := entry
	visited := make(map[string] struct{})

	sortToVisitSlice(toVisit)

	canVisit := func(node *graph.Node) bool {
		for _, edge := range node.Inbound {
			if _, found := visited[edge.Start.ID]; !found {
				return false
			}
		}

		return true
	}

	for len(toVisit) > 0 {
		// visit the node
		node := toVisit[0]
		toVisit = toVisit[1:]
		visited[node.ID] = struct{}{}
		order += node.ID

		// check which children can be added
		for _, edge := range node.Outbound {
			if canVisit(edge.End) {
				toVisit = append(toVisit, edge.End)
			}
		}

		// Resort the visit list based on alphabetical order
		sortToVisitSlice(toVisit)
	}

	return order
}

func part2(entry []*graph.Node, baseTime int, numWorkers int) (string, int) {
	order := ""
	toVisit := entry
	visited := make(map[string] struct{})

	canVisit := func(node *graph.Node) bool {
		for _, edge := range node.Inbound {
			if _, found := visited[edge.Start.ID]; !found {
				return false
			}
		}

		return true
	}

	clock := 0

	type worker struct {
		busyUntil int
		node *graph.Node
	}

	// A map of worker ID to when it's not busy
	workers := make(map[int]*worker)
	for i := 0; i < numWorkers; i++ {
		workers[i] = &worker {}
	}

	// Is a worker busy right now?
	isAWorkerBusy := func() bool {
		for _, worker := range workers {
			if worker.busyUntil >= clock {
				return true
			}
		}

		return false
	}

	sortToVisitSlice(toVisit)

	// Have we got nodes to visit or is a worker currently busy?
	for len(toVisit) > 0 || isAWorkerBusy() {

		for _, worker := range workers {
			// If this worker is busy skip over it
			if worker.busyUntil > clock {
				continue
			}

			// If this worker finished this clock then
			if worker.busyUntil == clock && worker.node != nil {
				visited[worker.node.ID] = struct{}{}
				order += worker.node.ID

				// check which children can be added
				for _, edge := range worker.node.Outbound {
					if canVisit(edge.End) {
						toVisit = append(toVisit, edge.End)
					}
				}

				// Resort the visit list based on alphabetical order
				sortToVisitSlice(toVisit)
			}

			if len(toVisit) > 0 {
				// visit the node
				worker.node = toVisit[0]
				toVisit = toVisit[1:]

				// The node "A" takes 1 second + baseTime
				time := int(worker.node.ID[0]) - 64
				worker.busyUntil = clock + baseTime + time
			}
		}

		clock++
	}

	return order, clock - 1
}

