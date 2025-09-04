package terminal

import "golang.org/x/sys/unix"

// setTermAttr applies terminal attributes on macOS (Darwin) using the TIOCSETA ioctl request.
// Unlike Linux, Darwin does not support TCSETS, so the equivalent request is TIOCSETA.
func setTermAttr(fd int, attr *unix.Termios) error {
	return unix.IoctlSetTermios(fd, unix.TIOCSETA, attr)
}
