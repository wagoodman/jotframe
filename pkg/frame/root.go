// +build !windows

package frame

const lineBreak = "\n"

var (
	terminalWidth  int
	terminalHeight int
)

func init() {
	// fetch initial values
	terminalWidth, terminalHeight = getTerminalSize()

	go pollSignals()
}