package main

import (
	"fmt"
	"github.com/DomBlack/advent-of-code-2018/lib"
	"log"
	"sort"
)

func main() {
	input := lib.InputAsStrings("day-04")

	guards := parseInput(input)

	fmt.Println("Part 1:", part1(guards))
	fmt.Println("Part 2:", part2(guards))
}

func parseInput(lines []string) []Guard {
	// Ensure the input is sorted
	sort.Strings(lines)

	// Create a struct to represent the l
	type Line struct {
		year, month, day, hour, minute int
		verb                           string
	}
	l := Line{}

	// Create our guard storage
	guardSlice := make([]Guard, 0)
	guards := make(map[int]Guard)
	currentGuard := Guard{}
	asleep := false
	fellAsleep := 0

	// Loop the input parsing it
	for _, line := range lines {
		num, err := fmt.Sscanf(
			line,
			"[%d-%d-%d %d:%d] %s",
			&l.year, &l.month, &l.day, &l.hour, &l.minute, &l.verb,
		)

		if err != nil {
			log.Fatal(err, line)
		}

		if num != 6 {
			log.Fatalf("Expected l to have 6 inputs, got %d", num)
		}

		switch l.verb {
		case "Guard":
			// Change of guard
			id := getGuardID(line)

			if asleep {
				log.Fatal("Guard change while still asleep", line)
			}

			guard, found := guards[id]
			if !found {
				guard = NewGuard(id)
				guards[id] = guard
				guardSlice = append(guardSlice, guard)
			}
			currentGuard = guard

		case "falls":
			fellAsleep = l.minute
			asleep = true

		case "wakes":
			asleep = false

			for i := fellAsleep; i < l.minute; i++ {
				if _, found := currentGuard.sleeping[i]; !found {
					currentGuard.sleeping[i] = 1
				} else {
					currentGuard.sleeping[i]++
				}
			}
		}
	}

	return guardSlice
}

// Out of the given guards, returns the result of the the guard ID * minute most asleep for the
// guard who slept the most during the midnight hour
func part1(guards []Guard) int {
	maxAmount := guards[0].TotalSleep()
	maxGuard := guards[0]

	for _, guard := range guards {
		amount := guard.TotalSleep()
		if amount > maxAmount {
			maxAmount = amount
			maxGuard = guard
		}
	}

	minuteMostAsleep, _ := maxGuard.MinuteMostAsleep()

	return maxGuard.id * minuteMostAsleep
}

// Which guard and which minute was the most slept minute in terms of frequency
func part2(guards []Guard) int {
	maxFrequency := 0
	maxSleepMin := 0
	maxGuardId := 0

	for _, guard := range guards {
		minuteMostAsleep, frequency := guard.MinuteMostAsleep()

		if frequency > maxFrequency {
			maxFrequency = frequency
			maxSleepMin = minuteMostAsleep
			maxGuardId = guard.id
		}
	}

	return maxGuardId * maxSleepMin
}

type Guard struct {
	id       int         // The guard ID
	sleeping map[int]int // Number of times the guard has slept at a given time
}

// How long has the guard slept for
func (guard Guard) TotalSleep() int {
	total := 0

	for _, amount := range guard.sleeping {
		total += amount
	}

	return total
}

// Which minute of the midnight shift is the guard most asleep and how many times that minute did they sleep for?
func (guard Guard) MinuteMostAsleep() (int, int) {
	maxAmount := 0
	maxMinute := 0

	for min, amount := range guard.sleeping {
		if amount > maxAmount {
			maxMinute = min
			maxAmount = amount
		}
	}

	return maxMinute, maxAmount
}

// Creates a guard
func NewGuard(id int) (res Guard) {
	res.id = id
	res.sleeping = make(map[int]int)
	return
}

// Extracts the guard ID from the string
func getGuardID(line string) (res int) {
	num, err := fmt.Sscanf(
		line[19:],
		"Guard #%d",
		&res,
	)

	if err != nil {
		log.Fatal(err)
	}

	if num != 1 {
		log.Fatalf("Expected line to have 1 input, got %d", num)
	}

	return
}
