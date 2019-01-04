package main

import "testing"

func Test_part2(t *testing.T) {
	tests := []struct {
		name     string
		want     int
		inputMap string
	}{
		{
			"First Summarized Example", 4988,
			`#######
#.G...#
#...EG#
#.#.#G#
#..G#E#
#.....#
#######`,
		},
		{
			"Second Example", 31284,
			`#######
#E..EG#
#.#G.E#
#E.##E#
#G..#.#
#..E#.#
#######`,
		},
		{
			"Third Example", 3478,
			`#######
#E.G#.#
#.#G..#
#G.#.G#
#G..#.#
#...E.#
#######`,
		},
		{
			"Forth Example", 6474,
			`#######
#.E...#
#.#..G#
#.###.#
#E#G#G#
#...#G#
#######`,
		},
		{
			"Last Example", 1140,
			`#########
#G......#
#.E.#...#
#..##..G#
#...##..#
#...#...#
#.G...G.#
#.....G.#
#########`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := part2(tt.inputMap); got != tt.want {
				t.Errorf("part2() = %v, want %v", got, tt.want)
			}
		})
	}
}
