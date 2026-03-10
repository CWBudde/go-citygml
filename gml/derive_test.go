package gml

import (
	"testing"

	"github.com/cwbudde/go-citygml/types"
)

func makeBox(minZ, maxZ float64) *types.Solid {
	// Simple box: bottom at minZ, top at maxZ.
	bottom := types.Polygon{
		Exterior: types.Ring{Points: []types.Point{
			{X: 0, Y: 0, Z: minZ},
			{X: 10, Y: 0, Z: minZ},
			{X: 10, Y: 10, Z: minZ},
			{X: 0, Y: 10, Z: minZ},
			{X: 0, Y: 0, Z: minZ},
		}},
	}
	top := types.Polygon{
		Exterior: types.Ring{Points: []types.Point{
			{X: 0, Y: 0, Z: maxZ},
			{X: 10, Y: 0, Z: maxZ},
			{X: 10, Y: 10, Z: maxZ},
			{X: 0, Y: 10, Z: maxZ},
			{X: 0, Y: 0, Z: maxZ},
		}},
	}

	return &types.Solid{
		Exterior: types.MultiSurface{Polygons: []types.Polygon{bottom, top}},
	}
}

func TestDeriveHeight_Solid(t *testing.T) {
	solid := makeBox(100, 112.5)

	h := DeriveHeight(solid, nil, nil)
	if h != 12.5 {
		t.Errorf("height = %g, want 12.5", h)
	}
}

func TestDeriveHeight_MultiSurface(t *testing.T) {
	ms := &types.MultiSurface{
		Polygons: []types.Polygon{{
			Exterior: types.Ring{Points: []types.Point{
				{X: 0, Y: 0, Z: 50},
				{X: 10, Y: 0, Z: 50},
				{X: 10, Y: 10, Z: 70},
				{X: 0, Y: 0, Z: 50},
			}},
		}},
	}

	h := DeriveHeight(nil, ms, nil)
	if h != 20 {
		t.Errorf("height = %g, want 20", h)
	}
}

func TestDeriveHeight_BoundedSurfaces(t *testing.T) {
	bounded := []types.Surface{
		{
			Type: "GroundSurface",
			Geometry: types.MultiSurface{Polygons: []types.Polygon{{
				Exterior: types.Ring{Points: []types.Point{
					{X: 0, Y: 0, Z: 0},
					{X: 10, Y: 0, Z: 0},
					{X: 10, Y: 10, Z: 0},
					{X: 0, Y: 0, Z: 0},
				}},
			}}},
		},
		{
			Type: "RoofSurface",
			Geometry: types.MultiSurface{Polygons: []types.Polygon{{
				Exterior: types.Ring{Points: []types.Point{
					{X: 0, Y: 0, Z: 15},
					{X: 10, Y: 0, Z: 15},
					{X: 10, Y: 10, Z: 15},
					{X: 0, Y: 0, Z: 15},
				}},
			}}},
		},
	}

	h := DeriveHeight(nil, nil, bounded)
	if h != 15 {
		t.Errorf("height = %g, want 15", h)
	}
}

func TestDeriveHeight_NoGeometry(t *testing.T) {
	h := DeriveHeight(nil, nil, nil)
	if h != 0 {
		t.Errorf("height = %g, want 0", h)
	}
}

func TestDeriveFootprint_GroundSurface(t *testing.T) {
	bounded := []types.Surface{
		{
			Type: "GroundSurface",
			Geometry: types.MultiSurface{Polygons: []types.Polygon{{
				Exterior: types.Ring{Points: []types.Point{
					{X: 0, Y: 0, Z: 100},
					{X: 10, Y: 0, Z: 100},
					{X: 10, Y: 10, Z: 100},
					{X: 0, Y: 0, Z: 100},
				}},
			}}},
		},
	}

	fp := DeriveFootprint(nil, nil, bounded)
	if fp == nil {
		t.Fatal("expected footprint")
	}
	// All Z should be 0 after projection.
	for _, pt := range fp.Exterior.Points {
		if pt.Z != 0 {
			t.Errorf("projected point has Z=%g, want 0", pt.Z)
		}
	}

	if len(fp.Exterior.Points) != 4 {
		t.Errorf("got %d points, want 4", len(fp.Exterior.Points))
	}
}

func TestDeriveFootprint_FromSolid(t *testing.T) {
	solid := makeBox(100, 112)

	fp := DeriveFootprint(solid, nil, nil)
	if fp == nil {
		t.Fatal("expected footprint")
	}
	// Should pick the bottom polygon (lower avgZ).
	for _, pt := range fp.Exterior.Points {
		if pt.Z != 0 {
			t.Errorf("projected point has Z=%g, want 0", pt.Z)
		}
	}
}

func TestDeriveFootprint_NoGeometry(t *testing.T) {
	fp := DeriveFootprint(nil, nil, nil)
	if fp != nil {
		t.Error("expected nil footprint")
	}
}

func TestDeriveFootprint_PreservesHoles(t *testing.T) {
	bounded := []types.Surface{{
		Type: "GroundSurface",
		Geometry: types.MultiSurface{Polygons: []types.Polygon{{
			Exterior: types.Ring{Points: []types.Point{
				{X: 0, Y: 0, Z: 5},
				{X: 20, Y: 0, Z: 5},
				{X: 20, Y: 20, Z: 5},
				{X: 0, Y: 0, Z: 5},
			}},
			Interior: []types.Ring{{Points: []types.Point{
				{X: 5, Y: 5, Z: 5},
				{X: 15, Y: 5, Z: 5},
				{X: 15, Y: 15, Z: 5},
				{X: 5, Y: 5, Z: 5},
			}}},
		}}},
	}}

	fp := DeriveFootprint(nil, nil, bounded)
	if fp == nil {
		t.Fatal("expected footprint")
	}

	if len(fp.Interior) != 1 {
		t.Errorf("got %d interior rings, want 1", len(fp.Interior))
	}
}
