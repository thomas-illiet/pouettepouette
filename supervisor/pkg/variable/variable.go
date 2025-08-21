package variable

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"common/log"
	"supervisor/pkg/config"
)

// AddDefault sets up the given command with a default set of environment variables.
func AddDefault(cmd *exec.Cmd, cfg *config.Config) *exec.Cmd {
	cmd.Env = buildChildProcEnv(cfg)
	return cmd
}

// buildChildProcEnv computes the environment variables to pass to a child process.
func buildChildProcEnv(cfg *config.Config) []string {
	// Start with the current process environment
	currentEnv := os.Environ()

	// Convert environment to a map for easy overriding
	envMap := make(map[string]string)
	for _, e := range currentEnv {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) < 2 {
			log.Printf("\"%s\" has invalid format, not including in IDE environment", e)
			continue
		}
		name, value := parts[0], parts[1]

		if isBlacklisted(name) {
			continue
		}
		envMap[name] = value
	}

	// Helper function for variable expansion
	getEnv := func(name string) string {
		return envMap[name]
	}

	// Overlay environment from configuration, with variable expansion
	for name, value := range cfg.Runtime.Environment {
		if isBlacklisted(name) {
			continue
		}
		envMap[name] = os.Expand(value, getEnv)
	}

	// Add supervisor-specific variables
	envMap["SUPERVISOR_ADDR"] = fmt.Sprintf("localhost:%d", cfg.APIEndpointPort)

	// Convert map back to slice for exec.Cmd
	var envList []string
	for name, value := range envMap {
		envList = append(envList, fmt.Sprintf("%s=%s", name, value))
	}

	return envList
}

// isBlacklisted checks whether an environment variable should be excluded from the child process.
func isBlacklisted(name string) bool {
	nameUpper := strings.ToUpper(name)

	prefixBlacklist := []string{
		"OPENCODER",
		"KUBERNETES",
	}

	for _, prefix := range prefixBlacklist {
		if strings.HasPrefix(nameUpper, prefix) {
			return true
		}
	}
	return false
}
