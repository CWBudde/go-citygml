package citygml

import (
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var updateGolden = flag.Bool("update-golden", false, "update golden snapshot files")

func TestSnapshot_MeasuredHeight(t *testing.T) {
	testSnapshot(t, "../testdata/building_measured_height.gml", "../testdata/golden/building_measured_height.json")
}

func TestSnapshot_ZExtents(t *testing.T) {
	testSnapshot(t, "../testdata/building_z_extents.gml", "../testdata/golden/building_z_extents.json")
}

func TestSnapshot_MultipleBuildings(t *testing.T) {
	testSnapshot(t, "../testdata/multiple_buildings.gml", "../testdata/golden/multiple_buildings.json")
}

func TestSnapshot_UnsupportedObjects(t *testing.T) {
	testSnapshot(t, "../testdata/unsupported_objects.gml", "../testdata/golden/unsupported_objects.json")
}

func TestSnapshot_CityGML30(t *testing.T) {
	testSnapshot(t, "../testdata/citygml30_building.gml", "../testdata/golden/citygml30_building.json")
}

func TestSnapshot_CityGML10BuildingParts(t *testing.T) {
	testSnapshot(t, "../testdata/citygml10_building_parts.gml", "../testdata/golden/citygml10_building_parts.json")
}

func TestSnapshot_URN_CRS(t *testing.T) {
	testSnapshot(t, "../testdata/citygml20_urn_crs.gml", "../testdata/golden/citygml20_urn_crs.json")
}

func testSnapshot(t *testing.T, inputPath, goldenPath string) {
	t.Helper()

	absInput, err := filepath.Abs(inputPath)
	if err != nil {
		t.Fatal(err)
	}

	absGolden, err := filepath.Abs(goldenPath)
	if err != nil {
		t.Fatal(err)
	}

	doc, err := ReadFile(absInput, Options{})
	if err != nil {
		t.Fatalf("ReadFile(%s): %v", inputPath, err)
	}

	got, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	if *updateGolden {
		dir := filepath.Dir(absGolden)

		err := os.MkdirAll(dir, 0o755)
		if err != nil {
			t.Fatal(err)
		}

		err = os.WriteFile(absGolden, got, 0o644)
		if err != nil {
			t.Fatal(err)
		}

		t.Logf("updated golden file: %s", absGolden)

		return
	}

	want, err := os.ReadFile(absGolden)
	if err != nil {
		t.Fatalf("reading golden file (run with -update-golden to create): %v", err)
	}

	// Compare newline-tolerantly: json.MarshalIndent emits no trailing
	// newline, but editors and formatters routinely add one to committed
	// golden files. Trimming trailing whitespace keeps the test robust
	// against that without weakening the structural comparison.
	if strings.TrimRight(string(got), "\n") != strings.TrimRight(string(want), "\n") {
		t.Errorf("snapshot mismatch for %s\ngot:\n%s\nwant:\n%s", inputPath, got, want)
	}
}
