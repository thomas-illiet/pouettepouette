package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Version holds the CLI version, set during build.
	Version = "dev"
)

// versionCmd prints the CLI version.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the CLI version",
	Long: "Prints the current version of the CLI. " +
		"This is useful for debugging or verifying the installed version.",
	Args: cobra.NoArgs, // No arguments are expected
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(Version)
		return nil
	},
}
