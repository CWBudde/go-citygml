package gml

import (
	"strings"
	"testing"

	"github.com/cwbudde/go-citygml/internal/xmlscan"
	"github.com/cwbudde/go-citygml/types"
)

// helper: create scanner, advance past wrapper, call parser on target element.
func scanAndParse[T any](t *testing.T, xml string, localName string, parser func(*xmlscan.Scanner) (T, types.Dimensionality, error)) (T, types.Dimensionality) {
	t.Helper()

	sc := xmlscan.NewScanner(strings.NewReader(xml))
	for {
		elem, err := sc.StartElement()
		if err != nil {
			t.Fatalf("scanning to %s: %v", localName, err)
		}

		if elem.LocalName() == localName {
			val, dim, err := parser(sc)
			if err != nil {
				t.Fatalf("parsing %s: %v", localName, err)
			}

			return val, dim
		}
	}
}

func TestParseLinearRing_PosList(t *testing.T) {
	input := `<root xmlns:gml="http://www.opengis.net/gml">
		<gml:LinearRing>
			<gml:posList>0 0 0 10 0 0 10 10 0 0 0 0</gml:posList>
		</gml:LinearRing>
	</root>`

	ring, dim := scanAndParse(t, input, "LinearRing", ParseLinearRing)
	if dim != types.Dim3D {
		t.Errorf("dim = %d, want 3", dim)
	}

	if len(ring.Points) != 4 {
		t.Fatalf("got %d points, want 4", len(ring.Points))
	}
	// Check closure: first == last
	if ring.Points[0] != ring.Points[3] {
		t.Error("ring not closed")
	}
}

func TestParseLinearRing_MultiplePos(t *testing.T) {
	input := `<root xmlns:gml="http://www.opengis.net/gml">
		<gml:LinearRing>
			<gml:pos>0 0 0</gml:pos>
			<gml:pos>10 0 0</gml:pos>
			<gml:pos>10 10 0</gml:pos>
			<gml:pos>0 0 0</gml:pos>
		</gml:LinearRing>
	</root>`

	ring, dim := scanAndParse(t, input, "LinearRing", ParseLinearRing)
	if dim != types.Dim3D {
		t.Errorf("dim = %d, want 3", dim)
	}

	if len(ring.Points) != 4 {
		t.Fatalf("got %d points, want 4", len(ring.Points))
	}
}

func TestParsePolygon_ExteriorOnly(t *testing.T) {
	input := `<root xmlns:gml="http://www.opengis.net/gml">
		<gml:Polygon>
			<gml:exterior>
				<gml:LinearRing>
					<gml:posList>0 0 0 10 0 0 10 10 0 0 0 0</gml:posList>
				</gml:LinearRing>
			</gml:exterior>
		</gml:Polygon>
	</root>`

	poly, dim := scanAndParse(t, input, "Polygon", ParsePolygon)
	if dim != types.Dim3D {
		t.Errorf("dim = %d, want 3", dim)
	}

	if len(poly.Exterior.Points) != 4 {
		t.Errorf("exterior has %d points, want 4", len(poly.Exterior.Points))
	}

	if len(poly.Interior) != 0 {
		t.Errorf("got %d interior rings, want 0", len(poly.Interior))
	}
}

func TestParsePolygon_WithHole(t *testing.T) {
	input := `<root xmlns:gml="http://www.opengis.net/gml">
		<gml:Polygon>
			<gml:exterior>
				<gml:LinearRing>
					<gml:posList>0 0 0 10 0 0 10 10 0 0 10 0 0 0 0</gml:posList>
				</gml:LinearRing>
			</gml:exterior>
			<gml:interior>
				<gml:LinearRing>
					<gml:posList>2 2 0 8 2 0 8 8 0 2 2 0</gml:posList>
				</gml:LinearRing>
			</gml:interior>
		</gml:Polygon>
	</root>`

	poly, _ := scanAndParse(t, input, "Polygon", ParsePolygon)
	if len(poly.Interior) != 1 {
		t.Fatalf("got %d interior rings, want 1", len(poly.Interior))
	}

	if len(poly.Interior[0].Points) != 4 {
		t.Errorf("interior ring has %d points, want 4", len(poly.Interior[0].Points))
	}
}

func TestParseMultiSurface(t *testing.T) {
	input := `<root xmlns:gml="http://www.opengis.net/gml">
		<gml:MultiSurface>
			<gml:surfaceMember>
				<gml:Polygon>
					<gml:exterior>
						<gml:LinearRing>
							<gml:posList>0 0 0 10 0 0 10 10 0 0 0 0</gml:posList>
						</gml:LinearRing>
					</gml:exterior>
				</gml:Polygon>
			</gml:surfaceMember>
			<gml:surfaceMember>
				<gml:Polygon>
					<gml:exterior>
						<gml:LinearRing>
							<gml:posList>0 0 5 10 0 5 10 10 5 0 0 5</gml:posList>
						</gml:LinearRing>
					</gml:exterior>
				</gml:Polygon>
			</gml:surfaceMember>
		</gml:MultiSurface>
	</root>`

	ms, dim := scanAndParse(t, input, "MultiSurface", ParseMultiSurface)
	if dim != types.Dim3D {
		t.Errorf("dim = %d, want 3", dim)
	}

	if len(ms.Polygons) != 2 {
		t.Fatalf("got %d polygons, want 2", len(ms.Polygons))
	}
}

func TestParseSolid(t *testing.T) {
	input := `<root xmlns:gml="http://www.opengis.net/gml">
		<gml:Solid>
			<gml:exterior>
				<gml:CompositeSurface>
					<gml:surfaceMember>
						<gml:Polygon>
							<gml:exterior>
								<gml:LinearRing>
									<gml:posList>0 0 0 10 0 0 10 10 0 0 0 0</gml:posList>
								</gml:LinearRing>
							</gml:exterior>
						</gml:Polygon>
					</gml:surfaceMember>
					<gml:surfaceMember>
						<gml:Polygon>
							<gml:exterior>
								<gml:LinearRing>
									<gml:posList>0 0 5 10 0 5 10 10 5 0 0 5</gml:posList>
								</gml:LinearRing>
							</gml:exterior>
						</gml:Polygon>
					</gml:surfaceMember>
				</gml:CompositeSurface>
			</gml:exterior>
		</gml:Solid>
	</root>`

	solid, dim := scanAndParse(t, input, "Solid", ParseSolid)
	if dim != types.Dim3D {
		t.Errorf("dim = %d, want 3", dim)
	}

	if len(solid.Exterior.Polygons) != 2 {
		t.Fatalf("got %d exterior polygons, want 2", len(solid.Exterior.Polygons))
	}
}

func TestParseCompositeSurface(t *testing.T) {
	input := `<root xmlns:gml="http://www.opengis.net/gml">
		<gml:CompositeSurface>
			<gml:surfaceMember>
				<gml:Polygon>
					<gml:exterior>
						<gml:LinearRing>
							<gml:posList>0 0 10 0 10 10 0 0</gml:posList>
						</gml:LinearRing>
					</gml:exterior>
				</gml:Polygon>
			</gml:surfaceMember>
		</gml:CompositeSurface>
	</root>`

	ms, dim := scanAndParse(t, input, "CompositeSurface", ParseCompositeSurface)
	if dim != types.Dim2D {
		t.Errorf("dim = %d, want 2", dim)
	}

	if len(ms.Polygons) != 1 {
		t.Fatalf("got %d polygons, want 1", len(ms.Polygons))
	}
}

func TestParseMultiSurface_Empty(t *testing.T) {
	input := `<root xmlns:gml="http://www.opengis.net/gml">
		<gml:MultiSurface/>
	</root>`

	sc := xmlscan.NewScanner(strings.NewReader(input))
	for {
		elem, err := sc.StartElement()
		if err != nil {
			t.Fatal(err)
		}

		if elem.LocalName() == "MultiSurface" {
			ms, _, err := ParseMultiSurface(sc)
			if err != nil {
				t.Fatal(err)
			}

			if len(ms.Polygons) != 0 {
				t.Errorf("got %d polygons, want 0", len(ms.Polygons))
			}

			return
		}
	}
}
