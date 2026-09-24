// Package crs provides parsing and interpretation of CRS (Coordinate Reference System)
// declarations found in CityGML and GML documents.
//
// The library preserves CRS metadata but does not perform coordinate reprojection.
package crs

import (
	"regexp"
	"strconv"

	"github.com/cwbudde/go-citygml/types"
)

// Supported srsName forms:
//
//   EPSG:25832
//   urn:ogc:def:crs:EPSG::25832
//   urn:ogc:def:crs:EPSG:6.12:25832
//   http://www.opengis.net/def/crs/EPSG/0/25832
//   urn:adv:crs:ETRS89_UTM32*DE_DHHN92_NH  (German ADV compound CRS)
//
// Compound (horizontal + vertical) forms are also recognised; they resolve
// to the horizontal component's code:
//
//   EPSG:25832+7837
//   urn:ogc:def:crs,crs:EPSG::25832,crs:EPSG::7837
//   http://www.opengis.net/def/crs-compound?1=http://www.opengis.net/def/crs/EPSG/0/25832&2=http://www.opengis.net/def/crs/EPSG/0/7837
//
// All are normalized to an EPSG integer code.

var (
	reEPSGShort = regexp.MustCompile(`^EPSG:(\d+)(?:\+\d+)?$`)
	reURN       = regexp.MustCompile(`^urn:ogc:def:crs:EPSG:[^:]*:(\d+)$`)
	reHTTP      = regexp.MustCompile(`^https?://www\.opengis\.net/def/crs/EPSG/\d+/(\d+)$`)
	// OGC compound URN: urn:ogc:def:crs,crs:EPSG::<horizontal>,crs:EPSG::<vertical>.
	reURNCompound = regexp.MustCompile(`^urn:ogc:def:crs,crs:EPSG:[^:,]*:(\d+)(?:,.*)?$`)
	// OGC compound HTTP URI: the first component (1=) is the horizontal CRS.
	reHTTPCompound = regexp.MustCompile(`^https?://www\.opengis\.net/def/crs-compound\?1=https?://www\.opengis\.net/def/crs/EPSG/\d+/(\d+)(?:&.*)?$`)
	// German ADV CRS: urn:adv:crs:ETRS89_UTM<zone>[*<vertical>]
	// ETRS89_UTM32 → EPSG:25832, ETRS89_UTM33 → EPSG:25833, etc.
	reADVUTM = regexp.MustCompile(`^urn:adv:crs:ETRS89_UTM(\d+)`)
)

// Parse interprets an srsName string and returns structured CRS metadata.
// For a compound CRS, Code is the horizontal component's EPSG code.
// Returns a CRS with Code=0 if the format is not recognized.
func Parse(srsName string) types.CRS {
	c := types.CRS{Raw: srsName}

	if code := extractCode(srsName); code > 0 {
		c.Code = code
		c.IsYXOrder = isYXOrder(code)
	}

	return c
}

func extractCode(s string) int {
	for _, re := range []*regexp.Regexp{reEPSGShort, reURN, reHTTP, reURNCompound, reHTTPCompound} {
		if m := re.FindStringSubmatch(s); m != nil {
			code, _ := strconv.Atoi(m[1])
			return code
		}
	}

	// German ADV UTM: zone number maps to EPSG 258XX (e.g. UTM32 → 25832)
	if m := reADVUTM.FindStringSubmatch(s); m != nil {
		zone, _ := strconv.Atoi(m[1])
		return 25800 + zone
	}

	return 0
}

// isYXOrder returns true for EPSG codes that define latitude (Y/northing) first.
// Geographic CRS (4326, 4258, etc.) and some projected CRS use Y,X natively.
//
// For CityGML, the most common case is UTM zones (EPSG:25832, 32632, etc.)
// which use easting/northing (X,Y) — so IsYXOrder is false for those.
//
// This covers the most common codes encountered in European CityGML data.
// URN-form srsNames (urn:ogc:def:crs:EPSG::*) follow strict EPSG axis order,
// while short-form (EPSG:*) typically implies X,Y regardless.
func isYXOrder(code int) bool {
	switch code {
	// WGS 84 geographic
	case 4326:
		return true
	// ETRS89 geographic
	case 4258:
		return true
	// Other common geographic CRS
	case 4269: // NAD83
		return true
	case 4167: // NZGD2000
		return true
	default:
		return false
	}
}
