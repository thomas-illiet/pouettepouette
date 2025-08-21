package workspace

import (
	"client/pkg/supervisor"
	"context"
	"fmt"
	"supervisor/api"
	"time"

	"github.com/spf13/cobra"
)

var UrlCmd = &cobra.Command{
	Use:   "url",
	Short: "Prints the URL of this workspace",
	Long:  "Prints the URL of this workspace",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Set a timeout for the request
		ctx, cancel := context.WithTimeout(cmd.Context(), 5*time.Second)
		defer cancel()

		// Create a supervisor client
		client, err := supervisor.New(ctx)
		if err != nil {
			return err
		}
		defer client.Close()

		// Fetch system information
		data, err := client.System.WorkspaceInfo(ctx, &api.WorkspaceInfoRequest{})
		if err != nil {
			return err
		}

		fmt.Println(data.WorkspaceUrl)
		return nil
	},
}
