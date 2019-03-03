package frame

import (
	"fmt"
	"strings"
	"sync"

	"github.com/k0kubun/go-ansi"
)

var (
	lockSync   sync.Once
	screenLock *sync.Mutex
)

func newScreenEvent(line *Line) *ScreenEvent {
	e := &ScreenEvent{
		row:   line.row,
		value: make([]byte, len(line.buffer)),
	}
	copy(e.value, line.buffer)
	return e
}

// TODO: this should be a ScreenEvent!
func advanceScreen(rows int) {
	setCursorRow(terminalHeight)
	fmt.Print(strings.Repeat(lineBreak, rows))
}

func writeAtRow(message string, row int) {
	setCursorRow(row)
	fmt.Print(strings.Replace(message, lineBreak, "", -1))
}

// todo: assumes VT100, not cross platform
func clearScreen() {
	fmt.Print("\x1b[2J")
}

// TODO: this should be a ScreenEvent!
func clearRow(row int) error {
	err := setCursorRow(row)
	if err != nil {
		panic(err)
	}
	ansi.EraseInLine(2)
	return nil
}

func getScreenLock() *sync.Mutex {
	lockSync.Do(func() {
		screenLock = &sync.Mutex{}
	})
	return screenLock
}
