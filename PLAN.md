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

- [x] Add optional helper package(s) for downstream mapping
  - [x] building footprints
  - [x] height extraction
  - [x] terrain mesh / polygon summaries
- [x] Keep these helpers generic and not Aconiq-specific
- [x] Decide whether GeoJSON conversion helpers belong here or in downstream apps
  - Decision: included as `geojson/` package — lightweight, no external deps, useful for quick visualization

---

## Phase 10 — Test corpus and fixtures

- [x] Add repo-safe synthetic CityGML fixtures
- [x] Add fixtures for
  - [x] measured height present
  - [x] height from Z extents
  - [x] multiple buildings
  - [x] unsupported object types
  - [x] malformed rings / malformed `posList`
  - [x] namespace/version variation
- [x] Add snapshot tests for normalized object output
- [x] Add property/fuzz tests for coordinate parsing and geometry robustness

---

## Phase 11 — Performance and memory

- [x] Benchmark large XML inputs
- [x] Avoid unnecessary DOM-style loading where possible
  - Already streaming: forward-only xml.Decoder, no DOM tree built
- [x] Evaluate streaming decode boundaries
  - Objects collected into Document struct; per-object streaming not needed for v1
  - Throughput ~45-48 MB/s, linear scaling confirmed up to 5000 buildings
- [x] Document expected memory behavior for large files

---

## Phase 12 — Release readiness

- [x] Finalize public API stability for `v0.x` or `v1`
  - Decision: release as `v0.1.0` (v0.x allows breaking changes per CONTRIBUTING.md)
- [x] Publish usage examples
- [x] Add changelog and migration notes
- [x] Tag the first release (`v0.1.0`)

---

## Phase 13 — Web demo (WASM + GitHub Pages)

- [ ] Create WASM entry point (`cmd/citygmlwasm/main.go`)
  - [ ] Build tag `//go:build js && wasm`
  - [ ] Register `parseCityGML(Uint8Array) → JSON` on `js.Global()`
  - [ ] Parse CityGML via `citygml.Read()` with height/footprint derivation
  - [ ] Run `citygml.Validate()` for findings
  - [ ] Convert to GeoJSON via `geojson.FromDocument()`
  - [ ] Compute bounding box via `helpers.BoundingBox()`, normalize 3D coords to centroid-origin
  - [ ] Return unified JSON: `geojson`, `scene`, `meta`, `objects`, `findings`, `bounds`
- [ ] Create web UI (`web/`)
  - [ ] `index.html` — single-page app, loads WASM + JS modules
  - [ ] `style.css` — dark/light theme, drop zone styling, responsive layout
  - [ ] `app.js` — WASM init, drop zone handlers, view transitions, state management
  - [ ] `modules/map.js` — MapLibre GL JS 2D map (buildings colored by height, terrain in green, click popups, auto-fit bounds)
  - [ ] `modules/scene.js` — Three.js 3D view (surfaces colored by type, OrbitControls, wireframe toggle, click-to-select)
  - [ ] `modules/sidebar.js` — tabs (Summary, Objects, Validation), object list with click-to-highlight
- [ ] UI flow
  - [ ] Landing: full-screen drop zone (drag-drop or click, accepts `.gml`/`.xml`/`.citygml`)
  - [ ] After parse: fade to visualization layout (70% map/3D + 30% sidebar)
  - [ ] Toggle button switches between 2D map and 3D scene
  - [ ] "New File" button returns to drop zone
  - [ ] Clicking object in sidebar highlights in view and vice versa
- [ ] GitHub Actions deployment (`.github/workflows/pages.yml`)
  - [ ] Trigger on push to `main`
  - [ ] Build WASM: `GOOS=js GOARCH=wasm go build -o web/citygml.wasm ./cmd/citygmlwasm`
  - [ ] Copy `wasm_exec.js` from Go runtime
  - [ ] Deploy `web/` to GitHub Pages
- [ ] Add `build-wasm` recipe to justfile

---

## Phase 14 — Quality & correctness remediation

Added after a full repository quality review (Aug 2026). Findings are grouped by
priority. Every item cites the file(s) to change. Check items off as they land.

### P0 — CI is red on `main`; unblock it first

The `Tests` workflow has failed on every push to `main` since March 2026. Two
jobs fail, and they impose contradictory requirements on the golden files, so
neither can be fixed in isolation:

- [ ] **Snapshot tests fail on a clean checkout.** `citygml/snapshot_test.go`
      compares `json.MarshalIndent` output (no trailing newline) byte-for-byte
      against `testdata/golden/*.json`, but `treefmt.toml`'s prettier step
      formats `*.json` and appends a trailing newline. Fix by excluding
      `testdata/golden/**` (or `testdata/**`) from prettier in `treefmt.toml`,
      then regenerate the goldens by running the snapshot test with the
      `-update-golden` flag. (Belt-and-braces: also `strings.TrimRight` in the
      comparison so the test is newline-tolerant.)
- [ ] **Format check fails.** `.golangci.yml` is not prettier-formatted (last
      edited by `98928ff` without a formatter run). Run `treefmt` / `just fmt`
      and commit the result.
- [ ] Verify locally that `just ci` (`check-formatted test lint check-tidy`)
      passes end-to-end before pushing.

### P1 — Correctness bugs (silent data corruption / crashes)

- [ ] **2D posList decoded as 3D.** `gml/coords.go:80 inferDimensionality`
      guesses 3D whenever the value count is divisible by 3; `srsDimension` is
      never read anywhere in the codebase. A 2D ring with a 3-divisible count is
      silently corrupted. Read the `srsDimension` attribute on
      `gml:pos`/`gml:posList` and thread it through `gml/parse.go` instead of
      guessing.
- [ ] **`cityObjectMember` sibling-skip.** `internal/xmlscan/document.go:92`
      calls `StartElement()`, which scans forward across element boundaries; an
      empty or xlink-only member consumes the _next_ member's start tag. Bound
      the scan to the current member's depth.
- [ ] **WASM entry point unsafe on bad input.** `cmd/citygmlwasm/main.go`:
      (a) add `defer/recover` in `parseCityGML` — one panic kills the instance
      for the page's lifetime; (b) validate `args[0]` is a `Uint8Array` before
      `ua.Get("length").Int()`; (c) guard `bbox.Empty` before emitting
      `bounds`, which currently leak `±math.MaxFloat64` for empty documents.
- [ ] **`measuredHeight` accepts NaN/Inf/negative.**
      `internal/decode/building.go:77` uses `strconv.ParseFloat` with no
      finiteness check (unlike `gml/coords.go`'s `parseFloat`). Reject
      non-finite values.
- [ ] **DOM-XSS in the web viewer.** `web/modules/map.js:146-163` interpolates
      unescaped user-file values (`id`, `class`, `function`, `lod`, …) into
      popup HTML via `setHTML`. `sidebar.js` already routes through
      `escapeHtml()`; map.js must too (or build nodes with `textContent`).

### P2 — API, conformance & correctness-adjacent

- [ ] **GeoJSON is not RFC 7946 conformant despite `doc.go`'s claim.**
      `geojson/geojson.go`: (a) a `Solid` shell is emitted as a 2D
      MultiPolygon (`:131`) — overlapping/degenerate polygons; derive a
      footprint instead; (b) coordinates are emitted in the source projected
      CRS, not WGS84, and `types.CRS.IsYXOrder` is ignored (`ringCoords`).
      Either reproject / honour axis order, or soften the "RFC 7946 compliant"
      claim and document the limitation.
- [ ] **Inconsistent `Options` API.** `citygml/options.go` mixes `bool`
      (`Strict`) and `*bool` (`DeriveHeights`, `DeriveFootprints`). Replace with
      inverted plain bools or functional options; document the partial-doc-on-
      error contract of `Read` (`read.go:86`).
- [ ] **CLI is thin and noisy.** `cmd/citygml/cli`: set
      `SilenceUsage`/`SilenceErrors` (validation failures currently dump full
      usage); remove double parse/validation error reporting; expose
      `--strict`/`--json` so library capabilities are reachable.
- [ ] **De-duplicate drifting logic.** Height selection is reimplemented three
      times (`helpers.BuildingHeight`, `geojson.BuildingFeature`, wasm
      `buildScene`/`buildObjects`) with subtly different rules (`>0` vs
      `HasMeasuredHeight`); geometry flattening is copied across `helpers`,
      `gml/derive`, and wasm. Collapse to shared helpers.
- [ ] **Extract the copy-pasted "drain-to-matching-end" loop** (~7 sites in
      `gml/parse.go` + `internal/decode`) into one documented `Scanner` helper;
      the manual `depth++/depth--` convention is the core fragility of the
      parser.
- [ ] **Three.js resource leaks.** `web/modules/scene.js`: `destroyScene` never
      disposes geometries/materials, never `disconnect()`s the `ResizeObserver`,
      and re-adds the canvas click listener on every re-init.
- [ ] **Unsupported CRS blanks the map silently.** `web/modules/map.js`
      reprojection only handles ETRS89/WGS84 UTM; anything else yields an empty
      map with no message. Surface a "cannot locate — unsupported CRS" notice.

### P3 — Tests, tooling, docs & hygiene

- [ ] **Zero coverage on `cmd/`** (CLI + WASM, the actual deliverables). Add
      cobra command tests (`SetArgs`/`SetOut`, exit codes) and extract+test the
      WASM conversion logic.
- [ ] **Dead malformed fixtures.** `testdata/malformed_ring.gml` and
      `malformed_poslist.gml` are referenced by no test. Wire them into
      error-path tests (covers the least-tested parser branches, ~65% in
      `internal/decode`) or delete them.
- [ ] **Add a coverage gate** (`-coverprofile` + floor, or Codecov) to
      `test-unit.yml`; nothing in CI observes the 78% coverage today.
- [ ] **Pin supply chain.** SHA-pin third-party actions
      (`softprops/action-gh-release`, `isbecker/treefmt-action`,
      `golangci/golangci-lint-action`) and pin `gofumpt/gci/shfmt/prettier`
      versions in `test-format.yml` (currently `@latest`). Add SRI hashes (or
      vendor) for the maplibre/three CDN scripts in `web/index.html`.
- [ ] **Module-path casing.** `go.mod` declares `github.com/cwbudde/...`
      (lowercase) but the canonical repo is `github.com/CWBudde/...`;
      `go install github.com/CWBudde/...@latest` fails. Pick one casing across
      `go.mod`, README, and LICENSE.
- [ ] **Release/version metadata is fictional.** No git tag exists, yet
      CHANGELOG and PLAN claim `v0.1.0` shipped and `release-binaries.yml` only
      fires on tags — so no release binary was ever built. Either tag `v0.1.0`
      or stop claiming it.
- [ ] **Docs drift.** Check off Phase 13 below (WASM demo is fully implemented);
      fix the `helpers.BoundingBox()` → `DocumentBBox()` reference at
      `PLAN.md:204`; update CHANGELOG for the WASM/web work; correct
      `CONTRIBUTING.md` (says Go 1.23+, repo requires 1.25; documents
      `gofmt`/`go vet` instead of the real `just`/`treefmt`/`golangci` flow).
- [ ] **Missing community files:** `SECURITY.md`, `CODE_OF_CONDUCT.md`, issue/PR
      templates.
- [ ] **Restrict `release-binaries.yml`** to `v*` tags instead of `"*"`.
- [ ] **Web a11y:** the drop zone (`div` + click), view-toggle, and sidebar
      tabs/list rows are not keyboard- or screen-reader-accessible.

### Deferred / decisions needed (not obviously "fix")

- Whether to actually reproject GeoJSON to WGS84 (contradicts the stated v1
  non-goal "CRS reprojection engine") or just document the limitation.
- Whether the `Dimensionality` return triple threaded through `gml/parse.go`
  (discarded by every caller) should be consumed or removed.
- Whether to prune dead API surface: `xmlscan.CityObjectMember`/`RawXML`,
  `Element.XLinkHref` (extracted, never read — no xlink resolution exists).

---

## Integration with Aconiq

- [x] Replace Aconiq’s local `citygmlimport` implementation with this library once the building scope is feature-complete
- [x] Keep Aconiq-specific mapping from generic city objects into the normalized noise model inside Aconiq
- [x] Add an adapter layer instead of leaking Aconiq types into this library
