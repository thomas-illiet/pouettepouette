package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:    "completion",
	Hidden: true,
	Short:  "Generate bash completion scripts",
	Long: `Generate Bash completion scripts for this CLI.

To enable completion for the current shell session, run:

    . <(oc completion)

To automatically load completion scripts in every new session,
add the following line to your shell configuration file (e.g., ~/.bashrc or ~/.profile):

    . <(oc completion)
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return rootCmd.GenBashCompletion(os.Stdout)
	},
}
