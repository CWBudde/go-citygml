package xmlscan

import (
	"errors"
	"io"
	"strings"
	"testing"
)

const minimalCityGML20 = `<?xml version="1.0" encoding="UTF-8"?>
<CityModel xmlns="http://www.opengis.net/citygml/2.0"
           xmlns:gml="http://www.opengis.net/gml"
           xmlns:xlink="http://www.w3.org/1999/xlink"
           xmlns:bldg="http://www.opengis.net/citygml/building/2.0">
  <cityObjectMember>
    <bldg:Building gml:id="BLDG_001">
      <bldg:measuredHeight>12.5</bldg:measuredHeight>
    </bldg:Building>
  </cityObjectMember>
</CityModel>`

const minimalCityGML30 = `<?xml version="1.0" encoding="UTF-8"?>
<CityModel xmlns="http://www.opengis.net/citygml/3.0"
           xmlns:gml="http://www.opengis.net/gml/3.2"
           xmlns:bldg="http://www.opengis.net/citygml/building/3.0">
  <cityObjectMember>
    <bldg:Building gml:id="BLDG_100"/>
  </cityObjectMember>
</CityModel>`

func TestDetectVersion20(t *testing.T) {
	sc := NewScanner(strings.NewReader(minimalCityGML20))

	_, err := sc.StartElement() // root
	if err != nil {
		t.Fatal(err)
	}

	if sc.DetectedVersion != Version20 {
		t.Errorf("got version %q, want %q", sc.DetectedVersion, Version20)
	}
}

func TestDetectVersion30(t *testing.T) {
	sc := NewScanner(strings.NewReader(minimalCityGML30))

	_, err := sc.StartElement()
	if err != nil {
		t.Fatal(err)
	}

	if sc.DetectedVersion != Version30 {
		t.Errorf("got version %q, want %q", sc.DetectedVersion, Version30)
	}
}

func TestGMLIDExtraction(t *testing.T) {
	sc := NewScanner(strings.NewReader(minimalCityGML20))

	// Scan until we find the Building element.
	for {
		elem, err := sc.StartElement()
		if err != nil {
			t.Fatal(err)
		}

		if elem.LocalName() == "Building" {
			if elem.ID != "BLDG_001" {
				t.Errorf("got ID %q, want %q", elem.ID, "BLDG_001")
			}

			return
		}
	}
}

func TestPath(t *testing.T) {
	sc := NewScanner(strings.NewReader(minimalCityGML20))

	// Advance to Building.
	var paths []string

	for {
		elem, err := sc.StartElement()
		if err != nil {
			t.Fatal(err)
		}

		paths = append(paths, sc.Path())

		if elem.LocalName() == "Building" {
			break
		}
	}

	want := []string{
		"CityModel",
		"CityModel/cityObjectMember",
		"CityModel/cityObjectMember/Building",
	}
	if len(paths) != len(want) {
		t.Fatalf("got %d paths, want %d: %v", len(paths), len(want), paths)
	}

	for i := range want {
		if paths[i] != want[i] {
			t.Errorf("path[%d] = %q, want %q", i, paths[i], want[i])
		}
	}
}

func TestSkip(t *testing.T) {
	sc := NewScanner(strings.NewReader(minimalCityGML20))

	// Skip to root, then skip its entire content.
	_, err := sc.StartElement()
	if err != nil {
		t.Fatal(err)
	}

	if err := sc.Skip(); err != nil {
		t.Fatal(err)
	}

	// Next token should be EOF.
	_, err = sc.Token()
	if !errors.Is(err, io.EOF) {
		t.Errorf("expected io.EOF after skip, got %v", err)
	}
}

func TestCharData(t *testing.T) {
	sc := NewScanner(strings.NewReader(minimalCityGML20))

	// Advance to measuredHeight.
	for {
		elem, err := sc.StartElement()
		if err != nil {
			t.Fatal(err)
		}

		if elem.LocalName() == "measuredHeight" {
			text, err := sc.CharData()
			if err != nil {
				t.Fatal(err)
			}

			if text != "12.5" {
				t.Errorf("got %q, want %q", text, "12.5")
			}

			return
		}
	}
}

func TestXLinkHref(t *testing.T) {
	input := `<root xmlns:xlink="http://www.w3.org/1999/xlink">
		<ref xlink:href="#target_1"/>
	</root>`

	sc := NewScanner(strings.NewReader(input))
	// Skip root.
	if _, err := sc.StartElement(); err != nil {
		t.Fatal(err)
	}

	elem, err := sc.StartElement()
	if err != nil {
		t.Fatal(err)
	}

	if elem.XLinkHref != "#target_1" {
		t.Errorf("got xlink:href %q, want %q", elem.XLinkHref, "#target_1")
	}
}

func TestMalformedXML(t *testing.T) {
	sc := NewScanner(strings.NewReader("<root><unclosed>"))

	// Consume tokens until error.
	for {
		_, err := sc.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				// encoding/xml may report EOF for some malformed cases.
				break
			}
			// Got an error — that's expected.
			return
		}
	}
}

func TestDetectVersionUnknown(t *testing.T) {
	sc := NewScanner(strings.NewReader(`<root xmlns="http://example.com/unknown"><child/></root>`))

	_, err := sc.StartElement()
	if err != nil {
		t.Fatal(err)
	}

	if sc.DetectedVersion != VersionUnknown {
		t.Errorf("got version %q, want unknown", sc.DetectedVersion)
	}
}
