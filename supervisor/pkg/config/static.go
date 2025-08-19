package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const supervisorConfigFile = "supervisor-config.json"

// StaticConfig is the supervisor-wide configuration.
type StaticConfig struct {
	// EditorConfigLocation is a path in the filesystem where to find the editor configuration
	EditorConfigLocation string `json:"editorConfigLocation"`

	// APIEndpointPort is the port where to serve the API endpoint on
	APIEndpointPort int `json:"apiEndpointPort"`
}

// loadStaticConfigFromFile loads the static supervisor configuration from
// a file named "supervisor-config.json" which is expected right next to the supervisor executable.
func loadStaticConfigFromFile() (*StaticConfig, error) {
	loc := "/Users/thomas-illiet/GolandProjects/Opencoder/supervisor/main.go"
	//loc, err := os.Executable()
	//if err != nil {
	//	return nil, fmt.Errorf("cannot get executable path: %w", err)
	//}

	loc = filepath.Join(filepath.Dir(loc), supervisorConfigFile)
	fc, err := os.ReadFile(loc)
	if err != nil {
		return nil, fmt.Errorf("cannot read supervisor config file %s: %w", loc, err)
	}

	var res StaticConfig
	err = json.Unmarshal(fc, &res)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal supervisor config file %s: %w", loc, err)
	}

	return &res, nil
}
