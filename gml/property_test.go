package gml

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"testing"

	"github.com/cwbudde/go-citygml/types"
)

// TestProperty_ParsePos_Roundtrip verifies that formatting and reparsing 3D coords preserves values.
func TestProperty_ParsePos_Roundtrip(t *testing.T) {
	cases := [][3]float64{
		{0, 0, 0},
		{1.5, -2.5, 3.5},
		{1e10, -1e10, 0},
		{0.001, 0.002, 0.003},
		{500000.123, 5700000.456, 100.789},
	}

	for _, c := range cases {
		text := fmt.Sprintf("%g %g %g", c[0], c[1], c[2])

		pt, dim, err := ParsePos(text)
		if err != nil {
			t.Errorf("ParsePos(%q): %v", text, err)
			continue
		}

		if dim != types.Dim3D {
			t.Errorf("ParsePos(%q): dim=%v, want 3D", text, dim)
		}

		if pt.X != c[0] || pt.Y != c[1] || pt.Z != c[2] {
			t.Errorf("ParsePos(%q): got {%g,%g,%g}, want {%g,%g,%g}",
				text, pt.X, pt.Y, pt.Z, c[0], c[1], c[2])
		}
	}
}

// TestProperty_ParsePosList_CountInvariant verifies that N*dim values produce N points.
func TestProperty_ParsePosList_CountInvariant(t *testing.T) {
	for n := 1; n <= 10; n++ {
		for _, dim := range []types.Dimensionality{types.Dim2D, types.Dim3D} {
			total := n * int(dim)

			vals := make([]string, total)
			for i := range vals {
				vals[i] = strconv.Itoa(i)
			}

			text := strings.Join(vals, " ")

			pts, _, err := ParsePosList(text, dim)
			if err != nil {
				t.Errorf("n=%d dim=%d: %v", n, dim, err)
				continue
			}

			if len(pts) != n {
				t.Errorf("n=%d dim=%d: got %d points, want %d", n, dim, len(pts), n)
			}
		}
	}
}

// TestProperty_ParsePosList_RejectsNonFinite verifies non-finite values are rejected.
func TestProperty_ParsePosList_RejectsNonFinite(t *testing.T) {
	nonFinite := []string{"NaN", "Inf", "-Inf", "+Inf"}
	for _, nf := range nonFinite {
		text := fmt.Sprintf("0 0 0 %s 0 0", nf)

		_, _, err := ParsePosList(text, types.Dim3D)
		if err == nil {
			t.Errorf("expected error for non-finite value %q in posList", nf)
		}
	}
}

// TestProperty_ValidateRing_ClosedRingAlwaysValid verifies that a properly closed ring with enough points passes.
func TestProperty_ValidateRing_ClosedRingAlwaysValid(t *testing.T) {
	for sides := 3; sides <= 20; sides++ {
		pts := make([]types.Point, sides+1)
		for i := range sides {
			angle := 2 * math.Pi * float64(i) / float64(sides)
			pts[i] = types.Point{X: math.Cos(angle), Y: math.Sin(angle), Z: 0}
		}

		pts[sides] = pts[0] // close ring

		ring := types.Ring{Points: pts}

		errs := ValidateRing(ring, fmt.Sprintf("ring_%d_sides", sides))
		if len(errs) != 0 {
			t.Errorf("ring with %d sides should be valid, got: %v", sides, errs)
		}
	}
}

// TestProperty_ValidatePolygon_EmptyInteriorOK verifies polygon with exterior-only is fine.
func TestProperty_ValidatePolygon_EmptyInteriorOK(t *testing.T) {
	poly := types.Polygon{
		Exterior: types.Ring{Points: []types.Point{
			{X: 0, Y: 0}, {X: 10, Y: 0}, {X: 10, Y: 10}, {X: 0, Y: 10}, {X: 0, Y: 0},
		}},
	}

	errs := ValidatePolygon(poly, "prop_test")
	if len(errs) != 0 {
		t.Errorf("valid polygon should have no errors: %v", errs)
	}
}
