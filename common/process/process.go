package process

import (
	"context"
	"errors"
	"golang.org/x/sys/unix"
	"os"
)

var ErrForceKilled = errors.New("process didn't terminate, so we sent SIGKILL")

// TerminateSync terminates the specified process using SIGTERM.
// It waits until the process has terminated or until the context is cancelled.
// If the context is cancelled, it sends a SIGKILL to forcefully terminate the process and returns an error.
func TerminateSync(ctx context.Context, pid int) error {
	// Find the process by its PID. This should never fail on UNIX systems.
	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	// Send a SIGTERM to the process.
	err = process.Signal(unix.SIGTERM)
	if err != nil {
		if errors.Is(err, os.ErrProcessDone) {
			return nil
		}
		return err
	}

	// Start a goroutine that waits for the process to terminate.
	terminated := make(chan error, 1)
	go func() {
		_, err := process.Wait()
		terminated <- err
	}()

	// Wait for either the termination channel or context cancellation.
	select {
	case err = <-terminated:
		return err
	case <-ctx.Done():
		err = process.Kill()
		if err != nil {
			return err
		}
		return ErrForceKilled
	}
}
