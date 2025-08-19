package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// RuntimeConfig defines the root structure of the runtime configuration file.
// It contains environment variables, Git configuration, tasks, and VS Code settings.
type RuntimeConfig struct {
	Environment      map[string]string      `yaml:"env"`       // Arbitrary environment variables
	GitConfiguration map[string]interface{} `yaml:"gitConfig"` // Git-related configuration values
	Tasks            []TaskConfig           `yaml:"tasks"`     // List of tasks to run in the workspace
	Vscode           VscodeConfig           `yaml:"vscode"`    // VS Code-specific settings
}

// VscodeConfig defines VS Code-related configuration such as required extensions.
type VscodeConfig struct {
	Extensions []string `yaml:"extensions"` // Extensions to be installed in VS Code
}

// TaskConfig represents the configuration of a single task that can be run
// within the workspace. Each field corresponds to a different execution phase.
type TaskConfig struct {
	Name     *string                 `yaml:"name"`     // Optional name of the task
	Before   *string                 `yaml:"before"`   // Command to run before other tasks
	Init     *string                 `yaml:"init"`     // Initialization command
	Prebuild *string                 `yaml:"prebuild"` // Command to run before build
	Command  *string                 `yaml:"command"`  // Main command for the task
	Env      *map[string]interface{} `yaml:"env"`      // Environment variables specific to the task
}

// NewRuntimeConfig creates a new RuntimeConfig with all properties initialized
func newRuntimeConfig() *RuntimeConfig {
	return &RuntimeConfig{
		Environment:      make(map[string]string),
		GitConfiguration: make(map[string]interface{}),
		Tasks:            []TaskConfig{},
		Vscode: VscodeConfig{
			Extensions: []string{},
		},
	}
}

// loadRuntimeConfig loads a runtime configuration from a given workspace file path.
func loadRuntimeConfig(workspaceLocation string) (*RuntimeConfig, error) {
	cfg := newRuntimeConfig()

	// Construct full path to .opencoder.yaml
	configPath := filepath.Join(workspaceLocation, ".opencoder.yml")

	// Attempt to open the configuration file
	file, err := os.Open(configPath)
	if err != nil {
		// If the file does not exist, return the default config without error
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	// Decode YAML into the RuntimeConfig struct
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(cfg); err != nil {
		return cfg, fmt.Errorf("failed to decode YAML: %w", err)
	}

	// Ensure all tasks have initialized Env maps to prevent nil dereferences
	for i := range cfg.Tasks {
		if cfg.Tasks[i].Env == nil {
			cfg.Tasks[i].Env = &map[string]interface{}{}
		}
	}

	return cfg, nil
}
