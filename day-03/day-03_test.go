package main

import (
	"reflect"
	"testing"
)

func Test_day3(t *testing.T) {
	tests := []struct {
		name   string
		pieces []Piece
		want1  int
		want2  int
	}{
		{
			"example",
			[]Piece{
				{1, Point{1, 3}, Point{4, 4}},
				{2, Point{3, 1}, Point{4, 4}},
				{3, Point{5, 5}, Point{2, 2}},
			},
			4,
			3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1, got2 := day3(tt.pieces)

			if got1 != tt.want1 {
				t.Errorf("part1() = %v, want %v", got1, tt.want1)
			}

			if got2 != tt.want2 {
				t.Errorf("part2() = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func TestNewPiece(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  Piece
	}{
		{"First Example", "#1 @ 1,3: 4x4", Piece{1, Point{1, 3}, Point{4, 4}}},
		{"Second Example", "#2 @ 3,1: 4x4", Piece{2, Point{3, 1}, Point{4, 4}}},
		{"Third Example", "#3 @ 5,5: 2x2", Piece{3, Point{5, 5}, Point{2, 2}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPiece(tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPiece() = %v, want %v", got, tt.want)
			}
		})
	}
}
