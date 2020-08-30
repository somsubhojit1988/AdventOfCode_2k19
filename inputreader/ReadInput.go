package inputreader

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type ReadInput struct {
	FileName string
}

func (r *ReadInput) GetLines() ([]string, error) {
	buf, err := os.Open(r.FileName)
	ret := []string{}
	if err != nil {
		return ret, err
	}

	defer func() {
		if err := buf.Close(); err != nil {
			panic(err)
		}
	}()

	snl := bufio.NewScanner(buf)

	for snl.Scan() {
		ret = append(ret, snl.Text())
	}
	err = snl.Err()
	return ret, err
}

func ProcessLines(lines []string) (int, []int) {
	// process i/p
	n, ret := 0, []int{}
	for _, line := range lines {
		for _, c := range strings.Split(strings.TrimSpace(line), ",") {
			if x, err := strconv.Atoi(c); err == nil {
				ret = append(ret, x)
				n++
			}
		}
	}
	return n, ret
}
