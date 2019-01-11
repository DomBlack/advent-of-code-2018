package main

import (
	"errors"
	"fmt"
	"github.com/DomBlack/advent-of-code-2018/lib"
	"golang.org/x/tools/container/intsets"
	"log"
	"strings"
)

func main() {
	input := lib.InputAsString("day-20")

	direction, err := NewDirectionsFrom(input)
	if err != nil {
		log.Fatal(err)
	}

	rooms := direction.CreateRoomMap()

	fmt.Println("Part 1", rooms.FurthestAwayRoom())
	fmt.Println("Part 2", rooms.RoomsOverNDoorsAway(1000))
}

const North complex64 = -1i
const East complex64 = +1
const South complex64 = +1i
const West complex64 = -1

type Room struct {
	coord                    complex64
	shortestDistance         int
	North, East, South, West *Room // Neighbouring Rooms
}

func (r *Room) OpenDoor(direction complex64, m *RoomMap) (newRoom *Room) {
	newRoom = m.GetRoom(r.coord + direction)

	switch direction {
	case North:
		r.North = newRoom
		newRoom.South = r
	case East:
		r.East = newRoom
		newRoom.West = r
	case South:
		r.South = newRoom
		newRoom.North = r
	case West:
		r.West = newRoom
		newRoom.East = r
	}

	return
}

type RoomMap struct {
	rooms map[complex64]*Room
}

func (rm *RoomMap) GetRoom(coord complex64) (room *Room) {
	room, found := rm.rooms[coord]

	if !found {
		room = &Room{
			coord,
			intsets.MaxInt,
			nil, nil, nil, nil,
		}
		rm.rooms[coord] = room
	}

	return room
}

func (rm *RoomMap) FurthestAwayRoom() int {
	distance := 0

	for _, room := range rm.rooms {
		if room.shortestDistance > distance {
			distance = room.shortestDistance
		}
	}

	return distance
}

func (rm *RoomMap) RoomsOverNDoorsAway(n int) (count int) {
	for _, room := range rm.rooms {
		if room.shortestDistance >= n {
			count ++
		}
	}

	return
}

type Direction struct {
	isBranch  bool         // Is this direction a branch
	direction complex64    // What is the direction if it's not a branch
	options   []*Direction // What are the options if it is a branch
	next      *Direction   // What's the next direction?
}

// Creates a direction structure from the given input string
func NewDirectionsFrom(input string) (res *Direction, err error) {
	if input[0] != '^' {
		err = errors.New("expect ^ at start of regex")
		return
	}

	reader := strings.NewReader(input[1:])
	res, err = newDirection(reader)
	if err != nil {
		return
	}

	ch, _, err := reader.ReadRune()
	if err != nil {
		return
	}

	if ch != '$' {
		return nil, errors.New(fmt.Sprintf("expected `$`, got `%v`", ch))
	}

	return
}

func (d Direction) CreateRoomMap() (m *RoomMap) {
	m = &RoomMap{make(map[complex64]*Room)}

	// Build the map from the directions
	rootRoom := m.GetRoom(0)
	d.buildMap([]*Room{rootRoom}, m)

	// Now flood fill the map for the distances
	rootRoom.shortestDistance = 0
	toVisit := []*Room{rootRoom}
	for len(toVisit) > 0 {
		room := toVisit[0]
		toVisit = toVisit[1:]

		costToVisit := room.shortestDistance + 1

		if room.North != nil && room.North.shortestDistance > costToVisit {
			room.North.shortestDistance = costToVisit
			toVisit = append(toVisit, room.North)
		}

		if room.East != nil && room.East.shortestDistance > costToVisit {
			room.East.shortestDistance = costToVisit
			toVisit = append(toVisit, room.East)
		}

		if room.South != nil && room.South.shortestDistance > costToVisit {
			room.South.shortestDistance = costToVisit
			toVisit = append(toVisit, room.South)
		}

		if room.West != nil && room.West.shortestDistance > costToVisit {
			room.West.shortestDistance = costToVisit
			toVisit = append(toVisit, room.West)
		}
	}

	return
}

func (d Direction) buildMap(rooms []*Room, m *RoomMap) []*Room {
	var res []*Room

	if d.isBranch {
		res = make([]*Room, 0)

		for _, option := range d.options {
			if option != nil {
				res = append(res, option.buildMap(rooms, m)...)
			}
		}
	} else {
		res = make([]*Room, len(rooms))

		for i, room := range rooms {
			res[i] = room.OpenDoor(d.direction, m)
		}
	}

	if d.next != nil {
		res = d.next.buildMap(res, m)
	}

	return res
}

// Creates a single direction or option from the reader at it's current point
func newDirection(reader *strings.Reader) (res *Direction, err error) {
	ch, _, err := reader.ReadRune()
	if err != nil {
		return nil, err
	}

	res = &Direction{
		false,
		0,
		make([]*Direction, 0),
		nil,
	}

	switch ch {
	case '(':
		res.isBranch = true
		for {
			// Read the option
			option, err := newDirection(reader)
			if err != nil {
				return nil, err
			}
			res.options = append(res.options, option)

			ch, _, err = reader.ReadRune()
			if err != nil {
				return nil, err
			}

			if ch == ')' {
				break
			} else if ch != '|' {
				return nil, errors.New(fmt.Sprintf("Unexpected rune `%v`", ch))
			}
		}
	case '|', ')', '$':
		err = reader.UnreadRune()
		return nil, err
	case 'N':
		res.direction = North
	case 'S':
		res.direction = South
	case 'E':
		res.direction = East
	case 'W':
		res.direction = West
	default:
		err = errors.New(fmt.Sprintf("Unknown rune `%v`", ch))
	}

	next, err := newDirection(reader)
	if err != nil {
		return nil, err
	}

	res.next = next

	return
}

func (d Direction) toString(str *strings.Builder) {
	if d.isBranch {
		str.WriteRune('(')
		for i, option := range d.options {
			if i > 0 {
				str.WriteRune('|')
			}

			if option != nil {
				option.toString(str)
			}
		}
		str.WriteRune(')')
	} else {
		switch d.direction {
		case North:
			str.WriteRune('N')
		case East:
			str.WriteRune('E')
		case South:
			str.WriteRune('S')
		case West:
			str.WriteRune('W')
		}
	}

	if d.next != nil {
		d.next.toString(str)
	}
}

// Writes the directions back out in the Regex form
func (d Direction) String() string {
	var str strings.Builder

	str.WriteRune('^')
	d.toString(&str)
	str.WriteRune('$')

	return str.String()
}
