package main

import (
	"github.com/DomBlack/advent-of-code-2018/lib"
	"github.com/DomBlack/advent-of-code-2018/lib/elf_code"
)

func main() {
	println(transpileToJS("day-19", []int{0}))
	println("===============")
	println(transpileToJS("day-21", []int {0 }))
}

func transpileToJS(day string, inputRegisters []int) string {
	input := lib.InputAsString(day)

	cpu, err := elf_code.NewCPUFromProgramFile(input)

	if err != nil {
		panic(err)
	}

	transpiler := cpu.StartTranspiler(elf_code.TranspileOptions{
		CompressConstants:           true,
		RemoveEmptyBlocks:           true,
		RemoveExtraJumps:            true,
		RemoveUnUsedRegisterWrites:  true,
		RewriteRecursionAsLoops:     true,
		InlineBlocksWherePossible:   true,
		DisplayBlockNumbersOnOutput: false,
		DisplayBlockRegisterUse:     false,
	})

	for _, registerNum := range inputRegisters {
		transpiler.Registers[registerNum].SetUnknownBool()
	}

	javaScript, err := transpiler.Run()
	if err != nil {
		println(err.Error())
	}

	return javaScript
}
