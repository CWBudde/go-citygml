package helpers

import (
	"testing"

	"github.com/cwbudde/go-citygml/types"
)

func TestSummarizeTerrain_Empty(t *testing.T) {
	doc := &types.Document{}
	s := SummarizeTerrain(doc)
	if s.Count != 0 || s.TotalPolygons != 0 || len(s.Polygons) != 0 {
		t.Errorf("expected empty summary: %+v", s)
	}
}

func TestSummarizeTerrain(t *testing.T) {
	doc := &types.Document{
		Terrains: []types.Terrain{
			{ID: "T1", Geometry: types.MultiSurface{Polygons: []types.Polygon{
				{Exterior: types.Ring{Points: []types.Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 0}}}},
				{Exterior: types.Ring{Points: []types.Point{{X: 1, Y: 0}, {X: 2, Y: 0}, {X: 1, Y: 1}, {X: 1, Y: 0}}}},
			}}},
			{ID: "T2", Geometry: types.MultiSurface{Polygons: []types.Polygon{
				{Exterior: types.Ring{Points: []types.Point{{X: 3, Y: 3}, {X: 4, Y: 3}, {X: 3, Y: 4}, {X: 3, Y: 3}}}},
			}}},
		},
	}
	s := SummarizeTerrain(doc)
	if s.Count != 2 {
		t.Errorf("count: got %d, want 2", s.Count)
	}
	if s.TotalPolygons != 3 {
		t.Errorf("total polygons: got %d, want 3", s.TotalPolygons)
	}
	if len(s.Polygons) != 3 {
		t.Errorf("polygons slice: got %d, want 3", len(s.Polygons))
	}
}
