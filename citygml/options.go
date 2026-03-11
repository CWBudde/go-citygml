package citygml

// Options configures the behavior of the CityGML decoder.
type Options struct {
	// Strict controls whether unsupported constructs cause errors (true)
	// or are silently skipped (false). Default is false.
	Strict bool

	// DeriveHeights controls whether the decoder computes building heights
	// from Z extents when MeasuredHeight is not present. Default is true.
	DeriveHeights *bool

	// DeriveFootprints controls whether the decoder computes 2D footprints
	// from 3D geometry. Default is true.
	DeriveFootprints *bool
}

func (o Options) deriveHeights() bool {
	if o.DeriveHeights == nil {
		return true
	}

	return *o.DeriveHeights
}

func (o Options) deriveFootprints() bool {
	if o.DeriveFootprints == nil {
		return true
	}

	return *o.DeriveFootprints
}
