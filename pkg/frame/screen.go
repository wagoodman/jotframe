package frame

import (
	"bufio"
	"fmt"
	"github.com/k0kubun/go-ansi"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
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
	fmt.Print(strings.Repeat("\n", rows))
}

func writeAtRow(message string, row int) {
	setCursorRow(row)
	fmt.Print(strings.Replace(message, "\n", "", -1))
}

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

// todo: will this be supported on windows?... https://github.com/nsf/termbox-go/blob/master/termbox_windows.go
// currently assumed VT100 compatible emulator
func setCursorRow(row int) error {
	// todo: is this "really" needed?
	// if isatty.IsTerminal(os.Stdin.Fd()) {
	// 	oldState, err := terminal.MakeRaw(0)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	defer terminal.Restore(0, oldState)
	// }

	// sets the cursor position where subsequent text will begin: <ESC>[{ROW};{COLUMN}H
	// great resource: http://www.termsys.demon.co.uk/vtansi.htm
	fmt.Printf("\x1b[%d;0H", row)
	return nil
}

// todo: will this be supported on windows?... https://github.com/nsf/termbox-go/blob/master/termbox_windows.go
// currently assumed VT100 compatible emulator
func GetCursorRow() (int, error) {
	var row int
	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		return -1, err
	}
	defer terminal.Restore(0, oldState)

	// capture keyboard output from echo
	reader := bufio.NewReader(os.Stdin)

	// request a "Report Cursor Position" response from the device: <ESC>[{ROW};{COLUMN}R
	// great resource: http://www.termsys.demon.co.uk/vtansi.htm
	fmt.Print("\x1b[6n")

	// capture the response up until the expected "R"
	text, err := reader.ReadSlice('R')
	if err != nil {
		return -1, fmt.Errorf("unable to read stdin")
	}

	// parse the row and column
	if strings.Contains(string(text), ";") {
		re := regexp.MustCompile(`\d+;\d+`)
		line := re.FindString(string(text))
		row, err = strconv.Atoi(strings.Split(line, ";")[0])

		if err != nil {
			return -1, fmt.Errorf("invalid row value: '%s'", line)
		}

	} else {
		return -1, fmt.Errorf("unable to fetch position")
	}

	return row, nil
}
