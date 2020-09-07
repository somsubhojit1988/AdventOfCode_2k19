package amplifier

import (
	"testing"

	"github.com/som.subhojit1988/aoc_2k19/intcomputer"
)

func TestCircuitNoFeedback(t *testing.T) {
	tt := []struct {
		instructions []int
		ps           []int
		expected     int
	}{
		{
			instructions: []int{3, 15, 3, 16, 1002, 16, 10, 16, 1, 16, 15, 15, 4, 15, 99, 0, 0},
			ps:           []int{4, 3, 2, 1, 0},
			expected:     43210,
		},

		{
			instructions: []int{3, 23, 3, 24, 1002, 24, 10, 24, 1002, 23, -1, 23,
				101, 5, 23, 23, 1, 24, 23, 23, 4, 23, 99, 0, 0},
			ps:       []int{0, 1, 2, 3, 4},
			expected: 54321,
		},

		{
			instructions: []int{3, 31, 3, 32, 1002, 32, 10, 32, 1001, 31,
				-2, 31, 1007, 31, 0, 33,
				1002, 33, 7, 33, 1, 33, 31, 31, 1,
				32, 31, 31, 4, 31, 99, 0, 0, 0},
			ps:       []int{1, 0, 4, 3, 2},
			expected: 65210,
		},
	}

	n := 5
	for _, tc := range tt {
		log := intcomputer.CreateLogger()
		c := CreateAmpCircuit(n, tc.ps, tc.instructions, log)

		ret, err := c.Run(0, false)
		if err != nil {
			t.Errorf("ERROR: %s", err)
			t.FailNow()
		}
		if ret != tc.expected {
			t.Errorf("Circuit out: %d expected %d", ret, tc.expected)
			t.FailNow()
		}
		t.Logf(" -------- [TestCircuitNoFeedback] o/p: %d --------- ", ret)
	}
}

func TestCircuitFeedback(t *testing.T) {
	tt := []struct {
		instructions []int
		ps           []int
		expected     int
	}{
		{
			instructions: []int{3, 26, 1001, 26, -4, 26, 3, 27, 1002, 27, 2, 27, 1, 27, 26,
				27, 4, 27, 1001, 28, -1, 28, 1005, 28, 6, 99, 0, 0, 5},
			ps:       []int{9, 8, 7, 6, 5},
			expected: 139629729,
		},
	}

	n := 5
	for _, tc := range tt {
		log := intcomputer.CreateLogger()
		c := CreateAmpCircuit(n, tc.ps, tc.instructions, log)

		ret, err := c.Run(0, true)
		if err != nil {
			t.Errorf("ERROR: %s", err)
			t.FailNow()
		}
		if ret != tc.expected {
			t.Errorf("Circuit out: %d expected %d", ret, tc.expected)
			t.FailNow()
		}
		t.Logf(" -------- [TestCircuitFeedback] o/p: %d --------- ", ret)
	}
}