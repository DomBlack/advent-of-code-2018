package elf_code

import (
	"fmt"
	"log"
	"sort"
	"strings"
)

type variableNameMap struct {
	m             map[int]string
	nextVarNumber uint32
}

func newVariableMap() (res variableNameMap) {
	res.m = make(map[int]string)
	res.nextVarNumber = 1
	return
}

func (m *variableNameMap) get(register int) string {
	name, found := m.m[register]
	if !found {
		return m.getNew(register)
	} else {
		return name
	}
}

func (m *variableNameMap) getNew(register int) string {
	name := ""
	n := m.nextVarNumber
	for n > 0 {
		name += string('A' + (n - 1%26))
		n /= 26
	}

	m.m[register] = name
	m.nextVarNumber++

	return name
}

// Re-orders program inputs when order doesn't matter,
// such that both `add A B A` and `add B A A` both become `add A B A`
// and on OpCodes which have an option for IR and RI inputs, normalises to RI
func (c *CPU) normaliseProgramForTranspile() {
	for i, instruction := range c.Program {
		// If input B is the same register as output C, reorder so it is always input A
		if instruction.B == instruction.C {
			switch instruction.OpCode {
			case AddR:
				instruction.A, instruction.B = instruction.B, instruction.A
			case MulR:
				instruction.A, instruction.B = instruction.B, instruction.A
			case EqIR:
				instruction.OpCode = EqRI
				instruction.A, instruction.B = instruction.B, instruction.A
			case GtIR:
				instruction.OpCode = GtRI
				instruction.A, instruction.B = instruction.B, instruction.A
			}
		}

		// Normalise IP modifications
		if instruction.OpCode == AddR {
			if instruction.A == c.InstructionPointerRegister {
				if instruction.B == c.InstructionPointerRegister {
					instruction.OpCode = SetI
					instruction.A = i + i
				} else {
					instruction.OpCode = AddI
					instruction.A, instruction.B = instruction.B, i
				}
			} else if instruction.B == c.InstructionPointerRegister {
				instruction.OpCode = AddI
				instruction.B = i
			}
		} else if instruction.OpCode == MulR {
			if instruction.A == c.InstructionPointerRegister {
				if instruction.B == c.InstructionPointerRegister {
					instruction.OpCode = SetI
					instruction.A = i * i
				} else {
					instruction.OpCode = MulI
					instruction.A, instruction.B = instruction.B, i

				}
			} else if instruction.B == c.InstructionPointerRegister {
				instruction.OpCode = MulI
				instruction.B = i
			}
		} else if instruction.OpCode == SetR && instruction.A == c.InstructionPointerRegister {
			instruction.OpCode = SetI
			instruction.A, instruction.B = i, 0
		}

		c.Program[i] = instruction
	}
}

type LineType int

const (
	Statement LineType = iota
	JumpStatement
	IfStatement
	DoWhileStatement
)

type ProgramLine struct {
	ip              int           // The original instruction pointer
	lineType        LineType      // The line type
	instruction     *Instruction  // The instruction
	jumpToBlock     *ProgramBlock // Jump to block (or true block for if statements)
	elseBlock       *ProgramBlock // false block for if statements
	parentBlock     *ProgramBlock // The parentBlock this line is within
	next            *ProgramLine  // The next line after this
	previous        *ProgramLine  // The previous line before this
	invertCondition bool          // Invert the condition of the while loop?
}

func NewProgramLine(ip int, instruction *Instruction, block *ProgramBlock) *ProgramLine {
	return &ProgramLine{
		ip,
		Statement,
		instruction,
		nil,
		nil,
		block,
		nil,
		nil,
		false,
	}
}

func (pl *ProgramLine) CanBeConst() bool {
	return pl.instruction.OpCode == AddI || pl.instruction.OpCode == MulI
}

func (pl *ProgramLine) forgetCallsOut() {
	if pl.jumpToBlock != nil {
		pl.jumpToBlock.removeCallFrom(pl)
		pl.jumpToBlock = nil
	}

	if pl.elseBlock != nil {
		pl.elseBlock.removeCallFrom(pl)
		pl.elseBlock = nil
	}
}

// Remove this line from the program as it's unused
func (pl *ProgramLine) RemoveUnusedLine() {
	if pl.previous != nil {
		pl.previous.next = pl.next
	} else {
		pl.parentBlock.firstLine = pl.next
	}

	if pl.next != nil {
		pl.next.previous = pl.previous
	} else {
		pl.parentBlock.lastLine = pl.previous
	}

	pl.forgetCallsOut()
}

func (pl *ProgramLine) moveToNewBlockAfter(line *ProgramLine) {
	// Update the block we used to be in to terminate on the line before us
	if pl.previous != nil {
		// Remove us from the previous line next
		pl.previous.next = nil
	} else {
		// of if we where the first line, set the block to empty
		pl.parentBlock.firstLine = nil
	}
	pl.parentBlock.lastLine = pl.previous

	// Find the last instruction we represent and updated all parent blocks to the new parent
	lastInChain := pl
	pl.parentBlock = line.parentBlock
	for next := pl.next; next != nil; next = next.next {
		next.parentBlock = line.parentBlock
		lastInChain = next
	}

	// Update the chain next reference
	if line.next != nil {
		line.next.previous = lastInChain
		lastInChain.next = line.next
	} else {
		line.parentBlock.lastLine = lastInChain
		lastInChain.next = nil
	}

	// Insert our chain start point into the existing block
	pl.previous = line
	line.next = pl
}

// Removes this line and all following lines from the block
func (pl *ProgramLine) RemoveAllFollowingLines() {
	pl.parentBlock.lastLine = pl

	for line := pl.next; line != nil; line = line.next {
		line.forgetCallsOut()
	}

	pl.next = nil
}

func (pl *ProgramLine) WriteInstructionWhenACSame(str *strings.Builder) {
	opCode := pl.instruction.OpCode
	isImmediate := OpCodeInputType[opCode]

	str.WriteString(fmt.Sprintf("R[%d]", pl.instruction.C))

	switch opCode {
	case AddR:
		str.WriteString(" += ")
	case AddI:
		switch pl.instruction.B {
		case 1:
			str.WriteString("++")
			return
		case -1:
			str.WriteString("--")
			return
		default:
			str.WriteString(" += ")
		}
	case MulR:
		str.WriteString(" *= ")
	case MulI:
		str.WriteString(" *= ")
	case BanR:
		str.WriteString(" &= ")
	case BanI:
		str.WriteString(" &= ")
	case BorR:
		str.WriteString(" |= ")
	case BorI:
		str.WriteString(" |= ")
	default:
		panic("Unable to write this!")
	}

	if isImmediate.B {
		str.WriteString(fmt.Sprintf("%d", pl.instruction.B))
	} else {
		str.WriteString(fmt.Sprintf("R[%d]", pl.instruction.B))
	}
}

func (pl *ProgramLine) WriteInstruction(str *strings.Builder) {
	opCode := pl.instruction.OpCode
	isImmediate := OpCodeInputType[opCode]

	if !isImmediate.A && pl.instruction.A == pl.instruction.C && opCode < SetR {
		pl.WriteInstructionWhenACSame(str)
		return
	}

	str.WriteString(fmt.Sprintf("R[%d] = ", pl.instruction.C))

	if isImmediate.A {
		str.WriteString(fmt.Sprintf("%d", pl.instruction.A))
	} else {
		str.WriteString(fmt.Sprintf("R[%d]", pl.instruction.A))
	}

	switch opCode {
	case AddR:
		str.WriteString(" + ")
	case AddI:
		str.WriteString(" + ")
	case MulR:
		str.WriteString(" * ")
	case MulI:
		str.WriteString(" * ")
	case BanR:
		str.WriteString(" & ")
	case BanI:
		str.WriteString(" & ")
	case BorR:
		str.WriteString(" | ")
	case BorI:
		str.WriteString(" | ")
	case SetR:
		return
	case SetI:
		return
	case GtIR:
		str.WriteString(" > ")
	case GtRI:
		str.WriteString(" > ")
	case GtRR:
		str.WriteString(" > ")
	case EqIR:
		str.WriteString(" == ")
	case EqRI:
		str.WriteString(" == ")
	case EqRR:
		str.WriteString(" == ")
	}

	if isImmediate.B {
		str.WriteString(fmt.Sprintf("%d", pl.instruction.B))
	} else {
		str.WriteString(fmt.Sprintf("R[%d]", pl.instruction.B))
	}
}

type BlockType int

const (
	FunctionBlock BlockType = iota
	InlineBlock
)

type ProgramBlock struct {
	blockNum           int            // The block number
	blockType          BlockType      // What type of block is this
	calledBy           []*ProgramLine // All lines which jump to this block
	firstLine          *ProgramLine   // The first line within this block
	lastLine           *ProgramLine   // The last line within this block
	modifiersRegisters []bool         // If a call to this block will modify the given registers
	readsRegisters     []bool         // List of registers which have their existing values used when coming into this block
}

func (b ProgramBlock) writeBlockAsJavaScript(state *TranspileState, str *strings.Builder, indent int) {
	indentStr := ""
	for i := 0; i < indent; i++ {
		indentStr += "    "
	}

	if state.options.DisplayBlockNumbersOnOutput {
		str.WriteString(fmt.Sprintf("%s// BEGIN block%d\n", indentStr, b.blockNum))
	}

	for line := b.firstLine; line != nil; line = line.next {
		switch line.lineType {
		case Statement:
			str.WriteString(indentStr)
			line.WriteInstruction(str)
		case JumpStatement:
			if line.jumpToBlock == nil {
				str.WriteString(indentStr + "// FIXME: jump to nil\n")
			} else if line.jumpToBlock.blockType == InlineBlock {
				line.jumpToBlock.writeBlockAsJavaScript(state, str, indent) // does not need extra indent
				continue                                                    // don't need to original comment
			} else {
				str.WriteString(fmt.Sprintf("%sblock%d()", indentStr, line.jumpToBlock.blockNum))
			}
		case IfStatement:
			if b.firstLine != line {
				// Add whitespace above if statements for readability
				str.WriteString("\n")
			}

			if line.jumpToBlock != nil {
				str.WriteString(fmt.Sprintf("%sif (R[%d] == 1) { // IP: %d (%s)\n", indentStr, line.instruction.A, line.ip, state.originalCode[line.ip]))

				if line.jumpToBlock.blockType == InlineBlock {
					line.jumpToBlock.writeBlockAsJavaScript(state, str, indent+1)
				} else {
					str.WriteString(fmt.Sprintf("%s    block%d()\n", indentStr, line.jumpToBlock.blockNum))
				}

				if line.elseBlock != nil {
					str.WriteString(indentStr + "} else {\n")
					if line.elseBlock.blockType == InlineBlock {
						line.elseBlock.writeBlockAsJavaScript(state, str, indent+1)
					} else {
						str.WriteString(fmt.Sprintf("%s    block%d()\n", indentStr, line.elseBlock.blockNum))
					}
				}

				str.WriteString(indentStr + "}\n")
			} else if line.elseBlock != nil {
				str.WriteString(fmt.Sprintf("%sif (R[%d] == 0) { // IP: %d (%s)\n", indentStr, line.instruction.A, line.ip, state.originalCode[line.ip]))
				if line.elseBlock.blockType == InlineBlock {
					line.elseBlock.writeBlockAsJavaScript(state, str, indent+1)
				} else {
					str.WriteString(fmt.Sprintf("%s    block%d()\n", indentStr, line.elseBlock.blockNum))
				}
				str.WriteString(indentStr + "}\n")
			} else {
				panic("JumpStatement with no jumps!")
			}

			if b.lastLine != line {
				// Add whitespace under if statements for readability
				str.WriteString("\n")
			}
			continue
		case DoWhileStatement:
			if b.firstLine != line {
				// Add whitespace above if statements for readability
				str.WriteString("\n")
			}

			str.WriteString(fmt.Sprintf("%sdo {\n", indentStr))

			line.jumpToBlock.writeBlockAsJavaScript(state, str, indent+1)

			targetValue := 1
			if line.invertCondition {
				targetValue = 0
			}
			str.WriteString(fmt.Sprintf("%s} while (R[%d] == %d); // IP: %d (%s)\n", indentStr, line.instruction.A, targetValue, line.ip, state.originalCode[line.ip]))

			if b.lastLine != line {
				// Add whitespace under if statements for readability
				str.WriteString("\n")
			}
			continue

		default:
			log.Println("Unknown line type", line.lineType)
		}

		str.WriteString(fmt.Sprintf(" // IP: %d (%s)\n", line.ip, state.originalCode[line.ip]))
	}

	if state.options.DisplayBlockNumbersOnOutput {
		str.WriteString(fmt.Sprintf("%s// END block%d\n", indentStr, b.blockNum))
	}
}

func (b *ProgramBlock) replaceJumpsToBlockWithBlock(newBlock *ProgramBlock) {
	for _, line := range b.calledBy {
		if line.jumpToBlock == b {
			line.jumpToBlock = newBlock
		} else if line.elseBlock == b {
			line.elseBlock = newBlock
		} else {
			panic(fmt.Sprintf("Line apparently calls block %d, but it's not jumpToBlock or elseBlock!", b.blockNum))
		}

		if newBlock != nil {
			newBlock.calledBy = append(newBlock.calledBy, line)
		}
	}

	b.calledBy = []*ProgramLine{}
}

func (b *ProgramBlock) removeCallFrom(line *ProgramLine) {
	calledBy := make([]*ProgramLine, 0)
	for _, other := range b.calledBy {
		if other != line {
			calledBy = append(calledBy, other)
		}
	}

	b.calledBy = calledBy
}

func (b *ProgramBlock) getLastJumpIfInlined() *ProgramLine {
	if b.blockType == InlineBlock {
		if b.lastLine != nil && b.lastLine.lineType == JumpStatement && b.lastLine.jumpToBlock != nil {
			if b.lastLine.jumpToBlock.blockType == InlineBlock {
				return b.lastLine.jumpToBlock.getLastJumpIfInlined()
			} else {
				return b.lastLine
			}
		} else {
			return nil
		}
	} else {
		return nil
	}
}

func (b ProgramBlock) debugPrintBlock(t *TranspileState) {
	var str strings.Builder
	b.writeBlockAsJavaScript(t, &str, 1)
	log.Println(str.String())
}

func (b *ProgramBlock) compressConstants(state *TranspileState) (changes int) {
	compressInlineBlock := func(jumpTo *ProgramBlock) {
		if jumpTo.blockType == InlineBlock {
			changes += jumpTo.compressConstants(state)
		} else {
			// The jump could have modified something, so let's remove any constant flags
			// we have for things it could have changed
			for i, m := range jumpTo.modifiersRegisters {
				if m {
					state.Registers[i].SetUnknownInt()
				}
			}

			for i, m := range jumpTo.readsRegisters {
				if m {
					state.Registers[i].MarkRead()
				}
			}
		}
	}

	for line := b.firstLine; line != nil; line = line.next {
		if line.lineType == Statement {
			instruction := line.instruction
			isImmediate := OpCodeInputType[instruction.OpCode]

			// Can we evaluate this expression at compile time?
			if (isImmediate.A || state.Registers[instruction.A].isConst) &&
				(isImmediate.B || state.Registers[instruction.B].isConst) {
				state.updateCPURegisters()
				value, err := OpCodeFunc[line.instruction.OpCode](instruction.A, instruction.B, state.cpu.Registers)
				if err != nil {
					panic(err)
				}

				if state.options.RemoveUnUsedRegisterWrites {
					lastLine := state.Registers[instruction.C].lastUnreadSet
					if lastLine != nil && lastLine.parentBlock == b {
						lastLine.RemoveUnusedLine()
					}
				}

				if instruction.OpCode != SetI {
					instruction.OpCode = SetI
					instruction.A = value
					instruction.B = 0
					changes++
				}
				state.Registers[instruction.C].SetInt(value, line)
			} else if iOp, found := OpCodeImmedateVersion[instruction.OpCode]; found &&
				state.Registers[instruction.B].isConst {

				if !isImmediate.A {
					state.Registers[instruction.A].MarkRead()
				}

				instruction.OpCode = iOp
				instruction.B = state.Registers[instruction.B].value

				state.Registers[instruction.C].SetUnknownInt()
				changes++
			} else {
				if !isImmediate.A {
					state.Registers[instruction.A].MarkRead()
				}

				if !isImmediate.B {
					state.Registers[instruction.B].MarkRead()
				}

				state.Registers[instruction.C].SetUnknownInt()
			}
		} else if line.lineType != DoWhileStatement {
			if line.jumpToBlock != nil {
				compressInlineBlock(line.jumpToBlock)
			}

			if line.elseBlock != nil {
				compressInlineBlock(line.elseBlock)
			}
		}
	}

	// If we're leaving a function, remove any unused writes
	if b.blockType == FunctionBlock && state.options.RemoveUnUsedRegisterWrites {
		for _, register := range state.Registers {
			if register.lastUnreadSet != nil {
				register.lastUnreadSet.RemoveUnusedLine()
			}
		}
	}

	return
}

func (b *ProgramBlock) buildRegisterUsageLists(visited map[*ProgramBlock]bool) {
	// Track if we've built it for this block yet
	if _, found := visited[b]; found {
		return
	}
	visited[b] = true

	for line := b.firstLine; line != nil; line = line.next {
		if line.jumpToBlock != nil {
			line.jumpToBlock.buildRegisterUsageLists(visited)

			for i := 0; i < len(b.modifiersRegisters); i++ {
				b.modifiersRegisters[i] = b.modifiersRegisters[i] || line.jumpToBlock.modifiersRegisters[i]
			}
		}

		if line.elseBlock != nil {
			line.elseBlock.buildRegisterUsageLists(visited)

			for i := 0; i < len(b.modifiersRegisters); i++ {
				b.modifiersRegisters[i] = b.modifiersRegisters[i] || line.elseBlock.modifiersRegisters[i]
			}
		}
	}
}

func (b *ProgramBlock) WriteBlockRegisterModificationList(str *strings.Builder) {
	str.WriteString("// Reads Registers: ")
	count := 0
	for index, b := range b.readsRegisters {
		if b {
			if count > 0 {
				str.WriteString(", ")
			}

			str.WriteString(fmt.Sprintf("R[%d]", index))
			count++
		}
	}
	str.WriteString("\n")

	str.WriteString("// Modifies Registers: ")
	count = 0
	for index, b := range b.modifiersRegisters {
		if b {
			if count > 0 {
				str.WriteString(", ")
			}

			str.WriteString(fmt.Sprintf("R[%d]", index))
			count++
		}
	}
	str.WriteString("\n")
}

type DataType int

const (
	IntType DataType = iota
	BoolType
)

type RegisterState struct {
	isConst       bool         // Is the register holding a constant? If so see `value`
	value         int          // The constant value the register is holding
	dataType      DataType     // Is the register holding a boolean result?
	lastUnreadSet *ProgramLine // The previous line setting the value
}

func (r *RegisterState) SetInt(value int, lineSettingValue *ProgramLine) {
	r.isConst = true
	r.value = value
	r.dataType = IntType
	r.lastUnreadSet = lineSettingValue
}

func (r *RegisterState) SetUnknownInt() {
	r.isConst = false
	r.value = 0
	r.dataType = IntType
}

func (r *RegisterState) SetUnknownBool() {
	r.isConst = false
	r.value = 0
	r.dataType = BoolType
}

func (r *RegisterState) MarkRead() {
	r.lastUnreadSet = nil
}

type TranspileOptions struct {
	CompressConstants           bool // Where there are multiple instructions in a block building a constant, compress into a single SetI
	RemoveEmptyBlocks           bool // Remove blocks with no instructions in them
	RemoveExtraJumps            bool // Remove blocks with a single jump instruction in them
	RemoveUnUsedRegisterWrites  bool // Remove unused register writes
	RewriteRecursionAsLoops     bool // If possible rewrite recursion as loops
	InlineBlocksWherePossible   bool // Inline Blocks which are only called once
	DisplayBlockNumbersOnOutput bool // Should the output include block information?
	DisplayBlockRegisterUse     bool // Should the output include which registers are modified by a block?
}

type TranspileState struct {
	cpu             *CPU                  // The CPU we are using as our basis
	originalCode    Program               // The original code
	blocks          map[int]*ProgramBlock // The program blocks found so far and their entry point
	ipToInstruction map[int]*ProgramLine  // The instruction pointer mapping to a program line
	Registers       []RegisterState       // The Registers
	options         TranspileOptions      // The options

	nextBlockNum int
}

func (cpu *CPU) StartTranspiler(options TranspileOptions) (state *TranspileState) {
	state = &TranspileState{
		cpu,
		make(Program, len(cpu.Program)),
		make(map[int]*ProgramBlock),
		make(map[int]*ProgramLine),
		make([]RegisterState, len(cpu.Registers)),
		options,
		0,
	}

	// Copy the original code across
	copy(state.originalCode, cpu.Program)

	// Set the instruction pointer const
	for i := 0; i < len(cpu.Registers); i++ {
		state.Registers[i].SetInt(cpu.Registers[i], nil)
	}

	// Normalise the CPU code for the transpile
	// this is mainly reordering input args
	cpu.normaliseProgramForTranspile()

	return
}

func (t *TranspileState) updateCPURegisters() {
	for i, register := range t.Registers {
		if register.isConst {
			t.cpu.Registers[i] = register.value
		} else {
			t.cpu.Registers[i] = -99999
		}
	}
}

func (t *TranspileState) resetRegistersToUnknownState() {
	for i := 0; i < len(t.Registers); i++ {
		// When entering a block, lose the constant flag as we dont know if it is constant
		// this block might be called a second time after what was constant has been changed!
		t.Registers[i].isConst = false
		t.Registers[i].dataType = IntType
		t.Registers[i].value = 0
		t.Registers[i].lastUnreadSet = nil
	}
}

func (t *TranspileState) processBlock(startIP int, stopIP int, callingLine *ProgramLine) (block *ProgramBlock) {
	block, found := t.blocks[startIP]
	if found {
		return
	}

	block = &ProgramBlock{
		t.nextBlockNum,
		FunctionBlock,
		make([]*ProgramLine, 0),
		nil, nil,
		make([]bool, len(t.Registers)), make([]bool, len(t.Registers)),
	}
	t.nextBlockNum++
	t.blocks[startIP] = block

	// Is this instruction already in a block? If so, then we need to move it into it's own block
	// and replace the original call site with a jump to this new version of the block
	if line, found := t.ipToInstruction[startIP]; found {
		line.RemoveAllFollowingLines()
		line.lineType = JumpStatement
		line.instruction = nil
		line.jumpToBlock = block
	}

	t.Registers[t.cpu.InstructionPointerRegister].SetInt(startIP, nil)
	for {
		ip := t.Registers[t.cpu.InstructionPointerRegister].value

		// Stop executing
		if ip >= stopIP || ip < 0 || ip >= len(t.cpu.Program) {
			break
		}

		// Convert the instruction to a line and link it into the parentBlock
		line := NewProgramLine(ip, &t.cpu.Program[ip], block)
		t.ipToInstruction[ip] = line
		if block.lastLine == nil {
			block.firstLine = line
		} else {
			block.lastLine.next = line
			line.previous = block.lastLine
		}
		block.lastLine = line

		// Is this a program line change?
		isImmediate := OpCodeInputType[line.instruction.OpCode]

		if line.instruction.C == t.cpu.InstructionPointerRegister {
			if (isImmediate.A || t.Registers[line.instruction.A].isConst) &&
				(isImmediate.B || t.Registers[line.instruction.B].isConst) {
				// Update the Registers to present our known state
				t.updateCPURegisters()

				value, err := OpCodeFunc[line.instruction.OpCode](line.instruction.A, line.instruction.B, t.cpu.Registers)
				if err != nil {
					panic(err)
				}

				line.lineType = JumpStatement
				line.instruction = nil
				line.jumpToBlock = t.processBlock(value+1, stopIP, line)

				for i := 0; i < len(block.modifiersRegisters); i++ {
					if !block.modifiersRegisters[i] {
						block.readsRegisters[i] = block.readsRegisters[i] || line.jumpToBlock.readsRegisters[i]
					}
				}
				return // This block ends because fo the jump
			} else if (line.instruction.OpCode == AddI) &&
				(t.Registers[line.instruction.A].dataType == BoolType) {

				line.lineType = IfStatement
				line.jumpToBlock = t.processBlock(ip+2, stopIP, line)
				for i := 0; i < len(block.modifiersRegisters); i++ {
					if !block.modifiersRegisters[i] {
						block.readsRegisters[i] = block.readsRegisters[i] || line.jumpToBlock.readsRegisters[i]
					}
				}

				line.elseBlock = t.processBlock(ip+1, stopIP, line)
				for i := 0; i < len(block.modifiersRegisters); i++ {
					if !block.modifiersRegisters[i] {
						block.readsRegisters[i] = block.readsRegisters[i] || line.elseBlock.readsRegisters[i]
					}
				}
				// If statement found!
				return
			} else {
				panic(fmt.Sprintf("Non constant $IP change at %d: %s", ip, line.instruction))
			}
		} else {
			// Flag if this block reads a register before it has written to it
			if !isImmediate.A && !block.modifiersRegisters[line.instruction.A] {
				block.readsRegisters[line.instruction.A] = true
			}

			if !isImmediate.B && !block.modifiersRegisters[line.instruction.B] {
				block.readsRegisters[line.instruction.B] = true
			}

			// Flag this block has written to a register
			block.modifiersRegisters[line.instruction.C] = true

			if line.instruction.OpCode.isComparator() {
				t.Registers[line.instruction.C].SetUnknownBool()
			} else {
				t.Registers[line.instruction.C].SetUnknownInt()
			}
		}

		// Increment the IP counter
		t.Registers[t.cpu.InstructionPointerRegister].value++
	}

	return
}

func (t *TranspileState) buildCalledByListsForBlocks() {
	for _, b := range t.blocks {
		b.calledBy = make([]*ProgramLine, 0)
	}

	for _, b := range t.blocks {
		for line := b.firstLine; line != nil; line = line.next {
			if line.jumpToBlock != nil {
				line.jumpToBlock.calledBy = append(line.jumpToBlock.calledBy, line)
			}

			if line.elseBlock != nil {
				line.elseBlock.calledBy = append(line.elseBlock.calledBy, line)
			}
		}
	}
}

// Some blocks, might just be calls to other blocks, so we remove them at this point by inling them
func (t *TranspileState) removeJumpOnlyBlocks() (changes int) {
	for _, b := range t.blocks {
		if b.firstLine != nil &&
			b.firstLine == b.lastLine &&
			b.firstLine.lineType == JumpStatement &&
			len(b.calledBy) > 0 {

			targetBlock := b.firstLine.jumpToBlock
			targetBlock.removeCallFrom(b.firstLine) // Remove our jump call
			b.replaceJumpsToBlockWithBlock(targetBlock)
			b.firstLine = nil
			b.lastLine = nil
			b.calledBy = []*ProgramLine{}
			changes++
		}
	}

	return
}

func (t *TranspileState) removeEmptyBlocks() (changes int) {
	for _, b := range t.blocks {
		if b.firstLine == nil && len(b.calledBy) > 0 {
			b.replaceJumpsToBlockWithBlock(nil)
			changes++
		}
	}

	return
}

func (t *TranspileState) inlineCallOnceBlocks() (changes int) {
	for _, b := range t.blocks {
		if len(b.calledBy) == 1 {
			if b.calledBy[0].lineType == JumpStatement && b.calledBy[0].previous != nil && b.firstLine != nil {
				b.firstLine.moveToNewBlockAfter(b.calledBy[0].previous)
				b.calledBy[0].RemoveUnusedLine()
				b.calledBy = []*ProgramLine{}
				changes++
			} else if b.blockType != InlineBlock {
				b.blockType = InlineBlock
				changes++
			}
		}
	}

	return
}

func (t *TranspileState) inlineBlocksAfterIfsWhichEndBothBranches() (changes int) {
	for _, b := range t.blocks {
		for line := b.firstLine; line != nil; line = line.next {
			if line.lineType == IfStatement && // If this is an if statement
				line.jumpToBlock != nil && line.elseBlock != nil { // which has two branches

				branchT := line.jumpToBlock
				branchF := line.elseBlock

				if branchT == branchF {
					line.lineType = JumpStatement
					line.elseBlock = nil
					branchT.removeCallFrom(line)                      // remove 2 calls (if and else branches)
					branchT.calledBy = append(branchT.calledBy, line) // re-add just 1 call

					changes++
					break
				} else if branchT.blockType == InlineBlock && branchF.blockType == InlineBlock {
					// Both are inline statements
					jumpT := branchT.getLastJumpIfInlined()
					jumpF := branchF.getLastJumpIfInlined()

					if jumpT != nil && jumpF != nil && jumpT.jumpToBlock == jumpF.jumpToBlock {
						jumpT.moveToNewBlockAfter(line)
						jumpF.RemoveUnusedLine()
						changes++
						break
					}

				} else if branchT.blockType == InlineBlock && branchF.blockType == FunctionBlock {
					// Just the true side is inlined
					jumpT := branchT.getLastJumpIfInlined()

					if jumpT != nil && jumpT.jumpToBlock == branchF {
						jumpT.moveToNewBlockAfter(line)
						line.elseBlock = nil
						branchF.removeCallFrom(line) // remove the else call
						changes++
						break
					}
				} else if branchT.blockType == FunctionBlock && branchF.blockType == InlineBlock {
					// Just the true side is inlined
					jumpF := branchF.getLastJumpIfInlined()

					if jumpF != nil && jumpF.jumpToBlock == branchT {
						jumpF.moveToNewBlockAfter(line)
						line.jumpToBlock = nil
						branchT.removeCallFrom(line) // remove the else call
						changes++
						break
					}
				}
			}
		}
	}

	return
}

func (t *TranspileState) compressConstants() (changes int) {
	for _, b := range t.blocks {
		if b.blockType == FunctionBlock {
			t.resetRegistersToUnknownState()
			changes += b.compressConstants(t)
		}
	}

	return
}

func (t *TranspileState) findAndRewriteLoops() (changes int) {
	for _, b := range t.blocks {
		if len(b.calledBy) == 2 && // It's only called twice
			(b.calledBy[0].parentBlock == b || b.calledBy[1].parentBlock == b) && // and at least one of those calls is the block (recursion)
			b.lastLine != nil && b.lastLine.lineType == IfStatement { // and the block ends with an IF
			// Then if one of the two branches of the if end with the call back to this block
			// it's then recursion!

			var contentToAddAfterLoop *ProgramBlock
			isLoop := false
			loopConditionInverse := false

			if b.lastLine.jumpToBlock == b {
				contentToAddAfterLoop = b.lastLine.elseBlock
				isLoop = true
			} else if b.lastLine.elseBlock == b {
				isLoop = true
				loopConditionInverse = true
				contentToAddAfterLoop = b.lastLine.jumpToBlock
			}

			if isLoop {
				// Find the calling line
				callingLine := b.calledBy[0]
				if callingLine.parentBlock == b {
					callingLine = b.calledBy[1]
				}

				callingLine.lineType = DoWhileStatement
				callingLine.invertCondition = loopConditionInverse
				// Change this block to an inline block
				b.blockType = InlineBlock
				b.calledBy = []*ProgramLine{callingLine}
				callingLine.instruction = b.lastLine.instruction

				if contentToAddAfterLoop != nil && contentToAddAfterLoop.firstLine != nil {
					contentToAddAfterLoop.firstLine.moveToNewBlockAfter(callingLine)
				}
				b.lastLine.RemoveUnusedLine()
				changes++
			}
		}
	}

	return
}

// Transpile's the program to pseudo code to easier reading
func (t *TranspileState) Run() (string, error) {
	var str strings.Builder

	initalRegisterState := make([]RegisterState, len(t.Registers))
	copy(initalRegisterState, t.Registers)

	// Build the program out into the minimum number of blocks it can exist as
	startingBlock := t.processBlock(0, 37, nil)
	t.buildCalledByListsForBlocks()
	startingBlock.buildRegisterUsageLists(make(map[*ProgramBlock]bool))

	// Start making more readable by humans, each set of changes could result in
	// the previous operations being able to do more, so keep looping until no more
	// changes are made to the structure of the program
	for changes := 1; changes > 0; {
		changes = 0
		if t.options.RemoveExtraJumps {
			changes += t.removeJumpOnlyBlocks()
		}

		if t.options.RemoveEmptyBlocks {
			changes += t.removeEmptyBlocks()
		}

		if t.options.InlineBlocksWherePossible {
			changes += t.inlineCallOnceBlocks()
			changes += t.inlineBlocksAfterIfsWhichEndBothBranches()
		}

		if t.options.CompressConstants {
			// We need initial register state here for the main func
			copy(t.Registers, initalRegisterState)
			startingBlock.compressConstants(t)

			// Then rest as normal
			t.resetRegistersToUnknownState()
			changes += t.compressConstants()
		}
	}

	// This occurs outside the other loops, as inlinng the blocks could cause other optimisations to
	// wrongly apply (mainly compress constants)
	if t.options.RewriteRecursionAsLoops {
		for t.findAndRewriteLoops() > 0 {}
	}

	// Copy the blocks to a slice, then sort
	blocks := make([]*ProgramBlock, len(t.blocks))
	i := 0
	for _, b := range t.blocks {
		blocks[i] = b
		i++
	}
	sort.Slice(blocks, func(i, j int) bool {
		if blocks[i] == startingBlock {
			return true
		} else if blocks[i].blockNum < blocks[j].blockNum {
			return true
		} else {
			return false
		}
	})

	str.WriteString("// Registers\nvar R = [")
	for i := 0; i < len(t.Registers); i++ {
		if i > 0 {
			str.WriteString(", ")
		}

		str.WriteRune('0')
	}
	str.WriteString("]\n\n")

	for _, b := range blocks {
		if b.blockType == InlineBlock {
			// This should be written inside a function block somewhere
			continue
		}

		if len(b.calledBy) == 0 && b != startingBlock {
			// These blocks have been optimised out of the program
			continue
		}

		indent := 1

		if b == startingBlock {
			if t.options.DisplayBlockRegisterUse {
				b.WriteBlockRegisterModificationList(&str)
			}
			str.WriteString("function main() {\n")
		} else {
			if t.options.DisplayBlockNumbersOnOutput {
				str.WriteString("// Called By: ")
				for index, line := range b.calledBy {
					if index > 0 {
						str.WriteString(", ")
					}
					str.WriteString(fmt.Sprintf("block%d", line.parentBlock.blockNum))
				}
				str.WriteString("\n")
			}

			if t.options.DisplayBlockRegisterUse {
				b.WriteBlockRegisterModificationList(&str)
			}

			str.WriteString(fmt.Sprintf("function block%d() {\n", b.blockNum))
		}

		b.writeBlockAsJavaScript(t, &str, indent)

		str.WriteString("}\n\n")
	}

	return str.String(), nil
}
