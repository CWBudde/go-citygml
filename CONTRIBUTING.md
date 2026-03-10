# Contributing to go-citygml

## Getting started

1. Clone the repository
2. Ensure you have Go 1.23+ installed
3. Run `go test ./...` to verify everything works

## Development workflow

- Format code with `gofmt` (CI enforces this)
- Run `go vet ./...` before committing
- Write tests for new functionality
- Keep commits focused and well-described

## Pull requests

- Open PRs against `main`
- Include tests for new features and bug fixes
- Ensure CI passes before requesting review
- Keep PRs focused on a single change

## Versioning policy

This project follows [Semantic Versioning](https://semver.org/):

- **v0.x** releases may include breaking API changes between minor versions
- **v1.0** will mark a stable public API with the following guarantees:
  - Patch releases: bug fixes only, no API changes
  - Minor releases: additive changes only (new types, functions, fields)
  - Major releases: breaking changes with migration notes

## Compatibility

- The library targets the two most recent stable Go releases
- CityGML version support is explicitly documented in the README
- Unsupported CityGML constructs produce clear errors rather than silent failures

## Reporting issues

- Include the CityGML input (or a minimal reproduction) when reporting parsing issues
- Include the Go version and library version
