package gml

import (
	"testing"

	"github.com/cwbudde/go-citygml/types"
)

func FuzzParsePos(f *testing.F) {
	// Seed corpus with representative inputs.
	f.Add("1.0 2.0 3.0")
	f.Add("1.0 2.0")
	f.Add("")
	f.Add("nan inf -inf")
	f.Add("1e308 -1e308 0")
	f.Add("   1.5   2.5   3.5   ")
	f.Add("abc def ghi")
	f.Add("1.0")
	f.Add("1.0 2.0 3.0 4.0")

	f.Fuzz(func(t *testing.T, input string) {
		// Must not panic.
		_, _, _ = ParsePos(input)
	})
}

func FuzzParsePosList(f *testing.F) {
	// Seed corpus with 3D coordinate lists.
	f.Add("0 0 0 10 0 0 10 10 0 0 10 0 0 0 0", 3)
	f.Add("0 0 10 0 10 10 0 0", 2)
	f.Add("", 3)
	f.Add("1.0 2.0", 3)
	f.Add("nan 0 0 0 nan 0", 3)
	f.Add("abc", 3)
	f.Add("1 2 3 4 5", 3)
	f.Add("1e-300 2e300 3", 3)

	f.Fuzz(func(t *testing.T, input string, dim int) {
		if dim < 2 || dim > 3 {
			return
		}
		// Must not panic.
		_, _, _ = ParsePosList(input, types.Dimensionality(dim))
	})
}

func FuzzValidateRing(f *testing.F) {
	f.Add(true) // Minimal seed; ring is constructed in the fuzz function.

	f.Fuzz(func(t *testing.T, closed bool) {
		// Construct a ring programmatically with variable closure.
		pts := []types.Point{
			{X: 0, Y: 0, Z: 0},
			{X: 10, Y: 0, Z: 0},
			{X: 10, Y: 10, Z: 0},
			{X: 0, Y: 10, Z: 0},
		}
		if closed {
			pts = append(pts, pts[0])
		}

		ring := types.Ring{Points: pts}
		// Must not panic.
		_ = ValidateRing(ring, "fuzz")
	})
}
