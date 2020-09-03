package intcomputer

import "testing"

func TestInstructionDecode(t *testing.T) {
	tt := []struct {
		code int
		ins  *Instruction
	}{
		{code: 10101, ins: &Instruction{Opcode: 01, ParamAddrModes: []int{1, 0, 1}}},
	}

	for _, tc := range tt {
		ins := decode(tc.code)
		e := tc.ins
		if ins.Opcode != e.Opcode {
			t.Errorf("FAIL: Opcode mismatch expected= %s actual= %s", e, ins)
		}

		for i, m := range e.ParamAddrModes {
			if m != ins.ParamAddrModes[i] {
				t.Errorf("FAIL: Parameter addressing-modes mismatch expected= %s actual= %s",
					e, ins)
			}
		}
	}
}
