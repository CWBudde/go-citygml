package helpers

import "github.com/cwbudde/go-citygml/types"

// BuildingHeight returns the effective height of a building.
// It prefers MeasuredHeight when available, falls back to DerivedHeight.
// Returns 0 and false if no height is available.
func BuildingHeight(b *types.Building) (float64, bool) {
	if b.HasMeasuredHeight {
		return b.MeasuredHeight, true
	}

	if b.DerivedHeight > 0 {
		return b.DerivedHeight, true
	}

	return 0, false
}

// HeightResult pairs a building ID with its effective height.
type HeightResult struct {
	ID         string
	Height     float64
	IsMeasured bool
	HasHeight  bool
}

// BuildingHeights extracts the effective height for every top-level building
// in the document. BuildingParts are not included; use FlattenBuildings to
// reach them.
func BuildingHeights(doc *types.Document) []HeightResult {
	results := make([]HeightResult, len(doc.Buildings))
	for i := range doc.Buildings {
		b := &doc.Buildings[i]
		h, ok := BuildingHeight(b)
		results[i] = HeightResult{
			ID:         b.ID,
			Height:     h,
			IsMeasured: b.HasMeasuredHeight,
			HasHeight:  ok,
		}
	}

	return results
}

// FootprintResult pairs a building ID with its footprint polygon.
type FootprintResult struct {
	ID        string
	Footprint *types.Polygon
}

// BuildingFootprints extracts the derived footprint for every top-level
// building in the document. BuildingParts are not included; use
// FlattenBuildings to reach them.
// Buildings without a footprint will have a nil Footprint field.
func BuildingFootprints(doc *types.Document) []FootprintResult {
	results := make([]FootprintResult, len(doc.Buildings))
	for i := range doc.Buildings {
		b := &doc.Buildings[i]
		results[i] = FootprintResult{
			ID:        b.ID,
			Footprint: b.Footprint,
		}
	}

	return results
}

// FlattenBuildings returns every building and building part in the document,
// depth-first in document order: each building is followed by its parts
// (and their nested parts). The pointers refer into doc and stay valid as
// long as doc.Buildings and the Parts slices are not modified.
//
// A building composed of parts often carries no geometry or height of its
// own, so consumers that want every volume typically keep the entries that
// have a Footprint and a height.
func FlattenBuildings(doc *types.Document) []*types.Building {
	var out []*types.Building

	var walk func(b *types.Building)

	walk = func(b *types.Building) {
		out = append(out, b)
		for i := range b.Parts {
			walk(&b.Parts[i])
		}
	}

	for i := range doc.Buildings {
		walk(&doc.Buildings[i])
	}

	return out
}
