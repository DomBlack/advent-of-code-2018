package main

import "testing"

func Test_part1(t *testing.T) {
	type args struct {
		numPlayers int
		lastMarble int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"Example 1", args { 9, 25 }, 32 },
		{"Example 2", args { 10, 1618 }, 8317 },
		{"Example 3", args { 13, 7999 }, 146373 },
		{"Example 4", args { 17, 1104 }, 2764 },
		{"Example 5", args { 21, 6111 }, 54718 },
		{"Example 6", args { 30, 5807 }, 37305 },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := part1(tt.args.numPlayers, tt.args.lastMarble); got != tt.want {
				t.Errorf("part1() = %v, want %v", got, tt.want)
			}
		})
	}
}
