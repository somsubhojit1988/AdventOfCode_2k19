package intcomputer

import "testing"

func TestIntComputer_ExecuteEqualPosition(t *testing.T) {

	instructions := []int{3, 9, 8, 9, 10, 9, 4, 9, 99, -1, 8}
	l := CreateLogger()
	input := 8

	c := CreateIntComputer(instructions, l, func() int {
		return input
	}, func(n int) {
		t.Logf("Output: %d", n)
		if n != 1 {
			t.Errorf("result: %d, expected 1", n)
		}
	})

	if err := c.Run(); err != nil {
		panic(err)
	}

	c.Reset()
	c.Program(instructions)
	input = 9
	c.OutFunc = func(n int) {
		t.Logf("Output: %d", n)
		if n != 0 {
			t.Errorf("result: %d, expected 0", n)
		}
	}

	if err := c.Run(); err != nil {
		panic(err)
	}

}

func TestIntComputer_ExecuteEqualImmediate(t *testing.T) {

	instructions := []int{3, 3, 1108, -1, 8, 3, 4, 3, 99}
	l := CreateLogger()
	input := 8
	// equal to 8; output 1 (if it is) or 0 (if it is not)
	c := CreateIntComputer(instructions, l, func() int {
		return input
	}, func(n int) {
		t.Logf("Output: %d", n)
		if n != 1 {
			t.Errorf("result: %d, expected 1", n)
		}
	})

	if err := c.Run(); err != nil {
		panic(err)
	}

	c.Reset()
	input = 9
	c.OutFunc = func(n int) {
		t.Logf("Output: %d", n)
		if n != 0 {
			t.Errorf("result: %d, expected 0", n)
		}
	}

	c.Program(instructions)
	if err := c.Run(); err != nil {
		panic(err)
	}
}

func TestIntComputer_ExecuteLtPosition(t *testing.T) {

	instructions := []int{3, 9, 7, 9, 10, 9, 4, 9, 99, -1, 8}
	l := CreateLogger()
	input := 7
	// less than 8; output 1 (if it is) or 0 (if it is not)
	c := CreateIntComputer(instructions, l, func() int {
		return input
	}, func(n int) {
		t.Logf("Output: %d", n)
		if n != 1 {
			t.Errorf("result: %d, expected 1", n)
		}
	})

	if err := c.Run(); err != nil {
		panic(err)
	}

	c.Reset()
	input = 9
	c.Program(instructions)
	c.OutFunc = func(n int) {
		t.Logf("Output: %d", n)
		if n != 0 {
			t.Errorf("result: %d, expected 0", n)
		}
	}

	if err := c.Run(); err != nil {
		panic(err)
	}
}

func TestIntComputer_ExecuteLtImmediate(t *testing.T) {

	instructions := []int{3, 3, 1107, -1, 8, 3, 4, 3, 99}
	l := CreateLogger()
	input := 7
	// less than 8; output 1 (if it is) or 0 (if it is not)
	c := CreateIntComputer(instructions, l, func() int {
		return input
	}, func(n int) {
		t.Logf("Output: %d", n)
		if n != 1 {
			t.Errorf("result: %d, expected 1", n)
		}
	})

	if err := c.Run(); err != nil {
		panic(err)
	}

	c.Reset()
	input = 9
	c.OutFunc = func(n int) {
		t.Logf("Output: %d", n)
		if n != 0 {
			t.Errorf("result: %d, expected 0", n)
		}
	}
	c.Program(instructions)
	if err := c.Run(); err != nil {
		panic(err)
	}
}

func TestIntComputer_JmpPosition(t *testing.T) {

	instructions := []int{3, 12, 6, 12, 15, 1, 13, 14, 13, 4, 13, 99, -1, 0, 1, 9}
	l := CreateLogger()

	input := 0
	// take an input, then output 0 if the input was zero
	// or 1 if the input was non-zero
	c := CreateIntComputer(instructions, l, func() int {
		return input
	}, func(n int) {
		t.Logf("Output: %d", n)
		if n != 0 {
			t.Errorf("result: %d, expected 0", n)
		}
	})

	if err := c.Run(); err != nil {
		panic(err)
	}

	input = 9
	c.Program(instructions)
	c.OutFunc = func(n int) {
		t.Logf("Output: %d", n)
		if n != 1 {
			t.Errorf("result: %d, expected 1", n)
		}
	}

	if err := c.Run(); err != nil {
		panic(err)
	}
}

func TestIntComputer_JmpImmediate(t *testing.T) {

	instructions := []int{3, 3, 1105, -1, 9, 1101, 0, 0, 12, 4, 12, 99, 1}
	l := CreateLogger()

	var input int

	// take an input, then output 0 if the input was zero
	// or 1 if the input was non-zero
	input = 0
	c := CreateIntComputer(instructions, l, func() int {
		return input
	}, func(n int) {
		t.Logf("Output: %d", n)
		if n != 0 {
			t.Errorf("result: %d, expected 0", n)
		}
	})

	if err := c.Run(); err != nil {
		panic(err)
	}

	input = 9
	c.Reset()
	c.OutFunc = func(n int) {
		t.Logf("Output: %d", n)
		if n != 1 {
			t.Errorf("result: %d, expected 1", n)
		}
	}
	c.Program(instructions)
	if err := c.Run(); err != nil {
		panic(err)
	}
}

func TestIntComputer_ExecuteJmp(t *testing.T) {

	instructions := []int{3, 21, 1008, 21, 8, 20, 1005, 20, 22, 107, 8,
		21, 20, 1006, 20, 31, 1106, 0, 36, 98, 0, 0, 1002, 21, 125, 20,
		4, 20, 1105, 1, 46, 104, 999, 1105, 1, 46, 1101, 1000, 1, 20, 4,
		20, 1105, 1, 46, 98, 99}
	l := CreateLogger()

	var input int
	var expected int
	// output 999 if the input value is below 8, output 1000 if
	// the input value is equal to 8, or output 1001 if the input
	//  value is greater than 8.
	c := CreateIntComputer(instructions, l, func() int {
		return input
	}, func(n int) {
		t.Logf("Output: %d", n)
		if n != expected {
			t.Errorf("result: %d, expected 0", n)
		}
	})

	input, expected = 7, 999
	if err := c.Run(); err != nil {
		panic(err)
	}

	c.Reset()
	c.Program(instructions)
	input, expected = 8, 1000
	if err := c.Run(); err != nil {
		panic(err)
	}

	c.Reset()
	c.Program(instructions)
	input, expected = 9, 1001
	if err := c.Run(); err != nil {
		panic(err)
	}
}
