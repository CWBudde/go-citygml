# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/),
and this project adheres to [Semantic Versioning](https://semver.org/).

## [Unreleased]

## [0.3.0] - 2026-09-25

### Added

- **CityGML 1.0**: documents in the CityGML 1.0 namespaces (core, building, relief, transportation, vegetation, generics, appearance, cityobjectgroup) are decoded like 2.0 and report `Version` `"1.0"`. This covers AdV/LGLN LoD2 tiles, which previously yielded no buildings.
- **BuildingPart**: `bldg:consistsOfBuildingPart` (1.0/2.0) and `bldg:buildingPart` (3.0) are decoded recursively into the new `types.Building.Parts` field. Each part carries its own ID, attributes, measured or derived height, geometry, semantic surfaces and footprint.
- `helpers.FlattenBuildings` returns every building and building part depth-first.
- **Compound CRS**: `EPSG:25832+7837`, `urn:ogc:def:crs,crs:EPSG::25832,crs:EPSG::7837` and `http://www.opengis.net/def/crs-compound?1=…&2=…` resolve to the horizontal EPSG code.
- **srsDimension**: an `srsDimension` on a `posList`/`pos` or on an enclosing geometry now decides the coordinate dimensionality; the root envelope's value is used as a hint. Guessing from divisibility by 3 is only the last fallback.

### Changed

- `Validate` no longer warns about missing height or geometry on a building that has parts, and validates each part under `Building[i](id)/Part[j](id)`.
- `helpers.DocumentBBox` includes building-part geometry.
- `geojson.FromDocument` emits one feature per building part (type `BuildingPart`, with a `parent` property) after its parent's feature.
- Golden snapshots include the new `Parts` field.

### Fixed

- `DeriveFootprint` returns the largest-area polygon over all GroundSurface polygons (holes included) instead of the first polygon of the first GroundSurface.

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
