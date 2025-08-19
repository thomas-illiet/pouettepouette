package git

import (
	"os"
	"os/exec"

	"common/log"
	"supervisor/pkg/config"
)

// ConfigureGit applies a set of global Git configurations based on defaults
// and values from the opencoder configuration.
func ConfigureGit(cfg *config.Config) {
	defaultSettings := [][]string{
		{"push.default", "simple"},
		{"alias.lg", `log --color --graph --pretty=format:'%Cred%h%Creset -%C(yellow)%d%Creset %s %Cgreen(%cr) %C(bold blue)<%an>%Creset' --abbrev-commit`},
		{"credential.helper", "/usr/bin/oc credential-helper"},
		{"safe.directory", "*"},
	}

	// Add user-specific settings if provided
	if cfg.Workspace.GitUsername != "" {
		defaultSettings = append(defaultSettings, []string{"user.name", cfg.Workspace.GitUsername})
	}
	if cfg.Workspace.GitEmail != "" {
		defaultSettings = append(defaultSettings, []string{"user.email", cfg.Workspace.GitEmail})
	}

	applyGitSettings(defaultSettings)
}

// applyGitSettings iterates over key/value pairs and applies them via `git config --global`.
func applyGitSettings(settings [][]string) {
	for _, setting := range settings {
		args := append([]string{"config", "--global"}, setting...)
		cmd := exec.Command("git", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			log.WithError(err).
				WithField("setting", setting).
				Warn("failed to apply git config")
		}
	}
}
