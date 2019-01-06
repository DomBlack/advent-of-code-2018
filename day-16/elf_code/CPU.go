package elf_code

// A program for the CPU to execute
type Program []Instruction

// A virtual CPU for Elf Code
type CPU struct {
	Registers Registers
	Program   Program
}

// Creates a new CPU and loads it with the given program
func NewCPU(program Program) (res *CPU) {
	return &CPU {
		NewRegisters(),
		program,
	}
}

// Executes the program loaded into the CPU
func (cpu *CPU) Execute() (err error) {
	for _, instruction := range cpu.Program {
		result, e := OpCodeFunc[instruction.OpCode](instruction.A, instruction.B, cpu.Registers)

		if e != nil {
			return e
		}

		e = cpu.Registers.Set(instruction.C, result)
		if e != nil {
			return e
		}
	}

	return nil
}
