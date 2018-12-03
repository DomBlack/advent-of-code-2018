package lib

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
)

// Reads the given input file as a slice of integers
func InputAsIntegers(folder string) []int {
	file, err := os.Open(folder + "/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer Close(file)

	var input []int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		i, err := strconv.Atoi(scanner.Text())

		if err != nil {
			log.Fatal(err)
		}

		input = append(input, i)
	}

	return input
}

// Reads the input line by line as strings
func InputAsStrings(folder string) []string {
	file, err := os.Open(folder + "/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer Close(file)

	var input []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		input = append(input, scanner.Text())
	}

	return input
}

// Closes a closer handling the error
func Close(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Fatal(err)
	}
}
