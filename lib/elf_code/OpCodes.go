package elf_code

import (
	"errors"
	"strconv"
	"strings"
)

type OpCode int

const (
	AddR OpCode = iota // `addr` (add register) stores into register `C` the result of adding register `A` and register `B`.
	AddI               // `addi` (add immediate) stores into register `C` the result of adding register `A` and value `B`.
	MulR               // `mulr` (multiply register) stores into register `C` the result of multiplying register `A` and register `B`.
	MulI               // `muli` (multiply immediate) stores into register `C` the result of multiplying register `A` and value `B`.
	BanR               // `banr` (bitwise AND register) stores into register `C` the result of the bitwise AND of register `A` and register `B`.
	BanI               // `bani` (bitwise AND immediate) stores into register `C` the result of the bitwise AND of register `A` and value `B`.
	BorR               // `borr` (bitwise OR register) stores into register `C` the result of the bitwise OR of register `A` and register `B`.
	BorI               // `bori` (bitwise OR immediate) stores into register `C` the result of the bitwise OR of register `A` and value `B`.
	SetR               // `setr` (set register) copies the contents of register `A` into register `C`. (Input `B` is ignored.)
	SetI               // `seti` (set immediate) stores value `A` into register `C`. (Input `B` is ignored.)
	GtIR               // `gtir` (greater-than immediate/register) sets register `C` to 1 if value `A` is greater than register `B`. Otherwise, register `C` is set to 0.
	GtRI               // `gtri` (greater-than register/immediate) sets register `C` to 1 if register `A` is greater than value `B`. Otherwise, register `C` is set to 0.
	GtRR               // `gtrr` (greater-than register/register) sets register `C` to 1 if register `A` is greater than register `B`. Otherwise, register `C` is set to 0.
	EqIR               // `eqir` (equal immediate/register) sets register `C` to 1 if value `A` is equal to register `B`. Otherwise, register `C` is set to 0.
	EqRI               // `eqri` (equal register/immediate) sets register `C` to 1 if register `A` is equal to value `B`. Otherwise, register `C` is set to 0.
	EqRR               // `eqrr` (equal register/register) sets register `C` to 1 if register `A` is equal to register `B`. Otherwise, register `C` is set to 0.
)

const NumOpCodes = OpCode(int(EqRR) + 1)

func (o OpCode) String() string {
	switch o {
	case AddR:
		return "addr"
	case AddI:
		return "addi"
	case MulR:
		return "mulr"
	case MulI:
		return "muli"
	case BanR:
		return "banr"
	case BanI:
		return "bani"
	case BorR:
		return "borr"
	case BorI:
		return "bori"
	case SetR:
		return "setr"
	case SetI:
		return "seti"
	case GtIR:
		return "gtir"
	case GtRI:
		return "gtri"
	case GtRR:
		return "gtrr"
	case EqIR:
		return "eqir"
	case EqRI:
		return "eqri"
	case EqRR:
		return "eqrr"
	default:
		return "Unknown OpCode: " + strconv.Itoa(int(o))
	}
}

// Parse an op code from string
func ParseOpCode(str string) (opCode OpCode, err error) {
	switch strings.ToLower(str) {
	case "addr":
		return AddR, nil
	case "addi":
		return AddI, nil
	case "mulr":
		return MulR, nil
	case "muli":
		return MulI, nil
	case "banr":
		return BanR, nil
	case "bani":
		return BanI, nil
	case "borr":
		return BorR, nil
	case "bori":
		return BorI, nil
	case "setr":
		return SetR, nil
	case "seti":
		return SetI, nil
	case "gtir":
		return GtIR, nil
	case "gtri":
		return GtRI, nil
	case "gtrr":
		return GtRR, nil
	case "eqir":
		return EqIR, nil
	case "eqri":
		return EqRI, nil
	case "eqrr":
		return EqRR, nil
	default:
		err = errors.New("unknown op code: " + str)
		return
	}
}

// All op codes without the given `filter` list
func OpCodesWithout(filter []OpCode) []OpCode {
	result := make([]OpCode, 0)

	for opCode := OpCode(0); opCode < NumOpCodes; opCode++ {
		found := false
		for i := 0; i < len(filter); i++ {
			if filter[i] == opCode {
				found = true
			}
		}

		if !found {
			result = append(result, opCode)
		}
	}

	return result
}

// The function which executes the given op code, with the A, B inputs and the current registers. Returns the value to
// put in output C and any error with the op code
type OpFunc = func(a int, b int, registers Registers) (value int, err error)

var OpCodeFunc = map[OpCode]OpFunc{
	// addr (add register) stores into register C the result of adding register A and register B.
	AddR: func(a int, b int, registers Registers) (value int, err error) {
		aV, bV, err := registers.GetTwo(a, b)
		if err != nil {
			return 0, err
		}

		return aV + bV, nil
	},

	// addi (add immediate) stores into register C the result of adding register A and value B.
	AddI: func(a int, b int, registers Registers) (value int, err error) {
		aV, err := registers.Get(a)
		if err != nil {
			return 0, err
		}

		return aV + b, nil
	},

	// mulr (multiply register) stores into register C the result of multiplying register A and register B.
	MulR: func(a int, b int, registers Registers) (value int, err error) {
		aV, bV, err := registers.GetTwo(a, b)
		if err != nil {
			return 0, err
		}

		return aV * bV, nil
	},

	// muli (multiply immediate) stores into register C the result of multiplying register A and value B.
	MulI: func(a int, b int, registers Registers) (value int, err error) {
		aV, err := registers.Get(a)
		if err != nil {
			return 0, err
		}

		return aV * b, nil
	},

	// banr (bitwise AND register) stores into register C the result of the bitwise AND of register A and register B.
	BanR: func(a int, b int, registers Registers) (value int, err error) {
		aV, bV, err := registers.GetTwo(a, b)
		if err != nil {
			return 0, err
		}

		return aV & bV, nil
	},

	// bani (bitwise AND immediate) stores into register C the result of the bitwise AND of register A and value B.
	BanI: func(a int, b int, registers Registers) (value int, err error) {
		aV, err := registers.Get(a)
		if err != nil {
			return 0, err
		}

		return aV & b, nil
	},

	// borr (bitwise OR register) stores into register C the result of the bitwise OR of register A and register B.
	BorR: func(a int, b int, registers Registers) (value int, err error) {
		aV, bV, err := registers.GetTwo(a, b)
		if err != nil {
			return 0, err
		}

		return aV | bV, nil
	},

	// bori (bitwise OR immediate) stores into register C the result of the bitwise OR of register A and value B.
	BorI: func(a int, b int, registers Registers) (value int, err error) {
		aV, err := registers.Get(a)
		if err != nil {
			return 0, err
		}

		return aV | b, nil
	},

	// setr (set register) copies the contents of register A into register C. (Input B is ignored.)
	SetR: func(a int, b int, registers Registers) (value int, err error) {
		aV, err := registers.Get(a)
		if err != nil {
			return 0, err
		}

		return aV, nil
	},

	// seti (set immediate) stores value A into register C. (Input B is ignored.)
	SetI: func(a int, b int, registers Registers) (value int, err error) {
		return a, nil
	},

	// gtir (greater-than immediate/register) sets register C to 1 if value A is greater than register B. Otherwise, register C is set to 0.
	GtIR: func(a int, b int, registers Registers) (value int, err error) {
		bV, err := registers.Get(b)
		if err != nil {
			return 0, err
		}

		if a > bV {
			return 1, nil
		} else {
			return 0, nil
		}
	},

	// gtri (greater-than register/immediate) sets register C to 1 if register A is greater than value B. Otherwise, register C is set to 0.
	GtRI: func(a int, b int, registers Registers) (value int, err error) {
		aV, err := registers.Get(a)
		if err != nil {
			return 0, err
		}

		if aV > b {
			return 1, nil
		} else {
			return 0, nil
		}
	},

	// gtrr (greater-than register/register) sets register C to 1 if register A is greater than register B. Otherwise, register C is set to 0.
	GtRR: func(a int, b int, registers Registers) (value int, err error) {
		aV, bV, err := registers.GetTwo(a, b)
		if err != nil {
			return 0, err
		}

		if aV > bV {
			return 1, nil
		} else {
			return 0, nil
		}
	},

	// eqir (equal immediate/register) sets register C to 1 if value A is equal to register B. Otherwise, register C is set to 0.
	EqIR: func(a int, b int, registers Registers) (value int, err error) {
		bV, err := registers.Get(b)
		if err != nil {
			return 0, err
		}

		if a == bV {
			return 1, nil
		} else {
			return 0, nil
		}
	},

	// eqri (equal register/immediate) sets register C to 1 if register A is equal to value B. Otherwise, register C is set to 0.
	EqRI: func(a int, b int, registers Registers) (value int, err error) {
		aV, err := registers.Get(a)
		if err != nil {
			return 0, err
		}

		if aV == b {
			return 1, nil
		} else {
			return 0, nil
		}
	},

	// eqrr (equal register/register) sets register C to 1 if register A is equal to register B. Otherwise, register C is set to 0.
	EqRR: func(a int, b int, registers Registers) (value int, err error) {
		aV, bV, err := registers.GetTwo(a, b)
		if err != nil {
			return 0, err
		}

		if aV == bV {
			return 1, nil
		} else {
			return 0, nil
		}
	},
}
