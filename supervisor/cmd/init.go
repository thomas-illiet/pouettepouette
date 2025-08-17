package cmd

import (
	"bytes"
	"common/log"
	"common/process"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"supervisor/pkg/supervisor"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/prometheus/procfs"
	"github.com/ramr/go-reaper"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init the supervisor",
	Long: "The init command initializes the supervisor by starting a child process, " +
		"handling graceful termination, and cleaning up zombie processes.",
	Run: func(cmd *cobra.Command, args []string) {
		// Open log file and defer closing it
		logFile := initLog(true)
		defer logFile.Close()

		// Retrieve supervisor configuration
		cfg, err := supervisor.GetConfig()
		if err != nil {
			log.WithError(err).Info("cannot load config")
		}

		// Set up signal handling
		var sigInput = make(chan os.Signal, 1)
		signal.Notify(sigInput, os.Interrupt, syscall.SIGTERM)

		// Check if git executable exists, supervisor will fail if it doesn't
		_, err = exec.LookPath("git")
		if err != nil {
			log.WithError(err).Fatal("cannot find git executable, make sure it is installed as part of gitpod image")
		}

		// Determine supervisor executable path
		supervisorPath, err := os.Executable()
		if err != nil {
			log.WithError(err).Warn("unable to find supervisor executable")
			supervisorPath = "/user/bin/supervisor"
		}

		// Execute the child process (supervisor run)
		runCommand := exec.Command(supervisorPath, "run")
		runCommand.Args[0] = "supervisor"
		runCommand.Stdin = os.Stdin
		runCommand.Stdout = os.Stdout
		runCommand.Stderr = os.Stderr
		runCommand.Env = os.Environ()
		err = runCommand.Start()
		if err != nil {
			log.WithError(err).Error("unable to start supervisor")
			return
		}

		// Channel to signal when the supervisor process is done
		supervisorDone := make(chan struct{})

		// Channel for handling exit codes from reaper
		handledByReaper := make(chan int)

		// Atomic boolean flag to ignore unexpected exit codes
		ignoreUnexpectedExitCode := atomic.Bool{}

		// Function to handle the exit of the supervisor process
		handleSupervisorExit := func(exitCode int) {
			if exitCode == 0 {
				return
			}

			logs := extractFailureFromRun()
			if exitCode == 2 {
				log.Fatal(logs)
			} else {
				if ignoreUnexpectedExitCode.Load() {
					return
				}
				log.WithError(fmt.Errorf("%s", logs)).Fatal("supervisor run error with unexpected exit code")
			}
		}

		// Goroutine to wait for the supervisor process to finish and handle its exit
		go func() {
			defer close(supervisorDone)

			err := runCommand.Wait()
			if err == nil {
				return
			}

			// Check if the error occurred due to a reaper exit
			if strings.Contains(err.Error(), "no child processes") {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
				defer cancel()

				select {
				case <-ctx.Done(): // Timeout reached
					log.Warn("timed out waiting for reaper to clean up the process.")
				case exitCode := <-handledByReaper:
					handleSupervisorExit(exitCode)
				}
			} else if !(strings.Contains(err.Error(), "signal: ")) {
				var exitErr *exec.ExitError
				if errors.As(err, &exitErr) && exitErr.ExitCode() != 0 {
					handleSupervisorExit(exitErr.ExitCode())
				}
				log.WithError(err).Error("supervisor run error")
			}
		}()

		// Configuration for the reaper process
		reaperChan := make(chan reaper.Status, 10)
		reaperConfig := reaper.Config{
			Pid:              -1,
			Options:          0,
			DisablePid1Check: false,
			StatusChannel:    reaperChan,
		}

		// Start the reaper to monitor and clean up zombie processes
		reaper.Start(reaperConfig)

		// Goroutine to listen for reaper statuses and forward them to handledByReaper channel
		go func() {
			for status := range reaperChan {
				if status.Pid != runCommand.Process.Pid {
					continue
				}
				exitCode := status.WaitStatus.ExitStatus()
				handledByReaper <- exitCode
			}
		}()

		// Manages the shutdown sequence for a system or process.
		select {
		case <-supervisorDone:
			// supervisor has ended - we're all done here
			defer log.Info("supervisor has ended (supervisorDone)")
			return
		case <-sigInput:
			// We received a terminating signal, pass it to the supervisor and wait for it to finish
			ignoreUnexpectedExitCode.Store(true)

			ctx, cancel := context.WithTimeout(context.Background(), cfg.Workspace.GetTerminationGracePeriod())
			defer cancel()

			slog := newShutdownLogger()
			defer slog.Close()

			slog.write("initiating shutdown...")

			// Start a goroutine to terminate all processes
			terminationDone := make(chan struct{})
			go func() {
				defer close(terminationDone)
				slog.TerminateSync(ctx, runCommand.Process.Pid)
				terminateAllProcesses(ctx, slog)
			}()

			// Wait for either successful termination or the timeout
			select {
			case <-ctx.Done():
				// Time is up; give goroutines a little time to react to this
				time.Sleep(time.Millisecond * 1000)
				defer log.Info("supervisor has ended (ctx.Done)")

			case <-terminationDone:
				defer log.Info("supervisor has ended (terminationDone)")
			}

			slog.write("all processes have been shut down")
		}
	},
}

// terminateAllProcesses terminates all processes but ours until there are none anymore or the context is cancelled
// on context cancellation any still running processes receive a SIGKILL
func terminateAllProcesses(ctx context.Context, slog shutdownLogger) {
	for {
		//
		processes, err := procfs.AllProcs()
		if err != nil {
			log.WithError(err).Error("Cannot list processes")
			slog.write(fmt.Sprintf("Cannot list processes: %s", err))
			return
		}

		// only one process (must be us)
		if len(processes) == 1 {
			return
		}

		// terminate all processes but ourself
		var wg sync.WaitGroup
		for _, proc := range processes {
			if proc.PID == os.Getpid() {
				continue
			}
			p := proc
			wg.Add(1)
			go func() {
				defer wg.Done()
				slog.TerminateSync(ctx, p.PID)
			}()
		}
		wg.Wait()
	}
}

// shutdownLogger defines an interface for logging shutdown operations and
// handling process termination in a synchronized way.
type shutdownLogger interface {
	write(s string)
	TerminateSync(ctx context.Context, pid int)
	io.Closer
}

// newShutdownLogger creates a new shutdownLogger implementation that logs shutdown
// events to a file under /tmp/opencoder/supervisor-termination.
func newShutdownLogger() shutdownLogger {
	file := "/tmp/opencoder/supervisor-termination"
	f, err := os.Create(file)
	if err != nil {
		log.WithError(err).WithField("file", file).Error("Couldn't create shutdown log file")
	}
	result := shutdownLoggerImpl{
		file:      f,
		startTime: time.Now(),
	}
	return &result
}

// shutdownLoggerImpl is the concrete implementation of shutdownLogger.
type shutdownLoggerImpl struct {
	file      *os.File
	startTime time.Time
}

// write logs a message to the shutdown log file (if available).
func (l *shutdownLoggerImpl) write(s string) {
	if l.file != nil {
		msg := fmt.Sprintf("[%s] %s \n", time.Since(l.startTime), s)
		_, err := l.file.WriteString(msg)
		if err != nil {
			log.WithError(err).Error("couldn't write to log file")
		}
		log.Infof("slog: %s", msg)
	} else {
		log.Debug(s)
	}
}

// Close closes the underlying log file if present.
func (l *shutdownLoggerImpl) Close() error {
	return l.file.Close()
}

// TerminateSync inspects and attempts to terminate the process with the given PID.
func (l *shutdownLoggerImpl) TerminateSync(ctx context.Context, pid int) {
	proc, err := procfs.NewProc(pid)
	if err != nil {
		l.write(fmt.Sprintf("Couldn't obtain process information for PID %d.", pid))
		return
	}
	stat, err := proc.Stat()
	if err != nil {
		l.write(fmt.Sprintf("Couldn't obtain process information for PID %d.", pid))
	} else if stat.State == "Z" {
		l.write(fmt.Sprintf("Process %s with PID %d is a zombie, skipping termination.", stat.Comm, pid))
		return
	} else {
		l.write(fmt.Sprintf("Terminating process %s with PID %d (state: %s, cmdlind: %s).", stat.Comm, pid, stat.State, fmt.Sprint(proc.CmdLine())))
	}
	err = process.TerminateSync(ctx, pid)
	if err != nil {
		if errors.Is(err, process.ErrForceKilled) {
			l.write("Terminating process didn't finish, but had to be force killed")
		} else {
			l.write(fmt.Sprintf("Terminating main process errored: %s", err))
		}
	}
}

// extractFailureFromRun attempts to extract the last error message from `supervisor run` command.
func extractFailureFromRun() string {
	logs, err := os.ReadFile("/tmp/opencoder/termination")
	if err != nil {
		// Return empty if unable to read the file
		return ""
	}

	var msg struct {
		Error   string `json:"error"`
		Message string `json:"message"`
	}

	var index int
	var sep = []byte("\n")
	for idx := bytes.LastIndex(logs, sep); idx > 0; idx = index {
		index = bytes.LastIndex(logs[:idx], sep)
		if index < 0 {
			index = 0
		}

		line := logs[index:idx]
		err := json.Unmarshal(line, &msg)
		if err != nil {
			// Skip to the next line if JSON unmarshalling fails
			continue
		}

		if msg.Message == "" {
			// Skip if there's no message in the entry
			continue
		}

		// If error is empty, return just the message; otherwise, concatenate them with a colon
		if msg.Error == "" {
			return msg.Message
		}

		return msg.Message + ": " + msg.Error
	}

	// Return all logs if no valid error entry is found
	return string(logs)
}
