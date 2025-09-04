package task

import (
	"supervisor/pkg/config"
	"supervisor/pkg/terminal"
)

type tasksManager struct {
	config          *config.Config
	storeLocation   string
	tasks           []*task
	terminalService *terminal.MuxTerminalService
}

type taskSuccess string

type task struct {
	//api.TaskStatus
	config      config.TaskConfig
	command     string
	successChan chan taskSuccess
	title       string
	lastOutput  string
}

func newTasksManager(config *config.Config, terminalService *terminal.MuxTerminalService) *tasksManager {
	return &tasksManager{
		config:          config,
		storeLocation:   "/tmp/.opencoder",
		terminalService: terminalService,
	}
}
