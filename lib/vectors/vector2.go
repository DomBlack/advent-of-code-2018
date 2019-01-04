package vectors

import (
	"fmt"
	"github.com/DomBlack/advent-of-code-2018/lib"
)

type Vec2 struct {
	X, Y int
}

// Returns a point presenting the largest X & Y of these two points
func (p1 Vec2) Max(p2 Vec2) Vec2 {
	return Vec2{
		lib.Max(p1.X, p2.X),
		lib.Max(p1.Y, p2.Y),
	}
}

// Returns a point presenting the smallest X & Y of these two points
func (p1 Vec2) Min(p2 Vec2) Vec2 {
	return Vec2{
		lib.Min(p1.X, p2.X),
		lib.Min(p1.Y, p2.Y),
	}
}

// Adds two vectors together
func (p1 Vec2) Add(p2 Vec2) Vec2 {
	return Vec2 { p1.X + p2.X, p1.Y + p2.Y }
}

// Gets the manhattan distance between two vectors
func (p1 Vec2) Distance(p2 Vec2) int {
	return lib.Abs(p1.X - p2.X) + lib.Abs(p1.Y + p2.Y)
}

func (i Vec2) IsReadingOrderLess(j Vec2) bool {
	return i.Y < j.Y || (i.X < j.X && i.Y == j.Y)
}

func (p Vec2) String() string {
	return fmt.Sprintf("%v,%v", p.X, p.Y)
}

func NewVec2(x, y int) Vec2 {
	return Vec2 { x, y }
}
