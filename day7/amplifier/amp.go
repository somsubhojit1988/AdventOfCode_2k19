package amplifier

import (
	"github.com/som.subhojit1988/aoc_2k19/intcomputer"
)

const (
	ampStateInit              = -1
	ampStatePhaseProvided     = 0
	ampStateInputProvided     = 1
	ampStateFinishedExecution = 2
)

type Amplifier struct {
	c                    *intcomputer.IntComputer
	state                int
	phase, input, output int
	ampProgram           []int
	isFeedbackMode       bool
}

func CreateAmp(instructions []int,
	logger *intcomputer.Logger, phase, input int,
	feedbackMode bool) *Amplifier {
	prog := make([]int, len(instructions))
	copy(prog, instructions)
	a := &Amplifier{
		c:              intcomputer.CreateIntComputer(instructions, logger, nil, nil),
		state:          ampStateInit,
		phase:          phase,
		input:          input,
		ampProgram:     prog,
		isFeedbackMode: feedbackMode,
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
	a.c.OutFunc = func(x int) {
		a.output = x
		if feedbackMode {
			a.c.Break()
		}
	}
	return a
}

func (a *Amplifier) Run(in int) error {
	a.input = in
	if a.c.IsHalted() {
		a.state = ampStateFinishedExecution
		return nil
	}
	a.c.Resume()
	return a.c.Run()
}

func (a *Amplifier) Reset() {
	a.state = ampStateInit
	a.c.Reset()
	a.c.Program(a.ampProgram)
}

type SeriesAmpCircuit struct {
	n  int
	as []*Amplifier
}

func CreateAmpCircuit(
	n int, pSettings []int,
	instructions []int, logger *intcomputer.Logger,
	feedbackMode bool) *SeriesAmpCircuit {
	as := make([]*Amplifier, n)

	circuit := &SeriesAmpCircuit{n: n}

	for i := 0; i < n; i++ {
		as[i] = CreateAmp(instructions, logger, pSettings[i], 0, feedbackMode)
	}

	circuit.as = as

	return circuit
}

func (ac *SeriesAmpCircuit) Run(circuitIn int, feedbackMode bool) (int, error) {
	i := 0
	for {
		amp, nxtAmp := ac.as[i%ac.n], ac.as[(i+1)%ac.n]
		if err := amp.Run(amp.input); err != nil {
			return amp.output, err
		}
		if amp.state == ampStateFinishedExecution {
			break
		}

		nxtAmp.input = amp.output
		i++
	}

	return ac.as[ac.n-1].output, nil
}
