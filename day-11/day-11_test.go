package main

import (
	"reflect"
	"testing"

	"github.com/DomBlack/advent-of-code-2018/lib/vectors"
)

func TestGetPowerLevel(t *testing.T) {
	type args struct {
		position      vectors.Vec2
		gridSerialNum int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"Example 1", args{vectors.NewVec2(3, 5), 8}, 4},
		{"Example 2", args{vectors.NewVec2(122, 79), 57}, -5},
		{"Example 3", args{vectors.NewVec2(217, 196), 39}, 0},
		{"Example 4", args{vectors.NewVec2(101, 153), 71}, 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetPowerLevel(tt.args.position, tt.args.gridSerialNum); got != tt.want {
				t.Errorf("GetPowerLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_FindBestFuelCellPatch(t *testing.T) {
	type args struct {
		gridSerialNum int
		minPatchSize int
		maxPatchSize int
	}
	tests := []struct {
		name string
		args args
		wantPos vectors.Vec2
		wantSize int
	}{
		{ "Example 1", args{ 18, 3, 3 }, vectors.NewVec2(33, 45), 3 },
		{ "Example 2", args{ 42, 3, 3 }, vectors.NewVec2(21, 61), 3 },
		{ "Example 3", args{ 18, 1, 300 }, vectors.NewVec2(90,  269), 16 },
		{ "Example 4", args{ 42, 1, 300 }, vectors.NewVec2(232,  251), 12 },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, size := FindBestFuelCellPatch(tt.args.gridSerialNum, tt.args.minPatchSize, tt.args.maxPatchSize)

			if !reflect.DeepEqual(got, tt.wantPos) {
				t.Errorf("FindBestFuelCellPatch() Position = %v, want %v", got, tt.wantPos)
			}
			if tt.wantSize != size {
				t.Errorf("FindBestFuelCellPatch() Size = %v, want %v", size, tt.wantSize)
			}

		})
	}
}
