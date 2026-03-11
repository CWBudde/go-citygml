package geojson

import (
	"encoding/json"
	"testing"

	"github.com/cwbudde/go-citygml/types"
)

func TestPolygonGeometry(t *testing.T) {
	poly := types.Polygon{
		Exterior: types.Ring{Points: []types.Point{
			{X: 0, Y: 0}, {X: 10, Y: 0}, {X: 10, Y: 10}, {X: 0, Y: 10}, {X: 0, Y: 0},
		}},
	}

	g := PolygonGeometry(poly)
	if g.Type != "Polygon" {
		t.Errorf("type: got %q, want Polygon", g.Type)
	}

	data, err := json.Marshal(g)
	if err != nil {
		t.Fatal(err)
	}
	// Ensure it's valid JSON.
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}

	if m["type"] != "Polygon" {
		t.Errorf("serialized type: %v", m["type"])
	}
}

func TestMultiPolygonGeometry(t *testing.T) {
	ms := types.MultiSurface{Polygons: []types.Polygon{
		{Exterior: types.Ring{Points: []types.Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 0}}}},
		{Exterior: types.Ring{Points: []types.Point{{X: 2, Y: 2}, {X: 3, Y: 2}, {X: 2, Y: 3}, {X: 2, Y: 2}}}},
	}}

	g := MultiPolygonGeometry(ms)
	if g.Type != "MultiPolygon" {
		t.Errorf("type: got %q, want MultiPolygon", g.Type)
	}
}

func TestBuildingFeature_WithFootprint(t *testing.T) {
	b := &types.Building{
		ID:                "B1",
		Class:             "residential",
		HasMeasuredHeight: true,
		MeasuredHeight:    12.5,
		LoD:               types.LoD1,
		Footprint: &types.Polygon{
			Exterior: types.Ring{Points: []types.Point{
				{X: 0, Y: 0}, {X: 10, Y: 0}, {X: 10, Y: 10}, {X: 0, Y: 10}, {X: 0, Y: 0},
			}},
		},
	}

	f := BuildingFeature(b)
	if f.ID != "B1" {
		t.Errorf("id: %q", f.ID)
	}

	if f.Type != "Feature" {
		t.Errorf("type: %q", f.Type)
	}

	if f.Geometry == nil {
		t.Fatal("expected geometry")
	}

	if f.Geometry.Type != "Polygon" {
		t.Errorf("geometry type: %q", f.Geometry.Type)
	}

	if f.Properties["class"] != "residential" {
		t.Errorf("class: %v", f.Properties["class"])
	}

	if f.Properties["measuredHeight"] != 12.5 {
		t.Errorf("measuredHeight: %v", f.Properties["measuredHeight"])
	}
}

func TestBuildingFeature_NoGeometry(t *testing.T) {
	b := &types.Building{ID: "B2"}

	f := BuildingFeature(b)
	if f.Geometry != nil {
		t.Error("expected nil geometry")
	}
}

func TestTerrainFeature(t *testing.T) {
	tr := &types.Terrain{
		ID: "T1",
		Geometry: types.MultiSurface{Polygons: []types.Polygon{
			{Exterior: types.Ring{Points: []types.Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 0}}}},
		}},
	}

	f := TerrainFeature(tr)
	if f.ID != "T1" {
		t.Errorf("id: %q", f.ID)
	}

	if f.Geometry == nil || f.Geometry.Type != "MultiPolygon" {
		t.Error("expected MultiPolygon geometry")
	}

	if f.Properties["type"] != "Terrain" {
		t.Errorf("type prop: %v", f.Properties["type"])
	}
}

func TestFromDocument(t *testing.T) {
	doc := &types.Document{
		Buildings: []types.Building{
			{ID: "B1", Footprint: &types.Polygon{
				Exterior: types.Ring{Points: []types.Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 0}}},
			}},
		},
		Terrains: []types.Terrain{
			{ID: "T1", Geometry: types.MultiSurface{Polygons: []types.Polygon{
				{Exterior: types.Ring{Points: []types.Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 0}}}},
			}}},
		},
	}

	fc := FromDocument(doc)
	if fc.Type != "FeatureCollection" {
		t.Errorf("type: %q", fc.Type)
	}

	if len(fc.Features) != 2 {
		t.Fatalf("features: got %d, want 2", len(fc.Features))
	}

	// Verify full JSON round-trip.
	data, err := json.Marshal(fc)
	if err != nil {
		t.Fatal(err)
	}

	var decoded FeatureCollection
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}

	if len(decoded.Features) != 2 {
		t.Errorf("decoded features: %d", len(decoded.Features))
	}
}

func TestFromDocument_Empty(t *testing.T) {
	doc := &types.Document{}

	fc := FromDocument(doc)
	if len(fc.Features) != 0 {
		t.Errorf("expected 0 features, got %d", len(fc.Features))
	}
}

func TestPolygonGeometry_WithHoles(t *testing.T) {
	poly := types.Polygon{
		Exterior: types.Ring{Points: []types.Point{
			{X: 0, Y: 0}, {X: 20, Y: 0}, {X: 20, Y: 20}, {X: 0, Y: 20}, {X: 0, Y: 0},
		}},
		Interior: []types.Ring{{Points: []types.Point{
			{X: 5, Y: 5}, {X: 15, Y: 5}, {X: 15, Y: 15}, {X: 5, Y: 15}, {X: 5, Y: 5},
		}}},
	}
	g := PolygonGeometry(poly)

	// Verify coordinates include two rings.
	var rings [][][2]float64

	err := json.Unmarshal(g.Coordinates, &rings)
	if err != nil {
		t.Fatal(err)
	}

	if len(rings) != 2 {
		t.Errorf("got %d rings, want 2", len(rings))
	}
}
