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

// square returns a closed axis-aligned square ring at height z.
func square(x0, y0, size, z float64) types.Ring {
	return types.Ring{Points: []types.Point{
		{X: x0, Y: y0, Z: z},
		{X: x0 + size, Y: y0, Z: z},
		{X: x0 + size, Y: y0 + size, Z: z},
		{X: x0, Y: y0 + size, Z: z},
		{X: x0, Y: y0, Z: z},
	}}
}

func TestDeriveFootprint_LargestGroundPolygon(t *testing.T) {
	tests := []struct {
		name    string
		bounded []types.Surface
		wantX0  float64 // X of the chosen polygon's first exterior point
		holes   int
	}{
		{
			name: "largest polygon within one GroundSurface",
			bounded: []types.Surface{{
				Type: groundSurfaceType,
				Geometry: types.MultiSurface{Polygons: []types.Polygon{
					{Exterior: square(0, 0, 2, 50)},
					{Exterior: square(100, 0, 10, 50)},
					{Exterior: square(200, 0, 5, 50)},
				}},
			}},
			wantX0: 100,
		},
		{
			name: "largest polygon across GroundSurfaces, other surfaces ignored",
			bounded: []types.Surface{
				{Type: "RoofSurface", Geometry: types.MultiSurface{Polygons: []types.Polygon{{Exterior: square(300, 0, 50, 60)}}}},
				{Type: groundSurfaceType, Geometry: types.MultiSurface{Polygons: []types.Polygon{{Exterior: square(0, 0, 3, 50)}}}},
				{Type: groundSurfaceType, Geometry: types.MultiSurface{Polygons: []types.Polygon{{Exterior: square(100, 0, 4, 50)}}}},
			},
			wantX0: 100,
		},
		{
			name: "holes reduce the area that is compared",
			bounded: []types.Surface{{
				Type: groundSurfaceType,
				Geometry: types.MultiSurface{Polygons: []types.Polygon{
					// 10x10 minus an 8x8 hole = 36.
					{Exterior: square(0, 0, 10, 50), Interior: []types.Ring{square(1, 1, 8, 50)}},
					// 7x7 = 49.
					{Exterior: square(100, 0, 7, 50)},
				}},
			}},
			wantX0: 100,
		},
		{
			name: "chosen polygon keeps its holes",
			bounded: []types.Surface{{
				Type: groundSurfaceType,
				Geometry: types.MultiSurface{Polygons: []types.Polygon{
					{Exterior: square(100, 0, 2, 50)},
					{Exterior: square(0, 0, 20, 50), Interior: []types.Ring{square(5, 5, 2, 50), square(10, 10, 2, 50)}},
				}},
			}},
			wantX0: 0,
			holes:  2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fp := DeriveFootprint(nil, nil, tt.bounded)
			if fp == nil {
				t.Fatal("expected footprint")
			}

			if got := fp.Exterior.Points[0].X; got != tt.wantX0 {
				t.Errorf("chosen polygon starts at X=%g, want %g", got, tt.wantX0)
			}

			if len(fp.Interior) != tt.holes {
				t.Errorf("got %d interior rings, want %d", len(fp.Interior), tt.holes)
			}

			for _, ring := range append([]types.Ring{fp.Exterior}, fp.Interior...) {
				for _, pt := range ring.Points {
					if pt.Z != 0 {
						t.Fatalf("projected point has Z=%g, want 0", pt.Z)
					}
				}
			}
		})
	}
}

func TestPlanarArea(t *testing.T) {
	tests := []struct {
		name string
		poly types.Polygon
		want float64
	}{
		{name: "square", poly: types.Polygon{Exterior: square(0, 0, 4, 7)}, want: 16},
		{name: "square with hole", poly: types.Polygon{Exterior: square(0, 0, 4, 0), Interior: []types.Ring{square(1, 1, 1, 0)}}, want: 15},
		{name: "clockwise ring", poly: types.Polygon{Exterior: types.Ring{Points: []types.Point{{X: 0, Y: 0}, {X: 0, Y: 2}, {X: 3, Y: 2}, {X: 3, Y: 0}, {X: 0, Y: 0}}}}, want: 6},
		{name: "degenerate", poly: types.Polygon{Exterior: types.Ring{Points: []types.Point{{X: 0, Y: 0}, {X: 1, Y: 1}}}}, want: 0},
		{name: "vertical wall has no planar area", poly: types.Polygon{Exterior: types.Ring{Points: []types.Point{{X: 0, Y: 0, Z: 0}, {X: 5, Y: 0, Z: 0}, {X: 5, Y: 0, Z: 3}, {X: 0, Y: 0, Z: 3}, {X: 0, Y: 0, Z: 0}}}}, want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := planarArea(tt.poly); got != tt.want {
				t.Errorf("planarArea = %g, want %g", got, tt.want)
			}
		})
	}
}
