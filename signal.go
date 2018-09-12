package jotframe

import (
	"os"
	"syscall"
	"unsafe"
	"os/signal"
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
	var obj terminalSize
	_, _, _ = syscall.Syscall(syscall.SYS_IOCTL, os.Stdout.Fd(), uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(&obj)))
	return int(obj.cols), int(obj.rows)
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
			clearScreen()
			refresh()
			lock.Unlock()
		}
	}
}