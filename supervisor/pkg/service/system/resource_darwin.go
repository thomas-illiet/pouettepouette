package system

import (
	"math/rand"
	"supervisor/api"
)

// ⚠️ Compatibility note:
// csgroup is not supported on macOS. This package, and the final product
// as a whole, will not run on macOS in any case.
// Only Linux environments are supported.

func resolveCPUStat() (*cpuStat, error) {
	// Fake CPU usage: 0–100%, uptime: up to ~30 days
	return &cpuStat{
		usage:  rand.Float64() * 100,                // % CPU usage
		uptime: float64(rand.Intn(30*24*60*60 + 1)), // seconds
	}, nil
}

func resolveMemoryStatus() (*api.ResourceStatus, error) {
	// Fake memory: random between 4–32GB total, 10–90% used
	limit := int64((4 + rand.Intn(29)) * 1024 * 1024 * 1024) // 4–32 GiB
	used := int64(float64(limit) * (0.1 + rand.Float64()*0.8))
	return &api.ResourceStatus{
		Limit: limit,
		Used:  used,
	}, nil
}

// readProcUptime returns random uptime in seconds (up to 7 days).
func readProcUptime() (float64, error) {
	return float64(rand.Intn(7 * 24 * 60 * 60)), nil
}

// resolveCPUStatus calculates fake CPU quota limits.
func resolveCPUStatus() (*api.ResourceStatus, error) {
	// Random between 1–8 CPUs, with 5–95% usage
	limit := int64((1 + rand.Intn(8)) * 1000) // millicores
	used := int64(float64(limit) * (0.05 + rand.Float64()*0.9))
	return &api.ResourceStatus{
		Limit: limit,
		Used:  used,
	}, nil
}
