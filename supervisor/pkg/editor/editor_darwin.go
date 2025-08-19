package editor

import (
	"os/exec"
	"syscall"
)

// prepareSysProc configures an *exec.Cmd to run in its own process group.
func prepareSysProc(cmd *exec.Cmd) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}

	// We need the child process to run in its own process group, s.t.
	// we can suspend and resume editor and its children.
	cmd.SysProcAttr.Setpgid = true
}
