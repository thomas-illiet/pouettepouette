package workspace

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "ports",
	Short: "Interact with workspace ports.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		return nil
	},
}

func init() {

}
