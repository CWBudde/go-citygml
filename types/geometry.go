package types

// Dimensionality indicates coordinate dimensionality.
type Dimensionality int

const (
	// Dim2D indicates 2D coordinates (X, Y).
	Dim2D Dimensionality = 2
	// Dim3D indicates 3D coordinates (X, Y, Z).
	Dim3D Dimensionality = 3
)

// Point represents a coordinate tuple.
type Point struct {
	X, Y, Z float64
}

// Ring represents a closed linear ring of coordinates.
type Ring struct {
	Points []Point
}

// Polygon represents a surface bounded by an exterior ring and optional interior rings (holes).
type Polygon struct {
	Exterior Ring
	Interior []Ring
}

// MultiSurface represents a collection of polygons.
type MultiSurface struct {
	Polygons []Polygon
}

// Solid represents a volumetric geometry bounded by surfaces.
type Solid struct {
	Exterior MultiSurface
}

// Surface represents a single semantic surface (e.g. wall, roof, ground)
// attached to a city object.
type Surface struct {
	// ID is the gml:id of the surface, if present.
	ID string

	// Type classifies the surface (e.g. "RoofSurface", "WallSurface", "GroundSurface").
	Type string

	// Geometry is the surface geometry.
	Geometry MultiSurface
}
