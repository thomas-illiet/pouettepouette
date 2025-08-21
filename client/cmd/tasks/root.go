package tasks

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "tasks",
	Short: "Interact with workspace tasks",
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
