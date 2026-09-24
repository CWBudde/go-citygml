package gml

import (
	"math"

	"github.com/cwbudde/go-citygml/types"
)

// groundSurfaceType is the CityGML bounded-surface type preferred for footprint derivation.
const groundSurfaceType = "GroundSurface"

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

// DeriveFootprint projects 3D geometry onto the XY plane and returns a 2D
// footprint candidate. It considers every polygon of every GroundSurface in
// bounded and returns the one with the largest planar area (exterior minus
// holes), interior rings included. Without a usable GroundSurface it falls
// back to the lowest-Z polygon from the solid or multi-surface.
func DeriveFootprint(solid *types.Solid, ms *types.MultiSurface, bounded []types.Surface) *types.Polygon {
	// Strategy 1: the largest GroundSurface polygon.
	if ground := largestGroundPolygon(bounded); ground != nil {
		proj := projectPolygon(*ground)
		return &proj
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

// largestGroundPolygon returns the GroundSurface polygon with the largest
// planar area, or nil if bounded has no GroundSurface polygon. Ties keep the
// first polygon in document order.
func largestGroundPolygon(bounded []types.Surface) *types.Polygon {
	var best *types.Polygon

	bestArea := -1.0

	for i := range bounded {
		if bounded[i].Type != groundSurfaceType {
			continue
		}

		polys := bounded[i].Geometry.Polygons
		for j := range polys {
			if a := planarArea(polys[j]); a > bestArea {
				bestArea = a
				best = &polys[j]
			}
		}
	}

	return best
}

// planarArea returns the area of poly projected onto the XY plane: the
// absolute shoelace area of the exterior ring minus that of each interior
// ring. Z is ignored.
func planarArea(poly types.Polygon) float64 {
	area := ringArea(poly.Exterior)
	for _, hole := range poly.Interior {
		area -= ringArea(hole)
	}

	return math.Max(area, 0)
}

// ringArea returns the absolute XY shoelace area of a ring, closed or not.
func ringArea(ring types.Ring) float64 {
	pts := ring.Points
	if len(pts) < 3 {
		return 0
	}

	sum := 0.0

	for i := range pts {
		j := (i + 1) % len(pts)
		sum += pts[i].X*pts[j].Y - pts[j].X*pts[i].Y
	}

	return math.Abs(sum) / 2
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
