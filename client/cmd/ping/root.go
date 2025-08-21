package ping

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "ping",
	Short: "Ping command for testing connectivity",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		return nil
	},
}

func init() {
	Cmd.AddCommand(pingSupervisorCmd)
}
