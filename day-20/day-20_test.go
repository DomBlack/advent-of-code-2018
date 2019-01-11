package main

import (
	"testing"
)

func TestExample(t *testing.T) {
	input := "^ENWWW(NEEE|SSE(EE|N|))$"
	directions, err := NewDirectionsFrom(input)

	if err != nil {
		t.Error(err)
		return
	}

	if got := directions.String(); got != input {
		t.Errorf("NewDirectionsFrom() = %v, want %v", got, input)
		return
	}

	m := directions.CreateRoomMap()

	const want = 10
	if got := m.FurthestAwayRoom(); got != want {
		t.Errorf("FurthestAwayRoom() = %v, want %v", got, want)
		return
	}
}

func TestRoomMap_FurthestAwayRoom(t *testing.T) {
	tests := []struct {
		name   string
		want   int
	}{
		{ "^WNE$", 3 },
		{ "^ENWWW(NEEE|SSE(EE|N))$", 10 },
		{ "^ENNWSWW(NEWS|)SSSEEN(WNSE|)EE(SWEN|)NNN$", 18 },
		{ "^ESSWWN(E|NNENN(EESS(WNSE|)SSS|WWWSSSSE(SW|NNNE)))$", 23 },
		{ "^WSSEESWWWNW(S|NENNEEEENN(ESSSSW(NWSW|SSEN)|WSWWN(E|WWS(E|SS))))$", 31 },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			directions, err := NewDirectionsFrom(tt.name)
			if err != nil {
				t.Error(err)
				return
			}

			if got := directions.CreateRoomMap().FurthestAwayRoom(); got != tt.want {
				t.Errorf("RoomMap.FurthestAwayRoom() = %v, want %v", got, tt.want)
			}
		})
	}
}
