package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/cwbudde/go-citygml/citygml"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate [files...]",
	Short: "Validate CityGML files",
	Long: `Validate one or more CityGML files for structural correctness.

Checks include:
  - Well-formed XML
  - Recognized CityGML namespace and version
  - Valid GML geometry structure
  - Coordinate and ring validity`,
	Args: cobra.MinimumNArgs(1),
	RunE: runValidate,
}

func init() {
	rootCmd.AddCommand(validateCmd)
}

func runValidate(_ *cobra.Command, args []string) error {
	hasErrors := false

	for _, path := range args {
		err := validateFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", path, err)

			hasErrors = true
		}
	}

	if hasErrors {
		return errors.New("validation failed for one or more files")
	}

	return nil
}

func validateFile(path string) error {
	doc, err := citygml.ReadFile(path, citygml.Options{})
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	findings := citygml.Validate(doc)

	errorCount := 0
	warningCount := 0

	for _, f := range findings {
		switch f.Severity {
		case citygml.SeverityError:
			fmt.Fprintf(os.Stderr, "%s: ERROR %s\n", path, f)

			errorCount++
		case citygml.SeverityWarning:
			fmt.Fprintf(os.Stderr, "%s: WARN  %s\n", path, f)

			warningCount++
		}
	}

	if errorCount > 0 {
		fmt.Fprintf(os.Stdout, "%s: FAILED (%d errors, %d warnings)\n", path, errorCount, warningCount)
		return fmt.Errorf("%d errors found", errorCount)
	}

	if warningCount > 0 {
		fmt.Fprintf(os.Stdout, "%s: OK (%d warnings)\n", path, warningCount)
	} else {
		fmt.Fprintf(os.Stdout, "%s: OK\n", path)
	}

	return nil
}
