package elf_code

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Instruction struct {
	OpCode OpCode // Which OpCode this instruction is
	A      int    // The first Input
	B      int    // The second Input
	C      int    // The output register
}

// Creates an instruction from a given string
func NewInstruction(code string) (res Instruction, err error) {
	num, err := fmt.Sscanf(code, "%d %d %d %d", &res.OpCode, &res.A, &res.B, &res.C)

	if err != nil {
		return
	}

	if num != 4 {
		err = errors.New("instruction had an invalid number of ints")
		return
	}

	return
}

// Convert the instruction back to a string
func (i Instruction) String() string {
	var str strings.Builder

	str.WriteString(i.OpCode.String())
	str.WriteRune(' ')
	str.WriteString(strconv.Itoa(i.A))
	str.WriteRune(' ')
	str.WriteString(strconv.Itoa(i.B))
	str.WriteRune(' ')
	str.WriteString(strconv.Itoa(i.C))

	return str.String()
}
