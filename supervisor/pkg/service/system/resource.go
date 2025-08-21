package system

import (
	"context"
	"os"
	"strconv"
	"strings"
	"supervisor/api"

	"golang.org/x/sys/unix"
)

// ResourcesStatus returns system resource usage (CPU, Memory).
func (is *SystemService) ResourcesStatus(ctx context.Context, request *api.ResourcesStatusRequest) (*api.ResourcesStatusResponse, error) {
	return is.GetResources(ctx)
}

// calcSeverity maps a percentage to a severity level.
func calcSeverity(percentage int64) api.ResourceStatusSeverity {
	switch {
	case percentage >= 95:
		return api.ResourceStatusSeverity_danger
	case percentage >= 80:
		return api.ResourceStatusSeverity_warning
	default:
		return api.ResourceStatusSeverity_normal
	}
}

// GetResources collects and returns current CPU and Memory usage with severity levels.
func (is *SystemService) GetResources(ctx context.Context) (*api.ResourcesStatusResponse, error) {
	memory, err := resolveMemoryStatus()
	if err != nil {
		return nil, err
	}

	cpu, err := resolveCPUStatus()
	if err != nil {
		return nil, err
	}

	disk, err := getDiskUsage(is.Cfg.WorkspaceLocation)
	if err != nil {
		return nil, err
	}

	cpuPct := int64(float64(cpu.Used) / float64(cpu.Limit) * 100)
	memPct := int64(float64(memory.Used) / float64(memory.Limit) * 100)
	diskPct := int64(float64(disk.Used) / float64(disk.Limit) * 100)

	cpu.Severity = calcSeverity(cpuPct)
	memory.Severity = calcSeverity(memPct)
	disk.Severity = calcSeverity(diskPct)

	return &api.ResourcesStatusResponse{
		Flavor: is.Cfg.FlavorName,
		Memory: memory,
		Cpu:    cpu,
		Disk:   disk,
	}, nil
}

// cpuStat represents a snapshot of CPU usage at a given uptime.
type cpuStat struct {
	usage  float64
	uptime float64
}

// readIntFile reads a file and parses it as an integer.
func readIntFile(path string) (int, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(string(content)))
}

func getDiskUsage(path string) (*api.ResourceStatus, error) {
	var stat unix.Statfs_t
	err := unix.Statfs(path, &stat)
	if err != nil {
		return nil, err
	}

	// Total size = total blocks * block size
	total := stat.Blocks * uint64(stat.Bsize)

	// Available size = available blocks * block size
	available := stat.Bavail * uint64(stat.Bsize)

	return &api.ResourceStatus{
		Limit: int64(total),
		Used:  int64(available),
	}, nil
}
