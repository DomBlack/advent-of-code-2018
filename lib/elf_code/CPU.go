package elf_code

import (
	"bufio"
	"errors"
	"fmt"
	"strings"
)

// A program for the CPU to execute
type Program []Instruction

// A virtual CPU for Elf Code
type CPU struct {
	Registers                  Registers // The register state of the CPU
	Program                    Program   // The loaded program within the CPU
	InstructionPointerRegister int       // The register that represents the instruction pointer
}

// Creates a new CPU and loads it with the given program
func NewCPU(program Program, ipRegister int, numRegisters int) (res *CPU) {
	if ipRegister < 0 || ipRegister >= numRegisters {
		panic("The IP Register must be within range of the Registers!")
	}

	return &CPU{
		NewRegisters(numRegisters),
		program,
		ipRegister,
	}
}

// Creates a new CPU using a program file
func NewCPUFromProgramFile(fileContents string) (res *CPU, err error) {
	ipRegister := 0
	program := make(Program, 0)

	scanner := bufio.NewScanner(strings.NewReader(fileContents))

	// Read the IP register binding
	scanner.Scan()
	num, e := fmt.Sscanf(scanner.Text(), "#ip %d", &ipRegister)
	if e != nil {
		return nil, e
	}
	if num != 1 {
		return nil, errors.New("invalid number of params found on ip line")
	}

	// Parse the program instructions
	for scanner.Scan() {
		instruction, e := NewInstruction(scanner.Text())
		if e != nil {
			return nil, e
		}

		program = append(program, instruction)
	}

	// Check the scanner for errors
	if err = scanner.Err(); err != nil {
		return
	}

	// Return the CPU
	return NewCPU(program, ipRegister, 6), nil
}

// Executes the program loaded into the CPU
func (cpu *CPU) Execute() (err error) {
	numInstructions := len(cpu.Program)

	// While the instruction pointer is within the program
	for cpu.Registers[cpu.InstructionPointerRegister] >= 0 &&
		cpu.Registers[cpu.InstructionPointerRegister] < numInstructions {

		// Grab the next instruction
		instruction := cpu.Program[cpu.Registers[cpu.InstructionPointerRegister]]

		// Execute it
		result, e := OpCodeFunc[instruction.OpCode](instruction.A, instruction.B, cpu.Registers)

		if e != nil {
			return e
		}

		// Store the result of it
		e = cpu.Registers.Set(instruction.C, result)
		if e != nil {
			return e
		}

		// Increment the instruction pointer
		cpu.Registers[cpu.InstructionPointerRegister]++
	}

	// Reduce the program counter back down
	cpu.Registers[cpu.InstructionPointerRegister]--

	return nil
}
