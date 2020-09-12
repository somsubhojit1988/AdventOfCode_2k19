package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/som.subhojit1988/aoc_2k19/day7/amplifier"
	"github.com/som.subhojit1988/aoc_2k19/inputreader"
	"github.com/som.subhojit1988/aoc_2k19/intcomputer"
)

func readInput(inFile string) []int {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	inreader := &inputreader.ReadInput{FileName: fmt.Sprintf("%s/%s", wd, inFile)}
	lines, err := inreader.GetLines()
	if err != nil {
		panic(err)
	}

	_, ret := inputreader.ProcessLines(lines)

	return ret
}

func generateInputs(n, l, h int) []string {

	ret, fringe := []string{}, []string{}

	for i := l; i <= h; i++ {
		fringe = append(fringe, fmt.Sprintf("%d", i))
	}

	for len(fringe) > 0 {
		s := fringe[0]
		fringe = fringe[1:]
		if len(s) < n {
			cmap := func(s string) map[int]int {
				cs := map[int]int{}
				for _, c := range s {
					q, err := strconv.Atoi(string(c))
					if err != nil {
						panic(err)
					}
					if n, ok := cs[q]; ok {
						fmt.Printf("Illegal string %d appears %d times\n",
							q, n)
					}
					cs[q] = 1
				}
				return cs
			}(s)
			for i := l; i <= h; i++ {
				if _, ok := cmap[i]; !ok {
					fringe = append(fringe, s+fmt.Sprintf("%d", i))
				}
			}
		}

		if len(s) == n {
			ret = append(ret, s)
		}
	}
	return ret
}

func maxOf(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func run(nAmps, pHigh, pLow int, feedback bool, inputFile string) int {
	instructions, inputs := readInput(inputFile), generateInputs(nAmps, pLow, pHigh)
	maxBoost := 0

	for _, in := range inputs {
		input, err := strconv.Atoi(in)
		if err != nil {
			panic(err)
		}
		ps := func() []int {
			ret := make([]int, nAmps)
			for i := 0; i < 5; i++ {
				ret[i] = input % 10
				input = input / 10
			}
			return ret
		}()

		logger := intcomputer.CreateLogger()
		c := amplifier.CreateAmpCircuit(nAmps, ps, instructions, logger, feedback)
		ret, err := c.Run(0, false)
		if err != nil {
			panic(err)
		}
		maxBoost = maxOf(maxBoost, ret)
		// fmt.Printf("Amp Circuit phases: %v o/p= %d [MaxThrust= %d] \n",
		// 	ps, ret, maxBoost)
	}

	return maxBoost
}

func part1() {
	n, pLow, pHigh, inputFile := 5, 0, 4, "day7-part1-input.txt"
	maxBoost := run(n, pHigh, pLow, false, inputFile)
	fmt.Printf("[Part-1] Max Boost: %d\n", maxBoost)
}

func part2() {
	n, pLow, pHigh, inputFile := 5, 5, 9, "day7-part1-input.txt"
	maxBoost := run(n, pHigh, pLow, true, inputFile)
	fmt.Printf("[Part-2] Max Boost: %d\n", maxBoost)
}

func main() {
	part1()
	part2()
}
