package gml

import (
	"math"
	"testing"

	"github.com/cwbudde/go-citygml/types"
)

func TestParsePos2D(t *testing.T) {
	pt, dim, err := ParsePos("3.5 7.2")
	if err != nil {
		t.Fatal(err)
	}

	if dim != types.Dim2D {
		t.Errorf("dim = %d, want 2", dim)
	}

	if pt.X != 3.5 || pt.Y != 7.2 || pt.Z != 0 {
		t.Errorf("point = %+v", pt)
	}
}

func TestParsePos3D(t *testing.T) {
	pt, dim, err := ParsePos("1.0 2.0 3.0")
	if err != nil {
		t.Fatal(err)
	}

	if dim != types.Dim3D {
		t.Errorf("dim = %d, want 3", dim)
	}

	if pt.X != 1.0 || pt.Y != 2.0 || pt.Z != 3.0 {
		t.Errorf("point = %+v", pt)
	}
}

func TestParsePosInvalid(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"single value", "1.0"},
		{"four values", "1.0 2.0 3.0 4.0"},
		{"non-numeric", "abc 2.0"},
		{"empty", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := ParsePos(tt.input)
			if err == nil {
				t.Error("expected error")
			}
		})
	}
}

func TestParsePosNonFinite(t *testing.T) {
	tests := []string{"NaN 1.0 2.0", "1.0 Inf 2.0", "1.0 2.0 -Inf"}
	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			_, _, err := ParsePos(input)
			if err == nil {
				t.Error("expected error for non-finite value")
			}
		})
	}
}

func TestParsePosList3D(t *testing.T) {
	pts, dim, err := ParsePosList("0 0 0 1 0 0 1 1 0 0 0 0", 0)
	if err != nil {
		t.Fatal(err)
	}

	if dim != types.Dim3D {
		t.Errorf("dim = %d, want 3", dim)
	}

	if len(pts) != 4 {
		t.Fatalf("got %d points, want 4", len(pts))
	}

	if pts[1].X != 1 || pts[1].Y != 0 || pts[1].Z != 0 {
		t.Errorf("pts[1] = %+v", pts[1])
	}
}

func TestParsePosList2DExplicit(t *testing.T) {
	pts, dim, err := ParsePosList("0 0 1 0 1 1 0 0", types.Dim2D)
	if err != nil {
		t.Fatal(err)
	}

	if dim != types.Dim2D {
		t.Errorf("dim = %d, want 2", dim)
	}

	if len(pts) != 4 {
		t.Fatalf("got %d points, want 4", len(pts))
	}
}

func TestParsePosListNotDivisible(t *testing.T) {
	_, _, err := ParsePosList("1 2 3 4 5", types.Dim3D)
	if err == nil {
		t.Error("expected error for non-divisible count")
	}
}

func TestParsePosListEmpty(t *testing.T) {
	_, _, err := ParsePosList("", 0)
	if err == nil {
		t.Error("expected error for empty input")
	}
}

func TestParseFloatNaN(t *testing.T) {
	_, err := parseFloat("NaN")
	if err == nil {
		t.Error("expected error for NaN")
	}
}

func TestParseFloatInf(t *testing.T) {
	for _, s := range []string{"Inf", "-Inf", "+Inf"} {
		_, err := parseFloat(s)
		if err == nil {
			t.Errorf("expected error for %s", s)
		}
	}
}

func TestInferDimensionality(t *testing.T) {
	tests := []struct {
		count int
		want  types.Dimensionality
	}{
		{6, types.Dim3D},
		{9, types.Dim3D},
		{12, types.Dim3D},
		{4, types.Dim2D},
		{8, types.Dim2D},
	}
	for _, tt := range tests {
		got := inferDimensionality(tt.count)
		if got != tt.want {
			t.Errorf("inferDimensionality(%d) = %d, want %d", tt.count, got, tt.want)
		}
	}
}

func TestParsePosWhitespace(t *testing.T) {
	pt, _, err := ParsePos("  1.0   2.0   3.0  ")
	if err != nil {
		t.Fatal(err)
	}

	if pt.X != 1.0 {
		t.Errorf("X = %f, want 1.0", pt.X)
	}
}

func TestParsePosListLargeValues(t *testing.T) {
	pts, _, err := ParsePosList("1e10 2e10 3e10", types.Dim3D)
	if err != nil {
		t.Fatal(err)
	}

	if len(pts) != 1 {
		t.Fatalf("got %d points", len(pts))
	}

	if math.Abs(pts[0].X-1e10) > 1 {
		t.Errorf("X = %g", pts[0].X)
	}
}
