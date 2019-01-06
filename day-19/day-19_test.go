package main

import "testing"

func Test_runGoLangVersion(t *testing.T) {
	tests := []struct {
		name string
		n int
		want int
	}{
		{ "Initial Example", 836 + 46, 2223 },
		{ "Extra #1", 34, 54 },
		{ "Extra #2", 423, 624 },
		{ "Extra #3", 528, 1488 },
		{ "Extra #4", 1234, 1854 },
		{ "Extra #5", 4301, 5184 },
		{ "Extra #6", 5431, 5432 },
		{ "Extra #7", 53491, 54000 },
		{ "Extra #8", 53441, 53442 },
		{ "Extra #9", 145239, 193656 },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sumFactors(tt.n); got != tt.want {
				t.Errorf("sumFactors(%d) = %v, want %v", tt.n, got, tt.want)
			}
		})
	}
}
