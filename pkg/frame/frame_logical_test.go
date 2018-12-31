package frame

import (
	"fmt"
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

func Test_NewLogicalFrame(t *testing.T) {

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
		frame := newLogicalFrame(Config{
			Lines:     table.rows,
			HasHeader: table.hasHeader,
			HasFooter: table.hasFooter,
			startRow:  table.destinationRow,
		})

		// verify the header
		if table.hasHeader && frame.header == nil {
			t.Errorf("NewLogicalFrame: expected a header but none was found")
		} else if !table.hasHeader && frame.header != nil {
			t.Errorf("NewLogicalFrame: expected no header but one was found (header:%+v)", frame.header)
		}

		// verify the footer
		if table.hasFooter && frame.footer == nil {
			t.Errorf("NewLogicalFrame: expected a footer but none was found")
		} else if !table.hasFooter && frame.footer != nil {
			t.Errorf("NewLogicalFrame: expected no footer but one was found (footer:%+v)", frame.footer)
		}

		// verify the number of lines in the frame list
		actualLineRowLen := len(frame.activeLines)
		expectedLineRowLen := len(table.expectedLineRows)
		if expectedLineRowLen != actualLineRowLen {
			t.Errorf("NewLogicalFrame: expected %d lines, found %d", expectedLineRowLen, actualLineRowLen)
		}

		// ensure the screen row values are correct relative to the given starting row
		var expectedRow, actualRow int

		if table.hasHeader {
			expectedRow = table.expectedHeaderRow
			actualRow = frame.header.row
			if expectedRow != actualRow {
				t.Errorf("NewLogicalFrame: expected header row to start at %d, but starts at %d", expectedRow, actualRow)
			}
		}

		if table.hasFooter {
			expectedRow = table.expectedFooterRow
			actualRow = frame.footer.row
			if expectedRow != actualRow {
				t.Errorf("NewLogicalFrame: expected footer row to start at %d, but starts at %d", expectedRow, actualRow)
			}
		}

		for idx, expectedRow := range table.expectedLineRows {
			actualRow = frame.activeLines[idx].row
			if expectedRow != actualRow {
				t.Errorf("NewLogicalFrame: expected line row to start at %d, but starts at %d", expectedRow, actualRow)
			}
		}

	}

}

func Test_LogicalFrame_Height(t *testing.T) {

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
		frame := newLogicalFrame(Config{
			Lines:     table.rows,
			HasHeader: table.hasHeader,
			HasFooter: table.hasFooter,
			startRow:  10,
		})
		actualHeight := frame.Height()
		if table.expectedHeight != actualHeight {
			t.Errorf("LogicalFrame.height(): expected a height of %d, but found %d", table.expectedHeight, actualHeight)
		}
	}

}

func Test_LogicalFrame_IsAtOrPastScreenBottom(t *testing.T) {
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
		frame := newLogicalFrame(Config{
			Lines:     frameRows,
			HasHeader: false,
			HasFooter: false,
			startRow:  table.destinationRow,
		})
		actualResult := frame.IsPastScreenBottom()
		if table.isAtBottom != actualResult {
			t.Errorf("LogicalFrame.IsPastScreenBottom(): expected result of %v, but found %v (startIdx:%d)", table.isAtBottom, actualResult, frame.topRow)
		}
	}

}

func Test_LogicalFrame_Append(t *testing.T) {

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
		frame := newLogicalFrame(Config{
			Lines:     0,
			HasHeader: table.hasHeader,
			HasFooter: table.hasFooter,
			startRow:  table.destinationRow,
		})

		// append rows...
		for idx := 0; idx < table.appendRows; idx++ {
			line, err := frame.Append()
			if err != nil {
				t.Errorf("LogicalFrame.Append(): expected no error on Append(), got %v", err)
			}
			// write out each index while appending, later check the order
			line.buffer = []byte(strconv.Itoa(idx))
		}

		// check if the number of rows matches
		if len(frame.activeLines) != table.appendRows {
			t.Errorf("LogicalFrame.Append(): expected %d number of lines, got %d", table.appendRows, len(frame.activeLines))
		}

		// check the contents of each line
		for idx := 0; idx < table.appendRows; idx++ {
			line := frame.activeLines[idx]
			actualNum, err := strconv.Atoi(string(line.buffer))
			if err != nil {
				t.Errorf("LogicalFrame.Append(): expected no error on line read (%v), got %v", actualNum, err)
			}
			if actualNum != idx {
				t.Errorf("LogicalFrame.Append(): expected %d, got %d", idx, actualNum)
			}
		}

		// ensure the screen row values are correct relative to the given starting row
		var actualRow int
		for idx, expectedRow := range table.expectedLineRows {
			actualRow = frame.activeLines[idx].row
			if expectedRow != actualRow {
				t.Errorf("LogicalFrame.Append(): expected line row to start at %d, but starts at %d", expectedRow, actualRow)
			}
		}
	}
}

func Test_LogicalFrame_Prepend(t *testing.T) {

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
		frame := newLogicalFrame(Config{
			Lines:     0,
			HasHeader: table.hasHeader,
			HasFooter: table.hasFooter,
			startRow:  table.destinationRow,
		})

		// append rows...
		for idx := 0; idx < table.prependRows; idx++ {
			line, err := frame.Prepend()
			if err != nil {
				t.Errorf("LogicalFrame.Prepend(): expected no error on Prepend(), got %v", err)
			}
			// write out each index while appending, later check the order
			line.buffer = []byte(strconv.Itoa(idx))
		}

		// check if the number of rows matches
		if len(frame.activeLines) != table.prependRows {
			t.Errorf("LogicalFrame.Prepend(): expected %d number of lines, got %d", table.prependRows, len(frame.activeLines))
		}

		// check the contents of each line
		for idx := 0; idx < table.prependRows; idx++ {
			// note: indexes should be in reverse
			line := frame.activeLines[table.prependRows-idx-1]
			actualNum, err := strconv.Atoi(string(line.buffer))
			if err != nil {
				t.Errorf("LogicalFrame.Prepend(): expected no error on line read (%v), got %v", actualNum, err)
			}
			if actualNum != idx {
				t.Errorf("LogicalFrame.Prepend(): expected %d, got %d", idx, actualNum)
			}
		}

		// ensure the screen row values are correct relative to the given starting row
		var actualRow int
		for idx, expectedRow := range table.expectedLineRows {
			actualRow = frame.activeLines[idx].row
			if expectedRow != actualRow {
				t.Errorf("LogicalFrame.Prepend(): expected line row to start at %d, but starts at %d", expectedRow, actualRow)
			}
		}
	}
}

func Test_LogicalFrame_Insert(t *testing.T) {

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
		frame := newLogicalFrame(Config{
			Lines:     4,
			HasHeader: table.hasHeader,
			HasFooter: table.hasFooter,
			startRow:  table.destinationRow,
		})
		finalRows := 5
		insertIdx := 2
		// append rows...

		line, err := frame.Insert(2)
		if err != nil {
			t.Errorf("LogicalFrame.Insert(): expected no error on Insert(), got %v", err)
		}
		// write out each index while appending, later check the order
		line.buffer = []byte(strconv.Itoa(insertIdx))

		// check if the number of rows matches
		if len(frame.activeLines) != finalRows {
			t.Errorf("LogicalFrame.Insert(): expected %d number of lines, got %d", finalRows, len(frame.activeLines))
		}

		// check the contents of the inserted line

		fetchedLine := frame.activeLines[insertIdx]
		actualNum, err := strconv.Atoi(string(fetchedLine.buffer))
		if err != nil {
			t.Errorf("LogicalFrame.Insert(): expected no error on line read (%v), got %v", actualNum, err)
		}
		if actualNum != insertIdx {
			t.Errorf("LogicalFrame.Insert(): expected %d, got %d", insertIdx, actualNum)
		}

		// ensure the screen row values are correct relative to the given starting row
		var actualRow int
		for idx, expectedRow := range table.expectedLineRows {
			actualRow = frame.activeLines[idx].row
			if expectedRow != actualRow {
				t.Errorf("LogicalFrame.Prepend(): expected line row to start at %d, but starts at %d", expectedRow, actualRow)
			}
		}
	}
}

func Test_LogicalFrame_Remove(t *testing.T) {

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
		frame := newLogicalFrame(Config{
			Lines:     table.startRows,
			HasHeader: table.hasHeader,
			HasFooter: table.hasFooter,
			startRow:  table.destinationRow,
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
			t.Errorf("LogicalFrame.Remove(): expected %d number of lines, got %d", table.startRows, len(frame.activeLines))
		}

		// check the contents of each line
		for idx := 0; idx < table.startRows-1; idx++ {
			line := frame.activeLines[idx]
			actualNum, err := strconv.Atoi(string(line.buffer))
			if err != nil {
				t.Errorf("LogicalFrame.Remove(): expected no error on line read (%v), got %v", actualNum, err)
			}

			if idx < rmIdx && actualNum != idx {
				t.Errorf("LogicalFrame.Remove(): expected %d, got %d", idx, actualNum)
			} else if idx >= rmIdx && actualNum != idx+1 {
				t.Errorf("LogicalFrame.Remove(): expected %d, got %d", idx, actualNum)
			}
		}

		// ensure the screen row values are correct relative to the given starting row
		var actualRow int
		for idx, expectedRow := range table.expectedLineRows {
			actualRow = frame.activeLines[idx].row
			if expectedRow != actualRow {
				t.Errorf("LogicalFrame.Remove(): expected line row to start at %d, but starts at %d (idx:%d)", expectedRow, actualRow, idx)
			}
		}
	}
}

func Test_LogicalFrame_Clear(t *testing.T) {

	tables := []struct {
		startRows         int
		hasHeader         bool
		hasFooter         bool
		destinationRow    int
		expectedClearRows []int
	}{
		{5, false, false, 10, []int{10, 11, 12, 13, 14}},
		{5, true, false, 10, []int{10, 11, 12, 13, 14, 15}},
		{5, false, true, 10, []int{10, 11, 12, 13, 14, 15}},
		{5, true, true, 10, []int{10, 11, 12, 13, 14, 15, 16}},
	}

	for _, table := range tables {
		frame := newLogicalFrame(Config{
			Lines:     table.startRows,
			HasHeader: table.hasHeader,
			HasFooter: table.hasFooter,
			startRow:  table.destinationRow,
		})
		err := frame.Clear()
		if err != nil {
			t.Errorf("LogicalFrame.Clear(): expected no error on Clear(), got %v", err)
		}

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
			t.Errorf("LogicalFrame.Clear(): expected number of lines cleared to be %d, but is %d", expectedRows, actualRows)
		}

		for _, expectedRow := range table.expectedClearRows {
			if !contains(frame.clearRows, expectedRow) {
				t.Errorf("LogicalFrame.Clear(): expected %d to be in Clear rows but is not", expectedRow)
			}
		}

	}
}

// todo: rewrite this
func Test_LogicalFrame_Close(t *testing.T) {

	tables := []struct {
		startRows      int
		hasHeader      bool
		hasFooter      bool
		destinationRow int
	}{
		{5, false, false, 10},
		{5, true, false, 10},
		{5, false, true, 10},
		{5, true, true, 10},
	}

	for _, table := range tables {
		frame := newLogicalFrame(Config{
			Lines:     table.startRows,
			HasHeader: table.hasHeader,
			HasFooter: table.hasFooter,
			startRow:  table.destinationRow,
		})
		err := frame.Close()
		if err != nil {
			t.Errorf("LogicalFrame.Close(): expected no error on Close(), got %v", err)
		}

		if !frame.closed {
			t.Errorf("LogicalFrame.Close(): expected frame to be closed but is not")
		}

		if table.hasHeader && !frame.header.closed {
			t.Errorf("LogicalFrame.Close(): expected header to be closed but is not")
		}
		if table.hasFooter && !frame.footer.closed {
			t.Errorf("LogicalFrame.Close(): expected footer to be closed but is not")
		}

		for idx, line := range frame.activeLines {
			if !line.closed {
				t.Errorf("LogicalFrame.Close(): expected line %d to be closed but is not", idx)
			}
		}
	}
}

func Test_LogicalFrame_Move(t *testing.T) {

	tables := []struct {
		startRows      int
		hasHeader      bool
		hasFooter      bool
		destinationRow int
		moveRows       int
	}{
		{5, false, false, 10, 5},
		{5, true, false, 10, -5},
		{5, false, true, 10, 5},
		{5, true, true, 10, -5},
	}

	for _, table := range tables {
		frame := newLogicalFrame(Config{
			Lines:     table.startRows,
			HasHeader: table.hasHeader,
			HasFooter: table.hasFooter,
			startRow:  table.destinationRow,
		})
		frame.Move(table.moveRows)

		expectedFrameRow := table.destinationRow + table.moveRows
		actualFrameRow := frame.topRow
		if expectedFrameRow != actualFrameRow {
			t.Errorf("LogicalFrame.Move(): expected frame to be at %d, but is at %d", expectedFrameRow, actualFrameRow)
		}

		headerOffset := 0
		if table.hasHeader {
			headerOffset += 1

			expectedFrameRow = table.destinationRow + table.moveRows
			actualFrameRow = frame.header.row
			if expectedFrameRow != actualFrameRow {
				t.Errorf("LogicalFrame.Move(): expected header to be at %d, but is at %d", expectedFrameRow, actualFrameRow)
			}
		}

		if table.hasFooter {
			expectedFrameRow = table.destinationRow + table.startRows + headerOffset + table.moveRows
			actualFrameRow = frame.footer.row
			if expectedFrameRow != actualFrameRow {
				t.Errorf("LogicalFrame.Move(): expected footer to be at %d, but is at %d", expectedFrameRow, actualFrameRow)
			}
		}

		for idx, line := range frame.activeLines {
			expectedFrameRow = table.destinationRow + idx + headerOffset + table.moveRows
			actualFrameRow = line.row
			if expectedFrameRow != actualFrameRow {
				t.Errorf("LogicalFrame.Move(): expected line to be at %d, but is at %d", expectedFrameRow, actualFrameRow)
			}
		}
	}
}

func Test_LogicalFrame_Update(t *testing.T) {
	terminalHeight = 100
	frameRows := 10

	fn := func(frame *logicalFrame) error {
		// if the frame has moved past the bottom of the screen, move it up a bit
		if frame.IsPastScreenBottom() {
			frameHeight := frame.visibleHeight()
			// offset is how many rows the frame needs to be adjusted to fit on the screen.
			// This is the same as how many rows past the edge of the screen this frame currently is.
			offset := (frame.topRow + frameHeight) - terminalHeight
			// offset += 1 // we want to move one line past the frame
			frame.Move(-offset)
			frame.rowAdvancements += offset
		}

		// if the frame has moved above the top of the screen, move it down a bit
		if frame.IsPastScreenTop() {
			offset := -1*frame.topRow + 1
			frame.Move(offset)
		}
		return nil
	}

	tables := []struct {
		destinationRow int
		adjustedRow    int
	}{
		{-20, 1},
		{-2, 1},
		{0, 1},
		{1, 1},
		{50, 50},
		{89, 89},
		{90, 90},
		{91, 90},
		{99, 90},
		{100, 90},
		{110, 90},
	}

	for _, table := range tables {
		frame := newLogicalFrame(Config{
			Lines:     frameRows,
			HasHeader: false,
			HasFooter: false,
			startRow:  table.destinationRow,
		})
		frame.updateFn = fn
		frame.Update()
		actualResult := frame.topRow
		if table.adjustedRow != actualResult {
			t.Errorf("LogicalFrame.Update(): expected Update row of %d, but is at row %d", table.adjustedRow, actualResult)
		}
	}

}

var drawTestCases = map[string]drawTestParams{
	"goCase": {3, false, false, 10, 40,
		[]ScreenEvent{
			{row: 10, value: []byte("LineIdx:0")},
			{row: 11, value: []byte("LineIdx:1")},
			{row: 12, value: []byte("LineIdx:2")},
		},
		[]string{},
	},
	"goCase_Header": {3, true, false, 10, 40,
		[]ScreenEvent{
			{row: 10, value: []byte("theHeader")},
			{row: 11, value: []byte("LineIdx:0")},
			{row: 12, value: []byte("LineIdx:1")},
			{row: 13, value: []byte("LineIdx:2")},
		},
		[]string{},
	},
	"goCase_Footer": {3, false, true, 10, 40,
		[]ScreenEvent{
			{row: 10, value: []byte("LineIdx:0")},
			{row: 11, value: []byte("LineIdx:1")},
			{row: 12, value: []byte("LineIdx:2")},
			{row: 13, value: []byte("theFooter")},
		},
		[]string{},
	},
	"goCase_HeaderFooter": {3, true, true, 10, 40,
		[]ScreenEvent{
			{row: 10, value: []byte("theHeader")},
			{row: 11, value: []byte("LineIdx:0")},
			{row: 12, value: []byte("LineIdx:1")},
			{row: 13, value: []byte("LineIdx:2")},
			{row: 14, value: []byte("theFooter")},
		},
		[]string{},
	},
	"termHeightSmall_Top": {3, false, false, 1, 2,
		[]ScreenEvent{
			{row: 1, value: []byte("LineIdx:0")},
			{row: 2, value: []byte("LineIdx:1")},
		},
		[]string{
			"line is out of bounds (row=3)",
		},
	},
	"termHeightSmall_Top_Header": {3, true, false, 1, 2,
		[]ScreenEvent{
			{row: 1, value: []byte("theHeader")},
			{row: 2, value: []byte("LineIdx:0")},
		},
		[]string{
			"line is out of bounds (row=3)",
			"line is out of bounds (row=4)",
		},
	},
	"termHeightSmall_Top_Footer": {3, false, true, 1, 2,
		[]ScreenEvent{
			{row: 1, value: []byte("LineIdx:0")},
			{row: 2, value: []byte("LineIdx:1")},
		},
		[]string{
			"line is out of bounds (row=3)",
			"line is out of bounds (row=4)",
		},
	},
	"termHeightSmall_Top_HeaderFooter": {3, true, true, 1, 2,
		[]ScreenEvent{
			{row: 1, value: []byte("theHeader")},
			{row: 2, value: []byte("LineIdx:0")},
		},
		[]string{
			"line is out of bounds (row=3)",
			"line is out of bounds (row=4)",
			"line is out of bounds (row=5)",
		},
	},
	"termHeightSmall_Bottom": {3, false, false, 49, 50,
		[]ScreenEvent{
			{row: 49, value: []byte("LineIdx:0")},
			{row: 50, value: []byte("LineIdx:1")},
		},
		[]string{
			"line is out of bounds (row=51)",
		},
	},
	"termHeightSmall_Bottom_Header": {3, true, false, 49, 50,
		[]ScreenEvent{
			{row: 49, value: []byte("theHeader")},
			{row: 50, value: []byte("LineIdx:0")},
		},
		[]string{
			"line is out of bounds (row=51)",
			"line is out of bounds (row=52)",
		},
	},
	"termHeightSmall_Bottom_Footer": {3, false, true, 49, 50,
		[]ScreenEvent{
			{row: 49, value: []byte("LineIdx:0")},
			{row: 50, value: []byte("LineIdx:1")},
		},
		[]string{
			"line is out of bounds (row=51)",
			"line is out of bounds (row=52)",
		},
	},
	"termHeightSmall_Bottom_HeaderFooter": {3, true, true, 49, 50,
		[]ScreenEvent{
			{row: 49, value: []byte("theHeader")},
			{row: 50, value: []byte("LineIdx:0")},
		},
		[]string{
			"line is out of bounds (row=51)",
			"line is out of bounds (row=52)",
			"line is out of bounds (row=53)",
		},
	},
}

func Test_LogicalFrame_Draw(t *testing.T) {

	for test, table := range drawTestCases {
		suppressOutput(func() {
			// setup...
			terminalHeight = table.terminalHeight
			handler := NewTestEventHandler(t)
			screenHandlers = make([]ScreenEventHandler, 0)
			addScreenHandler(handler)

			// run test...
			var errs []error
			frame := newLogicalFrame(Config{
				Lines:     table.rows,
				HasHeader: table.hasHeader,
				HasFooter: table.hasFooter,
				startRow:  table.startRow,
			})
			if table.hasHeader {
				frame.header.buffer = []byte("theHeader")
			}
			for idx, line := range frame.activeLines {
				line.buffer = []byte(fmt.Sprintf("LineIdx:%d", idx))
			}
			if table.hasFooter {
				frame.footer.buffer = []byte("theFooter")
			}
			errs = frame.Draw()

			// assert results...
			validateEvents(t, test, table, errs, frame, handler)

		})
	}

}

func Test_LogicalFrame_AdhocDraw(t *testing.T) {

	for test, table := range drawTestCases {
		suppressOutput(func() {
			// setup...
			terminalHeight = table.terminalHeight
			handler := NewTestEventHandler(t)
			screenHandlers = make([]ScreenEventHandler, 0)
			addScreenHandler(handler)

			// run test...
			var err error
			var errs = make([]error, 0)
			frame := newLogicalFrame(Config{
				Lines:     table.rows,
				HasHeader: table.hasHeader,
				HasFooter: table.hasFooter,
				startRow:  table.startRow,
			})
			if table.hasHeader {
				err = frame.header.WriteString("theHeader")
				if err != nil {
					errs = append(errs, err)
				}
			}
			for idx, line := range frame.activeLines {
				err = line.WriteString(fmt.Sprintf("LineIdx:%d", idx))
				if err != nil {
					errs = append(errs, err)
				}
			}
			if table.hasFooter {
				err = frame.footer.WriteString("theFooter")
				if err != nil {
					errs = append(errs, err)
				}
			}

			// assert results...
			validateEvents(t, test, table, errs, frame, handler)

		})
	}

}

func Test_LogicalFrame_UpdateDraw(t *testing.T) {

	for test, table := range drawTestCases {
		suppressOutput(func() {
			// setup...
			terminalHeight = table.terminalHeight
			handler := NewTestEventHandler(t)
			screenHandlers = make([]ScreenEventHandler, 0)
			addScreenHandler(handler)

			// run test...
			var errs []error
			frame := newLogicalFrame(Config{
				Lines:     table.rows,
				HasHeader: table.hasHeader,
				HasFooter: table.hasFooter,
				startRow:  table.startRow,
			})
			if table.hasHeader {
				frame.header.buffer = []byte("theHeader")
			}
			for idx, line := range frame.activeLines {
				line.buffer = []byte(fmt.Sprintf("LineIdx:%d", idx))
			}
			if table.hasFooter {
				frame.footer.buffer = []byte("theFooter")
			}
			errs = frame.updateAndDraw()

			// assert results...
			validateEvents(t, test, table, errs, frame, handler)

		})
	}

}
