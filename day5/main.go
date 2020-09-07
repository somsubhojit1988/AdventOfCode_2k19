package main

import (
	"fmt"
	"os"

	"github.com/som.subhojit1988/aoc_2k19/inputreader"
	"github.com/som.subhojit1988/aoc_2k19/intcomputer"
)

const partOneInputFileName = "day5-part1-input.txt"
const partTwoInputFileName = "day5-part2-input.txt"

func readInstructions(fname string) []int {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	inreader := &inputreader.ReadInput{FileName: fmt.Sprintf("%s/%s", wd, fname)}
	lines, err := inreader.GetLines()
	if err != nil {
		panic(err)
	}

	_, ret := inputreader.ProcessLines(lines)

	return ret
}

func readInput() int {
	var n int
	fmt.Print("Enter systemID: ")
	_, err := fmt.Scanf("%d", &n)
	if err != nil {
		panic(err)
	}
	return n
}

func printOutput(n int) {
	var status string
	if n == 0 {
		status = "PASS"
	} else {
		status = "FAIL (may be not if this is the last)"
	}
	fmt.Printf("Diagnostic code [%s] (Expected - output)= %v\n",
		status, n)
}

func part1() {
	input := readInstructions(partOneInputFileName)
	logger := intcomputer.CreateLogger()
	c := intcomputer.CreateIntComputer(input, logger, readInput, printOutput)

	err := c.Run() // ans: 9961446
	if err != nil {
		panic(err)
	}

	// for _, l := range logger.Logs() {
	// 	log.Printf(l)
	// }
}

func part2() {
	input := readInstructions(partTwoInputFileName)
	logger := intcomputer.CreateLogger()
	c := intcomputer.CreateIntComputer(input, logger, readInput, printOutput)

	err := c.Run()
	if err != nil {
		panic(err)
	}

	// for _, l := range logger.Logs() {
	// 	log.Printf(l)
	// }
}

func main() {
	part1()
	part2()
}
