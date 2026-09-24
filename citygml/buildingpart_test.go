package citygml

import (
	"strings"
	"testing"

	"github.com/cwbudde/go-citygml/types"
)

const citygml10PartsFixture = "../testdata/citygml10_building_parts.gml"

func TestRead_CityGML10BuildingParts(t *testing.T) {
	doc, err := ReadFile(citygml10PartsFixture, Options{})
	if err != nil {
		t.Fatal(err)
	}

	if doc.Version != "1.0" {
		t.Errorf("Version = %q, want 1.0", doc.Version)
	}

	if doc.CRS.Code != 25832 {
		t.Errorf("CRS.Code = %d, want 25832", doc.CRS.Code)
	}

	if len(doc.GenericObjects) != 0 {
		t.Errorf("got %d generic objects, want 0", len(doc.GenericObjects))
	}

	if len(doc.Buildings) != 2 {
		t.Fatalf("got %d buildings, want 2", len(doc.Buildings))
	}

	parent := doc.Buildings[0]
	if parent.ID != "BLDG_PARTS" || parent.Function != "31001_1000" {
		t.Errorf("parent = %q/%q, want BLDG_PARTS/31001_1000", parent.ID, parent.Function)
	}

	if parent.HasMeasuredHeight || parent.Solid != nil || len(parent.BoundedBy) != 0 || parent.Footprint != nil {
		t.Error("parent should carry no height, geometry or footprint of its own")
	}

	if len(parent.Parts) != 2 {
		t.Fatalf("got %d parts, want 2", len(parent.Parts))
	}

	checkPart(t, parent.Parts[0], wantPart{"PART_A", "31001_1000", 12, 550000, 550010})
	checkPart(t, parent.Parts[1], wantPart{"PART_B", "31001_2000", 6.5, 550010, 550016})

	plain := doc.Buildings[1]
	if plain.ID != "BLDG_PLAIN" || len(plain.Parts) != 0 {
		t.Errorf("plain = %q with %d parts, want BLDG_PLAIN with none", plain.ID, len(plain.Parts))
	}

	if !plain.HasMeasuredHeight || plain.MeasuredHeight != 8 {
		t.Errorf("plain MeasuredHeight = %g, want 8", plain.MeasuredHeight)
	}

	// The GroundSurface has a 2x2 polygon first and a 10x10 polygon with a
	// 2x2 hole second; the footprint must be the larger one, hole included.
	if plain.Footprint == nil {
		t.Fatal("plain building has no footprint")
	}

	if got := plain.Footprint.Exterior.Points[0]; got.X != 550040 || got.Y != 5803040 {
		t.Errorf("plain footprint starts at (%g, %g), want the larger polygon at (550040, 5803040)", got.X, got.Y)
	}

	if len(plain.Footprint.Interior) != 1 {
		t.Errorf("plain footprint has %d holes, want 1", len(plain.Footprint.Interior))
	}
}

type wantPart struct {
	id       string
	function string
	height   float64
	minX     float64
	maxX     float64
}

func checkPart(t *testing.T, p types.Building, want wantPart) {
	t.Helper()

	if p.ID != want.id || p.Function != want.function {
		t.Errorf("part = %q/%q, want %q/%q", p.ID, p.Function, want.id, want.function)
	}

	if !p.HasMeasuredHeight || p.MeasuredHeight != want.height {
		t.Errorf("part %s MeasuredHeight = %g (has=%v), want %g", p.ID, p.MeasuredHeight, p.HasMeasuredHeight, want.height)
	}

	if p.LoD != types.LoD2 {
		t.Errorf("part %s LoD = %q, want 2", p.ID, p.LoD)
	}

	// The lod2Solid holds only xlink:href members, which are not resolved.
	if p.Solid == nil || len(p.Solid.Exterior.Polygons) != 0 {
		t.Errorf("part %s: want an empty Solid from xlink-only members, got %+v", p.ID, p.Solid)
	}

	if len(p.BoundedBy) != 3 {
		t.Fatalf("part %s: got %d bounded surfaces, want 3", p.ID, len(p.BoundedBy))
	}

	if p.Footprint == nil || len(p.Footprint.Exterior.Points) != 5 {
		t.Fatalf("part %s: want a 5-point footprint, got %+v", p.ID, p.Footprint)
	}

	minX, maxX := p.Footprint.Exterior.Points[0].X, p.Footprint.Exterior.Points[0].X
	for _, pt := range p.Footprint.Exterior.Points {
		minX = min(minX, pt.X)
		maxX = max(maxX, pt.X)

		if pt.Z != 0 {
			t.Errorf("part %s footprint point has Z=%g, want 0", p.ID, pt.Z)
		}
	}

	if minX != want.minX || maxX != want.maxX {
		t.Errorf("part %s footprint X range = [%g, %g], want [%g, %g]", p.ID, minX, maxX, want.minX, want.maxX)
	}
}

func TestRead_BuildingPartsDerivedHeight(t *testing.T) {
	// A CityGML 3.0 bldg:buildingPart without measuredHeight: the part's height
	// must be derived from its own geometry, and an xlink-only part reference
	// must be ignored.
	input := `<?xml version="1.0" encoding="UTF-8"?>
<CityModel xmlns="http://www.opengis.net/citygml/3.0"
           xmlns:gml="http://www.opengis.net/gml/3.2"
           xmlns:xlink="http://www.w3.org/1999/xlink"
           xmlns:bldg="http://www.opengis.net/citygml/building/3.0">
  <cityObjectMember>
    <bldg:Building gml:id="B">
      <bldg:buildingPart xlink:href="#ELSEWHERE"/>
      <bldg:buildingPart>
        <bldg:BuildingPart gml:id="P">
          <bldg:buildingPart>
            <bldg:BuildingPart gml:id="PP">
              <bldg:measuredHeight>3</bldg:measuredHeight>
            </bldg:BuildingPart>
          </bldg:buildingPart>
          <bldg:lod1MultiSurface>
            <gml:MultiSurface>
              <gml:surfaceMember>
                <gml:Polygon>
                  <gml:exterior>
                    <gml:LinearRing>
                      <gml:posList srsDimension="3">0 0 10 4 0 10 4 4 10 0 4 10 0 0 10</gml:posList>
                    </gml:LinearRing>
                  </gml:exterior>
                </gml:Polygon>
              </gml:surfaceMember>
              <gml:surfaceMember>
                <gml:Polygon>
                  <gml:exterior>
                    <gml:LinearRing>
                      <gml:posList srsDimension="3">0 0 17 4 0 17 4 4 17 0 4 17 0 0 17</gml:posList>
                    </gml:LinearRing>
                  </gml:exterior>
                </gml:Polygon>
              </gml:surfaceMember>
            </gml:MultiSurface>
          </bldg:lod1MultiSurface>
        </bldg:BuildingPart>
      </bldg:buildingPart>
    </bldg:Building>
  </cityObjectMember>
</CityModel>`

	doc, err := Read(strings.NewReader(input), Options{})
	if err != nil {
		t.Fatal(err)
	}

	b := doc.Buildings[0]
	if len(b.Parts) != 1 {
		t.Fatalf("got %d parts, want 1 (the xlink-only reference is skipped)", len(b.Parts))
	}

	p := b.Parts[0]
	if p.ID != "P" || p.DerivedHeight != 7 {
		t.Errorf("part = %q with DerivedHeight %g, want P with 7", p.ID, p.DerivedHeight)
	}

	if p.Footprint == nil {
		t.Error("part footprint was not derived")
	}

	if len(p.Parts) != 1 || p.Parts[0].ID != "PP" || p.Parts[0].MeasuredHeight != 3 {
		t.Errorf("nested part = %+v, want PP with measuredHeight 3", p.Parts)
	}
}

func TestValidate_BuildingParts(t *testing.T) {
	doc, err := ReadFile(citygml10PartsFixture, Options{})
	if err != nil {
		t.Fatal(err)
	}

	for _, f := range Validate(doc) {
		t.Errorf("unexpected finding: %s", f)
	}

	// A part without height or geometry is reported under its parent's path.
	doc.Buildings[0].Parts[1] = types.Building{ID: "EMPTY_PART"}

	findings := Validate(doc)

	paths := make([]string, 0, len(findings))
	for _, f := range findings {
		paths = append(paths, f.Path)
	}

	want := "Building[0](BLDG_PARTS)/Part[1](EMPTY_PART)"
	if len(paths) != 2 || paths[0] != want || paths[1] != want {
		t.Errorf("finding paths = %q, want two findings at %q", paths, want)
	}
}
