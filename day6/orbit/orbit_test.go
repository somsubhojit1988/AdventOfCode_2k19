package orbit

import (
	"strings"
	"testing"
)

func TestOrbitCnt(t *testing.T) {
	data := strings.Split(`
		COM)B
		B)C
		C)D
		D)E
		E)F
		B)G
		G)H
		D)I
		E)J
		J)K
		K)L`, "\n")[1:]

	c := CreateOrbitChart(data)
	t.Logf("orbit map: %s", c)
	r := c.CalculateOrbitCntChecksum()
	t.Logf("result: %d", r)
}

func TestCalculateOrbitTransfer(t *testing.T) {
	data := strings.Split(`
	COM)B
	B)C
	C)D
	D)E
	E)F
	B)G
	G)H
	D)I
	E)J
	J)K
	K)L
	K)YOU
	I)SAN
	`, "\n")[1:]
	c := CreateOrbitChart(data)
	// t.Logf("orbit map: %s", c)

	n := c.CalculateOrbitTransfer("YOU", "SAN")
	t.Logf("# of orbit transfer: %d", n)
}
