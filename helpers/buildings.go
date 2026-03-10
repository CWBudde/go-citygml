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
	ID              string
	Height          float64
	IsMeasured      bool
	HasHeight       bool
}

// BuildingHeights extracts the effective height for every building in the document.
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

// BuildingFootprints extracts the derived footprint for every building in the document.
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
