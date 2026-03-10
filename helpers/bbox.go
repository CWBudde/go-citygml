package helpers

import (
	"math"

	"github.com/cwbudde/go-citygml/types"
)

// BBox represents an axis-aligned bounding box.
type BBox struct {
	MinX, MinY, MinZ float64
	MaxX, MaxY, MaxZ float64
	// Has3D is true if any Z coordinates were found.
	Has3D bool
	// Empty is true if no coordinates were found.
	Empty bool
}

// DocumentBBox computes the axis-aligned bounding box over all geometry in the document.
func DocumentBBox(doc *types.Document) BBox {
	bb := BBox{
		MinX:  math.MaxFloat64,
		MinY:  math.MaxFloat64,
		MinZ:  math.MaxFloat64,
		MaxX:  -math.MaxFloat64,
		MaxY:  -math.MaxFloat64,
		MaxZ:  -math.MaxFloat64,
		Empty: true,
	}

	visit := func(pt types.Point) {
		bb.Empty = false
		if pt.X < bb.MinX {
			bb.MinX = pt.X
		}
		if pt.X > bb.MaxX {
			bb.MaxX = pt.X
		}
		if pt.Y < bb.MinY {
			bb.MinY = pt.Y
		}
		if pt.Y > bb.MaxY {
			bb.MaxY = pt.Y
		}
		if pt.Z != 0 {
			bb.Has3D = true
		}
		if pt.Z < bb.MinZ {
			bb.MinZ = pt.Z
		}
		if pt.Z > bb.MaxZ {
			bb.MaxZ = pt.Z
		}
	}

	for i := range doc.Buildings {
		b := &doc.Buildings[i]
		if b.Solid != nil {
			visitMultiSurface(&b.Solid.Exterior, visit)
		}
		if b.MultiSurface != nil {
			visitMultiSurface(b.MultiSurface, visit)
		}
		for j := range b.BoundedBy {
			visitMultiSurface(&b.BoundedBy[j].Geometry, visit)
		}
	}

	for i := range doc.Terrains {
		visitMultiSurface(&doc.Terrains[i].Geometry, visit)
	}

	return bb
}

func visitMultiSurface(ms *types.MultiSurface, fn func(types.Point)) {
	for _, poly := range ms.Polygons {
		for _, pt := range poly.Exterior.Points {
			fn(pt)
		}
		for _, ring := range poly.Interior {
			for _, pt := range ring.Points {
				fn(pt)
			}
		}
	}
}
