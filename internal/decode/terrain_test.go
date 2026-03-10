package decode

import (
	"strings"
	"testing"

	"github.com/cwbudde/go-citygml/internal/xmlscan"
)

const tinRelief = `<?xml version="1.0" encoding="UTF-8"?>
<CityModel xmlns="http://www.opengis.net/citygml/2.0"
           xmlns:gml="http://www.opengis.net/gml"
           xmlns:dem="http://www.opengis.net/citygml/relief/2.0">
  <cityObjectMember>
    <dem:ReliefFeature gml:id="RF_1">
      <dem:reliefComponent>
        <dem:TINRelief gml:id="TIN_1">
          <dem:tin>
            <gml:MultiSurface>
              <gml:surfaceMember>
                <gml:Polygon>
                  <gml:exterior>
                    <gml:LinearRing>
                      <gml:posList>0 0 0 10 0 0 5 10 5 0 0 0</gml:posList>
                    </gml:LinearRing>
                  </gml:exterior>
                </gml:Polygon>
              </gml:surfaceMember>
            </gml:MultiSurface>
          </dem:tin>
        </dem:TINRelief>
      </dem:reliefComponent>
    </dem:ReliefFeature>
  </cityObjectMember>
</CityModel>`

func TestTerrain_ReliefFeature(t *testing.T) {
	sc := xmlscan.NewScanner(strings.NewReader(tinRelief))

	hdr, err := xmlscan.ReadHeader(sc)
	if err != nil {
		t.Fatal(err)
	}

	var terrainCount int

	err = xmlscan.EachCityObjectMember(sc, hdr, func(elem *xmlscan.Element, sc *xmlscan.Scanner) error {
		if !IsTerrainElement(elem) {
			t.Fatalf("expected terrain element, got %s", elem.LocalName())
		}

		terrains, err := Terrain(elem, sc)
		if err != nil {
			return err
		}

		terrainCount += len(terrains)

		if len(terrains) != 1 {
			t.Fatalf("got %d terrains, want 1", len(terrains))
		}

		tr := terrains[0]
		if tr.ID != "TIN_1" {
			t.Errorf("ID = %q, want TIN_1", tr.ID)
		}

		if len(tr.Geometry.Polygons) != 1 {
			t.Errorf("got %d polygons, want 1", len(tr.Geometry.Polygons))
		}

		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if terrainCount != 1 {
		t.Errorf("total terrains = %d, want 1", terrainCount)
	}
}

func TestIsTerrainElement(t *testing.T) {
	tests := []struct {
		ns    string
		local string
		want  bool
	}{
		{xmlscan.NSCityGML20Dem, "ReliefFeature", true},
		{xmlscan.NSCityGML20Dem, "TINRelief", true},
		{xmlscan.NSCityGML30Dem, "ReliefFeature", true},
		{xmlscan.NSCityGML20Bldg, "Building", false},
		{xmlscan.NSCityGML20Dem, "Unknown", false},
	}
	for _, tt := range tests {
		elem := &xmlscan.Element{}
		elem.Name.Space = tt.ns
		elem.Name.Local = tt.local

		got := IsTerrainElement(elem)
		if got != tt.want {
			t.Errorf("IsTerrainElement(%s:%s) = %v, want %v", tt.ns, tt.local, got, tt.want)
		}
	}
}
