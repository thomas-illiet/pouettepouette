package utils

import (
	"os"
	"syscall"

	"golang.org/x/term"
)

func isInteractiveTerminal() bool {
	return term.IsTerminal(syscall.Stdin) && term.IsTerminal(syscall.Stdout)
}

func userRequestsNoColor() bool {
	isDumbTerm := os.Getenv("TERM") == "dumb"
	_, noColorPresent := os.LookupEnv("NO_COLOR")
	_, gpNoColor := os.LookupEnv("GP_NO_COLOR")

	return isDumbTerm || noColorPresent || gpNoColor
}

func ColorsEnabled() bool {
	colorsDisabled := userRequestsNoColor() || !isInteractiveTerminal()
	return !colorsDisabled
}
