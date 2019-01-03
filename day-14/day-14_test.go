package main

import (
	"strconv"
	"testing"
)

func TestNewScoreBoard(t *testing.T) {
	sb := NewScoreBoard()
	want := "(3)[7]"
	if got := sb.String(); got != want {
		t.Errorf("NewScoreBoard() = %v, want %v", got, want)
	}
}

func TestScoreBoard_CreateRecipe(t *testing.T) {
	tests := []string{
		"(3)[7] 1  0 ",
		" 3  7  1 [0](1) 0 ",
		" 3  7  1  0 [1] 0 (1)",
		"(3) 7  1  0  1  0 [1] 2 ",
		" 3  7  1  0 (1) 0  1  2 [4]",
		" 3  7  1 [0] 1  0 (1) 2  4  5 ",
		" 3  7  1  0 [1] 0  1  2 (4) 5  1 ",
		" 3 (7) 1  0  1  0 [1] 2  4  5  1  5 ",
		" 3  7  1  0  1  0  1  2 [4](5) 1  5  8 ",
		" 3 (7) 1  0  1  0  1  2  4  5  1  5  8 [9]",
		" 3  7  1  0  1  0  1 [2] 4 (5) 1  5  8  9  1  6 ",
		" 3  7  1  0  1  0  1  2  4  5 [1] 5  8  9  1 (6) 7 ",
		" 3  7  1  0 (1) 0  1  2  4  5  1  5 [8] 9  1  6  7  7 ",
		" 3  7 [1] 0  1  0 (1) 2  4  5  1  5  8  9  1  6  7  7  9 ",
		" 3  7  1  0 [1] 0  1  2 (4) 5  1  5  8  9  1  6  7  7  9  2 ",
	}

	sb := NewScoreBoard()

	for index, expected := range tests {
		t.Run("Recipe "+strconv.Itoa(index), func(t *testing.T) {
			sb.CreateRecipe()
			got := sb.String()

			if got != expected {
				t.Errorf("ScoreBoard.CreateRecipe() = %v, want %v", got, expected)
			}
		})
	}
}

func TestScoreBoard_NextTenAfter(t *testing.T) {
	tests := []struct {
		number int
		want   string
	}{
		{ 9, "5158916779" },
		{ 5, "0124515891" },
		{ 18, "9251071085" },
		{ 2018, "5941429882" },
	}

	sb := NewScoreBoard()

	for _, tt := range tests {
		t.Run(strconv.Itoa(tt.number), func(t *testing.T) {
			if got := sb.NextTenAfter(tt.number); got != tt.want {
				t.Errorf("ScoreBoard.NextTenAfter(%v) = %v, want %v", tt.number, got, tt.want)
			}
		})
	}
}

func TestScoreBoard_NumberRecipesBeforeDigits(t *testing.T) {
	tests := []struct {
		number string
		want   int
	}{
		{ "51589", 9 },
		{ "01245", 5 },
		{ "92510", 18 },
		{ "59414", 2018 },
	}

	sb := NewScoreBoard()

	for _, tt := range tests {
		t.Run(tt.number, func(t *testing.T) {
			if got := sb.NumberRecipesBeforeDigits(tt.number); got != tt.want {
				t.Errorf("ScoreBoard.NumberRecipesBeforeDigits(%v) = %v, want %v", tt.number, got, tt.want)
			}
		})
	}
}
