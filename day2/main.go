package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/som.subhojit1988/aoc_2k19/inputreader"
	"github.com/som.subhojit1988/aoc_2k19/intcomputer"
)

func readInput() []int {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fptr := flag.String("fpath",
		fmt.Sprintf("%s/%s", wd, "day2-input.txt"),
		"file path to read from")

	inreader := &inputreader.ReadInput{FileName: *fptr}
	lines, err := inreader.GetLines()
	if err != nil {
		panic(err)
	}

	_, ret := inputreader.ProcessLines(lines)
	return ret
}

func main() {
	instructions := readInput()
	logger := intcomputer.CreateLogger()
	c := intcomputer.CreateIntComputer(instructions, logger, nil, nil)

	// replace position 1 with the value 12 and
	// replace position 2 with the value 2
	c.Store(12, 1)
	c.Store(2, 2)
	err := c.Run()
	if err != nil {
		panic(err)
	}

	for _, l := range logger.Logs() {
		log.Printf(l)
	}

	// Your puzzle answer should be:  4090689.
	v, err := c.ReadMemory(0, 1)
	if err != nil {
		panic(err)
	}
	log.Printf("result: %v", v)

}
