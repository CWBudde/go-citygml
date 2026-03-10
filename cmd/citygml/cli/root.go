package cli

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "citygml",
	Short: "CityGML toolkit",
	Long:  "A command-line tool for working with CityGML files.",
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}
