package pkg

import (
	"client/pkg/supervisor"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"supervisor/api"
	"time"

	"github.com/spf13/cobra"
)

// InstallCmd represents the install package command.
var InstallCmd = &cobra.Command{
	Use:  "install",
	Args: cobra.MinimumNArgs(1),
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

		// Request package installation
		data, err := client.Package.Install(ctx, &api.InstallPackageRequest{Id: int64(id)})
		if err != nil {
			return err
		}

		// Output in JSON or table format
		if jsonFormat {
			content, _ := json.Marshal(data)
			fmt.Println(string(content))
		} else {
			fmt.Println(data.Status)
		}

		return nil
	},
}
