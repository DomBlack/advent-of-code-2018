package main

import (
	"reflect"
	"testing"

	"github.com/DomBlack/advent-of-code-2018/lib/vectors"
)

func TestCart(t *testing.T) {
	cart := NewCart(34, 85, +1)

	testPos := func(expectedX, expectedY int) {
		p := cart.Position()

		if p.X != expectedX || p.Y != expectedY {
			t.Errorf("Cart.Position() = (%v, %v), want (%v, %v)", p.X, p.Y, expectedX, expectedY)
		}
	}

	// Check the position function actually works
	testPos(34, 85)

	// Check moving right (direction +1)
	cart.MoveForward()
	testPos(35, 85)

	// Check moving left (direction -1)
	cart.direction = -1
	cart.MoveForward()
	testPos(34, 85)

	// Check moving up (direction -1i)
	cart.direction = -1i
	cart.MoveForward()
	testPos(34, 84)

	// Check moving down (direction +1i)
	cart.direction = +1i
	cart.MoveForward()
	testPos(34, 85)

	// First turn will be left, which means we'll move right as we're currently pointing down
	cart.MoveThroughIntersection()
	testPos(35, 85)

	// Second turn will be straight, which means we'll move right again
	cart.MoveThroughIntersection()
	testPos(36, 85)

	// Third turn will be right, which as we're pointing right, means we'll go down
	cart.MoveThroughIntersection()
	testPos(36, 86)

	// Forth turn will be left again, which means we'll go right
	cart.MoveThroughIntersection()
	testPos(37, 86)

	// Test going around a corner from going right to going up
	cart.MoveThroughCorner('/')
	testPos(37, 85)

	// Test going around a corner from going up to going right
	cart.MoveThroughCorner('/')
	testPos(38, 85)

	// Test going around a corner from going right to going down
	cart.MoveThroughCorner('\\')
	testPos(38, 86)

	// Test going around a corner from going down to going left
	cart.MoveThroughCorner('/')
	testPos(37, 86)

	// Test going around a corner from going left to going down
	cart.MoveThroughCorner('/')
	testPos(37, 87)

	// Test going around a corner from going down to going right
	cart.MoveThroughCorner('\\')
	testPos(38, 87)
}

func TestNewMap(t *testing.T) {
	input :=
		`  />-\
/-/  |
v    |
\----/
`
	want := Map{
		6,
		[]rune{
			' ', ' ', '/', '-', '-', '\\',
			'/', '-', '/', ' ', ' ', '|',
			'-', ' ', ' ', ' ', ' ', '|',
			'\\', '-', '-', '-', '-', '/',
		},
		Carts{
			NewCart(3, 0, 1),
			NewCart(0, 2, 1i),
		},
	}

	if got := NewMap(input); !reflect.DeepEqual(got, want) {
		t.Errorf("NewMap() = %v, want %v", got, want)
	}
}

func TestMap_Tick(t *testing.T) {
	currentMap := NewMap(">------<-")

	expectedAfterTick := []struct {
		m       string
		crashes []vectors.Vec2
	}{
		{"->----<--", []vectors.Vec2{}},
		{"-->--<---", []vectors.Vec2{}},
		{"---><----", []vectors.Vec2{}},
		{"---------", []vectors.Vec2{vectors.NewVec2(4, 0)}},
	}

	for index, expected := range expectedAfterTick {
		gotCrashes := currentMap.Tick()

		if got := currentMap.String(); got != expected.m {
			t.Errorf("Tick(%v) = %v, wanted %v", index, got, expected.m)
		}

		if !reflect.DeepEqual(gotCrashes, expected.crashes) {
			t.Errorf("Tick(%v).crashes = %v, wanted %v", index, gotCrashes, expected.crashes)
		}
	}
}

func Test_part1(t *testing.T) {
	input := `/->-\        
|   |  /----\
| /-+--+-\  |
| | |  | v  |
\-+-/  \-+--/
  \------/   `

	want := vectors.NewVec2(7, 3)

	if got := part1(input); !reflect.DeepEqual(got, want) {
		t.Errorf("part1() = %v, wanted %v", got, want)
	}
}

func Test_part2(t *testing.T) {
	input := `/>-<\  
|   |  
| /<+-\
| | | v
\>+</ |
  |   ^
  \<->/`

	want := vectors.NewVec2(6, 4)

	if got := part2(input); !reflect.DeepEqual(got, want) {
		t.Errorf("part2() = %v, wanted %v", got, want)
	}
}
