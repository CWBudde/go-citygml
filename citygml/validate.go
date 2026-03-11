package citygml

import (
	"errors"
	"fmt"

	"github.com/cwbudde/go-citygml/gml"
	"github.com/cwbudde/go-citygml/types"
)

// Severity indicates the severity of a validation finding.
type Severity int

const (
	// SeverityWarning indicates a recoverable or informational issue.
	SeverityWarning Severity = iota
	// SeverityError indicates a structural problem that likely affects correctness.
	SeverityError
)

func (s Severity) String() string {
	switch s {
	case SeverityWarning:
		return "warning"
	case SeverityError:
		return "error"
	default:
		return "unknown"
	}
}

// Finding represents a single validation result with context.
type Finding struct {
	Severity Severity
	Path     string
	Message  string
}

func (f Finding) String() string {
	return fmt.Sprintf("[%s] %s: %s", f.Severity, f.Path, f.Message)
}

// Validate performs structural validation on a parsed Document.
// It returns a list of findings (warnings and errors) without modifying the document.
// This is separate from the decode step and can be called after Read().
func Validate(doc *types.Document) []Finding {
	findings := make([]Finding, 0, len(doc.Buildings)+len(doc.Terrains)+3)

	findings = append(findings, validateDocumentMeta(doc)...)

	for i := range doc.Buildings {
		findings = append(findings, validateBuilding(&doc.Buildings[i], fmt.Sprintf("Building[%d](%s)", i, doc.Buildings[i].ID))...)
	}

	for i := range doc.Terrains {
		findings = append(findings, validateTerrain(&doc.Terrains[i], fmt.Sprintf("Terrain[%d](%s)", i, doc.Terrains[i].ID))...)
	}

	return findings
}

func validateDocumentMeta(doc *types.Document) []Finding {
	var findings []Finding

	if doc.Version == "" {
		findings = append(findings, Finding{
			Severity: SeverityWarning,
			Path:     "Document",
			Message:  "CityGML version not detected",
		})
	}

	if doc.SRSName == "" {
		findings = append(findings, Finding{
			Severity: SeverityWarning,
			Path:     "Document",
			Message:  "no srsName (CRS) declared",
		})
	} else if doc.CRS.Code == 0 {
		findings = append(findings, Finding{
			Severity: SeverityWarning,
			Path:     "Document",
			Message:  fmt.Sprintf("srsName %q could not be parsed to an EPSG code", doc.SRSName),
		})
	}

	if len(doc.Buildings) == 0 && len(doc.Terrains) == 0 && len(doc.GenericObjects) == 0 {
		findings = append(findings, Finding{
			Severity: SeverityWarning,
			Path:     "Document",
			Message:  "document contains no city objects",
		})
	}

	return findings
}

func validateBuilding(b *types.Building, path string) []Finding {
	var findings []Finding

	if !b.HasMeasuredHeight && b.DerivedHeight == 0 {
		findings = append(findings, Finding{
			Severity: SeverityWarning,
			Path:     path,
			Message:  "no measured height and could not derive height from geometry",
		})
	}

	if b.Solid == nil && b.MultiSurface == nil && len(b.BoundedBy) == 0 {
		findings = append(findings, Finding{
			Severity: SeverityWarning,
			Path:     path,
			Message:  "building has no geometry",
		})
	}

	if b.Solid != nil {
		findings = append(findings, validateGeometryErrors(
			gml.ValidateSolid(*b.Solid, path+"/Solid"),
		)...)
	}

	if b.MultiSurface != nil {
		findings = append(findings, validateGeometryErrors(
			gml.ValidateMultiSurface(*b.MultiSurface, path+"/MultiSurface"),
		)...)
	}

	for i, surf := range b.BoundedBy {
		surfPath := fmt.Sprintf("%s/BoundedBy[%d](%s)", path, i, surf.Type)
		findings = append(findings, validateGeometryErrors(
			gml.ValidateMultiSurface(surf.Geometry, surfPath),
		)...)
	}

	return findings
}

func validateTerrain(t *types.Terrain, path string) []Finding {
	var findings []Finding

	if len(t.Geometry.Polygons) == 0 {
		findings = append(findings, Finding{
			Severity: SeverityWarning,
			Path:     path,
			Message:  "terrain has no geometry",
		})
	}

	findings = append(findings, validateGeometryErrors(
		gml.ValidateMultiSurface(t.Geometry, path+"/Geometry"),
	)...)

	return findings
}

// validateGeometryErrors converts gml validation errors into Findings with Error severity.
func validateGeometryErrors(errs []error) []Finding {
	findings := make([]Finding, 0, len(errs))
	for _, err := range errs {
		ve := &gml.ValidationError{}
		if errors.As(err, &ve) {
			findings = append(findings, Finding{
				Severity: SeverityError,
				Path:     ve.Path,
				Message:  ve.Message,
			})
		} else {
			findings = append(findings, Finding{
				Severity: SeverityError,
				Path:     "",
				Message:  err.Error(),
			})
		}
	}

	return findings
}
