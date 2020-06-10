// +build !windows

package frame

import (
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
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
			scr := getScreen()
			terminalWidth, terminalHeight = getTerminalSize()
			lock := scr.lock
			lock.Lock()
			scr.refresh()
			lock.Unlock()
		}
	}
}
