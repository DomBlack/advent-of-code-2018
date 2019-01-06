package elf_code

import (
	"reflect"
	"testing"
)

func TestNewCPUFromProgramFile(t *testing.T) {
	fileContents := `#ip 0
seti 5 0 1
seti 6 0 2
addi 0 1 0
addr 1 2 3
setr 1 0 0
seti 8 0 4
seti 9 0 5`

	cpu, err := NewCPUFromProgramFile(fileContents)
	if err != nil {
		t.Errorf("NewCPUFromProgramFile() error = %v", err)
		return
	}

	if len(cpu.Registers) != 6 {
		t.Errorf("NewCPUFromProgramFile() registerCount = %v, want %v", len(cpu.Registers), 6)
		return
	}

	if cpu.InstructionPointerRegister != 0 {
		t.Errorf("NewCPUFromProgramFile() InstructionPointerRegister = %v, want %v", cpu.InstructionPointerRegister, 0)
		return
	}

	err = cpu.Execute()
	if err != nil {
		t.Errorf("NewCPUFromProgramFile() Execute error = %v", err)
		return
	}

	want := Registers { 6, 5, 6, 0 , 0, 9 }
	if !reflect.DeepEqual(cpu.Registers, want) {
		t.Errorf("NewCPUFromProgramFile() = %v, want %v", cpu.Registers, want)
	}
}
