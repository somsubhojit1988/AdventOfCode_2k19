package orbit

import (
	"fmt"
	"strings"
)

type OrbitChart struct {
	orbitMap map[string]string
}

func CreateOrbitChart(input []string) *OrbitChart {

	orbitMap := func() map[string]string {
		m := make(map[string]string)
		for _, l := range input {
			r := strings.Split(l, ")")
			if len(r) == 2 {
				k, v := strings.TrimSpace(r[1]), strings.TrimSpace(r[0])
				m[k] = v
			}
		}
		return m
	}()

	return &OrbitChart{
		orbitMap: orbitMap,
	}
}

func (c *OrbitChart) String() string {
	strb := &strings.Builder{}
	for k, v := range c.orbitMap {
		strb.WriteString(fmt.Sprintf("%s orbits %s\n", k, v))
	}
	return strb.String()
}

func (o *OrbitChart) CalculateOrbitCntChecksum() int {
	cnt := 0
	for k := range o.orbitMap {
		path := []string{k}

		n := func(key string) int {
			ptr, n := key, 0
			for tmp, ok := o.orbitMap[ptr]; ok && tmp != ptr; n++ {
				ptr = tmp
				path = append(path, ptr)

				tmp, ok = o.orbitMap[ptr]
			}
			// log.Printf("%s : %s  ", key, strings.Join(path, " -> "))
			return n
		}(k)
		cnt += n
		// log.Printf("[Total count= %d ]%s orbit count= %d", cnt, k, n)
	}
	return cnt
}

func (o *OrbitChart) CalculateOrbitTransfer(ins ...string) int {

	src, dest := "", ""

	if len(ins) >= 2 {
		src, dest = ins[0], ins[1]
	} else {
		src, dest = "YOU", "SAN"
	}

	p2 := func(node string) map[string]int {
		ptr := node
		path, cnt := map[string]int{}, 0
		p := []string{ptr}

		for {
			tmp, ok := o.orbitMap[ptr]

			if !ok || ptr == tmp {
				break
			}
			p = append(p, tmp)
			cnt++
			path[tmp] = cnt
			ptr = tmp
		}

		return path
	}(dest)

	n, ptr := 0, o.orbitMap[src]
	for {
		tmp, ok := o.orbitMap[ptr]
		n++
		if !ok || tmp == ptr {
			break
		}
		if _, ok := p2[tmp]; ok {
			n = n + p2[tmp] - 1
			break
		}
		ptr = tmp
	}
	return n
}
