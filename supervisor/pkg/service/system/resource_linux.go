package system

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"supervisor/api"
	cgroups "supervisor/pkg/cgroups/v2"

	linuxproc "github.com/c9s/goprocinfo/linux"
	"golang.org/x/xerrors"
)

func resolveCPUStat() (*cpuStat, error) {
	cpu := cgroups.NewCpuController("/sys/fs/cgroups")

	stats, err := cpu.Stat()
	if err != nil {
		return nil, xerrors.Errorf("failed to parse cpu.stat: %w", err)
	}

	usage := float64(stats.UsageTotal) * 1e-9

	uptime, err := readProcUptime()
	if err != nil {
		return nil, xerrors.Errorf("failed to parse uptime: %w", err)
	}

	return &cpuStat{
		usage:  usage,
		uptime: uptime,
	}, nil
}

func resolveMemoryStatus() (*api.ResourceStatus, error) {
	memory := cgroups.NewMemoryController("/sys/fs/cgroup")

	limit, err := memory.Max()
	if err != nil {
		return nil, fmt.Errorf("failed to parse memory.max: %w", err)
	}

	memInfo, err := linuxproc.ReadMemInfo("/proc/meminfo")
	if err != nil {
		return nil, fmt.Errorf("failed to read meminfo: %w", err)
	}

	memTotal := memInfo.MemTotal * 1024
	if limit > memTotal && memTotal > 0 {
		limit = memTotal
	}

	used, err := memory.Current()
	if err != nil {
		return nil, fmt.Errorf("failed to parse memory.current: %w", err)
	}

	memstats, err := memory.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to parse memory.stats: %w", err)
	}

	if used < memstats.InactiveFileTotal {
		used = 0
	} else {
		used -= memstats.InactiveFileTotal
	}

	return &api.ResourceStatus{
		Limit: int64(limit),
		Used:  int64(used),
	}, nil
}

// readProcUptime returns system uptime in seconds.
func readProcUptime() (float64, error) {
	content, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return 0, fmt.Errorf("failed to read uptime: %w", err)
	}
	fields := strings.Fields(strings.TrimSpace(string(content)))
	uptime, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse uptime: %w", err)
	}
	return uptime, nil
}

// resolveCPUStatus calculates CPU usage and quota limits.
func resolveCPUStatus() (*api.ResourceStatus, error) {
	// Capture two snapshots with delay to measure CPU usage.
	sample1, err := resolveCPUStat()
	if err != nil {
		return nil, err
	}
	time.Sleep(time.Second)

	sample2, err := resolveCPUStat()
	if err != nil {
		return nil, err
	}

	// Calculate usage.
	cpuUsage := sample2.usage - sample1.usage
	totalTime := sample2.uptime - sample1.uptime
	used := cpuUsage / totalTime * 1000

	// Determine CPU quota.
	quota, err := readIntFile("/sys/fs/cgroup/cpu/cpu.cfs_quota_us")
	if err != nil {
		return nil, fmt.Errorf("failed to read cpu.cfs_quota_us: %w", err)
	}

	var limit int
	if quota > 0 {
		period, err := readIntFile("/sys/fs/cgroup/cpu/cpu.cfs_period_us")
		if err != nil {
			return nil, fmt.Errorf("failed to read cpu.cfs_period_us: %w", err)
		}
		limit = quota / period * 1000
	} else {
		content, err := os.ReadFile("/sys/fs/cgroup/cpu/cpuacct.usage_percpu")
		if err != nil {
			return nil, fmt.Errorf("failed to read cpuacct.usage_percpu: %w", err)
		}
		limit = len(strings.Fields(strings.TrimSpace(string(content)))) * 1000
	}

	return &api.ResourceStatus{
		Limit: int64(limit),
		Used:  int64(used),
	}, nil
}
