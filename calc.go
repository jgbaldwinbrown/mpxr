package mpxr

import (
	"fmt"
	"io"
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

// finalConc := ((v1 * c1) + (v2 * c2) + (v3 * c3) + (v4 * c4)) / totalVol
// fmt.Printf("finalConc: %v\n", finalConc)


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
