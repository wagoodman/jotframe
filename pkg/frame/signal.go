// +build !windows

package frame

import (
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"os/signal"
	"syscall"
)

var (
	sigwinch = make(chan os.Signal)
)

type terminalSize struct {
	rows    uint16
	cols    uint16
	xPixels uint16
	yPixels uint16
}

func GetTerminalSize() (int, int) {
	return terminalWidth, terminalHeight
}

func getTerminalSize() (int, int) {
	termWidth, termHeight, _ := terminal.GetSize(int(os.Stdout.Fd()))
	return termWidth, termHeight
}

func pollSignals() {
	// set signal handlers
	signal.Notify(sigwinch, syscall.SIGWINCH)

	// watch for events
	for {
		select {
		case <-sigwinch:
			terminalWidth, terminalHeight = getTerminalSize()
			lock := getScreenLock()
			lock.Lock()
			// the screen may have a trail, which is unmanaged at this point. Don't clear the screen
			// clearScreen()
			refresh()
			lock.Unlock()
		}
	}
}
