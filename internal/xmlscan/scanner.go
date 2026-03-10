package xmlscan

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"
)

// Element represents an XML start element with tracked context.
type Element struct {
	xml.StartElement

	// ID is the gml:id attribute value, if present.
	ID string

	// XLinkHref is the xlink:href attribute value, if present.
	XLinkHref string
}

// LocalName returns the local name of the element.
func (e Element) LocalName() string {
	return e.Name.Local
}

// Namespace returns the namespace URI of the element.
func (e Element) Namespace() string {
	return e.Name.Space
}

// Scanner provides namespace-aware, forward-only XML scanning with
// element path tracking and attribute extraction.
type Scanner struct {
	dec  *xml.Decoder
	path []string // element path stack for error context
	err  error    // sticky error

	// DetectedVersion is set when the root element's namespaces are inspected.
	DetectedVersion Version
}

// NewScanner creates a Scanner from the given reader.
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		dec: xml.NewDecoder(r),
	}
}

// Path returns the current element path as a slash-separated string
// (e.g. "CityModel/cityObjectMember/Building").
func (s *Scanner) Path() string {
	return strings.Join(s.path, "/")
}

// Err returns the first non-EOF error encountered during scanning.
func (s *Scanner) Err() error {
	return s.err
}

// Token advances the decoder and returns the next XML token.
// On the first StartElement it detects the CityGML version from namespaces.
// Returns nil, io.EOF when the document ends.
func (s *Scanner) Token() (xml.Token, error) {
	if s.err != nil {
		return nil, s.err
	}

	tok, err := s.dec.Token()
	if err != nil {
		if !errors.Is(err, io.EOF) {
			s.err = fmt.Errorf("xmlscan: %w", err)
		}

		return nil, err
	}

	switch t := tok.(type) {
	case xml.StartElement:
		s.path = append(s.path, t.Name.Local)

		// Detect version from root element namespaces.
		if len(s.path) == 1 && s.DetectedVersion == VersionUnknown {
			s.detectVersionFromElement(t)
		}

	case xml.EndElement:
		if len(s.path) > 0 {
			s.path = s.path[:len(s.path)-1]
		}
	}

	return tok, nil
}

// StartElement advances past non-element tokens and returns the next StartElement.
// Returns nil, io.EOF if no more start elements exist.
func (s *Scanner) StartElement() (*Element, error) {
	for {
		tok, err := s.Token()
		if err != nil {
			return nil, err
		}

		if se, ok := tok.(xml.StartElement); ok {
			return s.WrapElement(se), nil
		}
	}
}

// Skip consumes tokens until the current element's matching EndElement is consumed.
// Call this after receiving a StartElement to skip its entire subtree.
func (s *Scanner) Skip() error {
	depth := 1
	for depth > 0 {
		tok, err := s.Token()
		if err != nil {
			return err
		}

		switch tok.(type) {
		case xml.StartElement:
			depth++
		case xml.EndElement:
			depth--
		}
	}

	return nil
}

// CharData reads and returns the concatenated character data inside the
// current element. It stops at the matching EndElement.
func (s *Scanner) CharData() (string, error) {
	var buf strings.Builder

	depth := 1
	for depth > 0 {
		tok, err := s.Token()
		if err != nil {
			return "", err
		}

		switch t := tok.(type) {
		case xml.CharData:
			buf.Write(t)
		case xml.StartElement:
			depth++
		case xml.EndElement:
			depth--
		}
	}

	return buf.String(), nil
}

// WrapElement extracts gml:id and xlink:href from a StartElement.
func (s *Scanner) WrapElement(se xml.StartElement) *Element {
	e := &Element{StartElement: se}
	for _, attr := range se.Attr {
		switch {
		case attr.Name.Local == "id" && isGMLNamespace(attr.Name.Space):
			e.ID = attr.Value
		case attr.Name.Local == "href" && attr.Name.Space == NSXLink:
			e.XLinkHref = attr.Value
		}
	}

	return e
}

// detectVersionFromElement inspects namespace declarations on the root element.
func (s *Scanner) detectVersionFromElement(se xml.StartElement) {
	// Check the element's own namespace first.
	if v := DetectVersion(se.Name.Space); v != VersionUnknown {
		s.DetectedVersion = v
		return
	}

	// Check xmlns attributes for a core CityGML namespace.
	for _, attr := range se.Attr {
		if v := DetectVersion(attr.Value); v != VersionUnknown {
			s.DetectedVersion = v
			return
		}
	}
}

// isGMLNamespace returns true if the URI is a known GML namespace.
func isGMLNamespace(uri string) bool {
	return uri == NSGML31 || uri == NSGML32
}
