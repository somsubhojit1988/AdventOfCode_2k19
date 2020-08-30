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
	logger := intcomputer.CreateLooger()
	c := intcomputer.CreateIntComputer(instructions, logger)

	// replace position 1 with the value 12 and
	// replace position 2 with the value 2
	c.Program(1, 12)
	c.Program(2, 2)
	log.Printf("result: %v", c.Run())

	// Your puzzle answer should be:  4090689.
	// for _, l := range logger.Logs() {
	// 	log.Printf(l)
	// }

}
