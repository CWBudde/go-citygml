package xmlscan

import (
	"strings"
	"testing"
)

const twoBuildings = `<?xml version="1.0" encoding="UTF-8"?>
<CityModel xmlns="http://www.opengis.net/citygml/2.0"
           xmlns:gml="http://www.opengis.net/gml"
           xmlns:bldg="http://www.opengis.net/citygml/building/2.0">
  <gml:boundedBy>
    <gml:Envelope srsName="EPSG:25832"/>
  </gml:boundedBy>
  <cityObjectMember>
    <bldg:Building gml:id="B1">
      <bldg:measuredHeight>10.0</bldg:measuredHeight>
    </bldg:Building>
  </cityObjectMember>
  <cityObjectMember>
    <bldg:Building gml:id="B2">
      <bldg:measuredHeight>20.0</bldg:measuredHeight>
    </bldg:Building>
  </cityObjectMember>
</CityModel>`

func TestReadHeader(t *testing.T) {
	sc := NewScanner(strings.NewReader(twoBuildings))

	hdr, err := ReadHeader(sc)
	if err != nil {
		t.Fatal(err)
	}

	if hdr.Version != Version20 {
		t.Errorf("version = %q, want %q", hdr.Version, Version20)
	}
}

func TestReadHeaderNonCityModel(t *testing.T) {
	sc := NewScanner(strings.NewReader(`<Something xmlns="http://example.com"/>`))

	_, err := ReadHeader(sc)
	if err == nil {
		t.Error("expected error for non-CityModel root")
	}
}

func TestEachCityObjectMember_Order(t *testing.T) {
	sc := NewScanner(strings.NewReader(twoBuildings))

	hdr, err := ReadHeader(sc)
	if err != nil {
		t.Fatal(err)
	}

	var ids []string

	err = EachCityObjectMember(sc, hdr, func(elem *Element, sc *Scanner) error {
		ids = append(ids, elem.ID)
		return sc.Skip()
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(ids) != 2 {
		t.Fatalf("got %d objects, want 2", len(ids))
	}

	if ids[0] != "B1" || ids[1] != "B2" {
		t.Errorf("got IDs %v, want [B1, B2]", ids)
	}
}

func TestEachCityObjectMember_SkipsBoundedBy(t *testing.T) {
	sc := NewScanner(strings.NewReader(twoBuildings))

	hdr, err := ReadHeader(sc)
	if err != nil {
		t.Fatal(err)
	}

	count := 0

	err = EachCityObjectMember(sc, hdr, func(elem *Element, sc *Scanner) error {
		count++

		if elem.LocalName() == "Envelope" || elem.LocalName() == "boundedBy" {
			t.Error("should not yield boundedBy/Envelope as a city object")
		}

		return sc.Skip()
	})
	if err != nil {
		t.Fatal(err)
	}

	if count != 2 {
		t.Errorf("got %d city objects, want 2", count)
	}
}

func TestEachCityObjectMember_Empty(t *testing.T) {
	input := `<CityModel xmlns="http://www.opengis.net/citygml/2.0"></CityModel>`
	sc := NewScanner(strings.NewReader(input))

	hdr, err := ReadHeader(sc)
	if err != nil {
		t.Fatal(err)
	}

	count := 0

	err = EachCityObjectMember(sc, hdr, func(_ *Element, sc *Scanner) error {
		count++
		return sc.Skip()
	})
	if err != nil {
		t.Fatal(err)
	}

	if count != 0 {
		t.Errorf("got %d objects, want 0", count)
	}
}

func TestEachCityObjectMember_ExtractsEnvelopeSRS(t *testing.T) {
	sc := NewScanner(strings.NewReader(twoBuildings))

	hdr, err := ReadHeader(sc)
	if err != nil {
		t.Fatal(err)
	}

	err = EachCityObjectMember(sc, hdr, func(_ *Element, sc *Scanner) error {
		return sc.Skip()
	})
	if err != nil {
		t.Fatal(err)
	}

	if hdr.SRSName != "EPSG:25832" {
		t.Errorf("SRSName = %q, want EPSG:25832", hdr.SRSName)
	}
}
