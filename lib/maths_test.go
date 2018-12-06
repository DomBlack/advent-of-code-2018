package lib

import (
	"strconv"
	"testing"
)

func TestAbs(t *testing.T) {
	tests := []struct {
		a int
		want int
	}{
		{ -1, 1 },
		{ 0, 0 },
		{1, 1 },
		{ -42348, 42348 },
	}
	for _, tt := range tests {
		t.Run(strconv.Itoa(tt.a), func(t *testing.T) {
			if got := Abs(tt.a); got != tt.want {
				t.Errorf("Abs() = %v, want %v", got, tt.want)
			}
		})
	}
}
