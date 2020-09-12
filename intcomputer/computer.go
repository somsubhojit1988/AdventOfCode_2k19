package intcomputer

import (
	"fmt"
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
	ret := make([]string, len(l.buffer))
	copy(ret, l.buffer)
	return ret
}

func (l *Logger) clear() {
	l.buffer = []string{}
}

type IntComputer struct {
	Mem     *Memory
	InPtr   int
	InFunc  InputMethod
	OutFunc OutputMethod
	logger  *Logger

	//  xxxx xxxx xxxx b3 b2 b1 Halt
	flags int16
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

func (c *IntComputer) jmpIfTrue(ins *Instruction) error {
	params, err := c.readParams(ins)
	if err != nil {
		return err
	}
	v, ptr := params[0], params[1]
	if v != 0 {
		c.InPtr = ptr
	} else {
		c.InPtr += 3
	}
	return err
}

func (c *IntComputer) jmpIfFalse(ins *Instruction) error {
	params, err := c.readParams(ins)
	if err != nil {
		return err
	}
	v, ptr := params[0], params[1]
	if v == 0 {
		c.InPtr = ptr
	} else {
		c.InPtr += 3
	}
	return err
}

func (c *IntComputer) lt(ins *Instruction) error {
	params, err := c.readParams(ins)
	if err != nil {
		return err
	}
	if params[0] < params[1] {
		err = c.Mem.write(1, params[2])
	} else {
		err = c.Mem.write(0, params[2])
	}
	c.InPtr += 4
	return err
}

func (c *IntComputer) eq(ins *Instruction) error {
	params, err := c.readParams(ins)
	if err != nil {
		return err
	}
	if params[0] == params[1] {
		err = c.Mem.write(1, params[2])
	} else {
		err = c.Mem.write(0, params[2])
	}
	c.InPtr += 4
	return err
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
	// set flag-0th bit of
	c.flags |= 0x01
}

func (c *IntComputer) isHalted() bool {
	return ((c.flags & 0x01) != 0)
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
		c.jmpIfTrue(ins)
	case JmpIfFalse:
		c.jmpIfFalse(ins)
	case LessThan:
		c.lt(ins)
	case Equals:
		c.eq(ins)
	case Halt:
		c.halt()
	default:
		return fmt.Errorf("Unsupported opcode")
	}
	c.Mem.memPtr = c.InPtr
	return err
}

func CreateIntComputer(instructions []int, logger *Logger,
	in InputMethod, out OutputMethod) *IntComputer {
	ins := make([]int, len(instructions))
	copy(ins, instructions)
	return &IntComputer{
		Mem:     &Memory{storage: ins, memPtr: 0, logger: logger},
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
	return ret[:n], nil
}

func (c *IntComputer) Run() error {
	var err error
	for !c.isHalted() && err == nil {
		err = c.execute()
	}
	return err
}

func (c *IntComputer) Reset() {
	c.logger.clear()
	c.Mem = &Memory{storage: []int{99}, memPtr: 0, logger: c.logger}
	c.InPtr = 0
}

func (c *IntComputer) Program(instructions []int) {
	c.Reset()
	c.Mem = &Memory{storage: instructions, memPtr: 0, logger: c.logger}
}
