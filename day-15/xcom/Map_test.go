package xcom

import (
	"testing"
)

func TestNewMap(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantRes string
	}{
		{
			"Simple",
			`#######
#.G.E.#
#E.G.E#
#.G.E.#
#######`,
			`#######   
#.G.E.#   G(200), E(200)
#E.G.E#   E(200), G(200), E(200)
#.G.E.#   G(200), E(200)
#######   
`,
		},
		{
			"Basic Example",
			`#######
#.G...#
#...EG#
#.#.#G#
#..G#E#
#.....#
#######`,
			`#######   
#.G...#   G(200)
#...EG#   E(200), G(200)
#.#.#G#   G(200)
#..G#E#   G(200), E(200)
#.....#   
#######   
`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMap(tt.input, 3)
			if gotRes := m.String(); gotRes != tt.wantRes {
				t.Errorf("NewMap() = \n%v \n\n, want\n%v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestMap_Round(t *testing.T) {
	tests := []struct {
		name                string
		executeNumberRounds int
		wantCombatOver      bool
		wantMap             string
	}{
		{
			"After 1 round",
			1, false,
			`#######   
#..G..#   G(200)
#...EG#   E(197), G(197)
#.#G#G#   G(200), G(197)
#...#E#   E(197)
#.....#   
#######   
`,
		},
		{
			"After 2 rounds",
			1, false,
			`#######   
#...G.#   G(200)
#..GEG#   G(200), E(188), G(194)
#.#.#G#   G(194)
#...#E#   E(194)
#.....#   
#######   
`,
		},
		{
			"After 23 rounds",
			23 - 2, false,
			`#######   
#...G.#   G(200)
#..G.G#   G(200), G(131)
#.#.#G#   G(131)
#...#E#   E(131)
#.....#   
#######   
`,
		},
		{
			"After 24 rounds",
			1, false,
			`#######   
#..G..#   G(200)
#...G.#   G(131)
#.#G#G#   G(200), G(128)
#...#E#   E(128)
#.....#   
#######   
`,
		},
		{
			"After 25 rounds",
			1, false,
			`#######   
#.G...#   G(200)
#..G..#   G(131)
#.#.#G#   G(125)
#..G#E#   G(200), E(125)
#.....#   
#######   
`,
		},
		{
			"After 26 rounds",
			1, false,
			`#######   
#G....#   G(200)
#.G...#   G(131)
#.#.#G#   G(122)
#...#E#   E(122)
#..G..#   G(200)
#######   
`,
		},
		{
			"After 27 rounds",
			1, false,
			`#######   
#G....#   G(200)
#.G...#   G(131)
#.#.#G#   G(119)
#...#E#   E(119)
#...G.#   G(200)
#######   
`,
		},
		{
			"After 28 rounds",
			1, false,
			`#######   
#G....#   G(200)
#.G...#   G(131)
#.#.#G#   G(116)
#...#E#   E(113)
#....G#   G(200)
#######   
`,
		},
		{
			"After 47 rounds",
			47 - 28, false,
			`#######   
#G....#   G(200)
#.G...#   G(131)
#.#.#G#   G(59)
#...#.#   
#....G#   G(200)
#######   
`,
		},
		{
			"During Round 48 combat finishes",
			1, true,
			`#######   
#G....#   G(200)
#.G...#   G(131)
#.#.#G#   G(59)
#...#.#   
#....G#   G(200)
#######   
`,
		},
	}

	m := NewMap(`#######
#.G...#
#...EG#
#.#.#G#
#..G#E#
#.....#
#######`, 3)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotCombatOver bool
			for i := 0; i < tt.executeNumberRounds; i++ {
				gotCombatOver = m.Round()
			}

			gotMap := m.String()

			if gotMap != tt.wantMap {
				t.Errorf("gotMap = %v, want %v", gotMap, tt.wantMap)
			}

			if gotCombatOver != tt.wantCombatOver {
				t.Errorf("gotCombatOver = %v, want %v", gotCombatOver, tt.wantCombatOver)
			}
		})
	}
}

func TestMap_RunCombatSim(t *testing.T) {
	tests := []struct {
		name      string
		wantScore int
		inputMap  string
		wantMap   string
	}{
		{
			"Example 1", 27730,
			`#######
#.G...#
#...EG#
#.#.#G#
#..G#E#
#.....#
#######`,
			`#######   
#G....#   G(200)
#.G...#   G(131)
#.#.#G#   G(59)
#...#.#   
#....G#   G(200)
#######   
`,
		},

		{
			"Example 2", 36334,
			`#######
#G..#E#
#E#E.E#
#G.##.#
#...#E#
#...E.#
#######`,
			`#######   
#...#E#   E(200)
#E#...#   E(197)
#.E##.#   E(185)
#E..#E#   E(200), E(200)
#.....#   
#######   
`,
		},
		{
			"Example 3", 39514,
			`#######
#E..EG#
#.#G.E#
#E.##E#
#G..#.#
#..E#.#
#######`,
`#######   
#.E.E.#   E(164), E(197)
#.#E..#   E(200)
#E.##.#   E(98)
#.E.#.#   E(200)
#...#.#   
#######   
`,
		},
		{
			"Example 4", 27755,
			`#######
#E.G#.#
#.#G..#
#G.#.G#
#G..#.#
#...E.#
#######`,
`#######   
#G.G#.#   G(200), G(98)
#.#G..#   G(200)
#..#..#   
#...#G#   G(95)
#...G.#   G(200)
#######   
`,
		},
		{
			"Example 5", 28944,
			`#######
#.E...#
#.#..G#
#.###.#
#E#G#G#
#...#G#
#######`,
`#######   
#.....#   
#.#G..#   G(200)
#.###.#   
#.#.#.#   
#G.G#G#   G(98), G(38), G(200)
#######   
`,
		},
		{
			"Example 6", 18740,
			`#########
#G......#
#.E.#...#
#..##..G#
#...##..#
#...#...#
#.G...G.#
#.....G.#
#########`,
`#########   
#.G.....#   G(137)
#G.G#...#   G(200), G(200)
#.G##...#   G(200)
#...##..#   
#.G.#...#   G(200)
#.......#   
#.......#   
#########   
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMap(tt.inputMap, 3)

			if gotScore := m.RunCombatSim(); gotScore != tt.wantScore {
				t.Errorf("gotScore = %v, want %v\n\n%v", gotScore, tt.wantScore, m)
			}

			if gotMap := m.String(); gotMap != tt.wantMap {
				t.Errorf("gotMap:\n%v\n\nwantMap:\n%v\n\n", gotMap, tt.wantMap)
			}
		})
	}
}
