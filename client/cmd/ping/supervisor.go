package ping

import (
	"client/pkg/supervisor"
	"fmt"
	"supervisor/api"

	"github.com/spf13/cobra"
)

var pingSupervisorCmd = &cobra.Command{
	Use:   "supervisor",
	Short: "Ping a supervisor server",
	Long:  "Ping a supervisor server to check if it is available and responding.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create a new supervisor client with the command's context
		client, err := supervisor.New(cmd.Context())
		if err != nil {
			return fmt.Errorf("failed to create supervisor client: %w", err)
		}
		defer client.Close()

		// Ping the supervisor server
		result, err := client.Utility.Ping(cmd.Context(), &api.PingRequest{})
		if err != nil {
			return fmt.Errorf("failed to ping supervisor: %w", err)
		}

		// Print the response message
		fmt.Println(result.Message)
		return nil
	},
}
