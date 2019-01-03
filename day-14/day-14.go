package main

import (
	"container/ring"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
)

func main() {
	sb := NewScoreBoard()
	fmt.Println("Part 1", sb.NextTenAfter(74501))
	fmt.Println("Part 2", sb.NumberRecipesBeforeDigits("074501"))
}

type ScoreBoard struct {
	startingRecipe *ring.Ring // The first recipe generated
	lastRecipe     *ring.Ring // The last recipe generated
	count          int
	elf1           *ring.Ring // The position of elf 1
	elf2           *ring.Ring // The position of elf 2
}

// Creates a new scoreboard with the initial entries
func NewScoreBoard() (sb ScoreBoard) {
	sb.startingRecipe = ring.New(2) // init with 2 entries
	sb.lastRecipe = sb.startingRecipe.Next()
	sb.count = 2

	// Set initial values
	sb.startingRecipe.Value = 3
	sb.lastRecipe.Value = 7

	// Set elf pointers
	sb.elf1 = sb.startingRecipe
	sb.elf2 = sb.lastRecipe

	return
}

// Debug print the score board using the same formatting as the AoC site
func (sb ScoreBoard) String() string {
	var str strings.Builder

	currentEntry := sb.startingRecipe

	// do {} while (currentEntry != startingRecipe)
	for ok := true; ok; ok = currentEntry != sb.startingRecipe {
		value := rune(currentEntry.Value.(int) + '0')

		if currentEntry == sb.elf1 {
			str.WriteRune('(')
			str.WriteRune(value)
			str.WriteRune(')')
		} else if currentEntry == sb.elf2 {
			str.WriteRune('[')
			str.WriteRune(value)
			str.WriteRune(']')
		} else {
			str.WriteRune(' ')
			str.WriteRune(value)
			str.WriteRune(' ')
		}

		currentEntry = currentEntry.Next()
	}

	return str.String()
}

// Create a new recipe; Each digit of the sum of elf1 & elf2's current recipe become new recipes
func (sb *ScoreBoard) CreateRecipe() {
	newRecipes := sb.elf1.Value.(int) + sb.elf2.Value.(int)

	// Add the new recipes
	ringInsertPoint := sb.lastRecipe
	firstInsert := true

	for ok := true; ok; ok = newRecipes > 0 {
		// Work out the digit
		digit := newRecipes % 10

		// Create a ring entry for it
		entry := ring.New(1)
		entry.Value = digit

		// Insert after the original last entry, before we started looping
		// as we are adding them in the reverse order
		ringInsertPoint.Link(entry)
		sb.count++

		// If this was the first digit we added, it becomes the new last recipe
		if firstInsert {
			sb.lastRecipe = entry
			firstInsert = false
		}

		newRecipes /= 10
	}

	// Each elf now picks a new recipe
	sb.elf1 = sb.elf1.Move(sb.elf1.Value.(int) + 1)
	sb.elf2 = sb.elf2.Move(sb.elf2.Value.(int) + 1)
}

// Returns the next ten recipes after the given number of recipes (part 1)
func (sb *ScoreBoard) NextTenAfter(number int) string {
	var str strings.Builder

	// Ensure we have enough recipes, we need number + 10
	for ; sb.count < (number + 10); sb.CreateRecipe() {
	}

	// Add the next 10 entries to the list
	entry := sb.startingRecipe.Move(number)
	for i := 0; i < 10; i++ {
		str.WriteRune(rune('0' + entry.Value.(int)))
		entry = entry.Next()
	}

	return str.String()
}

// How many entries appear before the given string of digits
func (sb *ScoreBoard) NumberRecipesBeforeDigits(digits string) int {
	// Convert the digits to a required score
	requiredScore, err := strconv.Atoi(digits)

	if err != nil {
		log.Fatal(err)
	}

	mask := int(math.Pow10(len(digits)))
	currentScore := 0

	entry := sb.startingRecipe
	count := 0
	for {
		currentScore = ((currentScore * 10) % mask) + entry.Value.(int)

		if requiredScore == currentScore && count > len(digits) {
			return count - len(digits) + 1
		}

		// If we are on the last recipe, create more before moving to the next one
		if entry == sb.lastRecipe {
			sb.CreateRecipe()
		}
		entry = entry.Next()
		count++
	}
}
