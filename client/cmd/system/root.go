package system

import (
	"github.com/spf13/cobra"
)

var jsonFormat bool
var noColor bool

// Cmd represents the "system" command used to retrieve system information.
var Cmd = &cobra.Command{
	Use:   "system",
	Short: "Retrieve system information",
	Long: "The system command is used to retrieve various system information " +
		"such as OS details, resources, and configurations.",
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		return nil
	},
}

func init() {
	Cmd.AddCommand(InfoCmd)
	Cmd.AddCommand(ResourceCmd)
}
