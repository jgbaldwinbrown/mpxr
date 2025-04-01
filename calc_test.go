package mpxr

import (
	"testing"
)

func FullCalc4(conc1, conc2, conc3, conc4, totalVol float64) float64 {
	c1 := conc1
	c2 := conc2
	c3 := conc3
	c4 := conc4
	den := 1.0 + c1 * (1 / c2 + 1/c3 + 1/c4)
	v1 := totalVol / den

	v2 := c1 * v1 / c2
	v3 := c1 * v1 / c3
	v4 := c1 * v1 / c4

	finalConc := ((v1 * c1) + (v2 * c2) + (v3 * c3) + (v4 * c4)) / totalVol
	return finalConc
}

func TestFullCalc(t *testing.T) {
	cs := []NamedConc {
		NamedConc{Name: "a", Conc: 3.5},
		NamedConc{Name: "b", Conc: 2.2},
		NamedConc{Name: "c", Conc: 4.0},
		NamedConc{Name: "d", Conc: 9.0},
	}
	tot1 := FullCalc4(cs[0].Conc, cs[1].Conc, cs[2].Conc, cs[3].Conc, 12.0)
	tot2 := Total(FullCalc(12.0, cs...)...)
	if tot1 != tot2.Conc {
		t.Errorf("tot1 %v != tot2.Conc %v", tot1, tot2.Conc)
	}
}
