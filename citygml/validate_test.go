package citygml

import (
	"math"
	"testing"

	"github.com/cwbudde/go-citygml/types"
)

func TestValidate_ValidDocument(t *testing.T) {
	doc := &types.Document{
		Version: "2.0",
		SRSName: "EPSG:25832",
		CRS:     types.CRS{Raw: "EPSG:25832", Code: 25832},
		Buildings: []types.Building{{
			ID:                "B1",
			HasMeasuredHeight: true,
			MeasuredHeight:    10.0,
			Solid: &types.Solid{
				Exterior: types.MultiSurface{Polygons: []types.Polygon{{
					Exterior: types.Ring{Points: []types.Point{
						{X: 0, Y: 0, Z: 0},
						{X: 10, Y: 0, Z: 0},
						{X: 10, Y: 10, Z: 0},
						{X: 0, Y: 10, Z: 0},
						{X: 0, Y: 0, Z: 0},
					}},
				}}},
			},
		}},
	}

	findings := Validate(doc)
	for _, f := range findings {
		if f.Severity == SeverityError {
			t.Errorf("unexpected error: %s", f)
		}
	}
}

func TestValidate_MissingCRS(t *testing.T) {
	doc := &types.Document{
		Version: "2.0",
		Buildings: []types.Building{{
			ID:                "B1",
			HasMeasuredHeight: true,
			MeasuredHeight:    5,
		}},
	}

	findings := Validate(doc)
	found := false

	for _, f := range findings {
		if f.Path == "Document" && f.Severity == SeverityWarning && f.Message == "no srsName (CRS) declared" {
			found = true
		}
	}

	if !found {
		t.Error("expected warning about missing CRS")
	}
}

func TestValidate_UnparsableCRS(t *testing.T) {
	doc := &types.Document{
		Version: "2.0",
		SRSName: "some-unknown-crs",
		CRS:     types.CRS{Raw: "some-unknown-crs"},
		Buildings: []types.Building{{
			ID:                "B1",
			HasMeasuredHeight: true,
			MeasuredHeight:    5,
		}},
	}

	findings := Validate(doc)
	found := false

	for _, f := range findings {
		if f.Severity == SeverityWarning && f.Path == "Document" {
			found = true
		}
	}

	if !found {
		t.Error("expected warning about unparsable CRS")
	}
}

func TestValidate_NoVersion(t *testing.T) {
	doc := &types.Document{
		SRSName: "EPSG:25832",
		CRS:     types.CRS{Raw: "EPSG:25832", Code: 25832},
		Buildings: []types.Building{{
			ID:                "B1",
			HasMeasuredHeight: true,
			MeasuredHeight:    5,
		}},
	}

	findings := Validate(doc)
	found := false

	for _, f := range findings {
		if f.Message == "CityGML version not detected" {
			found = true
		}
	}

	if !found {
		t.Error("expected warning about missing version")
	}
}

func TestValidate_EmptyDocument(t *testing.T) {
	doc := &types.Document{Version: "2.0", SRSName: "EPSG:25832", CRS: types.CRS{Code: 25832}}

	findings := Validate(doc)
	found := false

	for _, f := range findings {
		if f.Message == "document contains no city objects" {
			found = true
		}
	}

	if !found {
		t.Error("expected warning about empty document")
	}
}

func TestValidate_BuildingNoGeometry(t *testing.T) {
	doc := &types.Document{
		Version: "2.0",
		SRSName: "EPSG:25832",
		CRS:     types.CRS{Code: 25832},
		Buildings: []types.Building{{
			ID:                "B1",
			HasMeasuredHeight: true,
			MeasuredHeight:    10,
		}},
	}

	findings := Validate(doc)
	found := false

	for _, f := range findings {
		if f.Message == "building has no geometry" {
			found = true
		}
	}

	if !found {
		t.Error("expected warning about missing geometry")
	}
}

func TestValidate_BuildingNoHeight(t *testing.T) {
	doc := &types.Document{
		Version: "2.0",
		SRSName: "EPSG:25832",
		CRS:     types.CRS{Code: 25832},
		Buildings: []types.Building{{
			ID: "B1",
		}},
	}

	findings := Validate(doc)
	found := false

	for _, f := range findings {
		if f.Message == "no measured height and could not derive height from geometry" {
			found = true
		}
	}

	if !found {
		t.Error("expected warning about missing height")
	}
}

func TestValidate_MalformedGeometry(t *testing.T) {
	doc := &types.Document{
		Version: "2.0",
		SRSName: "EPSG:25832",
		CRS:     types.CRS{Code: 25832},
		Buildings: []types.Building{{
			ID:                "B1",
			HasMeasuredHeight: true,
			MeasuredHeight:    10,
			Solid: &types.Solid{
				Exterior: types.MultiSurface{Polygons: []types.Polygon{{
					Exterior: types.Ring{Points: []types.Point{
						{X: 0, Y: 0, Z: 0},
						{X: 10, Y: 0, Z: 0},
						{X: 5, Y: 5, Z: 0}, // not closed, only 3 points
					}},
				}}},
			},
		}},
	}

	findings := Validate(doc)
	var errors []Finding

	for _, f := range findings {
		if f.Severity == SeverityError {
			errors = append(errors, f)
		}
	}

	if len(errors) == 0 {
		t.Error("expected geometry validation errors")
	}
}

func TestValidate_NonFiniteCoordinates(t *testing.T) {
	doc := &types.Document{
		Version: "2.0",
		SRSName: "EPSG:25832",
		CRS:     types.CRS{Code: 25832},
		Buildings: []types.Building{{
			ID:                "B1",
			HasMeasuredHeight: true,
			MeasuredHeight:    10,
			MultiSurface: &types.MultiSurface{Polygons: []types.Polygon{{
				Exterior: types.Ring{Points: []types.Point{
					{X: 0, Y: 0, Z: 0},
					{X: math.NaN(), Y: 0, Z: 0},
					{X: 10, Y: 10, Z: 0},
					{X: 0, Y: 10, Z: 0},
					{X: 0, Y: 0, Z: 0},
				}},
			}}},
		}},
	}

	findings := Validate(doc)
	var errors []Finding

	for _, f := range findings {
		if f.Severity == SeverityError {
			errors = append(errors, f)
		}
	}

	if len(errors) == 0 {
		t.Error("expected error for non-finite coordinates")
	}
}

func TestValidate_TerrainNoGeometry(t *testing.T) {
	doc := &types.Document{
		Version:  "2.0",
		SRSName:  "EPSG:25832",
		CRS:      types.CRS{Code: 25832},
		Terrains: []types.Terrain{{ID: "T1"}},
	}

	findings := Validate(doc)
	found := false

	for _, f := range findings {
		if f.Message == "terrain has no geometry" {
			found = true
		}
	}

	if !found {
		t.Error("expected warning about terrain with no geometry")
	}
}

func TestValidate_BoundedSurfaceGeometry(t *testing.T) {
	doc := &types.Document{
		Version: "2.0",
		SRSName: "EPSG:25832",
		CRS:     types.CRS{Code: 25832},
		Buildings: []types.Building{{
			ID:                "B1",
			HasMeasuredHeight: true,
			MeasuredHeight:    10,
			BoundedBy: []types.Surface{{
				ID:   "GS1",
				Type: "GroundSurface",
				Geometry: types.MultiSurface{Polygons: []types.Polygon{{
					Exterior: types.Ring{Points: []types.Point{
						{X: 0, Y: 0}, {X: 10, Y: 0}, {X: 0, Y: 0}, // too few points
					}},
				}}},
			}},
		}},
	}

	findings := Validate(doc)
	var errors []Finding

	for _, f := range findings {
		if f.Severity == SeverityError {
			errors = append(errors, f)
		}
	}

	if len(errors) == 0 {
		t.Error("expected geometry error for bounded surface")
	}
}

func TestSeverity_String(t *testing.T) {
	if SeverityWarning.String() != "warning" {
		t.Errorf("got %q", SeverityWarning.String())
	}

	if SeverityError.String() != "error" {
		t.Errorf("got %q", SeverityError.String())
	}
}

func TestFinding_String(t *testing.T) {
	f := Finding{Severity: SeverityError, Path: "Building[0]/Solid", Message: "ring not closed"}

	want := "[error] Building[0]/Solid: ring not closed"
	if got := f.String(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
