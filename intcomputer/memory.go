package intcomputer

import (
	"fmt"
	"strings"
)

type Memory struct {
	storage []int
	memPtr  int
	logger  *Logger
}

func (m *Memory) Size() int {
	return len(m.storage)
}

func (m *Memory) String() string {
	dump := func() string {
		sb := &strings.Builder{}
		cntr := 0
		for i, v := range m.storage {
			cntr++
			sb.WriteString(fmt.Sprintf("%d (%d) ", v, i))
			if cntr >= 10 {
				cntr = 0
				sb.WriteString("\n")
			}
		}
		return fmt.Sprintf("[ %s ]", sb.String())
	}
	return fmt.Sprintf("MEM: ptr= %d\n%s\n", m.memPtr, dump())
}

func (m *Memory) read(addrMode, d int) (int, error) {
	var ret int
	if m.memPtr+d >= m.Size() {
		return -1, fmt.Errorf("MEMREAD (addressing-mode= %d, addr = %d) Out of range",
			addrMode, m.memPtr+d)
	}
	x := m.storage[m.memPtr+d]
	switch addrMode {
	case Position:
		if x >= m.Size() {
			return -1, fmt.Errorf("MEMREAD (addressing-mode= %d, addr = %d) Out of range",
				addrMode, x)
		}
		ret = m.storage[x]
	case Immediate:
		ret = x
	}
	return ret, nil
}

func (m *Memory) readAddress(ptr int) (int, error) {
	if ptr < 0 || ptr >= m.Size() {
		return -1, fmt.Errorf("MEMREAD (addr = %d) Out of range", ptr)
	}
	return m.storage[ptr], nil
}

func (m *Memory) write(v, ptr int) error {
	if ptr < 0 || ptr >= m.Size() {
		return fmt.Errorf("MEMWRITE (addr = %d  v= %d) Out of range",
			ptr, v)
	}
	m.storage[ptr] = v
	return nil
}

func (m *Memory) opcodeFetch() (int, error) {
	return m.read(Immediate, 0)
}
