package intcomputer

import (
	"fmt"
	"strings"
)

const (
	// Addressing Modes
	Position  = 0
	Immediate = 1

	// Instructions
	Unsupported = -1
	Add         = 1
	Mul         = 2
	Input       = 3
	Output      = 4
	JmpIfTrue   = 5
	JmpIfFalse  = 6
	LessThan    = 7
	Equals      = 8
	Halt        = 99
)

type InputMethod func() int

type OutputMethod func(int)

type Instruction struct {
	Opcode         int
	ParamAddrModes []int
}

type Memory struct {
	Storage []int
	MemPtr  int
	logger  *Logger
}

type Logger struct {
	buffer []string
}

func CreateLogger() *Logger {
	return &Logger{buffer: []string{}}
}
func (l *Logger) log(msg string) {
	l.buffer = append(l.buffer, msg)
}

func (l *Logger) Logs() []string {
	return l.buffer
}

type IntComputer struct {
	Mem     *Memory
	InPtr   int
	InFunc  InputMethod
	OutFunc OutputMethod
	logger  *Logger
}

func decode(ins int) *Instruction {
	// M3 M2 M1 OP OP
	parse := func(in int) (int, []int) {
		op := in % 100
		in /= 100

		var ret []int
		switch op {
		case Add, Mul, LessThan, Equals:
			// 3 parameters
			ret = make([]int, 3)
			ret[0] = in % 10
			in /= 10
			ret[1] = in % 10
			in /= 10
			ret[2] = Immediate // param value gives address to store results to
		case Input, Output:
			// 1 parameter
			ret = make([]int, 1)
			if op == Input {
				ret[0] = Immediate
			} else {
				ret[0] = in % 10
			}
		case JmpIfTrue, JmpIfFalse:
			// 2 parameters
			ret = make([]int, 2)
			ret[0] = in % 10
			in /= 10
			ret[1] = Immediate
		default:
			// 0 params ex: Halt
			ret = make([]int, 0)
		}
		return op, ret
	}
	op, addrModes := parse(ins)
	return &Instruction{
		Opcode:         op,
		ParamAddrModes: addrModes,
	}
}

func (i *Instruction) String() string {

	opToString := func(op int) string {
		var ret string
		switch op {
		case Add:
			ret = "Add"
		case Mul:
			ret = "Mul"
		case Input:
			ret = "Input"
		case Output:
			ret = "Output"
		case JmpIfTrue:
			ret = "JumpIfTrue"
		case JmpIfFalse:
			ret = "JumpIfFalse"
		case LessThan:
			ret = "LessThan"
		case Equals:
			ret = "Equals"
		default:
			ret = "Unsupported instruction"
		}
		return ret
	}

	addrModeToString := func(m int) string {
		var ret string
		switch m {
		case Position:
			ret = "Position"
		case Immediate:
			ret = "Immediate"
		}
		return ret
	}

	strB := strings.Builder{}
	for _, m := range i.ParamAddrModes {
		s := fmt.Sprintf("%d (%s) ", m, addrModeToString(m))
		strB.WriteString(s)
	}

	return fmt.Sprintf("Opcode= %d {%s} AddressingModes [p1, p2, ...] = %s",
		i.Opcode, opToString(i.Opcode), strB.String())
}

func (m *Memory) Size() int {
	return len(m.Storage)
}

func (m *Memory) String() string {
	dump := func() string {
		sb := &strings.Builder{}
		cntr := 0
		for i, v := range m.Storage {
			cntr++
			sb.WriteString(fmt.Sprintf("%d (%d) ", v, i))
			if cntr >= 10 {
				cntr = 0
				sb.WriteString("\n")
			}
		}
		return fmt.Sprintf("[ %s ]", sb.String())
	}
	return fmt.Sprintf("MEM: ptr= %d\n%s\n", m.MemPtr, dump())
}

func (m *Memory) read(addrMode, d int) (int, error) {
	var ret int
	if m.MemPtr+d >= m.Size() {
		return -1, fmt.Errorf("MEMREAD (addressing-mode= %d, addr = %d) Out of range",
			addrMode, m.MemPtr+d)
	}
	x := m.Storage[m.MemPtr+d]
	switch addrMode {
	case Position:
		if x >= m.Size() {
			return -1, fmt.Errorf("MEMREAD (addressing-mode= %d, addr = %d) Out of range",
				addrMode, x)
		}
		ret = m.Storage[x]
	case Immediate:
		ret = x
	}
	return ret, nil
}

func (m *Memory) readAddress(ptr int) (int, error) {
	if ptr < 0 || ptr >= m.Size() {
		return -1, fmt.Errorf("MEMREAD (addr = %d) Out of range", ptr)
	}
	return m.Storage[ptr], nil
}

func (m *Memory) write(v, ptr int) error {
	if ptr >= m.Size() {
		return fmt.Errorf("MEMWRITE (addr = %d  v= %d) Out of range",
			ptr, v)
	}
	m.Storage[ptr] = v
	return nil
}

func (m *Memory) opcodeFetch() (int, error) {
	return m.read(Immediate, 0)
}

func (c *IntComputer) readParams(ins *Instruction) ([]int, error) {
	ret := make([]int, len(ins.ParamAddrModes))
	for i, m := range ins.ParamAddrModes {
		v, err := c.Mem.read(m, i+1)
		if err != nil {
			return ret, err
		}
		ret[i] = v
	}
	return ret, nil
}

func (c *IntComputer) storeResult(v, ptr int) error {
	return c.Mem.write(v, ptr)
}

func (c *IntComputer) add(ins *Instruction) error {
	params, err := c.readParams(ins)
	if err != nil {
		return err
	}
	if err := c.storeResult(params[0]+params[1], params[2]); err != nil {
		return err
	}
	c.InPtr += 4
	return nil
}

func (c *IntComputer) mul(ins *Instruction) error {
	params, err := c.readParams(ins)
	if err != nil {
		return err
	}
	if err := c.storeResult(params[0]*params[1], params[2]); err != nil {
		return err
	}
	c.InPtr += 4
	return nil
}

func (c *IntComputer) input() error {
	ptr, err := c.Mem.read(Immediate, 1)
	if err != nil {
		return err
	}
	v := c.InFunc()
	err = c.Mem.write(v, ptr)
	if err != nil {
		return err
	}
	c.InPtr += 2
	return err
}

func (c *IntComputer) output(ins *Instruction) error {
	v, err := c.Mem.read(ins.ParamAddrModes[0], 1)
	if err != nil {
		return err
	}
	c.OutFunc(v)
	c.InPtr += 2
	return err
}

func (c *IntComputer) halt() {
	c.InPtr = c.Mem.Size()
}

func (c *IntComputer) isHalted() bool {
	return c.InPtr >= c.Mem.Size()
}

func (c *IntComputer) execute() error {
	code, err := c.Mem.opcodeFetch()
	if err != nil {
		return err
	}
	ins := decode(code)

	c.logger.log(fmt.Sprintf("[IntComputer] {InsPtr: %d} Execute: %v",
		c.InPtr, ins))
	switch ins.Opcode {
	case Add:
		c.add(ins)
	case Mul:
		c.mul(ins)
	case Input:
		c.input()
	case Output:
		c.output(ins)
	case JmpIfTrue:
	case JmpIfFalse:
	case LessThan:
	case Equals:
	case Halt:
		c.halt()
	}
	c.Mem.MemPtr = c.InPtr
	return err
}

func CreateIntComputer(instructions []int, logger *Logger,
	in InputMethod, out OutputMethod) *IntComputer {
	return &IntComputer{
		Mem:     &Memory{Storage: instructions, MemPtr: 0, logger: logger},
		InPtr:   0,
		InFunc:  in,
		OutFunc: out,
		logger:  logger,
	}
}

func (c *IntComputer) Store(val, ptr int) error {
	return c.Mem.write(val, ptr)
}

func (c *IntComputer) ReadMemory(ptr, n int) ([]int, error) {
	ret := make([]int, 5)
	for i := 0; i < n; i++ {
		v, err := c.Mem.readAddress(ptr + i)
		if err != nil {
			return nil, err
		}
		ret[i] = v
	}
	return ret, nil
}

func (c *IntComputer) Run() error {
	var err error
	for !c.isHalted() && err == nil {
		err = c.execute()
	}
	return err
}
