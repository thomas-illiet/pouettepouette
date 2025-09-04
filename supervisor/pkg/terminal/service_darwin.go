package terminal

import "os/exec"

// setAmbientCaps is not supported on Darwin.
func (srv *MuxTerminalService) setAmbientCaps(cmd *exec.Cmd) {}
