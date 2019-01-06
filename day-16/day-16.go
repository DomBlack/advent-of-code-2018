package main

import (
	"fmt"
	"github.com/DomBlack/advent-of-code-2018/day-16/elf_code"
	"github.com/DomBlack/advent-of-code-2018/lib"
	"log"
	"reflect"
	"strings"
)

func main() {
	input := lib.InputAsStrings("day-16")
	samples, program := parseInput(input)

	fmt.Println("Part 1", part1(samples))
	fmt.Println("Part 2", part2(samples, program))
}

// How Many Samples have three or more matching op codes
func part1(samples []Sample) (count int) {
	for _, sample := range samples {
		if len(sample.MatchingOpCodes()) >= 3 {
			count++
		}
	}

	return
}

// Work out the number of each OpCode and then execute the sample program
func part2(samples []Sample, program elf_code.Program) int {
	// Create a map of all possibilities - OpCode => Possible Numbers
	possibleNumbers := make(map[elf_code.OpCode]map[int]bool, elf_code.NumOpCodes)
	knownNumbers := make(map[int]elf_code.OpCode)
	for opCode := elf_code.OpCode(0); opCode < elf_code.NumOpCodes; opCode++ {
		possibleNumbers[opCode] = make(map[int]bool)
	}

	// Loop over every sample and set all initial possibilities
	for _, sample := range samples {
		instructionNumber := int(sample.instruction.OpCode) // The OpCode number in the sample
		matching := sample.MatchingOpCodes()                // The OpCodes which match the result of the sample

		for _, opCode := range matching {
			possibleNumbers[elf_code.OpCode(opCode)][instructionNumber] = true
		}
	}

	// Minimise the matches until there are no unknown ones
	for len(knownNumbers) < int(elf_code.NumOpCodes) {
		for opCode, numbers := range possibleNumbers {
			if len(numbers) == 1 {

				// Remove this number from all others
				for o := elf_code.OpCode(0); o < elf_code.NumOpCodes; o++ {
					if o != opCode {
						for number := range numbers {
							knownNumbers[number] = opCode
							delete(possibleNumbers[o], number)
						}
					}
				}
			}
		}
	}

	// Now map the program from the original op codes to my op codes (defined by my enum)
	for index := range program {
		program[index].OpCode = knownNumbers[int(program[index].OpCode)]
	}

	// Create a CPU and run it with the given program
	cpu := elf_code.NewCPU(program)
	err := cpu.Execute()

	if err != nil {
		log.Fatal(err)
	}

	return cpu.Registers[0]
}

func parseInput(input []string) (samples []Sample, program elf_code.Program) {
	samples = make([]Sample, 0)
	program = make(elf_code.Program, 0)

	for i := 0; i < len(input); i++ {
		line := input[i]

		if strings.HasPrefix(line, "Before: [") {
			samples = append(
				samples,
				parseSample(line, input[i+1], input[i+2]),
			)
			i += 3
		} else if line != "" {
			instruction, err := elf_code.NewInstruction(line)
			if err != nil {
				log.Fatal("Unable to parse instruction " + line)
			}

			program = append(program, instruction)
		}
	}

	return
}

// The sample
type Sample struct {
	before, after elf_code.Registers   // The before and after register states
	instruction   elf_code.Instruction // The instruction to the sample was from
}

// Parse the sample
func parseSample(before, instruction, after string) (res Sample) {
	res.before = elf_code.NewRegisters()
	res.after = elf_code.NewRegisters()

	num, err := fmt.Sscanf(before, "Before: [%d, %d, %d, %d]", &res.before[0], &res.before[1], &res.before[2], &res.before[3])
	if num != 4 && err != nil {
		log.Fatal("Unable to parse: " + before)
	}

	res.instruction, err = elf_code.NewInstruction(instruction)
	if err != nil {
		log.Fatal("Unable to parse instruction: " + instruction)
	}

	num, err = fmt.Sscanf(after, "After: [%d, %d, %d, %d]", &res.after[0], &res.after[1], &res.after[2], &res.after[3])
	if num != 4 && err != nil {
		log.Fatal("Unable to parse: " + after)
	}

	return
}

// Tests if the given OpCode
func (s Sample) TestOpCode(opcode elf_code.OpCode) bool {
	f, found := elf_code.OpCodeFunc[opcode]

	if !found {
		log.Println("Unknown op code for TestOpCode: ", opcode)
		return false
	}

	registers := s.before.Copy()

	// Execute the op code and get the value of the result
	cValue, err := f(s.instruction.A, s.instruction.B, registers)
	if err != nil {
		return false
	}

	// Set the register as required
	err = registers.Set(s.instruction.C, cValue)
	if err != nil {
		return false
	}

	// Are the registers the same?
	return reflect.DeepEqual(registers, s.after)
}

// Finds all matching OpCodes
func (s Sample) MatchingOpCodes() (matching []elf_code.OpCode) {
	matching = make([]elf_code.OpCode, 0)

	// Test every op code
	for opCode := elf_code.OpCode(0); opCode < elf_code.NumOpCodes; opCode++ {
		if s.TestOpCode(opCode) {
			matching = append(matching, opCode)
		}
	}

	return
}

func (s Sample) String() string {
	var str strings.Builder

	str.WriteString("Before: ")
	str.WriteString(s.before.String())
	str.WriteRune('\n')
	str.WriteString(s.instruction.String())
	str.WriteRune('\n')
	str.WriteString("After: ")
	str.WriteString(s.after.String())

	return str.String()
}
