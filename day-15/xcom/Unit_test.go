package xcom

import (
	"reflect"
	"testing"

	"github.com/DomBlack/advent-of-code-2018/lib/vectors"
)

func TestUnit_FindTargets(t *testing.T) {
	tests := []struct {
		name   string
		inputMap string
		targetIndexes []int
	}{
		{ "No Targets", "#..E...E#\n#....E..#", []int{}, },
		{ "Two Targets", "#..E...G#\n#.G..E..#", []int{1, 2}, },
		{ "All Targets", "#..G...G#\n#.G..G..#", []int{0, 1, 2, 3}, },
	}

	u := NewElf(vectors.NewVec2(0, 0))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMap(tt.inputMap)

			want := make(Units, len(tt.targetIndexes))
			for i, t := range tt.targetIndexes {
				want[i] = m.Units[t]
			}

			if got := u.FindTargets(m); !reflect.DeepEqual(got, want) {
				t.Errorf("Unit.FindTargets() = %v, want %v", got, want)
			}
		})
	}
}
