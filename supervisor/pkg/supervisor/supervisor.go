package supervisor

import (
	"common/log"
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"supervisor/pkg/config"
	"supervisor/pkg/editor"
	"sync"
	"syscall"
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

	// Load supervisor configuration
	cfg, err := config.GetConfig()
	if err != nil {
		log.WithError(err).Fatal("configuration error")
	}

	// Check if the program is called with "run" as an argument to start the supervisor.
	if len(os.Args) < 2 || os.Args[1] != "run" {
		fmt.Println("supervisor makes sure your workspace/Editor keeps running smoothly.\n" +
			"You don't have to call this thing, Opencoder calls it for you.")
		return
	}

	// Set git credential helper configuration
	//configureGit(cfg)

	ctx, cancel := context.WithCancel(context.Background())

	// Start editor
	var ideWG sync.WaitGroup
	var ideReady = editor.NewEditorReadyState()
	//ideWG.Add(1)
	go editor.StartAndWatchEditor(ctx, cfg, &ideWG, ideReady)

	// to shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	select {
	case <-sigChan:
	}

	log.Info("received SIGTERM (or shutdown) - tearing down")

	cancel()
	ideWG.Wait()

}

func handleExit(ec *int) {
	exitCode := *ec
	log.WithField("exitCode", exitCode).Debug("supervisor exit")
	os.Exit(exitCode)
}
