package config

import "common/log"

// Config configures supervisor.
type Config struct {
	Static    StaticConfig
	Workspace WorkspaceConfig
	Editor    EditorConfig
	Runtime   RuntimeConfig
}

// GetConfig loads the supervisor configuration.
func GetConfig() (*Config, error) {
	log.Debug("loading static config from file...")
	static, err := loadStaticConfigFromFile()
	if err != nil {
		log.WithError(err).
			Error("failed to load static config")
		return nil, err
	}

	log.Debug("loading editor config...")
	editor, err := loadEditorConfig(static.EditorConfigLocation)
	if err != nil {
		log.WithError(err).
			WithField("path", static.EditorConfigLocation).
			Error("failed to load editor config")
		return nil, err
	}

	log.Debug("loading workspace config...")
	workspace, err := loadWorkspaceConfig()
	if err != nil {
		log.WithError(err).
			Error("failed to load workspace config")
		return nil, err
	}

	log.Debug("loading runtime config...")
	runtime, err := loadRuntimeConfig(workspace.WorkspaceLocation)
	if err != nil {
		log.WithError(err).
			WithField("path", workspace.WorkspaceLocation).
			Error("failed to load runtime config")
	}

	log.Debug("configuration loaded successfully")
	return &Config{
		Static:    *static,
		Workspace: *workspace,
		Editor:    *editor,
		Runtime:   *runtime,
	}, nil
}

// EditorLogRateLimit returns the log rate limit for the IDE process in kib/sec.
// If log rate limiting is disabled, this function returns 0.
func (c Config) EditorLogRateLimit() int {
	if c.Workspace.LogRateLimit == 0 {
		return c.Editor.LogRateLimit
	}
	if c.Workspace.LogRateLimit < c.Editor.LogRateLimit {
		return c.Workspace.LogRateLimit
	}
	return c.Editor.LogRateLimit
}
