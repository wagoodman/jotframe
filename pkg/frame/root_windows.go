package frame

import (
	"syscall"
)

const lineBreak = "\r\n"

var (
	outHandle      syscall.Handle
	terminalWidth  int
	terminalHeight int
)

func init() {
	outHandle, _ = syscall.Open("CONOUT$", syscall.O_RDWR, 0)

	// fetch initial values
	terminalWidth, terminalHeight = getTerminalSize()

	go pollSignals()
}
