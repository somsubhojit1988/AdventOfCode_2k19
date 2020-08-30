package intcomputer

import (
	"fmt"
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

func (p *Instruction) toString() string {
	return fmt.Sprintf("{opcode: %s addressingModes:%v}",
		p.getInstructionString(), p.AddrModes)
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
		m3 := n % 10 //param-3: addressing mode - 10_000s place
		return opcode, m1, m2, m3
	}(n)

	return &Instruction{
		Opcode:    opcode,
		AddrModes: []int{m1, m2, m3},
	}
}

func (p *Instruction) getInstructionString() string {
	var res string
	switch p.Opcode {
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
		res = fmt.Sprintf("Unsupported opcode: %v", p.Opcode)
	}
	return res
}

// Memory ... used to store instruction set and fetch parameter values
type Memory struct {
	InstructionSet   []int
	MemPtr           Ptr
	CurrentAddrModes []int
	logger           *Logger
}

// InitMemory ... initialize memory with instruction set, memory pointer to index-0
func InitMemory(ins, modes []int, logger *Logger) *Memory {
	return &Memory{
		InstructionSet:   ins,
		MemPtr:           Ptr(0),
		CurrentAddrModes: modes,
		logger:           logger,
	}
}
func (m *Memory) dump() string {
	return fmt.Sprintf("%v", m.InstructionSet)
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

// ReadParams ... reads the parameter values based on currentAddressingModes
// .. [Param:n, Param:n-1, ..,Param:0]
func (m *Memory) ReadParams() []int {
	return []int{
		m.read(m.CurrentAddrModes[0], m.MemPtr+1),
		m.read(m.CurrentAddrModes[1], m.MemPtr+2),
		m.read(ImmediateMode, m.MemPtr+3),
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
	m.logger.log(fmt.Sprintf("[Memory::ExecuteInstruction(%v)] memory: %v",
		opcode, m.toString()))

	params := m.ReadParams()
	res := -1
	switch opcode {

	case Add:
		res = params[0] + params[1]
		// m.logger.log(fmt.Sprintf("[Memory]Add: params: %v result %v at %v",
		// 	params, res, params[0]))
		m.Write(res, params[2])

	case Mul:
		res = params[0] * params[1]
		// m.logger.log(fmt.Sprintf("[Memory] storing Mul result %v at %v",
		// 	res, params[0]))
		m.Write(res, params[2])

	case Input:
		v := 1 // TODO: take user input
		// m.logger.log(fmt.Sprintf("[Memory] storing Input value %v at %v",
		// 	v, params[0]))
		m.Write(v, params[0])

	case Output:
		res = m.read(params[0], Ptr(params[0]))
		// m.logger.log(fmt.Sprintf("[Memory] read Output value %v from %v",
		// 	res, params[0]))
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
func CreateIntComputer(instructions []int, logger *Logger) *IntComputer {
	ins := CreateInstruction(instructions[0])
	return &IntComputer{
		State:          InitMemory(instructions, ins.AddrModes, logger),
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

// Run ... executes the provided instructions and produces an int result
func (c *IntComputer) Run() int {
	for !c.isHalted() {
		ins := CreateInstruction(c.State.read(ImmediateMode, c.InstructionPtr))
		c.logger.log(fmt.Sprintf("\n\n[IntComputer inPtr = %v] Executing:  instruction= %v \n",
			c.InstructionPtr, ins.toString()))
		c.State.CurrentAddrModes = ins.AddrModes
		c.State.ExecuteInstruction(ins.Opcode)

		d := 0
		switch ins.Opcode {
		case Add, Mul:
			d = 4
		case Input, Output:
			d = 2
		case Halt:
			d = c.MemorySize()
		default:
			panic(fmt.Sprintf("IntComputer: Unsupported opcode:%v\n", ins.Opcode))
		}

		c.InstructionPtr = Ptr(int(c.InstructionPtr) + d)
		c.State.IncrementPtr(d)
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
