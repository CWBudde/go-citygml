package types

// CRS represents parsed coordinate reference system metadata.
type CRS struct {
	// Raw is the original srsName string as found in the document.
	Raw string

	// Code is the extracted EPSG code (e.g. 25832), or 0 if not recognized.
	Code int

	// IsYXOrder indicates that the CRS uses latitude/northing first (Y,X)
	// as its native axis order per the EPSG definition.
	// When true, coordinates in the source may need swapping to reach X,Y (easting,northing).
	IsYXOrder bool
}
