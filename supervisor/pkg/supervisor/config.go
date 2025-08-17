package supervisor

import (
	"fmt"
	env "github.com/Netflix/go-env"
	"time"
)

// Config configures supervisor.
type Config struct {
	Workspace WorkspaceConfig
}

type WorkspaceConfig struct {
	// WorkspaceID is the ID of the workspace
	WorkspaceID int64 `env:"OPENCODER_WORKSPACE_ID"`

	// WorkspaceUrl is an URL for which workspace is accessed.
	WorkspaceUrl string `env:"OPENCODER_WORKSPACE_URL"`

	// GitUsername makes supervisor configure the global user.name Git setting.
	GitUsername string `env:"OPENCODER_GIT_USER_NAME"`

	// GitEmail makes supervisor configure the global user.email Git setting.
	GitEmail string `env:"OPENCODER_GIT_USER_EMAIL"`

	// DebugEnabled controls whether the supervisor debugging facilities (pprof, grpc tracing) should be enabled
	DebugEnable bool `env:"SUPERVISOR_DEBUG_ENABLE"`

	// OwnerId is the user id who owns the workspace
	OwnerId int64 `env:"OPENCODER_OWNER_ID"`

	// WorkspaceClusterHost is a host under which this workspace is served, e.g. fr42.oc.dev.echonet
	WorkspaceClusterHost string `env:"OPENCODER_WORKSPACE_CLUSTER_HOST"`

	// TerminationGracePeriodSeconds is the max number of seconds the workspace can take to shut down all its processes after SIGTERM was sent.
	TerminationGracePeriodSeconds *int `env:"OPENCODER_TERMINATION_GRACE_PERIOD_SECONDS"`
}

// GetConfig loads the supervisor configuration.
func GetConfig() (*Config, error) {
	workspace, err := loadWorkspaceConfigFromEnv()
	if err != nil {
		return nil, err
	}

	return &Config{
		Workspace: *workspace,
	}, nil
}

func (c WorkspaceConfig) GetTerminationGracePeriod() time.Duration {
	defaultGracePeriod := 15 * time.Second
	if c.TerminationGracePeriodSeconds == nil || *c.TerminationGracePeriodSeconds <= 0 {
		return defaultGracePeriod
	}
	return time.Duration(*c.TerminationGracePeriodSeconds) * time.Second
}

// loadWorkspaceConfigFromEnv loads the workspace configuration from environment variables.
func loadWorkspaceConfigFromEnv() (*WorkspaceConfig, error) {
	var res WorkspaceConfig
	_, err := env.UnmarshalFromEnviron(&res)
	if err != nil {
		return nil, fmt.Errorf("cannot load workspace config: %w", err)
	}

	return &res, nil
}
