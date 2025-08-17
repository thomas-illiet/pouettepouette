package supervisor

import (
	"common/log"
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime/debug"
)

var (
	Version = ""
)

// Run serves as main entrypoint to the supervisor.
func Run() {
	exitCode := 0
	defer handleExit(&exitCode)
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		log.WithField("cause", r).WithField("stack", string(debug.Stack())).Error("panicked")
		if exitCode == 0 {
			exitCode = 1
		}
	}()

	//
	cfg, err := GetConfig()
	if err != nil {
		log.WithError(err).Fatal("configuration error")
	}

	// Check if the program is called with "run" as an argument to start the supervisor.
	if len(os.Args) < 2 || os.Args[1] != "run" {
		fmt.Println("supervisor makes sure your workspace/IDE keeps running smoothly.\n" +
			"You don't have to call this thing, Opencoder calls it for you.")
		return
	}

	// Set git credential helper configuration
	//configureGit(cfg)

	ctx, cancel := context.WithCancel(context.Background())

	apiService = serverapi.NewServerApiService(ctx, &serverapi.ServiceConfig{
		Host:              host,
		Endpoint:          endpoint,
		InstanceID:        cfg.WorkspaceInstanceID,
		WorkspaceID:       cfg.WorkspaceID,
		OwnerID:           cfg.OwnerId,
		SupervisorVersion: Version,
		ConfigcatEnabled:  cfg.ConfigcatEnabled,
	}, tokenService)

	//

	_ = cfg
}

func handleExit(ec *int) {
	exitCode := *ec
	log.WithField("exitCode", exitCode).Debug("supervisor exit")
	os.Exit(exitCode)
}

func configureGit(cfg *Config) {
	settings := [][]string{
		{"push.default", "simple"},
		{"alias.lg", "log --color --graph --pretty=format:'%Cred%h%Creset -%C(yellow)%d%Creset %s %Cgreen(%cr) %C(bold blue)<%an>%Creset' --abbrev-commit"},
		{"credential.helper", "/usr/bin/oc credential-helper"},
		{"safe.directory", "*"},
	}
	if cfg.Workspace.GitUsername != "" {
		settings = append(settings, []string{"user.name", cfg.Workspace.GitUsername})
	}
	if cfg.Workspace.GitEmail != "" {
		settings = append(settings, []string{"user.email", cfg.Workspace.GitEmail})
	}

	for _, s := range settings {
		cmd := exec.Command("git", append([]string{"config", "--global"}, s...)...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.WithError(err).WithField("args", s).Warn("git config error")
		}
	}
}
