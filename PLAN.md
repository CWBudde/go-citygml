# PLAN.md — `go-citygml`

Status: draft, 10 March 2026

This plan defines a standalone Go library for reading and normalizing CityGML data without tying the parser to Aconiq-specific model semantics.

## Why this should exist

- CityGML parsing, namespace/version handling, and geometry extraction are a distinct problem from noise-model import.
- A dedicated library keeps XML/GML complexity out of Aconiq.
- The library can expose generic building/terrain/object abstractions that multiple tools can reuse.
- Version/profile support, validation, and test fixtures can evolve independently from application logic.

## Design goals

- Pure Go library, no CLI-first coupling
- Deterministic parsing and object ordering
- Explicit support boundaries by CityGML version, object type, and LoD
- Clear separation between XML parsing, GML geometry decoding, semantic object mapping, and validation
- Repo-safe fixtures and exhaustive tests

## Non-goals for v1

- Full CityGML conformance across every ADE/profile
- 3D rendering or visualization helpers
- CRS reprojection engine
- CityJSON support in the initial delivery

---

## Phase 0 — Repository foundation

- [x] Create Go module `github.com/cwbudde/go-citygml` (or final chosen path)
- [x] Add README with scope, supported versions, and non-goals
- [x] Add CI for formatting, tests, linting
- [x] Define contribution, versioning, and compatibility policy

---

## Phase 1 — Core API design

- [x] Define the public package layout
  - [x] `citygml/` high-level decode API
  - [x] `gml/` geometry parsing helpers
  - [x] `internal/xmlscan/` low-level token handling
  - [x] `types/` normalized semantic model or equivalent
- [x] Define the top-level decode API
  - [x] `Read(io.Reader, Options) (*Document, error)`
  - [x] `ReadFile(path string, Options) (*Document, error)`
  - [x] streaming/token-oriented API if needed
- [x] Define stable error types
  - [x] malformed XML
  - [x] unsupported CityGML version/profile
  - [x] unsupported geometry/object types
  - [x] invalid coordinate dimensionality / CRS metadata

---

## Phase 2 — XML and namespace foundation

- [x] Implement namespace-aware XML token scanning
- [x] Support the minimum namespace set for CityGML 2.0 and 3.0 detection
- [x] Detect document version/profile from namespaces and root structure
- [x] Preserve object IDs and xlink targets where present
- [x] Add deterministic traversal order for city object members and nested elements

---

## Phase 3 — GML geometry core

- [x] Parse `gml:pos`
- [x] Parse `gml:posList`
- [x] Parse `gml:LinearRing`
- [x] Parse `gml:Polygon`
- [x] Parse `gml:MultiSurface`
- [x] Parse `gml:Solid` / `CompositeSurface` for the initial supported scope
- [x] Track dimensionality (`2D`, `3D`) explicitly
- [x] Add geometry validation
  - [x] ring closure
  - [x] minimum coordinate counts
  - [x] finite numeric values
  - [x] supported dimensionality

---

## Phase 4 — Normalized semantic model

- [x] Define library-owned normalized types
  - [x] `Document`
  - [x] `Building`
  - [x] `Surface`
  - [x] `Terrain`
  - [x] generic `CityObject`
- [x] Define what metadata is preserved
  - [x] IDs
  - [x] class / function / usage
  - [x] measured height
  - [x] LoD markers
  - [x] CRS metadata
- [x] Define how raw source geometry relates to derived footprints / heights

---

## Phase 5 — Buildings v1

- [x] Support building extraction as the first shippable semantic object
- [x] Support common building geometry carriers
  - [x] `lod1Solid`
  - [x] `lod1MultiSurface`
  - [x] bounded surfaces where useful
- [x] Extract measured height when present
- [x] Derive height from Z extents when measured height is absent
- [x] Derive 2D footprint candidates deterministically from 3D geometry
- [x] Document exactly which building patterns are supported vs skipped

---

## Phase 6 — Terrain and context objects

- [x] Add terrain surface extraction
- [x] Add support for bridge / tunnel / transportation objects if kept in scope
- [x] Decide whether these belong in v1 or a later minor release
- [x] Document object-specific support boundaries clearly

---

## Phase 7 — CRS and axis-order handling

- [x] Parse `srsName` declarations from relevant geometry carriers
- [x] Define supported CRS declaration forms
- [x] Define axis-order normalization rules
- [x] Add error/warning behavior for missing CRS metadata
- [x] Decide whether the library only preserves CRS metadata or also reprojects

---

## Phase 8 — Validation layer

- [x] Add structural validation API separate from decode
- [x] Report unsupported but recoverable constructs as warnings
- [x] Report malformed required geometry/object structures as errors
- [x] Include object-path context in validation findings

---

## Phase 9 — Application integration helpers

- [ ] Add optional helper package(s) for downstream mapping
  - [ ] building footprints
  - [ ] height extraction
  - [ ] terrain mesh / polygon summaries
- [ ] Keep these helpers generic and not Aconiq-specific
- [ ] Decide whether GeoJSON conversion helpers belong here or in downstream apps

---

## Phase 10 — Test corpus and fixtures

- [ ] Add repo-safe synthetic CityGML fixtures
- [ ] Add fixtures for
  - [ ] measured height present
  - [ ] height from Z extents
  - [ ] multiple buildings
  - [ ] unsupported object types
  - [ ] malformed rings / malformed `posList`
  - [ ] namespace/version variation
- [ ] Add snapshot tests for normalized object output
- [ ] Add property/fuzz tests for coordinate parsing and geometry robustness

---

## Phase 11 — Performance and memory

- [ ] Benchmark large XML inputs
- [ ] Avoid unnecessary DOM-style loading where possible
- [ ] Evaluate streaming decode boundaries
- [ ] Document expected memory behavior for large files

---

## Phase 12 — Release readiness

- [ ] Finalize public API stability for `v0.x` or `v1`
- [ ] Publish usage examples
- [ ] Add changelog and migration notes
- [ ] Tag the first release

---

## Integration with Aconiq

- [ ] Replace Aconiq’s local `citygmlimport` implementation with this library once the building scope is feature-complete
- [ ] Keep Aconiq-specific mapping from generic city objects into the normalized noise model inside Aconiq
- [ ] Add an adapter layer instead of leaking Aconiq types into this library

---

## Suggested first delivery

The first worthwhile release should likely be:

- CityGML building-only support
- common `2.0` namespaces first
- `lod1Solid` / `lod1MultiSurface`
- measured height + Z-extent fallback
- deterministic footprint extraction
- repo-safe fixtures and strong validation errors

That is enough to let Aconiq depend on the library without forcing premature support for the entire CityGML ecosystem.
