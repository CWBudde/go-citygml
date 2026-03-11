package xmlscan

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
)

// CityObjectMember represents a single cityObjectMember entry
// containing one top-level city object.
type CityObjectMember struct {
	// Element is the city object's start element (e.g. bldg:Building).
	Element *Element

	// RawXML contains the inner XML of the city object for downstream parsing.
	RawXML []byte
}

// DocumentHeader holds metadata extracted from the CityModel root element.
type DocumentHeader struct {
	Version Version
	SRSName string
}

// ReadHeader reads the root CityModel element and extracts version and CRS metadata.
// After calling ReadHeader, the scanner is positioned inside the root element,
// ready to iterate over cityObjectMember children.
func ReadHeader(sc *Scanner) (*DocumentHeader, error) {
	elem, err := sc.StartElement()
	if err != nil {
		return nil, fmt.Errorf("xmlscan: reading root element: %w", err)
	}

	if elem.LocalName() != "CityModel" {
		return nil, fmt.Errorf("xmlscan: expected CityModel root, got %s", elem.LocalName())
	}

	hdr := &DocumentHeader{
		Version: sc.DetectedVersion,
	}

	// Extract srsName from envelope or root attributes if present.
	for _, attr := range elem.Attr {
		if attr.Name.Local == "srsName" {
			hdr.SRSName = attr.Value
			break
		}
	}

	return hdr, nil
}

// EachCityObjectMember iterates over cityObjectMember children of the root CityModel.
// It calls fn for each city object found inside a cityObjectMember wrapper.
// Elements are visited in document order (deterministic).
// If fn returns an error, iteration stops and that error is returned.
//
// If a gml:boundedBy element is encountered, it extracts srsName from the
// gml:Envelope inside it and stores it on the header (if not already set).
func EachCityObjectMember(sc *Scanner, hdr *DocumentHeader, fn func(elem *Element, sc *Scanner) error) error {
	for {
		tok, err := sc.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}

			return err
		}

		se, ok := tok.(xml.StartElement)
		if !ok {
			continue
		}

		if se.Name.Local == "boundedBy" {
			extractEnvelopeSRS(sc, hdr)
			continue
		}

		if se.Name.Local != "cityObjectMember" {
			err := sc.Skip()
			if err != nil {
				return err
			}

			continue
		}

		// Read the city object inside this cityObjectMember.
		inner, err := sc.StartElement()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}

			return fmt.Errorf("xmlscan: empty cityObjectMember at %s", sc.Path())
		}

		err = fn(inner, sc)
		if err != nil {
			return err
		}
	}
}

// extractEnvelopeSRS reads inside a boundedBy element looking for an Envelope
// with an srsName attribute. Sets hdr.SRSName if not already set.
func extractEnvelopeSRS(sc *Scanner, hdr *DocumentHeader) {
	depth := 1
	for depth > 0 {
		tok, err := sc.Token()
		if err != nil {
			return
		}

		switch t := tok.(type) {
		case xml.StartElement:
			depth++

			if t.Name.Local == "Envelope" {
				for _, attr := range t.Attr {
					if attr.Name.Local == "srsName" && hdr.SRSName == "" {
						hdr.SRSName = attr.Value
					}
				}
			}
		case xml.EndElement:
			depth--
		}
	}
}
