package main

import (
	"fmt"
	"os"

	"github.com/som.subhojit1988/aoc_2k19/day6/orbit"
	"github.com/som.subhojit1988/aoc_2k19/inputreader"
)

const (
	part1inputFile = "day6-part1-input.txt"
	part2inputFile = "day6-part2-input.txt"
)

func readInput(fname string) ([]string, error) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// fptr := flag.String("fpath", fmt.Sprintf("%s/%s", wd, fname),
	// 	"file path to read from")

	inreader := &inputreader.ReadInput{
		FileName: fmt.Sprintf("%s/%s", wd, fname),
	}
	return inreader.GetLines()
}

func part1() {
	input, err := readInput(part1inputFile)
	if err != nil {
		panic(err)
	}

	chrt := orbit.CreateOrbitChart(input)
	// fmt.Println(chrt)
	n := chrt.CalculateOrbitCntChecksum()
	fmt.Printf("Total orbit= %d\n", n)

}

func part2() {
	input, err := readInput(part2inputFile)
	if err != nil {
		panic(err)
	}

	chrt := orbit.CreateOrbitChart(input)
	// fmt.Println(chrt)
	n := chrt.CalculateOrbitTransfer("YOU", "SAN")
	fmt.Printf("# orbit transfers: %d\n", n)
}

func main() {
	part1()
	part2()
}
