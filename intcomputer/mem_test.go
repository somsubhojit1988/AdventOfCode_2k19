package intcomputer

import (
	"testing"
)

func TestCreateInstruction(t *testing.T) {
	tt := []struct {
		n   int
		ins Instruction
	}{
		{
			n: 1002,
			ins: Instruction{
				Opcode:    02,
				AddrModes: []int{0, 1, 0},
			},
		},
		{
			n: 3,
			ins: Instruction{
				Opcode:    03,
				AddrModes: []int{0, 0, 0},
			},
		},
	}

	for _, tc := range tt {
		expected := tc.ins
		res := CreateInstruction(tc.n)

		if res.Opcode != expected.Opcode {
			t.Errorf("opcode = %v expected %v", res.Opcode, expected.Opcode)
		}

		for i, m := range res.AddrModes {
			if m != expected.AddrModes[i] {
				t.Errorf("Param[%v]  addressingMode = %v expected %v",
					i+1, m, expected.AddrModes[i])
			}
		}
	}
}

func TestExecuteInstruction(t *testing.T) {
	tt := []struct {
		input  []int
		output []int
	}{
		{
			input:  []int{1, 0, 0, 0, 99},
			output: []int{2, 0, 0, 0, 99},
		},
		{
			input:  []int{2, 3, 0, 3, 99},
			output: []int{2, 3, 0, 6, 99},
		},
		{
			input:  []int{2, 4, 4, 5, 99, 0},
			output: []int{2, 4, 4, 5, 99, 9801},
		},
	}
	l := CreateLooger()
	for _, tc := range tt {
		l.Clear()
		ins := CreateInstruction(tc.input[0])
		m := InitMemory(tc.input, ins.AddrModes, l, nil, nil)
		m.ExecuteInstruction(ins.Opcode)

		for i, v := range m.InstructionSet {
			if v != tc.output[i] {

				for _, msg := range l.Logs() {
					t.Log(msg)
				}

				t.Errorf("[Test i/p: %v ]output mismatch res: %v expected: %v",
					tc.input, m.InstructionSet, tc.output)
				break
			}
		}
	}
}
