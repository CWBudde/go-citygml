package gml

import (
	"math"

	"github.com/cwbudde/go-citygml/types"
)

// DeriveHeight computes a building height from the Z extents of its geometry.
// It returns maxZ - minZ across all coordinates found in the geometry sources.
// Returns 0 if no 3D coordinates are available.
func DeriveHeight(solid *types.Solid, ms *types.MultiSurface, bounded []types.Surface) float64 {
	minZ := math.MaxFloat64
	maxZ := -math.MaxFloat64
	found := false

	visit := func(pt types.Point) {
		found = true

		if pt.Z < minZ {
			minZ = pt.Z
		}

		if pt.Z > maxZ {
			maxZ = pt.Z
		}
	}

	if solid != nil {
		visitMultiSurface(&solid.Exterior, visit)
	}

	if ms != nil {
		visitMultiSurface(ms, visit)
	}

	for i := range bounded {
		visitMultiSurface(&bounded[i].Geometry, visit)
	}

	if !found || maxZ <= minZ {
		return 0
	}

	return maxZ - minZ
}

// DeriveFootprint projects 3D geometry onto the XY plane and returns
// the polygon with the largest area as the footprint candidate.
// It considers ground surfaces first (from bounded), then falls back
// to the lowest-Z polygon from the solid or multi-surface.
func DeriveFootprint(solid *types.Solid, ms *types.MultiSurface, bounded []types.Surface) *types.Polygon {
	// Strategy 1: Use GroundSurface if available.
	for _, surf := range bounded {
		if surf.Type == "GroundSurface" && len(surf.Geometry.Polygons) > 0 {
			proj := projectPolygon(surf.Geometry.Polygons[0])
			return &proj
		}
	}

	// Strategy 2: Find the polygon with the lowest average Z (likely the footprint).
	var allPolygons []types.Polygon
	if solid != nil {
		allPolygons = append(allPolygons, solid.Exterior.Polygons...)
	}

	if ms != nil {
		allPolygons = append(allPolygons, ms.Polygons...)
	}

	if len(allPolygons) == 0 {
		return nil
	}

	bestIdx := 0

	bestAvgZ := avgZ(allPolygons[0])
	for i := 1; i < len(allPolygons); i++ {
		az := avgZ(allPolygons[i])
		if az < bestAvgZ {
			bestAvgZ = az
			bestIdx = i
		}
	}

	proj := projectPolygon(allPolygons[bestIdx])

	return &proj
}

// projectPolygon projects a polygon onto the XY plane (Z=0).
func projectPolygon(poly types.Polygon) types.Polygon {
	return types.Polygon{
		Exterior: projectRing(poly.Exterior),
		Interior: projectRings(poly.Interior),
	}
}

func projectRing(ring types.Ring) types.Ring {
	pts := make([]types.Point, len(ring.Points))
	for i, pt := range ring.Points {
		pts[i] = types.Point{X: pt.X, Y: pt.Y}
	}

	return types.Ring{Points: pts}
}

func projectRings(rings []types.Ring) []types.Ring {
	if len(rings) == 0 {
		return nil
	}

	out := make([]types.Ring, len(rings))
	for i, r := range rings {
		out[i] = projectRing(r)
	}

	return out
}

func avgZ(poly types.Polygon) float64 {
	sum := 0.0
	n := 0

	for _, pt := range poly.Exterior.Points {
		sum += pt.Z
		n++
	}

	if n == 0 {
		return 0
	}

	return sum / float64(n)
}

func visitMultiSurface(ms *types.MultiSurface, fn func(types.Point)) {
	for _, poly := range ms.Polygons {
		visitRing(poly.Exterior, fn)

		for _, ring := range poly.Interior {
			visitRing(ring, fn)
		}
	}
}

func visitRing(ring types.Ring, fn func(types.Point)) {
	for _, pt := range ring.Points {
		fn(pt)
	}
}
