package gml

import (
	"math"
	"testing"

	"github.com/cwbudde/go-citygml/types"
)

func closedRing(pts ...types.Point) types.Ring {
	return types.Ring{Points: pts}
}

func TestValidateRing_Valid(t *testing.T) {
	ring := closedRing(
		types.Point{X: 0, Y: 0, Z: 0},
		types.Point{X: 10, Y: 0, Z: 0},
		types.Point{X: 10, Y: 10, Z: 0},
		types.Point{X: 0, Y: 0, Z: 0},
	)

	errs := ValidateRing(ring, "test")
	if len(errs) != 0 {
		t.Errorf("expected no errors, got %v", errs)
	}
}

func TestValidateRing_TooFewPoints(t *testing.T) {
	ring := closedRing(
		types.Point{X: 0, Y: 0},
		types.Point{X: 1, Y: 0},
		types.Point{X: 0, Y: 0},
	)

	errs := ValidateRing(ring, "test")
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d: %v", len(errs), errs)
	}
}

func TestValidateRing_NotClosed(t *testing.T) {
	ring := closedRing(
		types.Point{X: 0, Y: 0, Z: 0},
		types.Point{X: 10, Y: 0, Z: 0},
		types.Point{X: 10, Y: 10, Z: 0},
		types.Point{X: 5, Y: 5, Z: 0}, // not equal to first
	)

	errs := ValidateRing(ring, "test")
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d: %v", len(errs), errs)
	}
}

func TestValidateRing_NonFinite(t *testing.T) {
	ring := closedRing(
		types.Point{X: 0, Y: 0, Z: 0},
		types.Point{X: math.NaN(), Y: 0, Z: 0},
		types.Point{X: 10, Y: math.Inf(1), Z: 0},
		types.Point{X: 0, Y: 0, Z: 0},
	)

	errs := ValidateRing(ring, "test")
	if len(errs) != 2 {
		t.Fatalf("expected 2 errors, got %d: %v", len(errs), errs)
	}
}

func TestValidatePolygon_WithHole(t *testing.T) {
	ext := closedRing(
		types.Point{X: 0, Y: 0}, types.Point{X: 10, Y: 0},
		types.Point{X: 10, Y: 10}, types.Point{X: 0, Y: 0},
	)
	// Interior ring with too few points
	hole := closedRing(
		types.Point{X: 1, Y: 1}, types.Point{X: 2, Y: 1}, types.Point{X: 1, Y: 1},
	)
	poly := types.Polygon{Exterior: ext, Interior: []types.Ring{hole}}

	errs := ValidatePolygon(poly, "poly")
	if len(errs) != 1 {
		t.Fatalf("expected 1 error from interior ring, got %d: %v", len(errs), errs)
	}
}

func TestValidateMultiSurface(t *testing.T) {
	good := types.Polygon{
		Exterior: closedRing(
			types.Point{X: 0, Y: 0}, types.Point{X: 10, Y: 0},
			types.Point{X: 10, Y: 10}, types.Point{X: 0, Y: 0},
		),
	}
	bad := types.Polygon{
		Exterior: closedRing(
			types.Point{X: 0, Y: 0}, types.Point{X: 1, Y: 0}, types.Point{X: 0, Y: 0},
		),
	}
	ms := types.MultiSurface{Polygons: []types.Polygon{good, bad}}

	errs := ValidateMultiSurface(ms, "ms")
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d: %v", len(errs), errs)
	}
}

func TestValidateSolid(t *testing.T) {
	solid := types.Solid{
		Exterior: types.MultiSurface{
			Polygons: []types.Polygon{{
				Exterior: closedRing(
					types.Point{X: 0, Y: 0, Z: 0}, types.Point{X: 10, Y: 0, Z: 0},
					types.Point{X: 10, Y: 10, Z: 0}, types.Point{X: 0, Y: 0, Z: 0},
				),
			}},
		},
	}

	errs := ValidateSolid(solid, "solid")
	if len(errs) != 0 {
		t.Errorf("expected no errors, got %v", errs)
	}
}

func TestValidationError_Format(t *testing.T) {
	err := &ValidationError{Path: "Polygon.exterior", Message: "ring not closed"}

	want := "gml: validation: Polygon.exterior: ring not closed"
	if got := err.Error(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestValidateRing_Empty(t *testing.T) {
	ring := types.Ring{}

	errs := ValidateRing(ring, "test")
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}
