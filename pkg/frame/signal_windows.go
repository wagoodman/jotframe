package frame

import (
	"time"

	"golang.org/x/term"
)

func GetTerminalSize() (int, int) {
	return terminalWidth, terminalHeight
}

func getTerminalSize() (int, int) {
	termWidth, termHeight, _ := term.GetSize(int(getScreen().output.Fd()))
	return termWidth, termHeight
}

func pollSignals() {

	// TODO: is there a way to make this event driven?
	for {
		terminalWidth, terminalHeight = getTerminalSize()

		lock := getScreenLock()
		lock.Lock()
		refresh()
		lock.Unlock()

		time.Sleep(1 * time.Second)
	}
}
