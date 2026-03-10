package helpers

import (
	"testing"

	"github.com/cwbudde/go-citygml/types"
)

func TestBuildingHeight_Measured(t *testing.T) {
	b := &types.Building{HasMeasuredHeight: true, MeasuredHeight: 15.5}
	h, ok := BuildingHeight(b)
	if !ok || h != 15.5 {
		t.Errorf("got %v, %v; want 15.5, true", h, ok)
	}
}

func TestBuildingHeight_Derived(t *testing.T) {
	b := &types.Building{DerivedHeight: 10.0}
	h, ok := BuildingHeight(b)
	if !ok || h != 10.0 {
		t.Errorf("got %v, %v; want 10.0, true", h, ok)
	}
}

func TestBuildingHeight_None(t *testing.T) {
	b := &types.Building{}
	_, ok := BuildingHeight(b)
	if ok {
		t.Error("expected no height")
	}
}

func TestBuildingHeights(t *testing.T) {
	doc := &types.Document{
		Buildings: []types.Building{
			{ID: "B1", HasMeasuredHeight: true, MeasuredHeight: 10},
			{ID: "B2", DerivedHeight: 5},
			{ID: "B3"},
		},
	}
	results := BuildingHeights(doc)
	if len(results) != 3 {
		t.Fatalf("got %d results, want 3", len(results))
	}
	if !results[0].HasHeight || !results[0].IsMeasured || results[0].Height != 10 {
		t.Errorf("B1: %+v", results[0])
	}
	if !results[1].HasHeight || results[1].IsMeasured || results[1].Height != 5 {
		t.Errorf("B2: %+v", results[1])
	}
	if results[2].HasHeight {
		t.Errorf("B3 should have no height: %+v", results[2])
	}
}

func TestBuildingFootprints(t *testing.T) {
	fp := &types.Polygon{Exterior: types.Ring{Points: []types.Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 1}, {X: 0, Y: 0}}}}
	doc := &types.Document{
		Buildings: []types.Building{
			{ID: "B1", Footprint: fp},
			{ID: "B2"},
		},
	}
	results := BuildingFootprints(doc)
	if len(results) != 2 {
		t.Fatalf("got %d results, want 2", len(results))
	}
	if results[0].Footprint == nil {
		t.Error("B1 should have footprint")
	}
	if results[1].Footprint != nil {
		t.Error("B2 should not have footprint")
	}
}
