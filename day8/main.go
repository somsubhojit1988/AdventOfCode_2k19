package main

import (
	"fmt"
	"os"

	"github.com/som.subhojit1988/aoc_2k19/day8/decoder"
	"github.com/som.subhojit1988/aoc_2k19/inputreader"
)

func main() {
	part1()
	part2()
}

func part1() {
	inputFile, h, w := "day8-part1-input.txt", 6, 25
	str := readInput(inputFile)
	img, err := decoder.CreateImage(str, h, w)
	if err != nil {
		panic(err)
	}
	fmt.Printf("CRC: %d\n", img.CRC)
}

func part2() {
	inputFile, h, w := "day8-part1-input.txt", 6, 25
	str := readInput(inputFile)
	mat, err := decoder.CreateImage(str, h, w)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", mat.Decode())
}

func readInput(inFile string) string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	inreader := &inputreader.ReadInput{FileName: fmt.Sprintf("%s/%s", wd, inFile)}
	lines, err := inreader.GetLines()
	if err != nil {
		panic(err)
	}
	if len(lines) != 1 {
		panic("Input file has unexpected no. of lines")
	}
	return lines[0]
}
