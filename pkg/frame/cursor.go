// +build !windows

package frame

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

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
