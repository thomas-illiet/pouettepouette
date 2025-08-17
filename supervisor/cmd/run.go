package cmd

import (
	"github.com/spf13/cobra"
	"supervisor/pkg/supervisor"

	common_grpc "common/grpc"
)

var runOpts struct {
	RunGP bool
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "starts the supervisor",

	Run: func(cmd *cobra.Command, args []string) {
		logFile := initLog(!runOpts.RunGP)
		defer logFile.Close()

		common_grpc.SetupLogging()
		supervisor.Version = Version
		supervisor.Run()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().BoolVar(&runOpts.RunGP, "rungp", false, "run ==supervisor in a run-gp context")
}
