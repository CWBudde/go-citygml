package cli

import (
	"errors"
	"fmt"
	"os"

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

func runValidate(cmd *cobra.Command, args []string) error {
	hasErrors := false

	for _, path := range args {
		err := validateFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", path, err)

			hasErrors = true
		} else {
			fmt.Printf("%s: ok\n", path)
		}
	}

	if hasErrors {
		return errors.New("validation failed for one or more files")
	}

	return nil
}

func validateFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// TODO: Wire up citygml.Read() + validation once the library core is implemented.
	return errors.New("validation not yet implemented")
}
