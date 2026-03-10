package types

// Document represents a parsed and normalized CityGML document.
type Document struct {
	// Version is the detected CityGML version (e.g. "2.0", "3.0").
	Version string

	// SRSName is the raw srsName string from the document (from root attributes or gml:Envelope).
	SRSName string

	// CRS is the parsed coordinate reference system metadata derived from SRSName.
	CRS CRS

	// Buildings contains all extracted building objects.
	Buildings []Building

	// Terrains contains all extracted terrain surfaces.
	Terrains []Terrain

	// GenericObjects contains city objects that don't map to a specific supported type.
	GenericObjects []CityObject
}
