package helpers

import (
	"testing"

	"github.com/cwbudde/go-citygml/types"
)

func partsDocument() *types.Document {
	pt := func(x, y, z float64) types.Point { return types.Point{X: x, Y: y, Z: z} }

	return &types.Document{Buildings: []types.Building{
		{
			ID: "A",
			Parts: []types.Building{
				{
					ID:           "A1",
					MultiSurface: &types.MultiSurface{Polygons: []types.Polygon{{Exterior: types.Ring{Points: []types.Point{pt(100, 200, 5), pt(110, 200, 25)}}}}},
					Parts:        []types.Building{{ID: "A1a"}},
				},
				{ID: "A2"},
			},
		},
		{ID: "B", MultiSurface: &types.MultiSurface{Polygons: []types.Polygon{{Exterior: types.Ring{Points: []types.Point{pt(0, 0, 10)}}}}}},
	}}
}

func TestFlattenBuildings(t *testing.T) {
	doc := partsDocument()

	got := FlattenBuildings(doc)

	want := []string{"A", "A1", "A1a", "A2", "B"}
	if len(got) != len(want) {
		t.Fatalf("got %d entries, want %d", len(got), len(want))
	}

	for i, b := range got {
		if b.ID != want[i] {
			t.Errorf("entry %d = %q, want %q", i, b.ID, want[i])
		}
	}

	// Entries point into the document, not at copies.
	if got[1] != &doc.Buildings[0].Parts[0] {
		t.Error("FlattenBuildings returned a copy instead of a pointer into the document")
	}
}

func TestFlattenBuildings_Empty(t *testing.T) {
	if got := FlattenBuildings(&types.Document{}); len(got) != 0 {
		t.Errorf("got %d entries, want 0", len(got))
	}
}

func TestDocumentBBox_IncludesParts(t *testing.T) {
	bb := DocumentBBox(partsDocument())

	if bb.Empty {
		t.Fatal("bbox is empty")
	}

	if bb.MaxX != 110 || bb.MaxY != 200 || bb.MinZ != 5 || bb.MaxZ != 25 {
		t.Errorf("bbox = %+v, want MaxX=110 MaxY=200 MinZ=5 MaxZ=25 from the part geometry", bb)
	}
}
