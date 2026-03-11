package decode

import (
	"strings"
	"testing"

	"github.com/cwbudde/go-citygml/internal/xmlscan"
	"github.com/cwbudde/go-citygml/types"
)

func TestIsBuildingElement(t *testing.T) {
	tests := []struct {
		ns    string
		local string
		want  bool
	}{
		{xmlscan.NSCityGML20Bldg, "Building", true},
		{xmlscan.NSCityGML30Bldg, "Building", true},
		{xmlscan.NSCityGML20Dem, "ReliefFeature", false},
		{"http://example.com", "Building", false},
	}
	for _, tt := range tests {
		elem := &xmlscan.Element{}
		elem.Name.Space = tt.ns
		elem.Name.Local = tt.local

		got := IsBuildingElement(elem)
		if got != tt.want {
			t.Errorf("IsBuildingElement(%s:%s) = %v, want %v", tt.ns, tt.local, got, tt.want)
		}
	}
}

const buildingWithAttributes = `<?xml version="1.0" encoding="UTF-8"?>
<CityModel xmlns="http://www.opengis.net/citygml/2.0"
           xmlns:gml="http://www.opengis.net/gml"
           xmlns:bldg="http://www.opengis.net/citygml/building/2.0">
  <cityObjectMember>
    <bldg:Building gml:id="B_ATTR">
      <bldg:class>1000</bldg:class>
      <bldg:function>1010</bldg:function>
      <bldg:usage>residential</bldg:usage>
      <bldg:measuredHeight>12.5</bldg:measuredHeight>
    </bldg:Building>
  </cityObjectMember>
</CityModel>`

func TestBuilding_Attributes(t *testing.T) {
	b := decodeTestBuilding(t, buildingWithAttributes)

	if b.ID != "B_ATTR" {
		t.Errorf("ID = %q, want B_ATTR", b.ID)
	}

	if b.Class != "1000" {
		t.Errorf("Class = %q, want 1000", b.Class)
	}

	if b.Function != "1010" {
		t.Errorf("Function = %q, want 1010", b.Function)
	}

	if b.Usage != "residential" {
		t.Errorf("Usage = %q, want residential", b.Usage)
	}

	if !b.HasMeasuredHeight || b.MeasuredHeight != 12.5 {
		t.Errorf("MeasuredHeight = %g (has=%v), want 12.5", b.MeasuredHeight, b.HasMeasuredHeight)
	}
}

const buildingWithLod1Solid = `<?xml version="1.0" encoding="UTF-8"?>
<CityModel xmlns="http://www.opengis.net/citygml/2.0"
           xmlns:gml="http://www.opengis.net/gml"
           xmlns:bldg="http://www.opengis.net/citygml/building/2.0">
  <cityObjectMember>
    <bldg:Building gml:id="B_SOLID1">
      <bldg:lod1Solid>
        <gml:Solid>
          <gml:exterior>
            <gml:CompositeSurface>
              <gml:surfaceMember>
                <gml:Polygon>
                  <gml:exterior>
                    <gml:LinearRing>
                      <gml:posList>0 0 0 10 0 0 10 10 0 0 10 0 0 0 0</gml:posList>
                    </gml:LinearRing>
                  </gml:exterior>
                </gml:Polygon>
              </gml:surfaceMember>
            </gml:CompositeSurface>
          </gml:exterior>
        </gml:Solid>
      </bldg:lod1Solid>
    </bldg:Building>
  </cityObjectMember>
</CityModel>`

func TestBuilding_Lod1Solid(t *testing.T) {
	b := decodeTestBuilding(t, buildingWithLod1Solid)

	if b.LoD != types.LoD1 {
		t.Errorf("LoD = %q, want 1", b.LoD)
	}

	if b.Solid == nil {
		t.Fatal("expected Solid geometry")
	}

	if len(b.Solid.Exterior.Polygons) != 1 {
		t.Errorf("got %d solid polygons, want 1", len(b.Solid.Exterior.Polygons))
	}

	if b.MultiSurface != nil {
		t.Error("expected nil MultiSurface")
	}
}

const buildingWithLod2Solid = `<?xml version="1.0" encoding="UTF-8"?>
<CityModel xmlns="http://www.opengis.net/citygml/2.0"
           xmlns:gml="http://www.opengis.net/gml"
           xmlns:bldg="http://www.opengis.net/citygml/building/2.0">
  <cityObjectMember>
    <bldg:Building gml:id="B_SOLID2">
      <bldg:lod2Solid>
        <gml:Solid>
          <gml:exterior>
            <gml:CompositeSurface>
              <gml:surfaceMember>
                <gml:Polygon>
                  <gml:exterior>
                    <gml:LinearRing>
                      <gml:posList>0 0 0 10 0 0 10 10 0 0 10 0 0 0 0</gml:posList>
                    </gml:LinearRing>
                  </gml:exterior>
                </gml:Polygon>
              </gml:surfaceMember>
            </gml:CompositeSurface>
          </gml:exterior>
        </gml:Solid>
      </bldg:lod2Solid>
    </bldg:Building>
  </cityObjectMember>
</CityModel>`

func TestBuilding_Lod2Solid(t *testing.T) {
	b := decodeTestBuilding(t, buildingWithLod2Solid)

	if b.LoD != types.LoD2 {
		t.Errorf("LoD = %q, want 2", b.LoD)
	}

	if b.Solid == nil {
		t.Fatal("expected Solid geometry")
	}
}

const buildingWithLod1MultiSurface = `<?xml version="1.0" encoding="UTF-8"?>
<CityModel xmlns="http://www.opengis.net/citygml/2.0"
           xmlns:gml="http://www.opengis.net/gml"
           xmlns:bldg="http://www.opengis.net/citygml/building/2.0">
  <cityObjectMember>
    <bldg:Building gml:id="B_MS1">
      <bldg:lod1MultiSurface>
        <gml:MultiSurface>
          <gml:surfaceMember>
            <gml:Polygon>
              <gml:exterior>
                <gml:LinearRing>
                  <gml:posList>0 0 100 20 0 100 20 15 100 0 15 100 0 0 100</gml:posList>
                </gml:LinearRing>
              </gml:exterior>
            </gml:Polygon>
          </gml:surfaceMember>
        </gml:MultiSurface>
      </bldg:lod1MultiSurface>
    </bldg:Building>
  </cityObjectMember>
</CityModel>`

func TestBuilding_Lod1MultiSurface(t *testing.T) {
	b := decodeTestBuilding(t, buildingWithLod1MultiSurface)

	if b.LoD != types.LoD1 {
		t.Errorf("LoD = %q, want 1", b.LoD)
	}

	if b.MultiSurface == nil {
		t.Fatal("expected MultiSurface geometry")
	}

	if len(b.MultiSurface.Polygons) != 1 {
		t.Errorf("got %d polygons, want 1", len(b.MultiSurface.Polygons))
	}

	if b.Solid != nil {
		t.Error("expected nil Solid")
	}
}

const buildingWithLod2MultiSurface = `<?xml version="1.0" encoding="UTF-8"?>
<CityModel xmlns="http://www.opengis.net/citygml/2.0"
           xmlns:gml="http://www.opengis.net/gml"
           xmlns:bldg="http://www.opengis.net/citygml/building/2.0">
  <cityObjectMember>
    <bldg:Building gml:id="B_MS2">
      <bldg:lod2MultiSurface>
        <gml:MultiSurface>
          <gml:surfaceMember>
            <gml:Polygon>
              <gml:exterior>
                <gml:LinearRing>
                  <gml:posList>0 0 0 10 0 0 10 10 0 0 10 0 0 0 0</gml:posList>
                </gml:LinearRing>
              </gml:exterior>
            </gml:Polygon>
          </gml:surfaceMember>
        </gml:MultiSurface>
      </bldg:lod2MultiSurface>
    </bldg:Building>
  </cityObjectMember>
</CityModel>`

func TestBuilding_Lod2MultiSurface(t *testing.T) {
	b := decodeTestBuilding(t, buildingWithLod2MultiSurface)

	if b.LoD != types.LoD2 {
		t.Errorf("LoD = %q, want 2", b.LoD)
	}

	if b.MultiSurface == nil {
		t.Fatal("expected MultiSurface geometry")
	}
}

const buildingWithBoundedSurfaces = `<?xml version="1.0" encoding="UTF-8"?>
<CityModel xmlns="http://www.opengis.net/citygml/2.0"
           xmlns:gml="http://www.opengis.net/gml"
           xmlns:bldg="http://www.opengis.net/citygml/building/2.0">
  <cityObjectMember>
    <bldg:Building gml:id="B_BOUNDED">
      <bldg:boundedBy>
        <bldg:GroundSurface gml:id="GS_1">
          <bldg:lod2MultiSurface>
            <gml:MultiSurface>
              <gml:surfaceMember>
                <gml:Polygon>
                  <gml:exterior>
                    <gml:LinearRing>
                      <gml:posList>0 0 0 10 0 0 10 10 0 0 10 0 0 0 0</gml:posList>
                    </gml:LinearRing>
                  </gml:exterior>
                </gml:Polygon>
              </gml:surfaceMember>
            </gml:MultiSurface>
          </bldg:lod2MultiSurface>
        </bldg:GroundSurface>
      </bldg:boundedBy>
      <bldg:boundedBy>
        <bldg:RoofSurface gml:id="RS_1">
          <bldg:lod2MultiSurface>
            <gml:MultiSurface>
              <gml:surfaceMember>
                <gml:Polygon>
                  <gml:exterior>
                    <gml:LinearRing>
                      <gml:posList>0 0 15 10 0 15 10 10 15 0 10 15 0 0 15</gml:posList>
                    </gml:LinearRing>
                  </gml:exterior>
                </gml:Polygon>
              </gml:surfaceMember>
            </gml:MultiSurface>
          </bldg:lod2MultiSurface>
        </bldg:RoofSurface>
      </bldg:boundedBy>
      <bldg:boundedBy>
        <bldg:WallSurface gml:id="WS_1">
          <bldg:lod2MultiSurface>
            <gml:MultiSurface>
              <gml:surfaceMember>
                <gml:Polygon>
                  <gml:exterior>
                    <gml:LinearRing>
                      <gml:posList>0 0 0 10 0 0 10 0 15 0 0 15 0 0 0</gml:posList>
                    </gml:LinearRing>
                  </gml:exterior>
                </gml:Polygon>
              </gml:surfaceMember>
            </gml:MultiSurface>
          </bldg:lod2MultiSurface>
        </bldg:WallSurface>
      </bldg:boundedBy>
    </bldg:Building>
  </cityObjectMember>
</CityModel>`

func TestBuilding_BoundedSurfaces(t *testing.T) {
	b := decodeTestBuilding(t, buildingWithBoundedSurfaces)

	if len(b.BoundedBy) != 3 {
		t.Fatalf("got %d bounded surfaces, want 3", len(b.BoundedBy))
	}

	want := []struct {
		id, typ string
	}{
		{"GS_1", "GroundSurface"},
		{"RS_1", "RoofSurface"},
		{"WS_1", "WallSurface"},
	}
	for i, w := range want {
		if b.BoundedBy[i].ID != w.id {
			t.Errorf("surface[%d].ID = %q, want %q", i, b.BoundedBy[i].ID, w.id)
		}

		if b.BoundedBy[i].Type != w.typ {
			t.Errorf("surface[%d].Type = %q, want %q", i, b.BoundedBy[i].Type, w.typ)
		}

		if len(b.BoundedBy[i].Geometry.Polygons) != 1 {
			t.Errorf("surface[%d] polygons = %d, want 1", i, len(b.BoundedBy[i].Geometry.Polygons))
		}
	}
}

const buildingMinimal = `<?xml version="1.0" encoding="UTF-8"?>
<CityModel xmlns="http://www.opengis.net/citygml/2.0"
           xmlns:gml="http://www.opengis.net/gml"
           xmlns:bldg="http://www.opengis.net/citygml/building/2.0">
  <cityObjectMember>
    <bldg:Building gml:id="B_MINIMAL"/>
  </cityObjectMember>
</CityModel>`

func TestBuilding_Minimal(t *testing.T) {
	b := decodeTestBuilding(t, buildingMinimal)

	if b.ID != "B_MINIMAL" {
		t.Errorf("ID = %q, want B_MINIMAL", b.ID)
	}

	if b.HasMeasuredHeight {
		t.Error("should not have measured height")
	}

	if b.Solid != nil {
		t.Error("expected nil Solid")
	}

	if b.MultiSurface != nil {
		t.Error("expected nil MultiSurface")
	}

	if len(b.BoundedBy) != 0 {
		t.Errorf("expected no bounded surfaces, got %d", len(b.BoundedBy))
	}
}

const buildingWithUnknownChildren = `<?xml version="1.0" encoding="UTF-8"?>
<CityModel xmlns="http://www.opengis.net/citygml/2.0"
           xmlns:gml="http://www.opengis.net/gml"
           xmlns:bldg="http://www.opengis.net/citygml/building/2.0">
  <cityObjectMember>
    <bldg:Building gml:id="B_UNK">
      <bldg:measuredHeight>5.0</bldg:measuredHeight>
      <bldg:yearOfConstruction>2020</bldg:yearOfConstruction>
      <bldg:storeysAboveGround>3</bldg:storeysAboveGround>
      <bldg:address><bldg:Street>Main St</bldg:Street></bldg:address>
    </bldg:Building>
  </cityObjectMember>
</CityModel>`

func TestBuilding_SkipsUnknownChildren(t *testing.T) {
	b := decodeTestBuilding(t, buildingWithUnknownChildren)

	if b.ID != "B_UNK" {
		t.Errorf("ID = %q, want B_UNK", b.ID)
	}

	if !b.HasMeasuredHeight || b.MeasuredHeight != 5.0 {
		t.Errorf("MeasuredHeight = %g (has=%v), want 5.0", b.MeasuredHeight, b.HasMeasuredHeight)
	}
}

const buildingCityGML30 = `<?xml version="1.0" encoding="UTF-8"?>
<CityModel xmlns="http://www.opengis.net/citygml/3.0"
           xmlns:gml="http://www.opengis.net/gml/3.2"
           xmlns:bldg="http://www.opengis.net/citygml/building/3.0">
  <cityObjectMember>
    <bldg:Building gml:id="B30">
      <bldg:class>office</bldg:class>
      <bldg:measuredHeight>20.0</bldg:measuredHeight>
    </bldg:Building>
  </cityObjectMember>
</CityModel>`

func TestBuilding_CityGML30(t *testing.T) {
	b := decodeTestBuilding(t, buildingCityGML30)

	if b.ID != "B30" {
		t.Errorf("ID = %q, want B30", b.ID)
	}

	if b.Class != "office" {
		t.Errorf("Class = %q, want office", b.Class)
	}

	if b.MeasuredHeight != 20.0 {
		t.Errorf("MeasuredHeight = %g, want 20.0", b.MeasuredHeight)
	}
}

// decodeTestBuilding is a helper that parses a full CityGML document and returns the first building.
func decodeTestBuilding(t *testing.T, input string) types.Building {
	t.Helper()

	sc := xmlscan.NewScanner(strings.NewReader(input))

	hdr, err := xmlscan.ReadHeader(sc)
	if err != nil {
		t.Fatal(err)
	}

	var result *types.Building

	err = xmlscan.EachCityObjectMember(sc, hdr, func(elem *xmlscan.Element, sc *xmlscan.Scanner) error {
		if !IsBuildingElement(elem) {
			return sc.Skip()
		}

		b, err := Building(elem, sc)
		if err != nil {
			return err
		}

		if result == nil {
			result = &b
		}

		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Fatal("no building found in input")
	}

	return *result
}
