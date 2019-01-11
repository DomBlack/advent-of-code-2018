package collections

import "github.com/DomBlack/advent-of-code-2018/lib/vectors"

type vecStackNode struct {
	value    vectors.Vec2  // The value at this point
	previous *vecStackNode // The previous node
}

type Vec2Stack struct {
	head *vecStackNode // The head
}

func NewVec2Stack() Vec2Stack {
	return Vec2Stack{}
}

func (s *Vec2Stack) Push(value vectors.Vec2) {
	newHead := &vecStackNode{
		value,
		s.head,
	}

	s.head = newHead
}

func (s *Vec2Stack) Pop() (value vectors.Vec2) {
	if s.head != nil {
		value = s.head.value
		s.head = s.head.previous
		return
	} else {
		panic("Pop on non empty stack")
	}
}

func (s *Vec2Stack) IsEmpty() bool {
	return s.head == nil
}
