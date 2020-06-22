// +build !windows

package frame

const lineBreak = "\n"

var (
	terminalWidth  int
	terminalHeight int
)

func init() {
	updateScreenDimensions()

	go pollSignals()
}

func updateScreenDimensions() {
	terminalWidth, terminalHeight = getTerminalSize()
}