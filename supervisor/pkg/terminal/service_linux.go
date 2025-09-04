package terminal

// setAmbientCaps applies DefaultAmbientCaps on Linux.
func (srv *MuxTerminalService) setAmbientCaps(cmd *exec.Cmd) {
	if srv.DefaultAmbientCaps != nil {
		if cmd.SysProcAttr == nil {
			cmd.SysProcAttr = &syscall.SysProcAttr{}
		}
		cmd.SysProcAttr.AmbientCaps = srv.DefaultAmbientCaps
	}
}
