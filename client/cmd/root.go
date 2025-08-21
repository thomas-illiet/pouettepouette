package cmd

import (
	"client/cmd/ping"
	"client/cmd/system"
	"client/cmd/tasks"
	"client/cmd/workspace"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

const rootCmdName = "oc"

// lastSignal stores the last OS signal received (e.g., SIGINT, SIGTERM).
var lastSignal os.Signal

var rootCmd = &cobra.Command{
	Use:           rootCmdName,
	Short:         "Command line interface for Opencoder",
	SilenceErrors: true,
	SilenceUsage:  true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Set up context with cancel
		ctx, cancel := context.WithCancel(cmd.Context())
		cmd.SetContext(ctx)

		// Listen for OS signals in a goroutine
		go func() {
			signals := make(chan os.Signal, 1)
			signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
			sig := <-signals
			lastSignal = sig
			cancel() // Cancel the command context
		}()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(ping.Cmd)
	rootCmd.AddCommand(tasks.Cmd)
	rootCmd.AddCommand(system.Cmd)
	rootCmd.AddCommand(workspace.Cmd)
	rootCmd.AddCommand(versionCmd)
}

// Execute runs the root command and exits with proper codes on error or signal.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// If a signal was received, exit with code 128 + signal number
	if sig, ok := lastSignal.(syscall.Signal); ok {
		os.Exit(128 + int(sig))
	}
}
