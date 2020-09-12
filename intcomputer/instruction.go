package intcomputer

import (
	"fmt"
	"strings"
)

type Instruction struct {
	Opcode         int
	ParamAddrModes []int
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
		case Halt:
			ret = "Halt"
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

func decode(ins int) *Instruction {
	// return opcode, {param-1-addrMode, parma-2-addrMode, ...}
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
			ret[1] = in % 10
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
