package cgroups_v2

import (
	"path/filepath"
	"supervisor/pkg/cgroups"
)

type IO struct {
	path string
}

func NewIOControllerWithMount(mountPoint, path string) *IO {
	fullPath := filepath.Join(mountPoint, path)
	return &IO{
		path: fullPath,
	}
}

func NewIOController(path string) *IO {
	return &IO{
		path: path,
	}
}

func (io *IO) PSI() (cgroups.PSI, error) {
	path := filepath.Join(io.path, "io.pressure")
	return cgroups.ReadPSIValue(path)
}

func (io *IO) Max() ([]cgroups.DeviceIOMax, error) {
	path := filepath.Join(io.path, "io.max")
	return cgroups.ReadIOMax(path)
}
