package config

import (
	"fmt"
	"time"

	"github.com/Netflix/go-env"
)

type WorkspaceConfig struct {
	// WorkspaceID is the ID of the workspace
	WorkspaceID int64 `env:"OPENCODER_WORKSPACE_ID"`

	// WorkspaceLocation is the location of the workspace
	WorkspaceLocation string `env:"OPENCODER_WORKSPACE_PATH"`

	// WorkspaceUrl is a URL for which workspace is accessed.
	WorkspaceUrl string `env:"OPENCODER_WORKSPACE_URL"`

	// GitUsername makes supervisor configure the global user.name Git setting.
	GitUsername string `env:"OPENCODER_GIT_USER_NAME"`

	// GitEmail makes supervisor configure the global user.email Git setting.
	GitEmail string `env:"OPENCODER_GIT_USER_EMAIL"`

	// DebugEnabled controls whether the supervisor debugging facilities should be enabled
	DebugEnable bool `env:"SUPERVISOR_DEBUG_ENABLE"`

	// OwnerId is the user id who owns the workspace
	OwnerId int64 `env:"OPENCODER_OWNER_ID"`

	// LogRateLimit limits the log output of the editor process.
	// Any output that exceeds this limit is silently dropped.
	// Expressed in kb/sec. Can be overridden by the IDE config (smallest value wins).
	LogRateLimit int `env:"OPENCODER_RATE_LIMIT_LOG"`

	// WorkspaceClusterHost is a host under which this workspace is served, e.g. fr42.oc.dev.local
	WorkspaceClusterHost string `env:"OPENCODER_WORKSPACE_CLUSTER_HOST"`

	// TerminationGracePeriodSeconds is the max number of seconds the workspace can take to shut down all its processes after SIGTERM was sent.
	TerminationGracePeriodSeconds *int `env:"OPENCODER_TERMINATION_GRACE_PERIOD_SECONDS"`
}

// loadWorkspaceConfig loads the workspace configuration from environment variables.
func loadWorkspaceConfig() (*WorkspaceConfig, error) {
	var res WorkspaceConfig
	_, err := env.UnmarshalFromEnviron(&res)
	if err != nil {
		return nil, fmt.Errorf("cannot load workspace config: %w", err)
	}

	return &res, nil
}

func (c WorkspaceConfig) GetTerminationGracePeriod() time.Duration {
	defaultGracePeriod := 15 * time.Second
	if c.TerminationGracePeriodSeconds == nil || *c.TerminationGracePeriodSeconds <= 0 {
		return defaultGracePeriod
	}
	return time.Duration(*c.TerminationGracePeriodSeconds) * time.Second
}
