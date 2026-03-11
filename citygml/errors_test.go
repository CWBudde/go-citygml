package citygml

import (
	"errors"
	"testing"
)

func TestParseError_Unwrap(t *testing.T) {
	err := &ParseError{
		Err:    ErrMalformedXML,
		Path:   "CityModel/cityObjectMember[0]",
		Detail: "unexpected EOF",
	}

	if !errors.Is(err, ErrMalformedXML) {
		t.Error("expected errors.Is to match ErrMalformedXML")
	}

	want := "citygml: malformed XML at CityModel/cityObjectMember[0]: unexpected EOF"
	if got := err.Error(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestParseError_PathOnly(t *testing.T) {
	err := &ParseError{Err: ErrUnsupportedGeometry, Path: "Building/lod2Solid"}

	want := "citygml: unsupported geometry type at Building/lod2Solid"
	if got := err.Error(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestParseError_DetailOnly(t *testing.T) {
	err := &ParseError{Err: ErrInvalidCoordinates, Detail: "expected 3D, got 2D"}

	want := "citygml: invalid coordinates: expected 3D, got 2D"
	if got := err.Error(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestParseError_Bare(t *testing.T) {
	err := &ParseError{Err: ErrInvalidCRS}
	if got := err.Error(); got != ErrInvalidCRS.Error() {
		t.Errorf("got %q, want %q", got, ErrInvalidCRS.Error())
	}
}
