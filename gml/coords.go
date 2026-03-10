package gml

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/cwbudde/go-citygml/types"
)

// ParsePos parses a gml:pos text content into a single Point.
// The text must contain 2 or 3 whitespace-separated float values.
func ParsePos(text string) (types.Point, types.Dimensionality, error) {
	fields := strings.Fields(text)
	switch len(fields) {
	case 2:
		x, y, err := parseXY(fields[0], fields[1])
		if err != nil {
			return types.Point{}, 0, err
		}

		return types.Point{X: x, Y: y}, types.Dim2D, nil
	case 3:
		x, y, z, err := parseXYZ(fields[0], fields[1], fields[2])
		if err != nil {
			return types.Point{}, 0, err
		}

		return types.Point{X: x, Y: y, Z: z}, types.Dim3D, nil
	default:
		return types.Point{}, 0, fmt.Errorf("gml: pos has %d values, expected 2 or 3", len(fields))
	}
}

// ParsePosList parses a gml:posList text content into a slice of Points.
// dim specifies the coordinate dimensionality (2 or 3). If dim is 0, it is
// inferred: 3 if the total count is divisible by 3, otherwise 2.
func ParsePosList(text string, dim types.Dimensionality) ([]types.Point, types.Dimensionality, error) {
	fields := strings.Fields(text)
	if len(fields) == 0 {
		return nil, 0, errors.New("gml: empty posList")
	}

	if dim == 0 {
		dim = inferDimensionality(len(fields))
	}

	d := int(dim)
	if len(fields)%d != 0 {
		return nil, 0, fmt.Errorf("gml: posList has %d values, not divisible by %d", len(fields), d)
	}

	n := len(fields) / d
	points := make([]types.Point, 0, n)

	for i := 0; i < len(fields); i += d {
		switch d {
		case 2:
			x, y, err := parseXY(fields[i], fields[i+1])
			if err != nil {
				return nil, 0, fmt.Errorf("gml: posList coordinate %d: %w", i/d, err)
			}

			points = append(points, types.Point{X: x, Y: y})
		case 3:
			x, y, z, err := parseXYZ(fields[i], fields[i+1], fields[i+2])
			if err != nil {
				return nil, 0, fmt.Errorf("gml: posList coordinate %d: %w", i/d, err)
			}

			points = append(points, types.Point{X: x, Y: y, Z: z})
		}
	}

	return points, dim, nil
}

func inferDimensionality(count int) types.Dimensionality {
	if count%3 == 0 {
		return types.Dim3D
	}

	return types.Dim2D
}

func parseXY(sx, sy string) (float64, float64, error) {
	x, err := parseFloat(sx)
	if err != nil {
		return 0, 0, err
	}

	y, err := parseFloat(sy)
	if err != nil {
		return 0, 0, err
	}

	return x, y, nil
}

func parseXYZ(sx, sy, sz string) (float64, float64, float64, error) {
	x, y, err := parseXY(sx, sy)
	if err != nil {
		return 0, 0, 0, err
	}

	z, err := parseFloat(sz)
	if err != nil {
		return 0, 0, 0, err
	}

	return x, y, z, nil
}

func parseFloat(s string) (float64, error) {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("gml: invalid coordinate %q: %w", s, err)
	}

	if math.IsNaN(v) || math.IsInf(v, 0) {
		return 0, fmt.Errorf("gml: non-finite coordinate %q", s)
	}

	return v, nil
}
