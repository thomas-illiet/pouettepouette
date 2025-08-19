package editor

import (
	"common/log"
	"context"
	"os"
	"os/exec"
	"runtime"
	"supervisor/pkg/config"
	"supervisor/pkg/dropwriter"
	"supervisor/pkg/variable"
	"sync"
	"syscall"
	"time"
)

const timeBudgetIDEShutdown = 15 * time.Second

// StartAndWatchEditor launches the configured editor process and continuously monitors it.
func StartAndWatchEditor(ctx context.Context, cfg *config.Config, wg *sync.WaitGroup, ideReady *ReadyState) {
	defer wg.Done()
	defer log.Debug("Editor supervisor stopped")

	//
	var cmd *exec.Cmd
	var ideStopped chan struct{}

	firstStart := true
	for {
		// Launch a new editor process
		ideStopped = make(chan struct{}, 1)
		cmd = prepareEditorLaunch(cfg)
		launchEditor(cfg, cmd, ideStopped, ideReady)

		// Only track readiness on the first start
		if firstStart {
			firstStart = false
			go monitorReadiness(ctx, ideReady)
		}

		// Wait until either the editor stops or the supervisor is cancelled
		select {
		case <-ideStopped:
			// Editor stopped unexpectedly -> cleanup and restart
			_ = syscall.Kill(-1*cmd.Process.Pid, syscall.SIGKILL)
			time.Sleep(1 * time.Second)

		case <-ctx.Done():
			// Supervisor shutdown requested
			log.Info("context cancelled, stopping editor")
			gracefulStop(cmd, ideStopped)
			return
		}
	}
}

// gracefulStop tries to stop the editor process cleanly. If the editor
// does not exit within `timeBudgetIDEShutdown`, it is force-killed.
func gracefulStop(cmd *exec.Cmd, ideStopped chan struct{}) {
	log.WithField("timeout", timeBudgetIDEShutdown).Info("waiting for editor shutdown")

	select {
	case <-ideStopped:
		log.Info("editor stopped gracefully")
	case <-time.After(timeBudgetIDEShutdown):
		log.Error("editor did not stop in time, sending SIGKILL")
		_ = cmd.Process.Signal(syscall.SIGKILL)
	}
}

// monitorReadiness waits for the editor to report readiness or a timeout.
// If the editor does not signal readiness within ~10s, a warning is logged.
func monitorReadiness(ctx context.Context, ideReady *ReadyState) {
	timer := time.NewTimer(10 * time.Second)
	defer timer.Stop()

	select {
	case <-timer.C:
		log.Warn("editor readiness timeout expired")
	case <-ideReady.Wait():
		log.Info("editor is ready")
	case <-ctx.Done():
		log.Debug("shutdown before editor became ready")
	}
}

// launchEditor starts the editor as a subprocess, runs readiness probes, and monitors for process exit.
func launchEditor(cfg *config.Config, cmd *exec.Cmd, ideStopped chan struct{}, ideReady *ReadyState) {
	go func() {
		// Lock thread to ensure Pdeathsig works correctly with SysProcAttr
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()

		log.Info("starting editor process")

		if err := cmd.Start(); err != nil {
			log.WithError(err).Fatal("Editor failed to start")
			return
		}

		// Run readiness probe in background
		go func() {
			runEditorReadinessProbe(cfg)
			ideReady.Set(true)
		}()

		// Block until the process exits
		if err := cmd.Wait(); err != nil && err.Error() != "signal: terminated" {
			log.WithError(err).Warn("Editor stopped unexpectedly")
			if !ideReady.Get() {
				log.WithError(err).Fatal("Editor failed before becoming ready")
				return
			}
		}

		// Reset readiness and signal stop
		ideReady.Set(false)
		close(ideStopped)
	}()
}

// prepareEditorLaunch configures the exec.Cmd used to start the editor.
func prepareEditorLaunch(cfg *config.Config) *exec.Cmd {
	log.WithField("args", cfg.Editor.EntrypointArgs).
		WithField("entrypoint", cfg.Editor.Entrypoint).
		Info("preparing editor launch")

	cmd := exec.Command(cfg.Editor.Entrypoint, cfg.Editor.EntrypointArgs...)
	variable.AddDefault(cmd, cfg)

	// Prepare operating system-specific attributes
	prepareSysProc(cmd)

	// Default: pass through raw stdout/stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Optionally wrap with a rate-limited writer
	if lrr := cfg.EditorLogRateLimit(); lrr > 0 {
		limit := int64(lrr)
		cmd.Stdout = dropwriter.Writer(cmd.Stdout, dropwriter.NewBucket(limit*1024*3, limit*1024))
		cmd.Stderr = dropwriter.Writer(cmd.Stderr, dropwriter.NewBucket(limit*1024*3, limit*1024))
		log.WithField("limit_kb_per_sec", limit).Info("rate limiting editor log output")
	}

	return cmd
}
