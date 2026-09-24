package gml

import (
	"strings"
	"testing"

	"github.com/cwbudde/go-citygml/internal/xmlscan"
	"github.com/cwbudde/go-citygml/types"
)

// ring2DSixPoints is a closed 2D ring of six points: 12 values, which the
// divisibility heuristic would misread as four 3D points.
const ring2DSixPoints = "0 0 10 0 20 0 20 10 0 10 0 0"

func TestParseLinearRing_SRSDimensionOnPosList(t *testing.T) {
	input := `<root xmlns:gml="http://www.opengis.net/gml">
		<gml:LinearRing>
			<gml:posList srsDimension="2">` + ring2DSixPoints + `</gml:posList>
		</gml:LinearRing>
	</root>`

	ring, dim := scanAndParse(t, input, "LinearRing", ParseLinearRing)
	if dim != types.Dim2D {
		t.Errorf("dim = %d, want 2", dim)
	}

	if len(ring.Points) != 6 {
		t.Fatalf("got %d points, want 6", len(ring.Points))
	}

	if ring.Points[2] != (types.Point{X: 20, Y: 0}) {
		t.Errorf("point[2] = %+v, want {20 0 0}", ring.Points[2])
	}
}

func TestParsePolygon_SRSDimensionInherited(t *testing.T) {
	input := `<root xmlns:gml="http://www.opengis.net/gml">
		<gml:Polygon srsDimension="2">
			<gml:exterior>
				<gml:LinearRing>
					<gml:posList>` + ring2DSixPoints + `</gml:posList>
				</gml:LinearRing>
			</gml:exterior>
		</gml:Polygon>
	</root>`

	poly, dim := scanAndParse(t, input, "Polygon", ParsePolygon)
	if dim != types.Dim2D {
		t.Errorf("dim = %d, want 2", dim)
	}

	if len(poly.Exterior.Points) != 6 {
		t.Fatalf("got %d points, want 6", len(poly.Exterior.Points))
	}
}

func TestParsePolygon_SRSDimensionScopeEnds(t *testing.T) {
	// The first Polygon declares 2D; the second declares nothing and must fall
	// back to the heuristic rather than inherit its closed sibling's value.
	input := `<root xmlns:gml="http://www.opengis.net/gml">
		<gml:MultiSurface>
			<gml:surfaceMember>
				<gml:Polygon srsDimension="2">
					<gml:exterior><gml:LinearRing><gml:posList>` + ring2DSixPoints + `</gml:posList></gml:LinearRing></gml:exterior>
				</gml:Polygon>
			</gml:surfaceMember>
			<gml:surfaceMember>
				<gml:Polygon>
					<gml:exterior><gml:LinearRing><gml:posList>` + ring2DSixPoints + `</gml:posList></gml:LinearRing></gml:exterior>
				</gml:Polygon>
			</gml:surfaceMember>
		</gml:MultiSurface>
	</root>`

	ms, _ := scanAndParse(t, input, "MultiSurface", ParseMultiSurface)
	if len(ms.Polygons) != 2 {
		t.Fatalf("got %d polygons, want 2", len(ms.Polygons))
	}

	if n := len(ms.Polygons[0].Exterior.Points); n != 6 {
		t.Errorf("declared polygon: got %d points, want 6", n)
	}

	if n := len(ms.Polygons[1].Exterior.Points); n != 4 {
		t.Errorf("undeclared polygon: got %d points, want 4 (heuristic)", n)
	}
}

func TestParseLinearRing_SRSDimensionMismatch(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "posList count not divisible by declared dimension",
			input: `<gml:LinearRing><gml:posList srsDimension="3">0 0 1 0 1 1 0 1 0 0</gml:posList></gml:LinearRing>`,
		},
		{
			name:  "unsupported srsDimension",
			input: `<gml:LinearRing><gml:posList srsDimension="4">0 0 0 0 1 0 0 0 1 1 0 0 0 0 0 0</gml:posList></gml:LinearRing>`,
		},
		{
			name:  "pos count differs from declared dimension",
			input: `<gml:LinearRing srsDimension="3"><gml:pos>0 0</gml:pos><gml:pos>1 0</gml:pos></gml:LinearRing>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := `<root xmlns:gml="http://www.opengis.net/gml">` + tt.input + `</root>`

			sc := xmlscan.NewScanner(strings.NewReader(input))
			for {
				elem, err := sc.StartElement()
				if err != nil {
					t.Fatalf("scanning: %v", err)
				}

				if elem.LocalName() == "LinearRing" {
					break
				}
			}

			_, _, err := ParseLinearRing(sc)
			if err == nil {
				t.Fatal("expected an error, got nil")
			}
		})
	}
}

func TestParsePosListFields_Hint(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		dim     types.Dimensionality
		hint    types.Dimensionality
		wantDim types.Dimensionality
		wantN   int
	}{
		{name: "hint applies when divisible", text: ring2DSixPoints, hint: types.Dim2D, wantDim: types.Dim2D, wantN: 6},
		{name: "hint ignored when not divisible", text: "0 0 0 1 1 1 2 2 2", hint: types.Dim2D, wantDim: types.Dim3D, wantN: 3},
		{name: "declared dimension beats hint", text: ring2DSixPoints, dim: types.Dim3D, hint: types.Dim2D, wantDim: types.Dim3D, wantN: 4},
		{name: "no hint falls back to heuristic", text: ring2DSixPoints, wantDim: types.Dim3D, wantN: 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pts, dim, err := parsePosListFields(strings.Fields(tt.text), tt.dim, tt.hint)
			if err != nil {
				t.Fatal(err)
			}

			if dim != tt.wantDim {
				t.Errorf("dim = %d, want %d", dim, tt.wantDim)
			}

			if len(pts) != tt.wantN {
				t.Errorf("got %d points, want %d", len(pts), tt.wantN)
			}
		})
	}
}
