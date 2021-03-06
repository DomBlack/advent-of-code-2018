package elf_code

import (
	"errors"
	"strconv"
	"strings"
)

type Registers []int

// Create the Registers
func NewRegisters(numRegisters int) (res Registers) {
	res = make(Registers, numRegisters)
	return
}

// Gets the register entry at register `a`
func (r Registers) Get(a int) (value int, err error) {
	if a >= len(r) || a < 0 {
		return 0, errors.New("register index out of bounds")
	}

	return r[a], nil
}

// Gets the value of the two given Registers `a` and `b`
func (r Registers) GetTwo(a, b int) (valueA int, valueB int, err error) {
	if a >= len(r) || a < 0 {
		return 0, 0, errors.New("register a out of bounds")
	}

	if b >= len(r) || b < 0 {
		return 0, 0, errors.New("register b out of bounds")
	}

	return r[a], r[b], nil
}

// Sets the register `a` to the given `value`
func (r Registers) Set(a int, value int) (err error) {
	if a >= len(r) || a < 0 {
		return errors.New("register index out of bounds")
	}

	r[a] = value
	return
}

// Creates a copy of the Registers
func (r Registers) Copy() (res Registers) {
	res = make(Registers, len(r))
	copy(res, r)
	return
}

// Writes the Registers out
func (r Registers) String() string {
	var str strings.Builder

	str.WriteRune('[')

	for index, value := range r {
		if index > 0 {
			str.WriteString(", ")
		}

		str.WriteString(strconv.Itoa(value))
	}

	str.WriteRune(']')

	return str.String()
}
