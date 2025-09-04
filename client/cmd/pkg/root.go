package pkg

import "github.com/spf13/cobra"

var jsonFormat bool

var Cmd = &cobra.Command{
	Use:   "pkg",
	Short: "Package command for application management",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		return nil
	},
}

func init() {
	Cmd.AddCommand(ListCmd)
	Cmd.AddCommand(GetCmd)
	Cmd.AddCommand(InstallCmd)
	Cmd.AddCommand(RemoveCmd)
}
