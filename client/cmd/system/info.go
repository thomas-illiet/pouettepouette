package system

import (
	"client/pkg/supervisor"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"supervisor/api"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

type infoCmd struct{}

func init() {
	InfoCmd.Flags().BoolVarP(&jsonFormat, "json", "j", false, "Output in JSON format")
}

// InfoCmd represents the info command.
var InfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Display workspace info, such as its ID, class, etc.",
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

		// Output in JSON or table format
		if jsonFormat {
			content, _ := json.Marshal(data)
			fmt.Println(string(content))
		} else {
			infoCmd{}.PrintTable(data)
		}

		return nil
	},
}

// PrintTable renders system information in a table format
func (rc infoCmd) PrintTable(resources *api.WorkspaceInfoResponse) {
	// Append rows
	table := tablewriter.NewWriter(os.Stdout)
	_ = table.Append([]string{"Workspace ID", strconv.FormatInt(resources.WorkspaceId, 10)})
	_ = table.Append([]string{"Checkout Location", resources.CheckoutLocation})
	_ = table.Append([]string{"User Home", resources.UserHome})
	_ = table.Append([]string{"Cluster host", resources.ClusterHost})
	_ = table.Append([]string{"Workspace URL", resources.WorkspaceUrl})
	_ = table.Append([]string{"Editor alias", resources.IdeAlias})
	_ = table.Append([]string{"Owner ID", strconv.FormatInt(resources.OwnerId, 10)})
	_ = table.Render()
}
