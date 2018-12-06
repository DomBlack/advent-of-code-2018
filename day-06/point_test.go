package main

import (
	"reflect"
	"testing"
)

func TestNewPoint(t *testing.T) {
	tests := []struct {
		coords  string
		wantRes Point
	}{
		{"1, 1", Point{1, 1}},
		{"-1, -1", Point{-1, -1}},
		{"1, 6", Point{1, 6}},
		{"8, 9", Point{8, 9}},
	}
	for _, tt := range tests {
		t.Run(tt.coords, func(t *testing.T) {
			if gotRes := NewPoint(tt.coords); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("NewPoint() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestPointsFromStrings(t *testing.T) {
	tests := []struct {
		name   string
		coords []string
		want   []Point
	}{
		{
			"Example",
			[]string{"1, -1", "-5, 23"},
			[]Point{{1, -1,}, {-5, 23}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PointsFromStrings(tt.coords); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PointsFromStrings() = %v, want %v", got, tt.want)
			}
		})
	}
}
