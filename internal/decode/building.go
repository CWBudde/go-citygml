// Package decode contains internal decoders for mapping XML elements
// to normalized semantic types.
package decode

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/cwbudde/go-citygml/gml"
	"github.com/cwbudde/go-citygml/internal/xmlscan"
	"github.com/cwbudde/go-citygml/types"
)

const (
	lod1MultiSurfaceElement = "lod1MultiSurface"
	lod2MultiSurfaceElement = "lod2MultiSurface"
	multiSurfaceElement     = "MultiSurface"
)

// IsBuildingElement returns true if the element is a recognized building element.
func IsBuildingElement(elem *xmlscan.Element) bool {
	ns := elem.Namespace()
	return ns == xmlscan.NSCityGML20Bldg || ns == xmlscan.NSCityGML30Bldg
}

// Building decodes a Building element from the scanner.
// The scanner must be positioned just after the Building StartElement.
func Building(elem *xmlscan.Element, sc *xmlscan.Scanner) (types.Building, error) {
	b := types.Building{
		ID: elem.ID,
	}

	depth := 1
	for depth > 0 {
		tok, err := sc.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return b, fmt.Errorf("decode: unexpected EOF in Building %s", b.ID)
			}

			return b, err
		}

		switch t := tok.(type) {
		case xml.StartElement:
			depth++

			err := decodeBuildingChild(t, sc, &b, &depth)
			if err != nil {
				return b, err
			}

		case xml.EndElement:
			depth--
		}
	}

	return b, nil
}

func decodeBuildingChild(se xml.StartElement, sc *xmlscan.Scanner, b *types.Building, depth *int) error {
	local := se.Name.Local

	switch local {
	case "measuredHeight":
		text, err := sc.CharData()
		if err != nil {
			return err
		}

		*depth--

		v, err := strconv.ParseFloat(strings.TrimSpace(text), 64)
		if err != nil {
			return fmt.Errorf("decode: measuredHeight: %w", err)
		}

		b.MeasuredHeight = v
		b.HasMeasuredHeight = true

	case "class":
		text, err := sc.CharData()
		if err != nil {
			return err
		}

		*depth--
		b.Class = strings.TrimSpace(text)

	case "function":
		text, err := sc.CharData()
		if err != nil {
			return err
		}

		*depth--
		b.Function = strings.TrimSpace(text)

	case "usage":
		text, err := sc.CharData()
		if err != nil {
			return err
		}

		*depth--
		b.Usage = strings.TrimSpace(text)

	case "lod1Solid":
		solid, err := decodeLodSolid(sc)
		if err != nil {
			return fmt.Errorf("decode: lod1Solid: %w", err)
		}

		*depth--
		b.LoD = types.LoD1
		b.Solid = solid

	case "lod2Solid":
		solid, err := decodeLodSolid(sc)
		if err != nil {
			return fmt.Errorf("decode: lod2Solid: %w", err)
		}

		*depth--
		b.LoD = types.LoD2
		b.Solid = solid

	case lod1MultiSurfaceElement:
		ms, err := decodeLodMultiSurface(sc)
		if err != nil {
			return fmt.Errorf("decode: %s: %w", lod1MultiSurfaceElement, err)
		}

		*depth--
		b.LoD = types.LoD1
		b.MultiSurface = ms

	case lod2MultiSurfaceElement:
		ms, err := decodeLodMultiSurface(sc)
		if err != nil {
			return fmt.Errorf("decode: %s: %w", lod2MultiSurfaceElement, err)
		}

		*depth--
		b.LoD = types.LoD2
		b.MultiSurface = ms

	case "boundedBy":
		surf, err := decodeBoundedBy(sc)
		if err != nil {
			return fmt.Errorf("decode: boundedBy: %w", err)
		}

		*depth--

		if surf != nil {
			b.BoundedBy = append(b.BoundedBy, *surf)
		}

	default:
		err := sc.Skip()
		if err != nil {
			return err
		}

		*depth--
	}

	return nil
}

// decodeLodSolid reads a Solid inside a lodXSolid wrapper.
func decodeLodSolid(sc *xmlscan.Scanner) (*types.Solid, error) {
	depth := 1
	for depth > 0 {
		tok, err := sc.Token()
		if err != nil {
			return nil, err
		}

		switch t := tok.(type) {
		case xml.StartElement:
			depth++

			if t.Name.Local == "Solid" {
				solid, _, err := gml.ParseSolid(sc)
				if err != nil {
					return nil, err
				}

				depth--
				// Drain remaining tokens in wrapper.
				for depth > 0 {
					tok2, err := sc.Token()
					if err != nil {
						return nil, err
					}

					switch tok2.(type) {
					case xml.StartElement:
						depth++
					case xml.EndElement:
						depth--
					}
				}

				return &solid, nil
			}

			err := sc.Skip()
			if err != nil {
				return nil, err
			}

			depth--

		case xml.EndElement:
			depth--
		}
	}

	//nolint:nilnil // Missing wrapped Solid is a valid absence case for optional LoD geometry.
	return nil, nil
}

// decodeLodMultiSurface reads a MultiSurface inside a lodXMultiSurface wrapper.
func decodeLodMultiSurface(sc *xmlscan.Scanner) (*types.MultiSurface, error) {
	depth := 1
	for depth > 0 {
		tok, err := sc.Token()
		if err != nil {
			return nil, err
		}

		switch t := tok.(type) {
		case xml.StartElement:
			depth++

			if t.Name.Local == multiSurfaceElement {
				ms, _, err := gml.ParseMultiSurface(sc)
				if err != nil {
					return nil, err
				}

				depth--
				for depth > 0 {
					tok2, err := sc.Token()
					if err != nil {
						return nil, err
					}

					switch tok2.(type) {
					case xml.StartElement:
						depth++
					case xml.EndElement:
						depth--
					}
				}

				return &ms, nil
			}

			err := sc.Skip()
			if err != nil {
				return nil, err
			}

			depth--

		case xml.EndElement:
			depth--
		}
	}

	//nolint:nilnil // Missing wrapped MultiSurface is a valid absence case for optional LoD geometry.
	return nil, nil
}

// decodeBoundedBy reads a semantic surface inside a boundedBy wrapper.
func decodeBoundedBy(sc *xmlscan.Scanner) (*types.Surface, error) {
	depth := 1
	for depth > 0 {
		tok, err := sc.Token()
		if err != nil {
			return nil, err
		}

		switch t := tok.(type) {
		case xml.StartElement:
			depth++
			local := t.Name.Local

			// Recognized bounded surface types.
			switch local {
			case "RoofSurface", "WallSurface", "GroundSurface",
				"ClosureSurface", "CeilingSurface", "FloorSurface",
				"OuterCeilingSurface", "OuterFloorSurface":
				wrapped := sc.WrapElement(t)

				surf, err := decodeSurfaceElement(local, wrapped, sc)
				if err != nil {
					return nil, err
				}

				depth--
				// Drain wrapper.
				for depth > 0 {
					tok2, err := sc.Token()
					if err != nil {
						return nil, err
					}

					switch tok2.(type) {
					case xml.StartElement:
						depth++
					case xml.EndElement:
						depth--
					}
				}

				return surf, nil

			default:
				err := sc.Skip()
				if err != nil {
					return nil, err
				}

				depth--
			}

		case xml.EndElement:
			depth--
		}
	}

	//nolint:nilnil // Missing recognized boundedBy surface is a valid absence case.
	return nil, nil
}

// decodeSurfaceElement reads geometry from a bounded surface element.
func decodeSurfaceElement(surfType string, elem *xmlscan.Element, sc *xmlscan.Scanner) (*types.Surface, error) {
	surf := &types.Surface{
		ID:   elem.ID,
		Type: surfType,
	}

	depth := 1
	for depth > 0 {
		tok, err := sc.Token()
		if err != nil {
			return nil, err
		}

		switch t := tok.(type) {
		case xml.StartElement:
			depth++
			// Look for lodXMultiSurface or MultiSurface directly.
			switch t.Name.Local {
			case "lod1MultiSurface", "lod2MultiSurface", "lod3MultiSurface", "lod4MultiSurface":
				ms, err := decodeLodMultiSurface(sc)
				if err != nil {
					return nil, err
				}

				depth--

				if ms != nil {
					surf.Geometry = *ms
				}

			case "MultiSurface":
				ms, _, err := gml.ParseMultiSurface(sc)
				if err != nil {
					return nil, err
				}

				depth--
				surf.Geometry = ms

			case "Polygon":
				poly, _, err := gml.ParsePolygon(sc)
				if err != nil {
					return nil, err
				}

				depth--
				surf.Geometry = types.MultiSurface{Polygons: []types.Polygon{poly}}

			default:
				err := sc.Skip()
				if err != nil {
					return nil, err
				}

				depth--
			}

		case xml.EndElement:
			depth--
		}
	}

	return surf, nil
}
