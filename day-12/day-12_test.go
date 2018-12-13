package main

import "testing"

var exampleInput = []string {
	"initial state: #..#.#..##......###...###",
	"",
	"...## => #",
	"..#.. => #",
	".#... => #",
	".#.#. => #",
	".#.## => #",
	".##.. => #",
	".#### => #",
	"#.#.# => #",
	"#.### => #",
	"##.#. => #",
	"##.## => #",
	"###.. => #",
	"###.# => #",
	"####. => #",
}

func Test_CalculatePotSum(t *testing.T) {
	const want = 325
	state, rules := parseInput(exampleInput)

	if got := CalculatePotSum(state, rules, 20); got != want {
		t.Errorf("CalculatePotSum() = %v, want %v", got, want)
	}
}
