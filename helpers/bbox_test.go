package helpers

import (
	"testing"

	"github.com/cwbudde/go-citygml/types"
)

func TestDocumentBBox_Empty(t *testing.T) {
	doc := &types.Document{}
	bb := DocumentBBox(doc)
	if !bb.Empty {
		t.Error("expected empty bbox")
	}
}

func TestDocumentBBox_Buildings(t *testing.T) {
	doc := &types.Document{
		Buildings: []types.Building{{
			ID: "B1",
			Solid: &types.Solid{
				Exterior: types.MultiSurface{Polygons: []types.Polygon{{
					Exterior: types.Ring{Points: []types.Point{
						{X: 10, Y: 20, Z: 30},
						{X: 40, Y: 50, Z: 60},
						{X: 10, Y: 20, Z: 30},
					}},
				}}},
			},
		}},
	}
	bb := DocumentBBox(doc)
	if bb.Empty {
		t.Fatal("expected non-empty bbox")
	}
	if bb.MinX != 10 || bb.MaxX != 40 {
		t.Errorf("X: got [%v, %v], want [10, 40]", bb.MinX, bb.MaxX)
	}
	if bb.MinY != 20 || bb.MaxY != 50 {
		t.Errorf("Y: got [%v, %v], want [20, 50]", bb.MinY, bb.MaxY)
	}
	if bb.MinZ != 30 || bb.MaxZ != 60 {
		t.Errorf("Z: got [%v, %v], want [30, 60]", bb.MinZ, bb.MaxZ)
	}
	if !bb.Has3D {
		t.Error("expected Has3D=true")
	}
}

func TestDocumentBBox_Terrain(t *testing.T) {
	doc := &types.Document{
		Terrains: []types.Terrain{{
			ID: "T1",
			Geometry: types.MultiSurface{Polygons: []types.Polygon{{
				Exterior: types.Ring{Points: []types.Point{
					{X: -5, Y: -10},
					{X: 5, Y: 10},
					{X: -5, Y: -10},
				}},
			}}},
		}},
	}
	bb := DocumentBBox(doc)
	if bb.Empty {
		t.Fatal("expected non-empty bbox")
	}
	if bb.MinX != -5 || bb.MaxX != 5 {
		t.Errorf("X: got [%v, %v], want [-5, 5]", bb.MinX, bb.MaxX)
	}
	if bb.Has3D {
		t.Error("expected Has3D=false for 2D-only terrain")
	}
}

func TestDocumentBBox_BoundedBy(t *testing.T) {
	doc := &types.Document{
		Buildings: []types.Building{{
			ID: "B1",
			BoundedBy: []types.Surface{{
				Type: "GroundSurface",
				Geometry: types.MultiSurface{Polygons: []types.Polygon{{
					Exterior: types.Ring{Points: []types.Point{
						{X: 100, Y: 200, Z: 0},
						{X: 300, Y: 400, Z: 0},
						{X: 100, Y: 200, Z: 0},
					}},
				}}},
			}},
		}},
	}
	bb := DocumentBBox(doc)
	if bb.MinX != 100 || bb.MaxX != 300 {
		t.Errorf("X: got [%v, %v], want [100, 300]", bb.MinX, bb.MaxX)
	}
	if bb.MinY != 200 || bb.MaxY != 400 {
		t.Errorf("Y: got [%v, %v], want [200, 400]", bb.MinY, bb.MaxY)
	}
}
