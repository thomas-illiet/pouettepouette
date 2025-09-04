package terminal

// setTermAttr applies terminal attributes on Linux using the TCSETS ioctl request.
// This is a thin wrapper around unix.IoctlSetTermios to allow per-OS differences
// (Darwin uses a different ioctl code).
func setTermAttr(fd int, attr *unix.Termios) error {
	return unix.IoctlSetTermios(fd, syscall.TCSETS, attr)
}
