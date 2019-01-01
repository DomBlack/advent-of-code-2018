package main

import (
	"bufio"
	"fmt"
	"github.com/DomBlack/advent-of-code-2018/lib/vectors"
	"io/ioutil"
	"log"
	"sort"
	"strings"
)

func main() {
	input, err := ioutil.ReadFile("day-13/input.txt")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Part 1:", part1(string(input)))
	fmt.Println("Part 2:", part2(string(input)))
}

// Where does the first crash occur
func part1(input string) vectors.Vec2 {
	m := NewMap(input)

	for {
		crashes := m.Tick()

		if len(crashes) > 0 {
			return crashes[0]
		}
	}
}

// Where does the final cart end up when all others have been removed?
func part2(input string) vectors.Vec2 {
	m := NewMap(input)

	for {
		m.Tick()

		if len(m.carts) == 1 {
			return m.carts[0].Position()
		}
	}
}

// A Cart
type Cart struct {
	// Our current position
	position complex64
	// Our direction of travel: +1, +1i, -1, -1i are the four cardinal directions (right, down, left, up)
	direction complex64
	// The next turn we'll make, multiply with the direction to make that turn (+1i = clockwise / -1i anti-clockwise)
	nextIntersectionTurn complex64
}

func NewCart(x, y int, direction complex64) *Cart {
	return &Cart{
		complex(float32(x), float32(y)),
		direction,
		-1i,
	}
}

// Gets real X,Y coords out of the complex number
func (c Cart) Position() vectors.Vec2 {
	return vectors.NewVec2(int(real(c.position)), int(imag(c.position)))
}

func (c Cart) MapIndex(width int) int {
	p := c.Position()
	return p.X + (p.Y * width)
}

// Continues moving in the existing direction
func (c *Cart) MoveForward() {
	c.position += c.direction
}

// Moves through the intersection, turning as per the rules
func (c *Cart) MoveThroughIntersection() {
	// Rotate our direction
	c.direction *= c.nextIntersectionTurn

	switch c.nextIntersectionTurn {
	case -1i:                      // was left
		c.nextIntersectionTurn = 1 // now straight
	case 1:                          // was straight
		c.nextIntersectionTurn = +1i // now right
	case +1i:
		c.nextIntersectionTurn = -1i // now left
	default:
		log.Fatal("Unknown intersection turn", c.nextIntersectionTurn)
	}

	// Now we've turned move forward
	c.position += c.direction
}

func (c *Cart) MoveThroughCorner(cornerRune rune) {
	switch cornerRune {
	case '/':
		c.direction = complex(-imag(c.direction), -real(c.direction))
	case '\\':
		c.direction = complex(imag(c.direction), real(c.direction))
	default:
		log.Fatal("Unknown corner type", cornerRune)
	}

	// Now we've turned around, continue forward
	c.position += c.direction
}

func (c Cart) String() string {
	switch c.direction {
	case -1i:
		return "^"
	case 1:
		return ">"
	case +1i:
		return "v"
	case -1:
		return "<"
	default:
		log.Fatal("Unknown direction", c.direction)
	}

	return "X"
}

// A slice of carts (for custom sorting)
type Carts []*Cart

func (s Carts) Len() int {
	return len(s)
}

func (s Carts) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Carts) Less(i, j int) bool {
	iP := s[i].Position()
	jP := s[j].Position()
	return iP.Y < jP.Y || (iP.Y == jP.Y && iP.X < jP.X)
}

// Our Map
type Map struct {
	width int
	data  []rune
	carts Carts
}

// Parses the map from the given string input
func NewMap(input string) (res Map) {
	// Read the map
	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		line := scanner.Text()

		if res.width == 0 {
			res.width = len(line)
		}

		res.data = append(res.data, []rune(line)...)
	}

	// Read all the carts
	for index, datum := range res.data {
		var direction complex64

		switch datum {
		case '^':
			direction = -1i
		case '>':
			direction = 1
		case 'v':
			direction = 1i
		case '<':
			direction = -1
		}

		if direction != 0 {
			res.data[index] = '-'
			res.carts = append(
				res.carts,
				NewCart(index%res.width, index/res.width, direction),
			)
		}
	}
	return
}

// Process a tick
func (m *Map) Tick() []vectors.Vec2 {
	crashes := make([]vectors.Vec2, 0)

	// First resort the carts so they are in the order we want to process them (top row first, left to right)
	sort.Sort(m.carts)

	cartMap := make(map[complex64]*Cart)
	for _, cart := range m.carts {
		cartMap[cart.position] = cart
	}

	// Move all carts
	for _, cart := range m.carts {
		// Has this cart already be removed this tick?
		if _, exists := cartMap[cart.position]; !exists {
			continue
		}

		delete(cartMap, cart.position)

		switch m.data[cart.MapIndex(m.width)] {
		case '+':
			cart.MoveThroughIntersection()
		case '/':
			cart.MoveThroughCorner('/')
		case '\\':
			cart.MoveThroughCorner('\\')
		default:
			cart.MoveForward()
		}

		if _, exists := cartMap[cart.position]; exists {
			// Remove the other cart, as we have collided
			delete(cartMap, cart.position)

			crashes = append(crashes, cart.Position())
		} else {
			cartMap[cart.position] = cart
		}
	}

	// Rebuild the slice of all remaining carts
	m.carts = make(Carts, 0)
	for _, cart := range cartMap {
		m.carts = append(m.carts, cart)
	}

	return crashes
}

// Convert the map back into a string format for visual debugging
func (m Map) String() string {
	var str strings.Builder

	// Put all carts into the map based on their datum index
	cartMap := make(map[int]*Cart)
	for _, cart := range m.carts {
		cartMap[cart.MapIndex(m.width)] = cart
	}

	// Build the map string
	for index, datum := range m.data {
		if index > 0 && index % m.width == 0 {
			str.WriteRune('\n')
		}

		if cart, exists := cartMap[index]; exists {
			str.WriteString(cart.String())
		} else {
			str.WriteRune(datum)
		}
	}

	return str.String()
}
