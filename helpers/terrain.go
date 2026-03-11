package helpers

import "github.com/cwbudde/go-citygml/types"

// TerrainSummary provides aggregated terrain information from a document.
type TerrainSummary struct {
	// Count is the number of terrain objects.
	Count int
	// TotalPolygons is the total number of polygons across all terrains.
	TotalPolygons int
	// Polygons is the flattened list of all terrain polygons.
	Polygons []types.Polygon
}

// SummarizeTerrain aggregates terrain geometry from all terrain objects in the document.
func SummarizeTerrain(doc *types.Document) TerrainSummary {
	var s TerrainSummary

	s.Count = len(doc.Terrains)
	for _, t := range doc.Terrains {
		s.TotalPolygons += len(t.Geometry.Polygons)
		s.Polygons = append(s.Polygons, t.Geometry.Polygons...)
	}

	return s
}
