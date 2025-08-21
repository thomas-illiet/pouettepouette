package system

import (
	"client/pkg/supervisor"
	"client/pkg/utils"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"supervisor/api"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

type resourceCmd struct{}

func init() {
	ResourceCmd.Flags().BoolVarP(&noColor, "no-color", "", false, "Disable output colorization")
	ResourceCmd.Flags().BoolVarP(&jsonFormat, "json", "j", false, "Output in JSON format")
}

// ResourceCmd defines the `resources` CLI command.
var ResourceCmd = &cobra.Command{
	Use:   "resources",
	Short: "Display workspace resource usage (CPU, Memory and Disk)",
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

		// Fetch resources usage
		data, err := client.System.ResourcesStatus(ctx, &api.ResourcesStatusRequest{})
		if err != nil {
			return err
		}

		// Output in JSON or table format
		if jsonFormat {
			content, _ := json.Marshal(data)
			fmt.Println(string(content))
		} else {
			resourceCmd{}.PrintTable(data)
		}
		return nil
	},
}

// PrintTable renders resource usage in a table format
func (rc resourceCmd) PrintTable(resources *api.ResourcesStatusResponse) {
	// Format Cpu and Memory data
	cpu := rc.formatCPU(resources)
	memory := rc.formatMemory(resources)
	disk := rc.formatDisk(resources)

	// Apply colorization if enabled
	if !noColor && utils.ColorsEnabled() {
		reset := "\033[0m"
		cpu = rc.getColor(resources.Cpu.Severity) + cpu + reset
		memory = rc.getColor(resources.Memory.Severity) + memory + reset
		disk = rc.getColor(resources.Disk.Severity) + disk + reset
	}

	// Append rows
	table := tablewriter.NewWriter(os.Stdout)
	_ = table.Append([]string{"Flavor name", resources.Flavor})
	_ = table.Append([]string{"CPU (millicores)", cpu})
	_ = table.Append([]string{"Memory (MiB)", memory})
	_ = table.Append([]string{"Disk (GiB)", disk})
	_ = table.Render()
}

// formatCPU returns a human-readable string for CPU usage
func (rc resourceCmd) formatCPU(r *api.ResourcesStatusResponse) string {
	used, limit := r.Cpu.Used, r.Cpu.Limit
	percent := int64((float64(used) / float64(limit)) * 100)
	return fmt.Sprintf("%dm/%dm (%d%%)", used, limit, percent)
}

// formatMemory returns a human-readable string for memory usage
func (rc resourceCmd) formatMemory(r *api.ResourcesStatusResponse) string {
	used, limit := r.Memory.Used, r.Memory.Limit
	usedMiB := used / (1024 * 1024)
	limitMiB := limit / (1024 * 1024)
	percent := int64((float64(used) / float64(limit)) * 100)
	return fmt.Sprintf("%dMi/%dMi (%d%%)", usedMiB, limitMiB, percent)
}

// formatDisk returns a human-readable string for disk usage in GiB
func (rc resourceCmd) formatDisk(r *api.ResourcesStatusResponse) string {
	used, limit := r.Disk.Used, r.Disk.Limit
	usedGiB := used / (1024 * 1024 * 1024)
	limitGiB := limit / (1024 * 1024 * 1024)
	percent := int64((float64(used) / float64(limit)) * 100)
	return fmt.Sprintf("%dGi/%dGi (%d%%)", usedGiB, limitGiB, percent)
}

// getColor returns the ANSI color code for a given severity.
func (rc resourceCmd) getColor(severity api.ResourceStatusSeverity) string {
	switch severity {
	case api.ResourceStatusSeverity_danger:
		return "\033[31m" // red
	case api.ResourceStatusSeverity_warning:
		return "\033[33m" // yellow
	default:
		return "\033[32m" // green
	}
}
