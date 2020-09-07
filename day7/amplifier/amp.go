package amplifier

import (
	"github.com/som.subhojit1988/aoc_2k19/intcomputer"
)

const (
	ampStateInit          = -1
	ampStatePhaseProvided = 0
	ampStateInputProvided = 1
)

type Amplifier struct {
	c            *intcomputer.IntComputer
	state        int
	phase, input int
	ampProgram   []int
}

func CreateAmp(instructions []int,
	logger *intcomputer.Logger, phase, input int,
	outFunc intcomputer.OutputMethod) *Amplifier {
	prog := make([]int, len(instructions))
	copy(prog, instructions)
	a := &Amplifier{
		c:          intcomputer.CreateIntComputer(instructions, logger, nil, outFunc),
		state:      ampStateInit,
		phase:      phase,
		input:      input,
		ampProgram: prog,
	}
	a.c.InFunc = func() int {
		ret := -1
		switch a.state {
		case ampStateInit:
			a.state = ampStatePhaseProvided
			ret = a.phase
		case ampStateInputProvided, ampStatePhaseProvided:
			a.state = ampStateInputProvided
			ret = a.input
		}
		return ret
	}
	return a
}

func (a *Amplifier) Run(in int) error {
	a.input = in
	return a.c.Run()
}

func (a *Amplifier) Reset() {
	a.state = ampStateInit
	a.c.Reset()
	a.c.Program(a.ampProgram)
}

type SeriesAmpCircuit struct {
	n       int
	as      []*Amplifier
	outFunc func(n int)
	currOut int
}

func CreateAmpCircuit(n int, pSettings []int,
	instructions []int, logger *intcomputer.Logger) *SeriesAmpCircuit {
	as := make([]*Amplifier, n)

	circuit := &SeriesAmpCircuit{
		n:       n,
		currOut: 0,
	}
	circuit.outFunc = func(n int) { circuit.currOut = n }

	for i := 0; i < n; i++ {
		as[i] = CreateAmp(instructions, logger, pSettings[i], 0, circuit.outFunc)
	}

	circuit.as = as

	return circuit
}

func (ac *SeriesAmpCircuit) Run(circuitIn int, feedbackMode bool) (int, error) {
	ac.currOut = circuitIn
	for i := 0; i < ac.n; i++ {
		if !feedbackMode {
			ac.as[i].Reset()
		}
		if err := ac.as[i].Run(ac.currOut); err != nil {
			return ac.currOut, err
		}
	}
	return ac.currOut, nil
}
