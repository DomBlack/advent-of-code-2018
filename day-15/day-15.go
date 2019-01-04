package main

import (
	"fmt"
	"github.com/DomBlack/advent-of-code-2018/day-15/xcom"
	"github.com/DomBlack/advent-of-code-2018/lib"
)

func main() {
	input := lib.InputAsString("day-15")

	fmt.Println("Part 1", part1(input))
	fmt.Println("Part 2", part2(input))
}

func part1(input string) int {
	m := xcom.NewMap(input, 3)
	return m.RunCombatSim()
}

func part2(input string) int {
	elfPower := 4

	for {
		m := xcom.NewMap(input, elfPower)
		score := m.RunCombatSim()

		// Check for any elf deaths
		elfDeaths := false
		for _, unit := range m.Units {
			if unit.IsElf && unit.IsDead() {
				elfDeaths = true
				break
			}
		}

		if !elfDeaths {
			return score
		}

		elfPower++
	}
}
