package mpxr

import (
	"fmt"
	"io"
	"strings"
)

type NamedConc struct {
	Name string
	Conc float64
}

type NamedConcVol struct {
	Name string
	Conc float64
	Vol float64
}

func Total(vols ...NamedConcVol) NamedConcVol {
	var t NamedConcVol
	t.Name = "total"
	for _, v := range vols {
		t.Vol += v.Vol
		t.Conc += v.Vol * v.Conc
	}
	t.Conc /= t.Vol
	return t
}

// This uses the following system of equations to calculate the volumes to add
// to produce a final volume with equal masses from all of the inputs (assuming 3 inputs here as an example):
//
// c1 * v1 = c2 * v2
// c1 * v1 = c3 * v3
// v1 + v2 + v3 = total
//
// solve the first 2 for v2 and v3, then substitute:
// v2 = (v1 * c1) / c2
// v3 = (v1 * c1) / c3
// v1 + ((v1 * c1) / c2) + ((v1 * c1) / c3) = total
//
// solve for v1:
// v1 * (1 + (c1 / c2) + (c1 / c3)) = total
// v1 * (1 + (c1 * (1 / c2 + 1 / c3))) = total
// v1 = total / (1 + (c1 * (1 / c2 + 1 / c3)))
func FullCalc(totalVol float64, concs ...NamedConc) []NamedConcVol {
	if len(concs) < 1 {
		return nil
	}
	c1 := concs[0]
	paren := 0.0
	for _, c := range concs[1:] {
		paren += 1.0 / c.Conc
	}
	den := 1.0 + c1.Conc * paren
	v1 := totalVol / den

	vols := make([]NamedConcVol, 0, len(concs))
	vols = append(vols, NamedConcVol{Name: c1.Name, Conc: c1.Conc, Vol: v1})

	for _, c := range concs[1:] {
		vols = append(vols, NamedConcVol{Name: c.Name, Conc: c.Conc, Vol: c1.Conc * v1 / c.Conc})
	}

	return vols
}

func FprintVols(w io.Writer, header bool, vols ...NamedConcVol) (n int, err error) {
	if header {
		nw, e := fmt.Fprintf(w, "name\tconc\tvol\tmass\n")
		n += nw
		if e != nil {
			return n, e
		}
	}
	for _, vol := range vols {
		nw, e := fmt.Fprintf(w, "%v\t%v\t%v\t%v\n", vol.Name, vol.Conc, vol.Vol, vol.Conc * vol.Vol)
		n += nw
		if e != nil {
			return n, e
		}
	}

	total := Total(vols...)
	nw, e := fmt.Fprintf(w, "%v\t%v\t%v\t%v\n", total.Name, total.Conc, total.Vol, total.Conc * total.Vol)
	n += nw
	return n, e
}

// this solves the following system of equations, assuming 3 inputs:
// v1 * c1 = v2 * c2
// v1 * c1 = v3 * c3
// ((v1 * c1) + (v2 * c2) + (v3 * c3)
// ...

func FullCalcFixedConc(totalVol, totalConc float64, concs ...NamedConc) []NamedConcVol {
	vols := make([]NamedConcVol, 0, len(concs))
	volsum := 0.0
	for _, c := range concs {
		vol := totalConc * totalVol / (float64(len(concs)) * c.Conc)
		volsum += vol
		vols = append(vols, NamedConcVol{Name: c.Name, Conc: c.Conc, Vol: vol})
	}
	vols = append(vols, NamedConcVol{Name: "water", Conc: 0.0, Vol: totalVol - volsum})
	return vols
}

func FprintVolsLatex(w io.Writer, header bool, digits int, vols ...NamedConcVol) (n int, err error) {
	if header {
		nw, e := fmt.Fprintf(w, "name & conc & vol & mass\\\\\n")
		n += nw
		if e != nil {
			return n, e
		}
	}
	for _, vol := range vols {
		nw, e := fmt.Fprintf(w, "%v & %.*f & %.*f & %.*f\\\\\n", strings.ReplaceAll(vol.Name, "_", "\\_"), digits, vol.Conc, digits, vol.Vol, digits, vol.Conc * vol.Vol)
		n += nw
		if e != nil {
			return n, e
		}
	}

	total := Total(vols...)
	nw, e := fmt.Fprintf(w, "%v & %.*f & %.*f & %.*f\n\\\\", strings.ReplaceAll(total.Name, "_", "\\_"), digits, total.Conc, digits, total.Vol, digits, total.Conc * total.Vol)
	n += nw
	return n, e
}

