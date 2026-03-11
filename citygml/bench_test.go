package citygml

import (
	"fmt"
	"strings"
	"testing"
)

// generateBuildings creates a CityGML 2.0 document with n buildings,
// each having a lod1Solid with 2 polygons.
func generateBuildings(n int) string {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8"?>
<CityModel xmlns="http://www.opengis.net/citygml/2.0"
           xmlns:gml="http://www.opengis.net/gml"
           xmlns:bldg="http://www.opengis.net/citygml/building/2.0">
  <gml:boundedBy>
    <gml:Envelope srsName="EPSG:25832">
      <gml:lowerCorner>500000 5700000 0</gml:lowerCorner>
      <gml:upperCorner>510000 5710000 100</gml:upperCorner>
    </gml:Envelope>
  </gml:boundedBy>
`)

	for i := range n {
		fmt.Fprintf(&b, `  <cityObjectMember>
    <bldg:Building gml:id="BLDG_%d">
      <bldg:class>1000</bldg:class>
      <bldg:measuredHeight>%.1f</bldg:measuredHeight>
      <bldg:lod1Solid>
        <gml:Solid>
          <gml:exterior>
            <gml:CompositeSurface>
              <gml:surfaceMember>
                <gml:Polygon>
                  <gml:exterior>
                    <gml:LinearRing>
                      <gml:posList>%d 5700000 0 %d 5700000 0 %d 5700010 0 %d 5700010 0 %d 5700000 0</gml:posList>
                    </gml:LinearRing>
                  </gml:exterior>
                </gml:Polygon>
              </gml:surfaceMember>
              <gml:surfaceMember>
                <gml:Polygon>
                  <gml:exterior>
                    <gml:LinearRing>
                      <gml:posList>%d 5700000 10 %d 5700000 10 %d 5700010 10 %d 5700010 10 %d 5700000 10</gml:posList>
                    </gml:LinearRing>
                  </gml:exterior>
                </gml:Polygon>
              </gml:surfaceMember>
            </gml:CompositeSurface>
          </gml:exterior>
        </gml:Solid>
      </bldg:lod1Solid>
    </bldg:Building>
  </cityObjectMember>
`, i, float64(i%30)+5.0,
			500000+i*10, 500010+i*10, 500010+i*10, 500000+i*10, 500000+i*10,
			500000+i*10, 500010+i*10, 500010+i*10, 500000+i*10, 500000+i*10)
	}

	b.WriteString(`</CityModel>`)

	return b.String()
}

func BenchmarkRead_1Building(b *testing.B) {
	input := generateBuildings(1)
	benchmarkRead(b, input)
}

func BenchmarkRead_10Buildings(b *testing.B) {
	input := generateBuildings(10)
	benchmarkRead(b, input)
}

func BenchmarkRead_100Buildings(b *testing.B) {
	input := generateBuildings(100)
	benchmarkRead(b, input)
}

func BenchmarkRead_1000Buildings(b *testing.B) {
	input := generateBuildings(1000)
	benchmarkRead(b, input)
}

func BenchmarkRead_5000Buildings(b *testing.B) {
	input := generateBuildings(5000)
	benchmarkRead(b, input)
}

func benchmarkRead(b *testing.B, input string) {
	b.Helper()
	b.ReportAllocs()
	b.SetBytes(int64(len(input)))

	for b.Loop() {
		_, err := Read(strings.NewReader(input), Options{})
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRead_NoDerive(b *testing.B) {
	input := generateBuildings(100)
	f := false
	opts := Options{DeriveHeights: &f, DeriveFootprints: &f}

	b.ReportAllocs()
	b.SetBytes(int64(len(input)))

	for b.Loop() {
		_, err := Read(strings.NewReader(input), opts)
		if err != nil {
			b.Fatal(err)
		}
	}
}
