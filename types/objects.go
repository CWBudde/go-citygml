package types

// LoD represents a Level of Detail indicator.
type LoD string

const (
	LoD0 LoD = "0"
	LoD1 LoD = "1"
	LoD2 LoD = "2"
	LoD3 LoD = "3"
	LoD4 LoD = "4"
)

// CityObject represents a generic city object that doesn't map
// to a specific supported semantic type.
type CityObject struct {
	// ID is the gml:id of the object.
	ID string

	// Type is the XML element local name (e.g. "Bridge", "Tunnel").
	Type string
}

// Building represents a parsed CityGML building or building part.
type Building struct {
	// ID is the gml:id of the building.
	ID string

	// Class is the CityGML class attribute, if present.
	Class string

	// Function is the CityGML function attribute, if present.
	Function string

	// Usage is the CityGML usage attribute, if present.
	Usage string

	// MeasuredHeight is the explicitly stated height in the source data.
	// Zero value means not present; use HasMeasuredHeight to distinguish.
	MeasuredHeight    float64
	HasMeasuredHeight bool

	// DerivedHeight is the height computed from Z extents when MeasuredHeight is absent.
	DerivedHeight float64

	// LoD indicates the Level of Detail of the geometry.
	LoD LoD

	// Solid is the building's solid geometry, if present.
	Solid *Solid

	// MultiSurface is the building's multi-surface geometry, if present.
	MultiSurface *MultiSurface

	// BoundedBy contains semantic surfaces (roof, wall, ground) attached to this building.
	BoundedBy []Surface

	// Footprint is a 2D polygon derived from the 3D geometry.
	Footprint *Polygon

	// Parts contains the building's BuildingParts (bldg:consistsOfBuildingPart
	// in CityGML 1.0/2.0, bldg:buildingPart in 3.0), in document order. Each
	// part is decoded like a building: it has its own ID, attributes,
	// measured or derived height, geometry, semantic surfaces, footprint and,
	// possibly, nested Parts. A building made of parts often carries no
	// geometry or height of its own, so consumers wanting every volume should
	// visit the parts too (see helpers.FlattenBuildings).
	Parts []Building
}

// Terrain represents a parsed CityGML terrain surface.
type Terrain struct {
	// ID is the gml:id of the terrain object.
	ID string

	// Geometry is the terrain surface geometry.
	Geometry MultiSurface
}
