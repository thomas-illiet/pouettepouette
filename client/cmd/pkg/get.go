package pkg

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

type getCmd struct{}

func init() {
	GetCmd.Flags().BoolVarP(&jsonFormat, "json", "j", false, "Output in JSON format")
}

// GetCmd represents the list package command.
var GetCmd = &cobra.Command{
	Use:   "get",
	Args:  cobra.MinimumNArgs(1),
	Short: "Get package info, such as its ID, name, etc.",
	RunE: func(cmd *cobra.Command, args []string) error {
		// string to int
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

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
		data, err := client.Package.Get(ctx, &api.GetPackageRequest{Id: int64(id)})
		if err != nil {
			return err
		}

		// Output in JSON or table format
		if jsonFormat {
			content, _ := json.Marshal(data)
			fmt.Println(string(content))
		} else {
			getCmd{}.PrintTable(data)
		}

		return nil
	},
}

// PrintTable renders resource usage in a table format
func (gc getCmd) PrintTable(resources *api.GetPackageResponse) {
	table := tablewriter.NewWriter(os.Stdout)
	_ = table.Append("Id", resources.Id)
	_ = table.Append("Name", resources.Name)
	_ = table.Append("Description", resources.Description)
	_ = table.Append("Status", resources.Status)
	_ = table.Append("Id", resources.Version)

	_ = table.Render()
}
