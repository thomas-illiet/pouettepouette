package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ReadinessProbeType determines the editor readiness probe type.
type ReadinessProbeType string

const (
	// ReadinessProcessProbe returns ready once the editor process has been started.
	ReadinessProcessProbe ReadinessProbeType = ""

	// ReadinessHTTPProbe returns ready once a single HTTP request against the editor was successful.
	ReadinessHTTPProbe ReadinessProbeType = "http"
)

// EditorConfig is the Editor specific configuration.
type EditorConfig struct {
	// Name is the unique identifier of the Editor.
	Name string `json:"name"`

	// Version is the version of the Editor.
	Version string `json:"version"`

	// Entrypoint is the command that gets executed by supervisor to start
	// the Editor process. If this command exits, supervisor will start it again.
	// If this command exits right after it was started with a non-zero exit
	// code the workspace is stopped.
	Entrypoint string `json:"entrypoint"`

	// EntrypointArgs
	EntrypointArgs []string `json:"entrypointArgs"`

	// LogRateLimit can be used to limit the log output of the editor process.
	// Any output that exceeds this limit is silently dropped.
	// Expressed in kb/sec. Can be overridden by the workspace config (smallest value wins).
	LogRateLimit int `json:"logRateLimit"`

	// ReadinessProbe configures the probe used to serve the editor status
	ReadinessProbe struct {
		// Type determines the type of readiness probe we'll use.
		// Defaults to process.
		Type ReadinessProbeType `json:"type"`

		// HTTPProbe configures the HTTP readiness probe.
		HTTPProbe struct {
			// Schema is either "http" or "https". Defaults to "http".
			Schema string `json:"schema"`

			// Host is the host to make requests to. Default to "localhost".
			Host string `json:"host"`

			// Port is the port to make requests to. Default it the editor port in the supervisor config.
			Port int `json:"port"`

			// Path is the path to make requests to. Defaults to "/".
			Path string `json:"path"`
		} `json:"http"`
	} `json:"readinessProbe"`
}

// loadEditorConfig reads and parses an Editor configuration from the given file path.
func loadEditorConfig(configPath string) (*EditorConfig, error) {
	f, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open editor config %q: %w", configPath, err)
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	var cfg EditorConfig
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse editor config %q: %w", configPath, err)
	}

	// If no name is provided, default to the parent directory's name.
	if cfg.Name == "" {
		cfg.Name = filepath.Base(filepath.Dir(configPath))
	}

	return &cfg, nil
}
