package decode

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"

	"github.com/cwbudde/go-citygml/gml"
	"github.com/cwbudde/go-citygml/internal/xmlscan"
	"github.com/cwbudde/go-citygml/types"
)

// IsTerrainElement returns true if the element is a recognized terrain/relief element.
func IsTerrainElement(elem *xmlscan.Element) bool {
	ns := elem.Namespace()
	local := elem.LocalName()

	if ns == xmlscan.NSCityGML20Dem || ns == xmlscan.NSCityGML30Dem {
		switch local {
		case "ReliefFeature", "TINRelief", "MassPointRelief", "BreaklineRelief", "RasterRelief":
			return true
		}
	}

	return false
}

// Terrain decodes a terrain element from the scanner.
// The scanner must be positioned just after the terrain StartElement.
// It handles ReliefFeature (which may contain nested relief components)
// and direct TINRelief elements.
func Terrain(elem *xmlscan.Element, sc *xmlscan.Scanner) ([]types.Terrain, error) {
	if elem.LocalName() == "ReliefFeature" {
		return decodeReliefFeature(elem, sc)
	}
	// Direct relief component (e.g. TINRelief at top level).
	t, err := decodeReliefComponent(elem, sc)
	if err != nil {
		return nil, err
	}

	if t != nil {
		return []types.Terrain{*t}, nil
	}

	return nil, nil
}

// decodeReliefFeature reads a ReliefFeature, which wraps one or more relief components.
func decodeReliefFeature(elem *xmlscan.Element, sc *xmlscan.Scanner) ([]types.Terrain, error) {
	var terrains []types.Terrain
	depth := 1

	for depth > 0 {
		tok, err := sc.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return terrains, fmt.Errorf("decode: unexpected EOF in ReliefFeature %s", elem.ID)
			}

			return terrains, err
		}

		switch t := tok.(type) {
		case xml.StartElement:
			depth++

			if t.Name.Local == "reliefComponent" {
				comps, err := decodeReliefComponentWrapper(sc)
				if err != nil {
					return terrains, err
				}

				depth--

				terrains = append(terrains, comps...)
			} else {
				err := sc.Skip()
				if err != nil {
					return terrains, err
				}

				depth--
			}

		case xml.EndElement:
			depth--
		}
	}

	return terrains, nil
}

// decodeReliefComponentWrapper reads the inner element of a reliefComponent wrapper.
func decodeReliefComponentWrapper(sc *xmlscan.Scanner) ([]types.Terrain, error) {
	depth := 1
	for depth > 0 {
		tok, err := sc.Token()
		if err != nil {
			return nil, err
		}

		switch t := tok.(type) {
		case xml.StartElement:
			depth++
			wrapped := sc.WrapElement(t)

			comp, err := decodeReliefComponent(wrapped, sc)
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

			if comp != nil {
				return []types.Terrain{*comp}, nil
			}

			return nil, nil

		case xml.EndElement:
			depth--
		}
	}

	return nil, nil
}

// decodeReliefComponent reads a TINRelief or similar relief component.
func decodeReliefComponent(elem *xmlscan.Element, sc *xmlscan.Scanner) (*types.Terrain, error) {
	t := &types.Terrain{
		ID: elem.ID,
	}

	depth := 1
	for depth > 0 {
		tok, err := sc.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return t, fmt.Errorf("decode: unexpected EOF in %s %s", elem.LocalName(), elem.ID)
			}

			return t, err
		}

		switch te := tok.(type) {
		case xml.StartElement:
			depth++

			switch te.Name.Local {
			case "tin", "triangulatedSurface", "lod1MultiSurface", "lod2MultiSurface":
				ms, err := decodeLodMultiSurface(sc)
				if err != nil {
					return nil, fmt.Errorf("decode: terrain geometry: %w", err)
				}

				depth--

				if ms != nil {
					t.Geometry = *ms
				}

			case "MultiSurface":
				ms, _, err := gml.ParseMultiSurface(sc)
				if err != nil {
					return nil, fmt.Errorf("decode: terrain MultiSurface: %w", err)
				}

				depth--
				t.Geometry = ms

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

	return t, nil
}
