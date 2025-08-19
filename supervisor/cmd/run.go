package cmd

import (
	"supervisor/pkg/supervisor"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "starts the supervisor",

	Run: func(cmd *cobra.Command, args []string) {
		logFile := initLog(false)
		defer logFile.Close()

		supervisor.Version = Version
		supervisor.Run()
	},
}
