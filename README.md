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

### CLI

A command-line tool is included for validating CityGML files:

```
go install github.com/cwbudde/go-citygml/cmd/citygml@latest

citygml validate building.gml
```

## Non-goals for v1

- Full CityGML conformance across every ADE/profile
- 3D rendering or visualization helpers
- CRS reprojection engine
- CityJSON support

## License

See [LICENSE](LICENSE) for details.
