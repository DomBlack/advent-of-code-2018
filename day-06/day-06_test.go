package main

import "testing"

func Test_part1(t *testing.T) {
	tests := []struct {
		name  string
		input []Point
		want  int
	}{
		{
			"Example",
			[]Point{{1, 1}, {1, 6}, {8, 3}, {3, 4}, {5, 5}, {8, 9}},
			17,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := part1(tt.input); got != tt.want {
				t.Errorf("part1() = %v, want %v", got, tt.want)
			}
		})
	}
}
