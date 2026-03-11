package gml

import (
	"fmt"
	"math"

	"github.com/cwbudde/go-citygml/types"
)

// ValidationError represents a geometry validation failure.
type ValidationError struct {
	// Path describes where the error occurred (e.g. "Polygon.exterior").
	Path string

	// Message describes the validation failure.
	Message string
}

func (e *ValidationError) Error() string {
	if e.Path != "" {
		return fmt.Sprintf("gml: validation: %s: %s", e.Path, e.Message)
	}

	return "gml: validation: " + e.Message
}

// ValidateRing checks that a ring satisfies geometry constraints:
//   - at least 4 points (3 distinct + closure)
//   - first and last points are identical (closed)
//   - all coordinates are finite
func ValidateRing(ring types.Ring, path string) []error {
	errs := make([]error, 0, 2)

	if len(ring.Points) < 4 {
		errs = append(errs, &ValidationError{
			Path:    path,
			Message: fmt.Sprintf("ring has %d points, minimum is 4", len(ring.Points)),
		})

		return errs // can't check further
	}

	first := ring.Points[0]

	last := ring.Points[len(ring.Points)-1]
	if first != last {
		errs = append(errs, &ValidationError{
			Path:    path,
			Message: fmt.Sprintf("ring not closed: first (%g,%g,%g) != last (%g,%g,%g)", first.X, first.Y, first.Z, last.X, last.Y, last.Z),
		})
	}

	for i, pt := range ring.Points {
		if !isFinitePoint(pt) {
			errs = append(errs, &ValidationError{
				Path:    fmt.Sprintf("%s[%d]", path, i),
				Message: fmt.Sprintf("non-finite coordinate (%g,%g,%g)", pt.X, pt.Y, pt.Z),
			})
		}
	}

	return errs
}

// ValidatePolygon validates a polygon's exterior and interior rings.
func ValidatePolygon(poly types.Polygon, path string) []error {
	errs := make([]error, 0, 1+len(poly.Interior))

	errs = append(errs, ValidateRing(poly.Exterior, path+".exterior")...)
	for i, ring := range poly.Interior {
		errs = append(errs, ValidateRing(ring, fmt.Sprintf("%s.interior[%d]", path, i))...)
	}

	return errs
}

// ValidateMultiSurface validates all polygons in a multi-surface.
func ValidateMultiSurface(ms types.MultiSurface, path string) []error {
	errs := make([]error, 0, len(ms.Polygons))
	for i, poly := range ms.Polygons {
		errs = append(errs, ValidatePolygon(poly, fmt.Sprintf("%s.polygon[%d]", path, i))...)
	}

	return errs
}

// ValidateSolid validates a solid's exterior shell.
func ValidateSolid(solid types.Solid, path string) []error {
	return ValidateMultiSurface(solid.Exterior, path+".exterior")
}

func isFinitePoint(pt types.Point) bool {
	return isFinite(pt.X) && isFinite(pt.Y) && isFinite(pt.Z)
}

func isFinite(v float64) bool {
	return !math.IsNaN(v) && !math.IsInf(v, 0)
}
