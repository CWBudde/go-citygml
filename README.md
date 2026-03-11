# go-citygml

A pure Go library for reading and normalizing [CityGML](https://www.ogc.org/standard/citygml/) data.

## What it does

- Parses CityGML XML with full namespace awareness
- Decodes GML geometry (pos, posList, LinearRing, Polygon, MultiSurface, Solid)
- Maps semantic objects (buildings, terrain, surfaces) into a normalized Go model
- Extracts heights (measured or Z-extent derived) and 2D footprints
- Validates geometry and structure with clear error reporting

## Supported versions

| CityGML version | Status    |
| --------------- | --------- |
| 2.0             | Supported |
| 3.0             | Supported |

## Supported building patterns

| Pattern                                                          | Status            |
| ---------------------------------------------------------------- | ----------------- |
| `lod1Solid`                                                      | Supported         |
| `lod2Solid`                                                      | Supported         |
| `lod1MultiSurface`                                               | Supported         |
| `lod2MultiSurface`                                               | Supported         |
| Bounded surfaces (GroundSurface, RoofSurface, WallSurface, etc.) | Supported         |
| `measuredHeight`                                                 | Supported         |
| Height from Z extents (fallback)                                 | Supported         |
| 2D footprint derivation                                          | Supported         |
| `class`, `function`, `usage` attributes                          | Supported         |
| BuildingPart                                                     | Not yet supported |
| lod3/lod4 geometry                                               | Not yet supported |
| BuildingInstallation                                             | Not yet supported |

## Object type support

| Object type                             | Status                                                  |
| --------------------------------------- | ------------------------------------------------------- |
| **Building**                            | Fully decoded (geometry, attributes, height, footprint) |
| **Terrain** (ReliefFeature / TINRelief) | Decoded (MultiSurface geometry)                         |
| Bridge                                  | Collected as GenericObject (ID + type only)             |
| Tunnel                                  | Collected as GenericObject (ID + type only)             |
| Transportation (Road, etc.)             | Collected as GenericObject (ID + type only)             |
| Vegetation                              | Collected as GenericObject (ID + type only)             |
| WaterBody                               | Collected as GenericObject (ID + type only)             |
| CityFurniture                           | Collected as GenericObject (ID + type only)             |
| LandUse                                 | Collected as GenericObject (ID + type only)             |

Unsupported object types are silently collected as `GenericObject` (with ID and type name) in non-strict mode. In strict mode (`Options{Strict: true}`), they produce an error.

## Installation

```
go get github.com/cwbudde/go-citygml
```

## Usage

### As a library

```go
import (
    "os"
    "fmt"
    "github.com/cwbudde/go-citygml/citygml"
)

doc, err := citygml.ReadFile("building.gml", citygml.Options{})
if err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
}

for _, b := range doc.Buildings {
    fmt.Printf("Building %s: height=%.1f\n", b.ID, b.MeasuredHeight)
}
```

### Validation

```go
import "github.com/cwbudde/go-citygml/citygml"

doc, _ := citygml.ReadFile("building.gml", citygml.Options{})

findings := citygml.Validate(doc)
for _, f := range findings {
    fmt.Println(f) // e.g. "[error] Building[0]/Solid: ring not closed"
}
```

### GeoJSON export

```go
import (
    "encoding/json"
    "github.com/cwbudde/go-citygml/citygml"
    "github.com/cwbudde/go-citygml/geojson"
)

doc, _ := citygml.ReadFile("city.gml", citygml.Options{})
fc := geojson.FromDocument(doc)

data, _ := json.MarshalIndent(fc, "", "  ")
os.WriteFile("output.geojson", data, 0o644)
```

### Helper utilities

```go
import "github.com/cwbudde/go-citygml/helpers"

// Bounding box over all geometry
bbox := helpers.DocumentBBox(doc)
fmt.Printf("Extent: [%.1f, %.1f] to [%.1f, %.1f]\n", bbox.MinX, bbox.MinY, bbox.MaxX, bbox.MaxY)

// Building heights (measured or derived)
for _, h := range helpers.BuildingHeights(doc) {
    fmt.Printf("%s: %.1fm (measured=%v)\n", h.ID, h.Height, h.IsMeasured)
}

// Terrain summary
ts := helpers.SummarizeTerrain(doc)
fmt.Printf("%d terrain objects, %d polygons\n", ts.Count, ts.TotalPolygons)
```

### Strict mode

```go
// Reject unsupported objects and require CRS metadata
doc, err := citygml.ReadFile("building.gml", citygml.Options{Strict: true})
if errors.Is(err, citygml.ErrUnsupportedObject) {
    // file contains object types not fully supported
}
```

### CLI

A command-line tool is included for validating CityGML files:

```
go install github.com/cwbudde/go-citygml/cmd/citygml@latest

citygml validate building.gml
```

## Performance

The parser uses forward-only streaming via Go's `encoding/xml.Decoder` — no DOM tree is built. Memory usage scales linearly with the number of city objects, not with raw XML size.

| Input size      | Time    | Throughput | Memory  |
| --------------- | ------- | ---------- | ------- |
| 1 building      | ~36 µs  | ~45 MB/s   | ~12 KB  |
| 100 buildings   | ~2.5 ms | ~48 MB/s   | ~735 KB |
| 1,000 buildings | ~30 ms  | ~45 MB/s   | ~7.3 MB |
| 5,000 buildings | ~130 ms | ~47 MB/s   | ~42 MB  |

Memory per building is approximately 7 KB (including geometry, attributes, and derived footprint). Disabling height/footprint derivation (`DeriveHeights`/`DeriveFootprints` options) reduces allocations slightly.

For very large files, consider processing in batches at the application level since the library currently collects all objects into a single `Document` struct.

## Non-goals for v1

- Full CityGML conformance across every ADE/profile
- 3D rendering or visualization helpers
- CRS reprojection engine
- CityJSON support

## License

See [LICENSE](LICENSE) for details.
