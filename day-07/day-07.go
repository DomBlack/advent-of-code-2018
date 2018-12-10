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

func part1(entry []*graph.Node) string {
	order := ""
	toVisit := entry
	visited := make(map[string] struct{})

	sort.Slice(toVisit, func(i, j int) bool {
		return toVisit[i].ID < toVisit[j].ID
	})

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
		sort.Slice(toVisit, func(i, j int) bool {
			return toVisit[i].ID < toVisit[j].ID
		})
	}

	return order
}
