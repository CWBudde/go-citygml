package xmlscan

// Known namespace URIs for CityGML and related standards.
const (
	// CityGML 2.0 core namespaces.
	NSCityGML20Core = "http://www.opengis.net/citygml/2.0"
	NSCityGML20Bldg = "http://www.opengis.net/citygml/building/2.0"
	NSCityGML20Dem  = "http://www.opengis.net/citygml/relief/2.0"
	NSCityGML20Tran = "http://www.opengis.net/citygml/transportation/2.0"
	NSCityGML20Veg  = "http://www.opengis.net/citygml/vegetation/2.0"
	NSCityGML20Gen  = "http://www.opengis.net/citygml/generics/2.0"

	// CityGML 3.0 core namespaces.
	NSCityGML30Core = "http://www.opengis.net/citygml/3.0"
	NSCityGML30Bldg = "http://www.opengis.net/citygml/building/3.0"
	NSCityGML30Dem  = "http://www.opengis.net/citygml/relief/3.0"
	NSCityGML30Con  = "http://www.opengis.net/citygml/construction/3.0"

	// GML namespaces.
	NSGML31 = "http://www.opengis.net/gml"
	NSGML32 = "http://www.opengis.net/gml/3.2"

	// XLink namespace.
	NSXLink = "http://www.w3.org/1999/xlink"
)

// Version represents a detected CityGML version.
type Version string

const (
	VersionUnknown Version = ""
	Version20      Version = "2.0"
	Version30      Version = "3.0"
)

// coreNamespaceVersion maps core CityGML namespace URIs to their version.
var coreNamespaceVersion = map[string]Version{
	NSCityGML20Core: Version20,
	NSCityGML30Core: Version30,
}

// DetectVersion returns the CityGML version based on a namespace URI.
// Returns VersionUnknown if the namespace is not a recognized core CityGML namespace.
func DetectVersion(nsURI string) Version {
	if v, ok := coreNamespaceVersion[nsURI]; ok {
		return v
	}

	return VersionUnknown
}
