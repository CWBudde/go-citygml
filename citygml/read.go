package citygml

import (
	"fmt"
	"io"
	"os"

	"github.com/cwbudde/go-citygml/crs"
	"github.com/cwbudde/go-citygml/gml"
	"github.com/cwbudde/go-citygml/internal/decode"
	"github.com/cwbudde/go-citygml/internal/xmlscan"
	"github.com/cwbudde/go-citygml/types"
)

// Read decodes a CityGML document from r and returns a normalized Document.
func Read(r io.Reader, opts Options) (*types.Document, error) {
	sc := xmlscan.NewScanner(r)

	hdr, err := xmlscan.ReadHeader(sc)
	if err != nil {
		return nil, &ParseError{Err: ErrMalformedXML, Detail: err.Error()}
	}

	if hdr.Version == xmlscan.VersionUnknown {
		if opts.Strict {
			return nil, &ParseError{Err: ErrUnsupportedVersion, Detail: "could not detect CityGML version from namespaces"}
		}
	}

	doc := &types.Document{
		Version: string(hdr.Version),
	}

	err = xmlscan.EachCityObjectMember(sc, hdr, func(elem *xmlscan.Element, sc *xmlscan.Scanner) error {
		if decode.IsBuildingElement(elem) {
			b, err := decode.Building(elem, sc)
			if err != nil {
				return fmt.Errorf("citygml: decode building %q: %w", elem.ID, err)
			}

			postProcessBuilding(&b, opts)
			doc.Buildings = append(doc.Buildings, b)

			return nil
		}

		if decode.IsTerrainElement(elem) {
			terrains, err := decode.Terrain(elem, sc)
			if err != nil {
				return fmt.Errorf("citygml: decode terrain %q: %w", elem.ID, err)
			}

			doc.Terrains = append(doc.Terrains, terrains...)

			return nil
		}

		// Unsupported object type.
		if opts.Strict {
			return &ParseError{
				Err:    ErrUnsupportedObject,
				Path:   sc.Path(),
				Detail: fmt.Sprintf("unsupported element: %s (ns: %s)", elem.LocalName(), elem.Namespace()),
			}
		}

		doc.GenericObjects = append(doc.GenericObjects, types.CityObject{
			ID:   elem.ID,
			Type: elem.LocalName(),
		})

		err := sc.Skip()
		if err != nil {
			return fmt.Errorf("citygml: skip unsupported element %s: %w", elem.LocalName(), err)
		}

		return nil
	})

	// Set CRS after iteration since boundedBy/Envelope may have been parsed during it.
	doc.SRSName = hdr.SRSName
	if hdr.SRSName != "" {
		doc.CRS = crs.Parse(hdr.SRSName)
	}

	if err != nil {
		return doc, fmt.Errorf("citygml: iterate city objects: %w", err)
	}

	// In strict mode, missing CRS is an error.
	if opts.Strict && doc.SRSName == "" {
		return doc, &ParseError{Err: ErrInvalidCRS, Detail: "no srsName found in document"}
	}

	return doc, nil
}

// ReadFile is a convenience wrapper that opens a file and calls Read.
func ReadFile(path string, opts Options) (*types.Document, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("citygml: open %s: %w", path, err)
	}
	defer f.Close()

	return Read(f, opts)
}

func postProcessBuilding(b *types.Building, opts Options) {
	if opts.deriveHeights() && !b.HasMeasuredHeight {
		b.DerivedHeight = gml.DeriveHeight(b.Solid, b.MultiSurface, b.BoundedBy)
	}

	if opts.deriveFootprints() {
		b.Footprint = gml.DeriveFootprint(b.Solid, b.MultiSurface, b.BoundedBy)
	}
}
