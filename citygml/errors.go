package citygml

import (
	"errors"
	"fmt"
)

// Sentinel errors for type-checking with errors.Is.
var (
	// ErrMalformedXML indicates the input is not well-formed XML.
	ErrMalformedXML = errors.New("citygml: malformed XML")

	// ErrUnsupportedVersion indicates the CityGML version or profile is not supported.
	ErrUnsupportedVersion = errors.New("citygml: unsupported CityGML version")

	// ErrUnsupportedGeometry indicates an unrecognized or unsupported geometry type.
	ErrUnsupportedGeometry = errors.New("citygml: unsupported geometry type")

	// ErrUnsupportedObject indicates an unrecognized or unsupported city object type.
	ErrUnsupportedObject = errors.New("citygml: unsupported object type")

	// ErrInvalidCoordinates indicates invalid coordinate data (wrong dimensionality,
	// non-finite values, or insufficient points).
	ErrInvalidCoordinates = errors.New("citygml: invalid coordinates")

	// ErrInvalidCRS indicates missing or unrecognized CRS metadata.
	ErrInvalidCRS = errors.New("citygml: invalid CRS metadata")
)

// ParseError provides structured context for a parsing failure.
type ParseError struct {
	// Err is the underlying sentinel or wrapped error.
	Err error

	// Path is the element path within the document where the error occurred
	// (e.g. "CityModel/cityObjectMember[2]/Building/lod1Solid").
	Path string

	// Detail provides additional context about the failure.
	Detail string
}

func (e *ParseError) Error() string {
	if e.Path != "" && e.Detail != "" {
		return fmt.Sprintf("%v at %s: %s", e.Err, e.Path, e.Detail)
	}

	if e.Path != "" {
		return fmt.Sprintf("%v at %s", e.Err, e.Path)
	}

	if e.Detail != "" {
		return fmt.Sprintf("%v: %s", e.Err, e.Detail)
	}

	return e.Err.Error()
}

func (e *ParseError) Unwrap() error {
	return e.Err
}
