package utils

import (
	"os"
	"runtime"

	isatty "github.com/mattn/go-isatty"
)

func DetectTerminal() bool {
	var session bool
	if runtime.GOOS == "windows" {
		// Detect the newer Windows Terminal
		_, session = os.LookupEnv("WT_SESSION")
	} else {
		session = false
	}
	// Detect Unix and Cygwin Terminals
	if isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd()) || session == true {
		return true
	} else {
		return false
	}
}
