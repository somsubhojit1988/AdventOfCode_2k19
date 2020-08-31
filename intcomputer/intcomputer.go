package intcomputer

import (
	"fmt"
	"strings"
)

const (

	// PositionMode ... indicates parameter addressing mode: position mode
	PositionMode = 0

	// ImmediateMode ... indicates parameter addressing mode: immediate mode
	ImmediateMode = 1

	// Add ... opcode for add
	Add = 1
	// Mul ... opcode for multiply
	Mul = 2
	// Input ... opcode for writting to memory; spec: input
	// ... takes a single integer as input and saves it to
	// ... the position given by its only parameter
	Input = 3
	// Output ... opcode for read memory;
	// ... spec: outputs the value of its only parameter.
	Output = 4
	// Halt ... opcode to halt execution
	Halt = 99
)

//Ptr ... unsigned int for array indexing
type Ptr uint

// Instruction ... represents a single instruction pointed to by the current
// ... instruction ptr
type Instruction struct {
	Opcode    int
	AddrModes []int
}

// CreateInstruction ... creates a new Instruction with supplied int instruction
// ...extract opcode and parameter addressing modes
func CreateInstruction(n int) *Instruction {

	opcode, m1, m2, m3 := func(n int) (int, int, int, int) {
		opcode := n % 100 // opcode: unit, 10s
		n /= 100
		m1 := n % 10 // param-1: addressing mode - 100s place
		n /= 10
		m2 := n % 10 //param-2: addressing mode - 1000s place
		n /= 10
		// m3 := n % 10 //param-3: addressing mode - 10_000s place
		return opcode, m1, m2, ImmediateMode
	}(n)

	return &Instruction{
		Opcode:    opcode,
		AddrModes: []int{m1, m2, m3},
	}
}

func getInstructionString(opcode int) string {
	var res string
	switch opcode {
	case Add:
		res = "Add"
	case Mul:
		res = "Multiply"
	case Input:
		res = "Input"
	case Output:
		res = "Output"
	case Halt:
		res = "Halt"
	default:
		res = fmt.Sprintf("Unsupported opcode: %v", opcode)
	}
	return res
}

func (p *Instruction) toString() string {
	return fmt.Sprintf("{opcode: %s addressingModes:%v}",
		getInstructionString(p.Opcode), p.AddrModes)
}

type InputFunc func() int

type OutputFunc func(int)

// Memory ... used to store instruction set and fetch parameter values
type Memory struct {
	InstructionSet   []int
	MemPtr           Ptr
	CurrentAddrModes []int
	logger           *Logger
	inFunc           InputFunc
	outFunc          OutputFunc
	dumpEnabled      bool
}

// InitMemory ... initialize memory with instruction set, memory pointer to index-0
func InitMemory(ins, modes []int, logger *Logger, inFunc InputFunc, outFunc OutputFunc) *Memory {
	return &Memory{
		InstructionSet:   ins,
		MemPtr:           Ptr(0),
		CurrentAddrModes: modes,
		logger:           logger,
		inFunc:           inFunc,
		outFunc:          outFunc,
		dumpEnabled:      true,
	}
}
func (m *Memory) dump() string {
	if !m.dumpEnabled {
		return ""
	}

	var b strings.Builder

	for i, x := range m.InstructionSet {
		fmt.Fprintf(&b, "%d: %d, ", i, x)
	}

	s := b.String()   // no copying
	s = s[:b.Len()-2] // no copying (removes trailing ", ")
	return fmt.Sprintf("[ %s ]", s)
}

func (m *Memory) toString() string {
	return fmt.Sprintf("{memPtr: %v, currentAddressingModes: %v  }",
		m.MemPtr, m.CurrentAddrModes)
}

func (m *Memory) read(mode int, ptr Ptr) int {
	switch mode {
	case PositionMode:
		return m.InstructionSet[m.InstructionSet[ptr]]
	case ImmediateMode:
		return m.InstructionSet[ptr]
	default:
		panic("mem-read: Unsupported addressing mode!")
	}
}

// Write ... writes value: `v` at position: `pos`
func (m *Memory) Write(v, pos int) {
	m.InstructionSet[pos] = v
}

// IncrementPtr ... increments the memory pointer by supplied displacement
func (m *Memory) IncrementPtr(d int) {
	if d >= 0 {
		m.MemPtr += Ptr(d)
	}
}

// ExecuteInstruction ... executes the supplied instruction
func (m *Memory) ExecuteInstruction(opcode int) int {
	m.logger.log(fmt.Sprintf("[Memory::ExecuteInstruction(%v{%v})] memory: %v",
		opcode, getInstructionString(opcode), m.toString()))

	readParams := func(addrModes []int) []int {
		ret := make([]int, len(addrModes))
		for i := 0; i < len(addrModes); i++ {
			ret[i] = m.read(m.CurrentAddrModes[i], Ptr(int(m.MemPtr)+i+1))
		}
		return ret
	}

	res := -1
	switch opcode {

	case Add:
		ps := readParams(m.CurrentAddrModes[:])
		res := ps[0] + ps[1]
		m.logger.log(fmt.Sprintf("[Memory]Add: result (%v) -> Memory[ %v ]",
			res, ps[2]))
		m.Write(res, ps[2])
	case Mul:
		ps := readParams(m.CurrentAddrModes[:])
		res := ps[0] * ps[1]
		m.logger.log(fmt.Sprintf("[Memory]Mul: result( %v) -> Memory[ %v]",
			res, ps[2]))
		m.Write(res, ps[2])

	case Input:
		if m.CurrentAddrModes[0] != PositionMode {
			panic(fmt.Sprintf("[Memory] Execute:i/p  Illegal addressingMode: %v",
				m.CurrentAddrModes[0]))
		}

		v, ptr := m.inFunc(), m.read(ImmediateMode, m.MemPtr+1)
		m.logger.log(fmt.Sprintf("[Memory]Input %v -> Memory[%v]", v, ptr))
		m.Write(v, ptr)

	case Output:
		res := m.read(m.CurrentAddrModes[0], m.MemPtr+1)
		m.logger.log(fmt.Sprintf("[Memory] Output <- Memory[addr mode: %v, ptr: %v] (%v) ",
			m.CurrentAddrModes[0], m.MemPtr+1, res))
		m.outFunc(res)
	}
	m.logger.log(fmt.Sprintf("[Memory-dump]\n%v\n", m.dump()))
	return res
}

// IntComputer ... intcomputer structure that received instructions and produces o/p
type IntComputer struct {
	State              *Memory
	InstructionPtr     Ptr
	CurrentInstruction *Instruction
	logger             *Logger
}

// CreateIntComputer ... create new int computer
func CreateIntComputer(instructions []int,
	logger *Logger,
	inFunc InputFunc,
	outFunc OutputFunc) *IntComputer {

	ins := CreateInstruction(instructions[0])
	return &IntComputer{
		State: InitMemory(instructions,
			ins.AddrModes, logger,
			inFunc, outFunc),
		InstructionPtr: Ptr(0),
		logger:         logger,
	}
}

func (c *IntComputer) isHalted() bool {
	return int(c.InstructionPtr) >= c.MemorySize()
}

// MemorySize ... gives the size of the int computer's memory
func (c *IntComputer) MemorySize() int {
	return len(c.State.InstructionSet)
}

func (c *IntComputer) execute(currInstruction *Instruction) {
	c.State.CurrentAddrModes = currInstruction.AddrModes
	c.State.ExecuteInstruction(currInstruction.Opcode)
}

func (c *IntComputer) updateState() {
	ins := CreateInstruction(c.State.read(ImmediateMode, c.InstructionPtr))
	c.logger.log(fmt.Sprintf("[IntComputer]execyuting instrunction: %v", ins.toString()))
	c.execute(ins)

	d := 0
	switch ins.Opcode {
	case Add, Mul:
		d = 4
	case Input, Output:
		d = 2
	case Halt:
		d = c.MemorySize()
	default:
		panic(fmt.Sprintf("IntComputer [InstructionPtr= %v]: Unsupported opcode:%v\n",
			c.InstructionPtr, ins.Opcode))
	}
	c.InstructionPtr = Ptr(int(c.InstructionPtr) + d)
	c.State.IncrementPtr(d)
}

// Run ... executes the provided instructions and produces an int result
func (c *IntComputer) Run() int {
	for !c.isHalted() {
		c.updateState()
	}

	return c.State.read(ImmediateMode, 0)
}

// Program ... store val at position= pos, in memory
func (c *IntComputer) Program(pos, val int) {
	c.State.Write(val, pos)
}

// Looger ... for intcomputer
type Logger struct {
	logBuffer []string
}

func CreateLooger() *Logger {
	return &Logger{
		logBuffer: []string{},
	}
}

func (l *Logger) log(msg string) {
	l.logBuffer = append(l.logBuffer, msg)
}

func (l *Logger) Logs() []string {
	return l.logBuffer
}

func (l *Logger) Clear() {
	l.logBuffer = []string{}
}
