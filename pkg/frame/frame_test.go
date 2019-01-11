package frame

import (
	"strconv"
	"testing"
)

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func Test_NewFrame(t *testing.T) {

	tables := []struct {
		rows              int
		hasHeader         bool
		hasFooter         bool
		destinationRow    int
		expectedHeaderRow int
		expectedFooterRow int
		expectedLineRows  []int
	}{
		{5, false, false, 10, -1, -1, []int{10, 11, 12, 13, 14}},
		{5, true, false, 10, 10, -1, []int{11, 12, 13, 14, 15}},
		{5, false, true, 10, -1, 15, []int{10, 11, 12, 13, 14}},
		{5, true, true, 10, 10, 16, []int{11, 12, 13, 14, 15}},
	}

	for _, table := range tables {
		suppressOutput(func() {
			frame := New(Config{
				Lines:          table.rows,
				HasHeader:      table.hasHeader,
				HasFooter:      table.hasFooter,
				startRow:       table.destinationRow,
				PositionPolicy: FloatFree,
			})

			// verify the header
			if table.hasHeader && frame.header == nil {
				t.Errorf("NewFrame: expected a header but none was found")
			} else if !table.hasHeader && frame.header != nil {
				t.Errorf("NewFrame: expected no header but one was found (header:%+v)", frame.header)
			}

			// verify the footer
			if table.hasFooter && frame.footer == nil {
				t.Errorf("NewFrame: expected a footer but none was found")
			} else if !table.hasFooter && frame.footer != nil {
				t.Errorf("NewFrame: expected no footer but one was found (footer:%+v)", frame.footer)
			}

			// verify the number of lines in the frame list
			actualLineRowLen := len(frame.activeLines)
			expectedLineRowLen := len(table.expectedLineRows)
			if expectedLineRowLen != actualLineRowLen {
				t.Errorf("NewFrame: expected %d lines, found %d", expectedLineRowLen, actualLineRowLen)
			}

			// ensure the screen row values are correct relative to the given starting row
			var expectedRow, actualRow int

			if table.hasHeader {
				expectedRow = table.expectedHeaderRow
				actualRow = frame.header.row
				if expectedRow != actualRow {
					t.Errorf("NewFrame: expected header row to start at %d, but starts at %d", expectedRow, actualRow)
				}
			}

			if table.hasFooter {
				expectedRow = table.expectedFooterRow
				actualRow = frame.footer.row
				if expectedRow != actualRow {
					t.Errorf("NewFrame: expected footer row to start at %d, but starts at %d", expectedRow, actualRow)
				}
			}

			for idx, expectedRow := range table.expectedLineRows {
				actualRow = frame.activeLines[idx].row
				if expectedRow != actualRow {
					t.Errorf("NewFrame: expected line row to start at %d, but starts at %d", expectedRow, actualRow)
				}
			}
		})

	}

}

func Test_Frame_Height(t *testing.T) {

	tables := []struct {
		rows           int
		hasHeader      bool
		hasFooter      bool
		expectedHeight int
	}{
		{5, false, false, 5},
		{5, true, false, 6},
		{5, false, true, 6},
		{5, true, true, 7},
	}

	for _, table := range tables {
		suppressOutput(func() {
			frame := New(Config{
				Lines:          table.rows,
				HasHeader:      table.hasHeader,
				HasFooter:      table.hasFooter,
				startRow:       10,
				PositionPolicy: FloatFree,
			})
			actualHeight := frame.Height()
			if table.expectedHeight != actualHeight {
				t.Errorf("Frame.height(): expected a height of %d, but found %d", table.expectedHeight, actualHeight)
			}
		})
	}

}

func Test_Frame_IsAtOrPastScreenBottom(t *testing.T) {
	terminalHeight = 100
	frameRows := 10

	tables := []struct {
		destinationRow int
		isAtBottom     bool
	}{
		{50, false},
		{89, false},
		{90, false},
		{91, true},
		{99, true},
		{100, true},
		{101, true},
	}

	for _, table := range tables {
		suppressOutput(func() {
			frame := New(Config{
				Lines:          frameRows,
				HasHeader:      false,
				HasFooter:      false,
				startRow:       table.destinationRow,
				PositionPolicy: FloatFree,
			})
			actualResult := frame.IsPastScreenBottom()
			if table.isAtBottom != actualResult {
				t.Errorf("Frame.IsPastScreenBottom(): expected result of %v, but found %v (startIdx:%d)", table.isAtBottom, actualResult, frame.topRow)
			}
		})
	}

}

func Test_Frame_Append(t *testing.T) {

	tables := []struct {
		appendRows        int
		hasHeader         bool
		hasFooter         bool
		destinationRow    int
		expectedHeaderRow int
		expectedFooterRow int
		expectedLineRows  []int
	}{
		{1, false, false, 10, -1, -1, []int{10}},
		{1, true, false, 10, 10, -1, []int{11}},
		{1, false, true, 10, -1, 11, []int{10}},
		{1, true, true, 10, 10, 12, []int{11}},
		{5, false, false, 10, -1, -1, []int{10, 11, 12, 13, 14}},
		{5, true, false, 10, 10, -1, []int{11, 12, 13, 14, 15}},
		{5, false, true, 10, -1, 15, []int{10, 11, 12, 13, 14}},
		{5, true, true, 10, 10, 16, []int{11, 12, 13, 14, 15}},
	}

	for _, table := range tables {
		suppressOutput(func() {
			frame := New(Config{
				Lines:          0,
				HasHeader:      table.hasHeader,
				HasFooter:      table.hasFooter,
				startRow:       table.destinationRow,
				PositionPolicy: FloatFree,
			})

			// append rows...
			for idx := 0; idx < table.appendRows; idx++ {
				line, err := frame.Append()
				if err != nil {
					t.Errorf("Frame.Append(): expected no error on Append(), got %v", err)
				}
				// write out each index while appending, later check the order
				line.buffer = []byte(strconv.Itoa(idx))
			}

			// check if the number of rows matches
			if len(frame.activeLines) != table.appendRows {
				t.Errorf("Frame.Append(): expected %d number of lines, got %d", table.appendRows, len(frame.activeLines))
			}

			// check the contents of each line
			for idx := 0; idx < table.appendRows; idx++ {
				line := frame.activeLines[idx]
				actualNum, err := strconv.Atoi(string(line.buffer))
				if err != nil {
					t.Errorf("Frame.Append(): expected no error on line read (%v), got %v", actualNum, err)
				}
				if actualNum != idx {
					t.Errorf("Frame.Append(): expected %d, got %d", idx, actualNum)
				}
			}

			// ensure the screen row values are correct relative to the given starting row
			var actualRow int
			for idx, expectedRow := range table.expectedLineRows {
				actualRow = frame.activeLines[idx].row
				if expectedRow != actualRow {
					t.Errorf("Frame.Append(): expected line row to start at %d, but starts at %d", expectedRow, actualRow)
				}
			}
		})
	}
}

func Test_Frame_Prepend(t *testing.T) {

	tables := []struct {
		prependRows       int
		hasHeader         bool
		hasFooter         bool
		destinationRow    int
		expectedHeaderRow int
		expectedFooterRow int
		expectedLineRows  []int
	}{
		{1, false, false, 10, -1, -1, []int{10}},
		{1, true, false, 10, 10, -1, []int{11}},
		{1, false, true, 10, -1, 11, []int{10}},
		{1, true, true, 10, 10, 12, []int{11}},
		{5, false, false, 10, -1, -1, []int{10, 11, 12, 13, 14}},
		{5, true, false, 10, 10, -1, []int{11, 12, 13, 14, 15}},
		{5, false, true, 10, -1, 15, []int{10, 11, 12, 13, 14}},
		{5, true, true, 10, 10, 16, []int{11, 12, 13, 14, 15}},
	}

	for _, table := range tables {
		suppressOutput(func() {
			frame := New(Config{
				Lines:          0,
				HasHeader:      table.hasHeader,
				HasFooter:      table.hasFooter,
				startRow:       table.destinationRow,
				PositionPolicy: FloatFree,
			})

			// append rows...
			for idx := 0; idx < table.prependRows; idx++ {
				line, err := frame.Prepend()
				if err != nil {
					t.Errorf("Frame.Prepend(): expected no error on Prepend(), got %v", err)
				}
				// write out each index while appending, later check the order
				line.buffer = []byte(strconv.Itoa(idx))
			}

			// check if the number of rows matches
			if len(frame.activeLines) != table.prependRows {
				t.Errorf("Frame.Prepend(): expected %d number of lines, got %d", table.prependRows, len(frame.activeLines))
			}

			// check the contents of each line
			for idx := 0; idx < table.prependRows; idx++ {
				// note: indexes should be in reverse
				line := frame.activeLines[table.prependRows-idx-1]
				actualNum, err := strconv.Atoi(string(line.buffer))
				if err != nil {
					t.Errorf("Frame.Prepend(): expected no error on line read (%v), got %v", actualNum, err)
				}
				if actualNum != idx {
					t.Errorf("Frame.Prepend(): expected %d, got %d", idx, actualNum)
				}
			}

			// ensure the screen row values are correct relative to the given starting row
			var actualRow int
			for idx, expectedRow := range table.expectedLineRows {
				actualRow = frame.activeLines[idx].row
				if expectedRow != actualRow {
					t.Errorf("Frame.Prepend(): expected line row to start at %d, but starts at %d", expectedRow, actualRow)
				}
			}
		})
	}
}

func Test_Frame_Insert(t *testing.T) {

	tables := []struct {
		hasHeader         bool
		hasFooter         bool
		destinationRow    int
		expectedHeaderRow int
		expectedFooterRow int
		expectedLineRows  []int
	}{
		{false, false, 10, -1, -1, []int{10, 11, 12, 13, 14}},
		{true, false, 10, 10, -1, []int{11, 12, 13, 14, 15}},
		{false, true, 10, -1, 15, []int{10, 11, 12, 13, 14}},
		{true, true, 10, 10, 16, []int{11, 12, 13, 14, 15}},
	}

	for _, table := range tables {
		suppressOutput(func() {
			frame := New(Config{
				Lines:          4,
				HasHeader:      table.hasHeader,
				HasFooter:      table.hasFooter,
				startRow:       table.destinationRow,
				PositionPolicy: FloatFree,
			})
			finalRows := 5
			insertIdx := 2
			// append rows...

			line, err := frame.Insert(2)
			if err != nil {
				t.Errorf("Frame.Insert(): expected no error on Insert(), got %v", err)
			}
			// write out each index while appending, later check the order
			line.buffer = []byte(strconv.Itoa(insertIdx))

			// check if the number of rows matches
			if len(frame.activeLines) != finalRows {
				t.Errorf("Frame.Insert(): expected %d number of lines, got %d", finalRows, len(frame.activeLines))
			}

			// check the contents of the inserted line

			fetchedLine := frame.activeLines[insertIdx]
			actualNum, err := strconv.Atoi(string(fetchedLine.buffer))
			if err != nil {
				t.Errorf("Frame.Insert(): expected no error on line read (%v), got %v", actualNum, err)
			}
			if actualNum != insertIdx {
				t.Errorf("Frame.Insert(): expected %d, got %d", insertIdx, actualNum)
			}

			// ensure the screen row values are correct relative to the given starting row
			var actualRow int
			for idx, expectedRow := range table.expectedLineRows {
				actualRow = frame.activeLines[idx].row
				if expectedRow != actualRow {
					t.Errorf("Frame.Prepend(): expected line row to start at %d, but starts at %d", expectedRow, actualRow)
				}
			}
		})
	}
}

func Test_Frame_Remove(t *testing.T) {

	tables := []struct {
		startRows        int
		hasHeader        bool
		hasFooter        bool
		destinationRow   int
		expectedLineRows []int
	}{
		{6, false, false, 10, []int{10, 11, 12, 13, 14}},
		{6, true, false, 10, []int{11, 12, 13, 14, 15}},
		{6, false, true, 10, []int{10, 11, 12, 13, 14}},
		{6, true, true, 10, []int{11, 12, 13, 14, 15}},
	}

	for _, table := range tables {
		suppressOutput(func() {
			frame := New(Config{
				Lines:          table.startRows,
				HasHeader:      table.hasHeader,
				HasFooter:      table.hasFooter,
				startRow:       table.destinationRow,
				PositionPolicy: FloatFree,
			})
			rmIdx := 2

			// add content to rows
			for idx, line := range frame.activeLines {
				line.buffer = []byte(strconv.Itoa(idx))
			}

			// Remove a single index
			frame.Remove(frame.activeLines[rmIdx])

			// check if the number of rows matches
			if len(frame.activeLines) != table.startRows-1 {
				t.Errorf("Frame.Remove(): expected %d number of lines, got %d", table.startRows, len(frame.activeLines))
			}

			// check the contents of each line
			for idx := 0; idx < table.startRows-1; idx++ {
				line := frame.activeLines[idx]
				actualNum, err := strconv.Atoi(string(line.buffer))
				if err != nil {
					t.Errorf("Frame.Remove(): expected no error on line read (%v), got %v", actualNum, err)
				}

				if idx < rmIdx && actualNum != idx {
					t.Errorf("Frame.Remove(): expected %d, got %d", idx, actualNum)
				} else if idx >= rmIdx && actualNum != idx+1 {
					t.Errorf("Frame.Remove(): expected %d, got %d", idx, actualNum)
				}
			}

			// ensure the screen row values are correct relative to the given starting row
			var actualRow int
			for idx, expectedRow := range table.expectedLineRows {
				actualRow = frame.activeLines[idx].row
				if expectedRow != actualRow {
					t.Errorf("Frame.Remove(): expected line row to start at %d, but starts at %d (idx:%d)", expectedRow, actualRow, idx)
				}
			}
		})
	}
}

func Test_Frame_Clear(t *testing.T) {

	tables := map[string]struct {
		startRows         int
		hasHeader         bool
		hasFooter         bool
		destinationRow    int
		expectedClearRows []int
	}{
		"goCase":        {5, false, false, 10, []int{10, 11, 12, 13, 14}},
		"Header":        {5, true, false, 10, []int{10, 11, 12, 13, 14, 15}},
		"Footer":        {5, false, true, 10, []int{10, 11, 12, 13, 14, 15}},
		"Header+Footer": {5, true, true, 10, []int{10, 11, 12, 13, 14, 15, 16}},
	}

	for test, table := range tables {
		suppressOutput(func() {
			frame := New(Config{
				Lines:          table.startRows,
				HasHeader:      table.hasHeader,
				HasFooter:      table.hasFooter,
				startRow:       table.destinationRow,
				PositionPolicy: FloatFree,
			})
			frame.clear()

			height := table.startRows
			if table.hasHeader {
				height++
			}
			if table.hasFooter {
				height++
			}

			expectedRows := height
			actualRows := len(frame.clearRows)
			if expectedRows != actualRows {
				t.Errorf("Frame.clear(): [case=%s] expected number of lines cleared to be %d, but is %d", test, expectedRows, actualRows)
			}

			for _, expectedRow := range table.expectedClearRows {
				if !contains(frame.clearRows, expectedRow) {
					t.Errorf("Frame.clear(): [case=%s] expected %d to be in clear rows but is not", test, expectedRow)
				}
			}
		})
	}
}

// todo: rewrite this
func Test_Frame_Close(t *testing.T) {

	tables := map[string]struct {
		startRows      int
		hasHeader      bool
		hasFooter      bool
		destinationRow int
	}{
		"goCase":        {5, false, false, 10},
		"Header":        {5, true, false, 10},
		"Footer":        {5, false, true, 10},
		"Header+Footer": {5, true, true, 10},
	}

	for test, table := range tables {
		suppressOutput(func() {
			frame := New(Config{
				Lines:          table.startRows,
				HasHeader:      table.hasHeader,
				HasFooter:      table.hasFooter,
				startRow:       table.destinationRow,
				PositionPolicy: FloatFree,
			})
			frame.Close()

			if !frame.closed {
				t.Errorf("Frame.Close(): [case=%s] expected frame to be closed but is not", test)
			}

			if table.hasHeader && !frame.header.closed {
				t.Errorf("Frame.Close(): [case=%s] expected header to be closed but is not", test)
			}
			if table.hasFooter && !frame.footer.closed {
				t.Errorf("Frame.Close(): [case=%s] expected footer to be closed but is not", test)
			}

			for idx, line := range frame.activeLines {
				if !line.closed {
					t.Errorf("Frame.Close(): [case=%s] expected line %d to be closed but is not", test, idx)
				}
			}
		})
	}
}

func Test_Frame_Move(t *testing.T) {

	tables := map[string]struct {
		startRows      int
		hasHeader      bool
		hasFooter      bool
		destinationRow int
		moveRows       int
	}{
		"goCase":        {5, false, false, 10, 5},
		"Header":        {5, true, false, 10, -5},
		"Footer":        {5, false, true, 10, 5},
		"Header+Footer": {5, true, true, 10, -5},
	}

	for test, table := range tables {
		suppressOutput(func() {
			frame := New(Config{
				Lines:          table.startRows,
				HasHeader:      table.hasHeader,
				HasFooter:      table.hasFooter,
				startRow:       table.destinationRow,
				PositionPolicy: FloatFree,
			})
			frame.Move(table.moveRows)

			expectedFrameRow := table.destinationRow + table.moveRows
			actualFrameRow := frame.topRow
			if expectedFrameRow != actualFrameRow {
				t.Errorf("Frame.Move(): [case=%s] expected frame to be at %d, but is at %d", test, expectedFrameRow, actualFrameRow)
			}

			headerOffset := 0
			if table.hasHeader {
				headerOffset += 1

				expectedFrameRow = table.destinationRow + table.moveRows
				actualFrameRow = frame.header.row
				if expectedFrameRow != actualFrameRow {
					t.Errorf("Frame.Move(): [case=%s] expected header to be at %d, but is at %d", test, expectedFrameRow, actualFrameRow)
				}
			}

			if table.hasFooter {
				expectedFrameRow = table.destinationRow + table.startRows + headerOffset + table.moveRows
				actualFrameRow = frame.footer.row
				if expectedFrameRow != actualFrameRow {
					t.Errorf("Frame.Move(): [case=%s] expected footer to be at %d, but is at %d", test, expectedFrameRow, actualFrameRow)
				}
			}

			for idx, line := range frame.activeLines {
				expectedFrameRow = table.destinationRow + idx + headerOffset + table.moveRows
				actualFrameRow = line.row
				if expectedFrameRow != actualFrameRow {
					t.Errorf("Frame.Move(): [case=%s] expected line to be at %d, but is at %d", test, expectedFrameRow, actualFrameRow)
				}
			}
		})
	}
}

// func Test_Frame_Update(t *testing.T) {
// 	terminalHeight = 100
// 	frameRows := 10
//
// 	fn := func(frame *Frame) error {
// 		// if the frame has moved past the bottom of the screen, move it up a bit
// 		if frame.IsPastScreenBottom() {
// 			frameHeight := frame.visibleHeight()
// 			// offset is how many rows the frame needs to be adjusted to fit on the screen.
// 			// This is the same as how many rows past the edge of the screen this frame currently is.
// 			offset := (frame.topRow + frameHeight) - terminalHeight
// 			// offset += 1 // we want to move one line past the frame
// 			frame.Move(-offset)
// 			frame.rowAdvancements += offset
// 		}
//
// 		// if the frame has moved above the top of the screen, move it down a bit
// 		if frame.IsPastScreenTop() {
// 			offset := -1*frame.topRow + 1
// 			frame.Move(offset)
// 		}
// 		return nil
// 	}
//
// 	tables := []struct {
// 		destinationRow int
// 		adjustedRow    int
// 	}{
// 		{-20, 1},
// 		{-2, 1},
// 		{0, 1},
// 		{1, 1},
// 		{50, 50},
// 		{89, 89},
// 		{90, 90},
// 		{91, 90},
// 		{99, 90},
// 		{100, 90},
// 		{110, 90},
// 	}
//
// 	for _, table := range tables {
// 		frame := New(Config{
// 			Lines:     frameRows,
// 			HasHeader: false,
// 			HasFooter: false,
// 			startRow:  table.destinationRow,
// 		})
// 		frame.updateFn = fn
// 		frame.Update()
// 		actualResult := frame.topRow
// 		if table.adjustedRow != actualResult {
// 			t.Errorf("Frame.Update(): expected Update row of %d, but is at row %d", table.adjustedRow, actualResult)
// 		}
// 	}
//
// }
