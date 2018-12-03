package main

import (
	"testing"
)

func Test_part1(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		want  int
	}{
		{"example 1", []int{+1, +1, +1}, 3},
		{"example 2", []int{+1, +1, -2}, 0},
		{"example 3", []int{-1, -2, -3}, -6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := part1(tt.input); got != tt.want {
				t.Errorf("part2() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_part2(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		want  int
	}{
		{"example 1", []int{+1, -1}, 0},
		{"example 2", []int{+3, +3, +4, -2, -4}, 10},
		{"example 3", []int{-6, +3, +8, +5, -6}, 5},
		{"example 4", []int{+7, +7, -2, -7, -4}, 14},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := part2(tt.input); got != tt.want {
				t.Errorf("part2() = %v, want %v", got, tt.want)
			}
		})
	}
}
