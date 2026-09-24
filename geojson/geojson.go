package geojson

import (
	"encoding/json"

	"github.com/cwbudde/go-citygml/types"
)

// GeoJSON object type tags (RFC 7946).
const (
	typeFeature           = "Feature"
	typeFeatureCollection = "FeatureCollection"
	typePolygon           = "Polygon"
	typeMultiPolygon      = "MultiPolygon"
	typeBuildingPart      = "BuildingPart"
)

// Geometry represents a GeoJSON geometry object.
type Geometry struct {
	Type        string          `json:"type"`
	Coordinates json.RawMessage `json:"coordinates"`
}

// Feature represents a GeoJSON Feature.
type Feature struct {
	Type       string         `json:"type"`
	ID         string         `json:"id,omitempty"`
	Geometry   *Geometry      `json:"geometry"`
	Properties map[string]any `json:"properties"`
}

// FeatureCollection represents a GeoJSON FeatureCollection.
type FeatureCollection struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}

// NewFeatureCollection creates an empty FeatureCollection.
func NewFeatureCollection() *FeatureCollection {
	return &FeatureCollection{
		Type:     typeFeatureCollection,
		Features: []Feature{},
	}
}

// PolygonGeometry converts a types.Polygon to a GeoJSON Polygon geometry.
func PolygonGeometry(poly types.Polygon) *Geometry {
	rings := make([][][2]float64, 0, 1+len(poly.Interior))

	rings = append(rings, ringCoords(poly.Exterior))
	for _, inner := range poly.Interior {
		rings = append(rings, ringCoords(inner))
	}

	coords := marshalCoordinates(rings)

	return &Geometry{
		Type:        typePolygon,
		Coordinates: coords,
	}
}

// MultiPolygonGeometry converts a types.MultiSurface to a GeoJSON MultiPolygon geometry.
func MultiPolygonGeometry(ms types.MultiSurface) *Geometry {
	polys := make([][][][2]float64, len(ms.Polygons))
	for i, poly := range ms.Polygons {
		rings := make([][][2]float64, 0, 1+len(poly.Interior))

		rings = append(rings, ringCoords(poly.Exterior))
		for _, inner := range poly.Interior {
			rings = append(rings, ringCoords(inner))
		}

		polys[i] = rings
	}

	coords := marshalCoordinates(polys)

	return &Geometry{
		Type:        typeMultiPolygon,
		Coordinates: coords,
	}
}

func marshalCoordinates(coords any) json.RawMessage {
	data, err := json.Marshal(coords)
	if err != nil {
		return json.RawMessage("null")
	}

	return data
}

func ringCoords(ring types.Ring) [][2]float64 {
	coords := make([][2]float64, len(ring.Points))
	for i, pt := range ring.Points {
		coords[i] = [2]float64{pt.X, pt.Y}
	}

	return coords
}

// BuildingFeature creates a GeoJSON Feature from a Building.
// It uses the footprint if available, otherwise falls back to the first available geometry.
func BuildingFeature(b *types.Building) Feature {
	props := map[string]any{
		"type": "Building",
	}
	if b.Class != "" {
		props["class"] = b.Class
	}

	if b.Function != "" {
		props["function"] = b.Function
	}

	if b.Usage != "" {
		props["usage"] = b.Usage
	}

	if b.HasMeasuredHeight {
		props["measuredHeight"] = b.MeasuredHeight
	}

	if b.DerivedHeight > 0 {
		props["derivedHeight"] = b.DerivedHeight
	}

	if b.LoD != "" {
		props["lod"] = string(b.LoD)
	}

	var geom *Geometry

	switch {
	case b.Footprint != nil:
		geom = PolygonGeometry(*b.Footprint)
	case b.MultiSurface != nil && len(b.MultiSurface.Polygons) > 0:
		geom = MultiPolygonGeometry(*b.MultiSurface)
	case b.Solid != nil && len(b.Solid.Exterior.Polygons) > 0:
		geom = MultiPolygonGeometry(b.Solid.Exterior)
	}

	return Feature{
		Type:       typeFeature,
		ID:         b.ID,
		Geometry:   geom,
		Properties: props,
	}
}

// TerrainFeature creates a GeoJSON Feature from a Terrain.
func TerrainFeature(t *types.Terrain) Feature {
	props := map[string]any{
		"type": "Terrain",
	}

	var geom *Geometry
	if len(t.Geometry.Polygons) > 0 {
		geom = MultiPolygonGeometry(t.Geometry)
	}

	return Feature{
		Type:       typeFeature,
		ID:         t.ID,
		Geometry:   geom,
		Properties: props,
	}
}

// FromDocument converts a full Document to a GeoJSON FeatureCollection.
// Buildings, building parts and terrain objects are all included as features.
// Each building part follows its parent as its own feature, with properties
// "type": "BuildingPart" and "parent" set to the parent's ID.
func FromDocument(doc *types.Document) *FeatureCollection {
	fc := NewFeatureCollection()

	for i := range doc.Buildings {
		fc.Features = appendBuildingFeatures(fc.Features, &doc.Buildings[i], "")
	}

	for i := range doc.Terrains {
		fc.Features = append(fc.Features, TerrainFeature(&doc.Terrains[i]))
	}

	return fc
}

// appendBuildingFeatures appends the feature for b followed by the features
// of its parts, depth-first. parentID is empty for a top-level building.
func appendBuildingFeatures(features []Feature, b *types.Building, parentID string) []Feature {
	f := BuildingFeature(b)
	if parentID != "" {
		f.Properties["type"] = typeBuildingPart
		f.Properties["parent"] = parentID
	}

	features = append(features, f)
	for i := range b.Parts {
		features = appendBuildingFeatures(features, &b.Parts[i], b.ID)
	}

	return features
}
