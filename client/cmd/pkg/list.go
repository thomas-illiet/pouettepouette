package pkg

import (
	"client/pkg/supervisor"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"supervisor/api"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

type listCmd struct{}

func init() {
	ListCmd.Flags().BoolVarP(&jsonFormat, "json", "j", false, "Output in JSON format")
}

// ListCmd represents the list package command.
var ListCmd = &cobra.Command{
	Use:   "list",
	Args:  cobra.NoArgs,
	Short: "List package info, such as its ID, name, etc.",
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
		data, err := client.Package.List(ctx, &api.ListPackageRequest{})
		if err != nil {
			return err
		}

		// Output in JSON or table format
		if jsonFormat {
			content, _ := json.Marshal(data)
			fmt.Println(string(content))
		} else {
			listCmd{}.PrintTable(data)
		}

		return nil
	},
}

// PrintTable renders system information in a table format
func (lc listCmd) PrintTable(resources *api.ListPackageResponse) {
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"Id", "Name", "Description", "Status", "Version"})
	for _, pkg := range resources.Packages {
		_ = table.Append(pkg.Id, pkg.Name, pkg.Description, pkg.Status, pkg.Version)
	}
	_ = table.Render()
}
