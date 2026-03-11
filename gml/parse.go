package gml

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"

	"github.com/cwbudde/go-citygml/internal/xmlscan"
	"github.com/cwbudde/go-citygml/types"
)

func wrapScannerError(context string, err error) error {
	return fmt.Errorf("gml: %s: %w", context, err)
}

// ParseLinearRing parses a gml:LinearRing element from the scanner.
// The scanner must be positioned just after the LinearRing StartElement.
func ParseLinearRing(sc *xmlscan.Scanner) (types.Ring, types.Dimensionality, error) {
	var points []types.Point
	var dim types.Dimensionality
	depth := 1

	for depth > 0 {
		tok, err := sc.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return types.Ring{}, 0, errors.New("gml: unexpected EOF in LinearRing")
			}

			return types.Ring{}, 0, wrapScannerError("LinearRing token", err)
		}

		switch t := tok.(type) {
		case xml.StartElement:
			depth++

			switch t.Name.Local {
			case "pos":
				text, err := sc.CharData()
				if err != nil {
					return types.Ring{}, 0, wrapScannerError("LinearRing pos char data", err)
				}

				depth-- // CharData consumed the EndElement

				pt, d, err := ParsePos(text)
				if err != nil {
					return types.Ring{}, 0, fmt.Errorf("gml: LinearRing pos: %w", err)
				}

				if dim == 0 {
					dim = d
				}

				points = append(points, pt)

			case "posList":
				text, err := sc.CharData()
				if err != nil {
					return types.Ring{}, 0, wrapScannerError("LinearRing posList char data", err)
				}

				depth-- // CharData consumed the EndElement

				pts, d, err := ParsePosList(text, 0)
				if err != nil {
					return types.Ring{}, 0, fmt.Errorf("gml: LinearRing posList: %w", err)
				}

				if dim == 0 {
					dim = d
				}

				points = append(points, pts...)

			default:
				// Skip unknown children.
				err := sc.Skip()
				if err != nil {
					return types.Ring{}, 0, wrapScannerError("LinearRing skip child", err)
				}

				depth--
			}

		case xml.EndElement:
			depth--
		}
	}

	return types.Ring{Points: points}, dim, nil
}

// ParsePolygon parses a gml:Polygon element from the scanner.
// The scanner must be positioned just after the Polygon StartElement.
func ParsePolygon(sc *xmlscan.Scanner) (types.Polygon, types.Dimensionality, error) {
	var poly types.Polygon
	var dim types.Dimensionality
	depth := 1

	for depth > 0 {
		tok, err := sc.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return types.Polygon{}, 0, errors.New("gml: unexpected EOF in Polygon")
			}

			return types.Polygon{}, 0, wrapScannerError("Polygon token", err)
		}

		switch t := tok.(type) {
		case xml.StartElement:
			depth++

			switch t.Name.Local {
			case "exterior":
				ring, d, err := parseRingWrapper(sc)
				if err != nil {
					return types.Polygon{}, 0, fmt.Errorf("gml: Polygon exterior: %w", err)
				}

				depth-- // parseRingWrapper consumed exterior's EndElement
				poly.Exterior = ring

				if dim == 0 {
					dim = d
				}

			case "interior":
				ring, d, err := parseRingWrapper(sc)
				if err != nil {
					return types.Polygon{}, 0, fmt.Errorf("gml: Polygon interior: %w", err)
				}

				depth--

				poly.Interior = append(poly.Interior, ring)

				if dim == 0 {
					dim = d
				}

			default:
				err := sc.Skip()
				if err != nil {
					return types.Polygon{}, 0, wrapScannerError("Polygon skip child", err)
				}

				depth--
			}

		case xml.EndElement:
			depth--
		}
	}

	return poly, dim, nil
}

// ParseMultiSurface parses a gml:MultiSurface element from the scanner.
// The scanner must be positioned just after the MultiSurface StartElement.
func ParseMultiSurface(sc *xmlscan.Scanner) (types.MultiSurface, types.Dimensionality, error) {
	return parseSurfaceCollection(sc, "MultiSurface")
}

// ParseSolid parses a gml:Solid element from the scanner.
// The scanner must be positioned just after the Solid StartElement.
func ParseSolid(sc *xmlscan.Scanner) (types.Solid, types.Dimensionality, error) {
	var solid types.Solid
	var dim types.Dimensionality
	depth := 1

	for depth > 0 {
		tok, err := sc.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return types.Solid{}, 0, errors.New("gml: unexpected EOF in Solid")
			}

			return types.Solid{}, 0, wrapScannerError("Solid token", err)
		}

		switch t := tok.(type) {
		case xml.StartElement:
			depth++

			switch t.Name.Local {
			case "exterior":
				ms, d, err := parseSolidExterior(sc)
				if err != nil {
					return types.Solid{}, 0, fmt.Errorf("gml: Solid exterior: %w", err)
				}

				depth--
				solid.Exterior = ms

				if dim == 0 {
					dim = d
				}

			default:
				err := sc.Skip()
				if err != nil {
					return types.Solid{}, 0, wrapScannerError("Solid skip child", err)
				}

				depth--
			}

		case xml.EndElement:
			depth--
		}
	}

	return solid, dim, nil
}

// ParseCompositeSurface parses a gml:CompositeSurface element.
// It is treated as equivalent to a MultiSurface for our purposes.
func ParseCompositeSurface(sc *xmlscan.Scanner) (types.MultiSurface, types.Dimensionality, error) {
	return parseSurfaceCollection(sc, "CompositeSurface")
}

func parseSurfaceCollection(sc *xmlscan.Scanner, elementName string) (types.MultiSurface, types.Dimensionality, error) {
	var ms types.MultiSurface
	var dim types.Dimensionality
	depth := 1

	for depth > 0 {
		tok, err := sc.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return types.MultiSurface{}, 0, fmt.Errorf("gml: unexpected EOF in %s", elementName)
			}

			return types.MultiSurface{}, 0, wrapScannerError(elementName+" token", err)
		}

		switch t := tok.(type) {
		case xml.StartElement:
			depth++

			switch t.Name.Local {
			case "surfaceMember":
				poly, d, err := parseSurfaceMember(sc)
				if err != nil {
					return types.MultiSurface{}, 0, fmt.Errorf("gml: %s surfaceMember: %w", elementName, err)
				}

				depth--

				if poly != nil {
					ms.Polygons = append(ms.Polygons, *poly)
				}

				if dim == 0 {
					dim = d
				}

			default:
				err := sc.Skip()
				if err != nil {
					return types.MultiSurface{}, 0, wrapScannerError(elementName+" skip child", err)
				}

				depth--
			}

		case xml.EndElement:
			depth--
		}
	}

	return ms, dim, nil
}

// parseRingWrapper reads a LinearRing inside an exterior/interior wrapper element.
// Consumes tokens up to and including the wrapper's EndElement.
func parseRingWrapper(sc *xmlscan.Scanner) (types.Ring, types.Dimensionality, error) {
	depth := 1
	for depth > 0 {
		tok, err := sc.Token()
		if err != nil {
			return types.Ring{}, 0, wrapScannerError("ring wrapper token", err)
		}

		switch t := tok.(type) {
		case xml.StartElement:
			depth++

			if t.Name.Local == "LinearRing" {
				ring, dim, err := ParseLinearRing(sc)
				if err != nil {
					return types.Ring{}, 0, fmt.Errorf("gml: ring wrapper LinearRing: %w", err)
				}

				depth--
				// Continue to consume the wrapper's EndElement.
				for depth > 0 {
					tok2, err := sc.Token()
					if err != nil {
						return types.Ring{}, 0, wrapScannerError("ring wrapper drain", err)
					}

					switch tok2.(type) {
					case xml.StartElement:
						depth++
					case xml.EndElement:
						depth--
					}
				}

				return ring, dim, nil
			}

			err := sc.Skip()
			if err != nil {
				return types.Ring{}, 0, wrapScannerError("ring wrapper skip child", err)
			}

			depth--

		case xml.EndElement:
			depth--
		}
	}

	return types.Ring{}, 0, errors.New("gml: no LinearRing found in exterior/interior")
}

// parseSurfaceMember reads a Polygon inside a surfaceMember wrapper.
// Consumes tokens up to and including the surfaceMember's EndElement.
func parseSurfaceMember(sc *xmlscan.Scanner) (*types.Polygon, types.Dimensionality, error) {
	depth := 1
	for depth > 0 {
		tok, err := sc.Token()
		if err != nil {
			return nil, 0, wrapScannerError("surfaceMember token", err)
		}

		switch t := tok.(type) {
		case xml.StartElement:
			depth++

			if t.Name.Local == "Polygon" {
				poly, dim, err := ParsePolygon(sc)
				if err != nil {
					return nil, 0, fmt.Errorf("gml: surfaceMember Polygon: %w", err)
				}

				depth--
				for depth > 0 {
					tok2, err := sc.Token()
					if err != nil {
						return nil, 0, wrapScannerError("surfaceMember drain", err)
					}

					switch tok2.(type) {
					case xml.StartElement:
						depth++
					case xml.EndElement:
						depth--
					}
				}

				return &poly, dim, nil
			}

			err := sc.Skip()
			if err != nil {
				return nil, 0, wrapScannerError("surfaceMember skip child", err)
			}

			depth--

		case xml.EndElement:
			depth--
		}
	}

	return nil, 0, nil // empty surfaceMember
}

// parseSolidExterior reads a CompositeSurface or Shell inside a Solid's exterior wrapper.
func parseSolidExterior(sc *xmlscan.Scanner) (types.MultiSurface, types.Dimensionality, error) {
	depth := 1
	for depth > 0 {
		tok, err := sc.Token()
		if err != nil {
			return types.MultiSurface{}, 0, wrapScannerError("Solid exterior token", err)
		}

		switch t := tok.(type) {
		case xml.StartElement:
			depth++

			switch t.Name.Local {
			case "CompositeSurface", "Shell":
				ms, dim, err := ParseCompositeSurface(sc)
				if err != nil {
					return types.MultiSurface{}, 0, fmt.Errorf("gml: Solid exterior %s: %w", t.Name.Local, err)
				}

				depth--
				for depth > 0 {
					tok2, err := sc.Token()
					if err != nil {
						return types.MultiSurface{}, 0, wrapScannerError("Solid exterior drain", err)
					}

					switch tok2.(type) {
					case xml.StartElement:
						depth++
					case xml.EndElement:
						depth--
					}
				}

				return ms, dim, nil

			default:
				err := sc.Skip()
				if err != nil {
					return types.MultiSurface{}, 0, wrapScannerError("Solid exterior skip child", err)
				}

				depth--
			}

		case xml.EndElement:
			depth--
		}
	}

	return types.MultiSurface{}, 0, errors.New("gml: no CompositeSurface/Shell in Solid exterior")
}
