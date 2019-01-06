package main

import (
	"fmt"
	"github.com/DomBlack/advent-of-code-2018/lib"
	"github.com/DomBlack/advent-of-code-2018/lib/elf_code"
	"log"
)

func main() {
	input := lib.InputAsString("day-19")

	fmt.Println("Part 1", runElfCodeVersion(input, 0))
	fmt.Println("Part 2", runGoLangVersion(1))
}

func runElfCodeVersion(input string, register0StartingValue int) int {
	cpu, err := elf_code.NewCPUFromProgramFile(input)
	if err != nil {
		log.Fatal(err)
	}

	cpu.Registers[0] = register0StartingValue

	err = cpu.Execute()
	if err != nil {
		log.Fatal(err)
	}

	return cpu.Registers[0]
}

// Reversed engineered GoLang version of input.txt
func runGoLangVersion(register0StartingValue int) int {
	var sumOfFactor, n int // sum = R[0], n = R[2]
	sumOfFactor = register0StartingValue

	// Init always
	n += 836 + 46 // (18-21) 25 (22 - 24)

	// Extra init when sumOfFactor starts at 1
	if sumOfFactor == 1 { // 26-27
		n += 10550400   // 34, 28-33
		sumOfFactor = 0 // 35
	}

	// Move into another function to unit test
	return sumFactors(n)
}

func sumFactors(n int) (sumOfFactor int) {
	// Main program - Returns the number of factors
	// <editor-fold desc="Original Code">
	//for i := 1; i <= n; i++ { // 2; 14-16; 13
	//	for j := 1; j <= n; j++ { // 3; 10-12; 9
	//		if (i * j) == n { // 4-7
	//			sumOfFactor += i // 8
	//		}
	//	}
	//}
	// </editor-fold>

	// <editor-fold desc="Optimized Code">
	for i := 1; i <= n; i++ {
		if n % i == 0 {
			sumOfFactor += i
		}
	}
	// </editor-fold>

	return
}
