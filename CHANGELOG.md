# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/),
and this project adheres to [Semantic Versioning](https://semver.org/).

## [0.1.0] - 2026-03-11

### Added

- **Core parsing**: `citygml.Read()` and `citygml.ReadFile()` for decoding CityGML 2.0 and 3.0 documents
- **GML geometry**: parsing for `gml:pos`, `gml:posList`, `LinearRing`, `Polygon`, `MultiSurface`, `Solid`, `CompositeSurface`
- **Building support**: extraction of `lod1Solid`, `lod2Solid`, `lod1MultiSurface`, `lod2MultiSurface`, bounded surfaces (GroundSurface, RoofSurface, WallSurface, etc.), class/function/usage attributes, and measuredHeight
- **Height derivation**: automatic height computation from Z extents when measuredHeight is absent
- **Footprint derivation**: 2D footprint extraction from 3D geometry (prefers GroundSurface, falls back to lowest-Z polygon)
- **Terrain support**: ReliefFeature and TINRelief decoding with MultiSurface geometry
- **Generic objects**: unsupported object types (Bridge, Tunnel, Road, etc.) collected as `GenericObject` with ID and type
- **CRS handling**: parsing of `srsName` in EPSG short form, URN form, and HTTP form with axis-order detection
- **Validation**: `citygml.Validate()` API for structural validation with severity levels and object-path context
- **CLI**: `citygml validate` command for validating CityGML files from the command line
- **Helpers package**: `BuildingHeight()`, `BuildingHeights()`, `BuildingFootprints()`, `SummarizeTerrain()`, `DocumentBBox()`
- **GeoJSON package**: `FromDocument()`, `BuildingFeature()`, `TerrainFeature()` for RFC 7946 GeoJSON conversion
- **Error types**: `ErrMalformedXML`, `ErrUnsupportedVersion`, `ErrUnsupportedGeometry`, `ErrUnsupportedObject`, `ErrInvalidCoordinates`, `ErrInvalidCRS`
- **Strict mode**: `Options{Strict: true}` for rejecting unsupported objects and missing CRS
- **Test fixtures**: synthetic CityGML files covering measured height, Z extents, multiple buildings, unsupported objects, malformed geometry, namespace variations
- **Snapshot tests**: golden-file tests for normalized parser output
- **Fuzz tests**: fuzz targets for `ParsePos`, `ParsePosList`, and `ValidateRing`
- **Benchmarks**: performance tests from 1 to 5000 buildings with allocation tracking
