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
