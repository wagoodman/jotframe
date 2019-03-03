package frame

import (
	"os"
	"time"

	"golang.org/x/crypto/ssh/terminal"
)

func GetTerminalSize() (int, int) {
	return terminalWidth, terminalHeight
}

func getTerminalSize() (int, int) {
	termWidth, termHeight, _ := terminal.GetSize(int(os.Stdout.Fd()))
	return termWidth, termHeight
}

func pollSignals() {

	// poll window size (how do you so this event)
	for {
		terminalWidth, terminalHeight = getTerminalSize()

		lock := getScreenLock()
		lock.Lock()
		// the screen may have a trail, which is unmanaged at this point. Don't clear the screen
		// clearScreen()
		refresh()
		lock.Unlock()

		time.Sleep(1 * time.Second)
	}
}
