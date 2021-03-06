package ui

import (
	"io"
	"os"

	"github.com/fatih/color"
	colorable "github.com/mattn/go-colorable"
	isatty "github.com/mattn/go-isatty"
)

// windowsIO is the Windows-specific implementation of the IO interface.
type windowsIO struct {
	standardIO
	coloredOutput io.Writer
}

// NewUserIO creates a new windowsIO.
func NewUserIO() IO {
	return windowsIO{
		standardIO:    newStdUserIO(),
		coloredOutput: colorable.NewColorable(os.Stdout),
	}
}

// Stdout returns the standardIO's Output.
func (o windowsIO) Output() io.Writer {
	if !color.NoColor {
		return o.coloredOutput
	}
	return o.output
}

// eofKey returns the key(s) that should be pressed to enter an EOF.
func eofKey() string {
	return "CTRL-Z + ENTER"
}

// isPiped checks whether the file is a pipe.
// If the file does not exist, it returns false.
func isPiped(file *os.File) bool {
	_, err := file.Stat()
	if err != nil {
		return false
	}

	return os.Getenv("TERM") == "dumb" ||
		(!isatty.IsTerminal(file.Fd()) && !isatty.IsCygwinTerminal(file.Fd()))
}
