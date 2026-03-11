package citygml

import (
	"errors"
	"strings"
	"testing"

	"github.com/cwbudde/go-citygml/types"
)

const buildingWithLod1Solid = `<?xml version="1.0" encoding="UTF-8"?>
<CityModel xmlns="http://www.opengis.net/citygml/2.0"
           xmlns:gml="http://www.opengis.net/gml"
           xmlns:bldg="http://www.opengis.net/citygml/building/2.0">
  <cityObjectMember>
    <bldg:Building gml:id="BLDG_001">
      <bldg:class>1000</bldg:class>
      <bldg:function>1010</bldg:function>
      <bldg:usage>residential</bldg:usage>
      <bldg:measuredHeight>12.5</bldg:measuredHeight>
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
              <gml:surfaceMember>
                <gml:Polygon>
                  <gml:exterior>
                    <gml:LinearRing>
                      <gml:posList>0 0 12.5 10 0 12.5 10 10 12.5 0 10 12.5 0 0 12.5</gml:posList>
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

func TestRead_BuildingLod1Solid(t *testing.T) {
	doc, err := Read(strings.NewReader(buildingWithLod1Solid), Options{})
	if err != nil {
		t.Fatal(err)
	}

	if doc.Version != "2.0" {
		t.Errorf("version = %q, want 2.0", doc.Version)
	}

	if len(doc.Buildings) != 1 {
		t.Fatalf("got %d buildings, want 1", len(doc.Buildings))
	}

	b := doc.Buildings[0]
	if b.ID != "BLDG_001" {
		t.Errorf("ID = %q, want BLDG_001", b.ID)
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

	if b.LoD != types.LoD1 {
		t.Errorf("LoD = %q, want 1", b.LoD)
	}

	if b.Solid == nil {
		t.Fatal("expected Solid geometry")
	}

	if len(b.Solid.Exterior.Polygons) != 2 {
		t.Errorf("got %d solid polygons, want 2", len(b.Solid.Exterior.Polygons))
	}
	// Footprint should be derived.
	if b.Footprint == nil {
		t.Error("expected derived footprint")
	}
	// DerivedHeight should be 0 since MeasuredHeight is present.
	if b.DerivedHeight != 0 {
		t.Errorf("DerivedHeight = %g, want 0 (measuredHeight is present)", b.DerivedHeight)
	}
}

const buildingWithLod1MultiSurface = `<?xml version="1.0" encoding="UTF-8"?>
<CityModel xmlns="http://www.opengis.net/citygml/2.0"
           xmlns:gml="http://www.opengis.net/gml"
           xmlns:bldg="http://www.opengis.net/citygml/building/2.0">
  <cityObjectMember>
    <bldg:Building gml:id="BLDG_002">
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
          <gml:surfaceMember>
            <gml:Polygon>
              <gml:exterior>
                <gml:LinearRing>
                  <gml:posList>0 0 108 20 0 108 20 15 108 0 15 108 0 0 108</gml:posList>
                </gml:LinearRing>
              </gml:exterior>
            </gml:Polygon>
          </gml:surfaceMember>
        </gml:MultiSurface>
      </bldg:lod1MultiSurface>
    </bldg:Building>
  </cityObjectMember>
</CityModel>`

func TestRead_BuildingLod1MultiSurface_DerivedHeight(t *testing.T) {
	doc, err := Read(strings.NewReader(buildingWithLod1MultiSurface), Options{})
	if err != nil {
		t.Fatal(err)
	}

	if len(doc.Buildings) != 1 {
		t.Fatalf("got %d buildings, want 1", len(doc.Buildings))
	}

	b := doc.Buildings[0]
	if b.ID != "BLDG_002" {
		t.Errorf("ID = %q", b.ID)
	}

	if b.HasMeasuredHeight {
		t.Error("should not have measured height")
	}

	if b.DerivedHeight != 8 {
		t.Errorf("DerivedHeight = %g, want 8", b.DerivedHeight)
	}

	if b.MultiSurface == nil {
		t.Fatal("expected MultiSurface geometry")
	}

	if len(b.MultiSurface.Polygons) != 2 {
		t.Errorf("got %d polygons, want 2", len(b.MultiSurface.Polygons))
	}

	if b.Footprint == nil {
		t.Error("expected derived footprint")
	}
}

const twoBuildings = `<?xml version="1.0" encoding="UTF-8"?>
<CityModel xmlns="http://www.opengis.net/citygml/2.0"
           xmlns:gml="http://www.opengis.net/gml"
           xmlns:bldg="http://www.opengis.net/citygml/building/2.0">
  <cityObjectMember>
    <bldg:Building gml:id="A">
      <bldg:measuredHeight>5.0</bldg:measuredHeight>
    </bldg:Building>
  </cityObjectMember>
  <cityObjectMember>
    <bldg:Building gml:id="B">
      <bldg:measuredHeight>10.0</bldg:measuredHeight>
    </bldg:Building>
  </cityObjectMember>
</CityModel>`

func TestRead_MultipleBuildings(t *testing.T) {
	doc, err := Read(strings.NewReader(twoBuildings), Options{})
	if err != nil {
		t.Fatal(err)
	}

	if len(doc.Buildings) != 2 {
		t.Fatalf("got %d buildings, want 2", len(doc.Buildings))
	}

	if doc.Buildings[0].ID != "A" || doc.Buildings[1].ID != "B" {
		t.Errorf("IDs = [%s, %s], want [A, B]", doc.Buildings[0].ID, doc.Buildings[1].ID)
	}

	if doc.Buildings[0].MeasuredHeight != 5.0 || doc.Buildings[1].MeasuredHeight != 10.0 {
		t.Error("measured heights mismatch")
	}
}

const buildingWithBoundedSurfaces = `<?xml version="1.0" encoding="UTF-8"?>
<CityModel xmlns="http://www.opengis.net/citygml/2.0"
           xmlns:gml="http://www.opengis.net/gml"
           xmlns:bldg="http://www.opengis.net/citygml/building/2.0">
  <cityObjectMember>
    <bldg:Building gml:id="BLDG_003">
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
    </bldg:Building>
  </cityObjectMember>
</CityModel>`

func TestRead_BoundedSurfaces(t *testing.T) {
	doc, err := Read(strings.NewReader(buildingWithBoundedSurfaces), Options{})
	if err != nil {
		t.Fatal(err)
	}

	if len(doc.Buildings) != 1 {
		t.Fatalf("got %d buildings, want 1", len(doc.Buildings))
	}

	b := doc.Buildings[0]
	if len(b.BoundedBy) != 2 {
		t.Fatalf("got %d bounded surfaces, want 2", len(b.BoundedBy))
	}

	if b.BoundedBy[0].Type != "GroundSurface" {
		t.Errorf("surface[0].Type = %q, want GroundSurface", b.BoundedBy[0].Type)
	}

	if b.BoundedBy[0].ID != "GS_1" {
		t.Errorf("surface[0].ID = %q, want GS_1", b.BoundedBy[0].ID)
	}

	if b.BoundedBy[1].Type != "RoofSurface" {
		t.Errorf("surface[1].Type = %q, want RoofSurface", b.BoundedBy[1].Type)
	}

	// Height derived from bounded surfaces: 15 - 0 = 15.
	if b.DerivedHeight != 15 {
		t.Errorf("DerivedHeight = %g, want 15", b.DerivedHeight)
	}

	// Footprint from GroundSurface.
	if b.Footprint == nil {
		t.Fatal("expected footprint from GroundSurface")
	}

	for _, pt := range b.Footprint.Exterior.Points {
		if pt.Z != 0 {
			t.Errorf("footprint point Z = %g, want 0", pt.Z)
		}
	}
}

const mixedObjects = `<?xml version="1.0" encoding="UTF-8"?>
<CityModel xmlns="http://www.opengis.net/citygml/2.0"
           xmlns:gml="http://www.opengis.net/gml"
           xmlns:bldg="http://www.opengis.net/citygml/building/2.0"
           xmlns:tran="http://www.opengis.net/citygml/transportation/2.0">
  <cityObjectMember>
    <bldg:Building gml:id="B1">
      <bldg:measuredHeight>5.0</bldg:measuredHeight>
    </bldg:Building>
  </cityObjectMember>
  <cityObjectMember>
    <tran:Road gml:id="R1">
      <tran:name>Main Street</tran:name>
    </tran:Road>
  </cityObjectMember>
</CityModel>`

func TestRead_UnsupportedObject_NonStrict(t *testing.T) {
	doc, err := Read(strings.NewReader(mixedObjects), Options{})
	if err != nil {
		t.Fatal(err)
	}

	if len(doc.Buildings) != 1 {
		t.Errorf("got %d buildings, want 1", len(doc.Buildings))
	}

	if len(doc.GenericObjects) != 1 {
		t.Fatalf("got %d generic objects, want 1", len(doc.GenericObjects))
	}

	if doc.GenericObjects[0].Type != "Road" {
		t.Errorf("generic type = %q, want Road", doc.GenericObjects[0].Type)
	}

	if doc.GenericObjects[0].ID != "R1" {
		t.Errorf("generic ID = %q, want R1", doc.GenericObjects[0].ID)
	}
}

func TestRead_UnsupportedObject_Strict(t *testing.T) {
	_, err := Read(strings.NewReader(mixedObjects), Options{Strict: true})
	if err == nil {
		t.Fatal("expected error in strict mode")
	}

	if !errors.Is(err, ErrUnsupportedObject) {
		t.Errorf("expected ErrUnsupportedObject, got %v", err)
	}
}

func TestRead_DeriveHeightsDisabled(t *testing.T) {
	f := false

	doc, err := Read(strings.NewReader(buildingWithLod1MultiSurface), Options{DeriveHeights: &f})
	if err != nil {
		t.Fatal(err)
	}

	if doc.Buildings[0].DerivedHeight != 0 {
		t.Errorf("DerivedHeight = %g, want 0 (derivation disabled)", doc.Buildings[0].DerivedHeight)
	}
}

func TestRead_DeriveFootprintsDisabled(t *testing.T) {
	f := false

	doc, err := Read(strings.NewReader(buildingWithLod1Solid), Options{DeriveFootprints: &f})
	if err != nil {
		t.Fatal(err)
	}

	if doc.Buildings[0].Footprint != nil {
		t.Error("expected nil footprint when derivation is disabled")
	}
}

func TestRead_EmptyDocument(t *testing.T) {
	input := `<CityModel xmlns="http://www.opengis.net/citygml/2.0"></CityModel>`

	doc, err := Read(strings.NewReader(input), Options{})
	if err != nil {
		t.Fatal(err)
	}

	if len(doc.Buildings) != 0 {
		t.Errorf("got %d buildings, want 0", len(doc.Buildings))
	}
}

func TestRead_MalformedXML(t *testing.T) {
	_, err := Read(strings.NewReader("<not valid xml"), Options{})
	if err == nil {
		t.Fatal("expected error for malformed XML")
	}
}

const cityGML30Building = `<?xml version="1.0" encoding="UTF-8"?>
<CityModel xmlns="http://www.opengis.net/citygml/3.0"
           xmlns:gml="http://www.opengis.net/gml/3.2"
           xmlns:bldg="http://www.opengis.net/citygml/building/3.0">
  <cityObjectMember>
    <bldg:Building gml:id="B30">
      <bldg:measuredHeight>20.0</bldg:measuredHeight>
    </bldg:Building>
  </cityObjectMember>
</CityModel>`

func TestRead_CityGML30(t *testing.T) {
	doc, err := Read(strings.NewReader(cityGML30Building), Options{})
	if err != nil {
		t.Fatal(err)
	}

	if doc.Version != "3.0" {
		t.Errorf("version = %q, want 3.0", doc.Version)
	}

	if len(doc.Buildings) != 1 {
		t.Fatalf("got %d buildings, want 1", len(doc.Buildings))
	}

	if doc.Buildings[0].MeasuredHeight != 20.0 {
		t.Errorf("measuredHeight = %g, want 20.0", doc.Buildings[0].MeasuredHeight)
	}
}

const terrainDocument = `<?xml version="1.0" encoding="UTF-8"?>
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

func TestRead_Terrain(t *testing.T) {
	doc, err := Read(strings.NewReader(terrainDocument), Options{})
	if err != nil {
		t.Fatal(err)
	}

	if len(doc.Terrains) != 1 {
		t.Fatalf("got %d terrains, want 1", len(doc.Terrains))
	}

	tr := doc.Terrains[0]
	if tr.ID != "TIN_1" {
		t.Errorf("ID = %q, want TIN_1", tr.ID)
	}

	if len(tr.Geometry.Polygons) != 1 {
		t.Errorf("got %d polygons, want 1", len(tr.Geometry.Polygons))
	}
}

const buildingAndTerrain = `<?xml version="1.0" encoding="UTF-8"?>
<CityModel xmlns="http://www.opengis.net/citygml/2.0"
           xmlns:gml="http://www.opengis.net/gml"
           xmlns:bldg="http://www.opengis.net/citygml/building/2.0"
           xmlns:dem="http://www.opengis.net/citygml/relief/2.0">
  <cityObjectMember>
    <bldg:Building gml:id="B1">
      <bldg:measuredHeight>10.0</bldg:measuredHeight>
    </bldg:Building>
  </cityObjectMember>
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

const bridgeAndTunnel = `<?xml version="1.0" encoding="UTF-8"?>
<CityModel xmlns="http://www.opengis.net/citygml/2.0"
           xmlns:gml="http://www.opengis.net/gml"
           xmlns:brid="http://www.opengis.net/citygml/bridge/2.0"
           xmlns:tun="http://www.opengis.net/citygml/tunnel/2.0"
           xmlns:tran="http://www.opengis.net/citygml/transportation/2.0">
  <cityObjectMember>
    <brid:Bridge gml:id="BR1"><brid:class>1000</brid:class></brid:Bridge>
  </cityObjectMember>
  <cityObjectMember>
    <tun:Tunnel gml:id="TN1"><tun:class>2000</tun:class></tun:Tunnel>
  </cityObjectMember>
  <cityObjectMember>
    <tran:Road gml:id="RD1"><tran:function>highway</tran:function></tran:Road>
  </cityObjectMember>
</CityModel>`

func TestRead_ContextObjects_AsGeneric(t *testing.T) {
	doc, err := Read(strings.NewReader(bridgeAndTunnel), Options{})
	if err != nil {
		t.Fatal(err)
	}

	if len(doc.GenericObjects) != 3 {
		t.Fatalf("got %d generic objects, want 3", len(doc.GenericObjects))
	}

	want := []struct{ id, typ string }{
		{"BR1", "Bridge"},
		{"TN1", "Tunnel"},
		{"RD1", "Road"},
	}
	for i, w := range want {
		if doc.GenericObjects[i].ID != w.id || doc.GenericObjects[i].Type != w.typ {
			t.Errorf("generic[%d] = {%s, %s}, want {%s, %s}",
				i, doc.GenericObjects[i].ID, doc.GenericObjects[i].Type, w.id, w.typ)
		}
	}
}

func TestRead_BuildingAndTerrain(t *testing.T) {
	doc, err := Read(strings.NewReader(buildingAndTerrain), Options{})
	if err != nil {
		t.Fatal(err)
	}

	if len(doc.Buildings) != 1 {
		t.Errorf("got %d buildings, want 1", len(doc.Buildings))
	}

	if len(doc.Terrains) != 1 {
		t.Errorf("got %d terrains, want 1", len(doc.Terrains))
	}
}

const documentWithEnvelopeSRS = `<?xml version="1.0" encoding="UTF-8"?>
<CityModel xmlns="http://www.opengis.net/citygml/2.0"
           xmlns:gml="http://www.opengis.net/gml"
           xmlns:bldg="http://www.opengis.net/citygml/building/2.0">
  <gml:boundedBy>
    <gml:Envelope srsName="EPSG:25832">
      <gml:lowerCorner>360000 5600000</gml:lowerCorner>
      <gml:upperCorner>370000 5610000</gml:upperCorner>
    </gml:Envelope>
  </gml:boundedBy>
  <cityObjectMember>
    <bldg:Building gml:id="B1">
      <bldg:measuredHeight>5.0</bldg:measuredHeight>
    </bldg:Building>
  </cityObjectMember>
</CityModel>`

func TestRead_CRS_FromEnvelope(t *testing.T) {
	doc, err := Read(strings.NewReader(documentWithEnvelopeSRS), Options{})
	if err != nil {
		t.Fatal(err)
	}

	if doc.SRSName != "EPSG:25832" {
		t.Errorf("SRSName = %q, want EPSG:25832", doc.SRSName)
	}

	if doc.CRS.Code != 25832 {
		t.Errorf("CRS.Code = %d, want 25832", doc.CRS.Code)
	}

	if doc.CRS.IsYXOrder {
		t.Error("EPSG:25832 should not be Y,X order")
	}
}

const documentWithURNSRS = `<?xml version="1.0" encoding="UTF-8"?>
<CityModel xmlns="http://www.opengis.net/citygml/2.0"
           xmlns:gml="http://www.opengis.net/gml"
           xmlns:bldg="http://www.opengis.net/citygml/building/2.0">
  <gml:boundedBy>
    <gml:Envelope srsName="urn:ogc:def:crs:EPSG::4326"/>
  </gml:boundedBy>
  <cityObjectMember>
    <bldg:Building gml:id="B1">
      <bldg:measuredHeight>5.0</bldg:measuredHeight>
    </bldg:Building>
  </cityObjectMember>
</CityModel>`

func TestRead_CRS_URN_YXOrder(t *testing.T) {
	doc, err := Read(strings.NewReader(documentWithURNSRS), Options{})
	if err != nil {
		t.Fatal(err)
	}

	if doc.CRS.Code != 4326 {
		t.Errorf("CRS.Code = %d, want 4326", doc.CRS.Code)
	}

	if !doc.CRS.IsYXOrder {
		t.Error("EPSG:4326 should be Y,X order")
	}
}

func TestRead_CRS_MissingStrict(t *testing.T) {
	input := `<CityModel xmlns="http://www.opengis.net/citygml/2.0">
		<cityObjectMember>
			<Building xmlns="http://www.opengis.net/citygml/building/2.0"/>
		</cityObjectMember>
	</CityModel>`

	_, err := Read(strings.NewReader(input), Options{Strict: true})
	if err == nil {
		t.Fatal("expected error for missing CRS in strict mode")
	}

	if !errors.Is(err, ErrInvalidCRS) {
		t.Errorf("expected ErrInvalidCRS, got %v", err)
	}
}

func TestRead_CRS_MissingNonStrict(t *testing.T) {
	input := `<CityModel xmlns="http://www.opengis.net/citygml/2.0">
		<cityObjectMember>
			<Building xmlns="http://www.opengis.net/citygml/building/2.0"/>
		</cityObjectMember>
	</CityModel>`

	doc, err := Read(strings.NewReader(input), Options{})
	if err != nil {
		t.Fatal(err)
	}

	if doc.SRSName != "" {
		t.Errorf("SRSName = %q, want empty", doc.SRSName)
	}

	if doc.CRS.Code != 0 {
		t.Errorf("CRS.Code = %d, want 0", doc.CRS.Code)
	}
}
