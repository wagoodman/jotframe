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
		headers           int
		footers           int
		destinationRow    int
		expectedHeaderRow int
		expectedFooterRow int
		expectedLineRows  []int
	}{
		{5, 0, 0, 10, -1, -1, []int{10, 11, 12, 13, 14}},
		{5, 1, 0, 10, 10, -1, []int{11, 12, 13, 14, 15}},
		{5, 0, 1, 10, -1, 15, []int{10, 11, 12, 13, 14}},
		{5, 1, 1, 10, 10, 16, []int{11, 12, 13, 14, 15}},
	}

	for _, table := range tables {
		getScreen().reset()

		frame, _, _, _ := New(Config{
			test: true,
			Lines:          table.rows,
			HeaderRows:     table.headers,
			FooterRows:     table.footers,
			startRow:       table.destinationRow,
			PositionPolicy: PolicyFloatOverflow,
		})

		// verify the Headers
		if table.headers != len(frame.Headers) {
			t.Errorf("NewFrame: expected headers to be %d but got %d", table.headers, len(frame.Headers))
		}

		// verify the Footers
		if table.footers != len(frame.Footers) {
			t.Errorf("NewFrame: expected footers to be %d but got %d", table.footers, len(frame.Footers))
		}

		// verify the number of lines in the frame list
		actualLineRowLen := len(frame.Lines)
		expectedLineRowLen := len(table.expectedLineRows)
		if expectedLineRowLen != actualLineRowLen {
			t.Errorf("NewFrame: expected %d lines, found %d", expectedLineRowLen, actualLineRowLen)
		}

		// ensure the screen row values are correct relative to the given starting row
		var expectedRow, actualRow int

		if table.headers > 0 {
			expectedRow = table.expectedHeaderRow
			actualRow = frame.Headers[0].row
			if expectedRow != actualRow {
				t.Errorf("NewFrame: expected Headers row to start at %d, but starts at %d", expectedRow, actualRow)
			}
		}

		if table.footers > 0 {
			expectedRow = table.expectedFooterRow
			actualRow = frame.Footers[0].row
			if expectedRow != actualRow {
				t.Errorf("NewFrame: expected Footers row to start at %d, but starts at %d", expectedRow, actualRow)
			}
		}

		for idx, expectedRow := range table.expectedLineRows {
			actualRow = frame.Lines[idx].row
			if expectedRow != actualRow {
				t.Errorf("NewFrame: expected line row to start at %d, but starts at %d", expectedRow, actualRow)
			}
		}


	}

}

func Test_Frame_Height(t *testing.T) {

	tables := []struct {
		rows           int
		headers        int
		footers        int
		expectedHeight int
	}{
		{5, 0, 0, 5},
		{5, 1, 0, 6},
		{5, 0, 1, 6},
		{5, 1, 1, 7},
	}

	for _, table := range tables {
		getScreen().reset()

		frame, _, _, _ := New(Config{
			test: true,
			Lines:          table.rows,
			HeaderRows:     table.headers,
			FooterRows:     table.footers,
			startRow:       10,
			PositionPolicy: PolicyFloatOverflow,
		})
		actualHeight := frame.Height()
		if table.expectedHeight != actualHeight {
			t.Errorf("Frame.height(): expected a height of %d, but found %d", table.expectedHeight, actualHeight)
		}

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
		getScreen().reset()

		frame, _, _, _ := New(Config{
			test: true,
			Lines:          frameRows,
			HeaderRows:     0,
			FooterRows:     0,
			startRow:       table.destinationRow,
			PositionPolicy: PolicyFloatOverflow,
		})
		actualResult := frame.IsPastScreenBottom()
		if table.isAtBottom != actualResult {
			t.Errorf("Frame.IsPastScreenBottom(): expected result of %v, but found %v (startIdx:%d)", table.isAtBottom, actualResult, frame.StartIdx)
		}

	}

}

func Test_Frame_Append(t *testing.T) {

	tables := []struct {
		appendRows        int
		headers           int
		footers           int
		destinationRow    int
		expectedHeaderRow int
		expectedFooterRow int
		expectedLineRows  []int
	}{
		{1, 0, 0, 10, -1, -1, []int{10}},
		{1, 1, 0, 10, 10, -1, []int{11}},
		{1, 0, 1, 10, -1, 11, []int{10}},
		{1, 1, 1, 10, 10, 12, []int{11}},
		{5, 0, 0, 10, -1, -1, []int{10, 11, 12, 13, 14}},
		{5, 1, 0, 10, 10, -1, []int{11, 12, 13, 14, 15}},
		{5, 0, 1, 10, -1, 15, []int{10, 11, 12, 13, 14}},
		{5, 1, 1, 10, 10, 16, []int{11, 12, 13, 14, 15}},
	}

	for _, table := range tables {
		getScreen().reset()

		frame, _, _, _ := New(Config{
			test: true,
			Lines:          0,
			HeaderRows:     table.headers,
			FooterRows:     table.footers,
			startRow:       table.destinationRow,
			PositionPolicy: PolicyFloatOverflow,
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
		if len(frame.Lines) != table.appendRows {
			t.Errorf("Frame.Append(): expected %d number of lines, got %d", table.appendRows, len(frame.Lines))
		}

		// check the contents of each line
		for idx := 0; idx < table.appendRows; idx++ {
			line := frame.Lines[idx]
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
			actualRow = frame.Lines[idx].row
			if expectedRow != actualRow {
				t.Errorf("Frame.Append(): expected line row to start at %d, but starts at %d", expectedRow, actualRow)
			}
		}

	}
}

func Test_Frame_Prepend(t *testing.T) {

	tables := []struct {
		prependRows       int
		headers           int
		footers           int
		destinationRow    int
		expectedHeaderRow int
		expectedFooterRow int
		expectedLineRows  []int
	}{
		{1, 0, 0, 10, -1, -1, []int{10}},
		{1, 1, 0, 10, 10, -1, []int{11}},
		{1, 0, 1, 10, -1, 11, []int{10}},
		{1, 1, 1, 10, 10, 12, []int{11}},
		{5, 0, 0, 10, -1, -1, []int{10, 11, 12, 13, 14}},
		{5, 1, 0, 10, 10, -1, []int{11, 12, 13, 14, 15}},
		{5, 0, 1, 10, -1, 15, []int{10, 11, 12, 13, 14}},
		{5, 1, 1, 10, 10, 16, []int{11, 12, 13, 14, 15}},
	}

	for _, table := range tables {
		getScreen().reset()

		frame, _, _, _ := New(Config{
			test: true,
			Lines:          0,
			HeaderRows:     table.headers,
			FooterRows:     table.footers,
			startRow:       table.destinationRow,
			PositionPolicy: PolicyFloatOverflow,
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
		if len(frame.Lines) != table.prependRows {
			t.Errorf("Frame.Prepend(): expected %d number of lines, got %d", table.prependRows, len(frame.Lines))
		}

		// check the contents of each line
		for idx := 0; idx < table.prependRows; idx++ {
			// note: indexes should be in reverse
			line := frame.Lines[table.prependRows-idx-1]
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
			actualRow = frame.Lines[idx].row
			if expectedRow != actualRow {
				t.Errorf("Frame.Prepend(): expected line row to start at %d, but starts at %d", expectedRow, actualRow)
			}
		}

	}
}

func Test_Frame_Insert(t *testing.T) {

	tables := []struct {
		header            int
		footer            int
		destinationRow    int
		expectedHeaderRow int
		expectedFooterRow int
		expectedLineRows  []int
	}{
		{0, 0, 10, -1, -1, []int{10, 11, 12, 13, 14}},
		{1, 0, 10, 10, -1, []int{11, 12, 13, 14, 15}},
		{0, 1, 10, -1, 15, []int{10, 11, 12, 13, 14}},
		{1, 1, 10, 10, 16, []int{11, 12, 13, 14, 15}},
	}

	for _, table := range tables {
		getScreen().reset()

		frame, _, _, _ := New(Config{
			test: true,
			Lines:          4,
			HeaderRows:     table.header,
			FooterRows:     table.footer,
			startRow:       table.destinationRow,
			PositionPolicy: PolicyFloatOverflow,
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
		if len(frame.Lines) != finalRows {
			t.Errorf("Frame.Insert(): expected %d number of lines, got %d", finalRows, len(frame.Lines))
		}

		// check the contents of the inserted line

		fetchedLine := frame.Lines[insertIdx]
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
			actualRow = frame.Lines[idx].row
			if expectedRow != actualRow {
				t.Errorf("Frame.Prepend(): expected line row to start at %d, but starts at %d", expectedRow, actualRow)
			}
		}

	}
}

func Test_Frame_Remove(t *testing.T) {

	tables := []struct {
		startRows        int
		headers          int
		footers          int
		destinationRow   int
		expectedLineRows []int
	}{
		{6, 0, 0, 10, []int{10, 11, 12, 13, 14}},
		{6, 1, 0, 10, []int{11, 12, 13, 14, 15}},
		{6, 0, 1, 10, []int{10, 11, 12, 13, 14}},
		{6, 1, 1, 10, []int{11, 12, 13, 14, 15}},
	}

	for _, table := range tables {
		getScreen().reset()

		frame, _, _, _ := New(Config{
			test: true,
			Lines:          table.startRows,
			HeaderRows:     table.headers,
			FooterRows:     table.footers,
			startRow:       table.destinationRow,
			PositionPolicy: PolicyFloatOverflow,
		})
		rmIdx := 2

		// add content to rows
		for idx, line := range frame.Lines {
			line.buffer = []byte(strconv.Itoa(idx))
		}

		// Remove a single index
		frame.Remove(frame.Lines[rmIdx])

		// check if the number of rows matches
		if len(frame.Lines) != table.startRows-1 {
			t.Errorf("Frame.Remove(): expected %d number of lines, got %d", table.startRows, len(frame.Lines))
		}

		// check the contents of each line
		for idx := 0; idx < table.startRows-1; idx++ {
			line := frame.Lines[idx]
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
			actualRow = frame.Lines[idx].row
			if expectedRow != actualRow {
				t.Errorf("Frame.Remove(): expected line row to start at %d, but starts at %d (idx:%d)", expectedRow, actualRow, idx)
			}
		}

	}
}

func Test_Frame_Clear(t *testing.T) {

	tables := map[string]struct {
		startRows         int
		headers           int
		footers           int
		destinationRow    int
		expectedClearRows []int
	}{
		"goCase":          {5, 0, 0, 10, []int{10, 11, 12, 13, 14}},
		"XHeader":         {5, 1, 0, 10, []int{10, 11, 12, 13, 14, 15}},
		"XFooter":         {5, 0, 1, 10, []int{10, 11, 12, 13, 14, 15}},
		"XHeader+XFooter": {5, 1, 1, 10, []int{10, 11, 12, 13, 14, 15, 16}},
	}

	for test, table := range tables {
		getScreen().reset()

		frame, _, _, _ := New(Config{
			test: true,
			Lines:          table.startRows,
			HeaderRows:     table.headers,
			FooterRows:     table.footers,
			startRow:       table.destinationRow,
			PositionPolicy: PolicyFloatOverflow,
		})
		frame.clear()

		height := table.startRows + table.headers + table.footers

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

	}
}

// todo: rewrite this
func Test_Frame_Close(t *testing.T) {

	tables := map[string]struct {
		startRows      int
		headers        int
		footers        int
		destinationRow int
	}{
		"goCase":          {5, 0, 0, 10},
		"XHeader":         {5, 1, 0, 10},
		"XFooter":         {5, 0, 1, 10},
		"XHeader+XFooter": {5, 1, 1, 10},
	}

	for test, table := range tables {
		getScreen().reset()

		frame, _, _, _ := New(Config{
			test: true,
			Lines:          table.startRows,
			HeaderRows:     table.headers,
			FooterRows:     table.footers,
			startRow:       table.destinationRow,
			PositionPolicy: PolicyFloatOverflow,
		})
		frame.Close()

		if !frame.closed {
			t.Errorf("Frame.Close(): [case=%s] expected frame to be closed but is not", test)
		}

		if table.headers > 0 && !frame.Headers[0].closed {
			t.Errorf("Frame.Close(): [case=%s] expected Headers to be closed but is not", test)
		}
		if table.footers > 0 && !frame.Footers[0].closed {
			t.Errorf("Frame.Close(): [case=%s] expected Footers to be closed but is not", test)
		}

		for idx, line := range frame.Lines {
			if !line.closed {
				t.Errorf("Frame.Close(): [case=%s] expected line %d to be closed but is not", test, idx)
			}
		}

	}
}

func Test_Frame_Move(t *testing.T) {

	tables := map[string]struct {
		startRows int
		header    int
		footer    int
		startRow  int
		moveRows  int
	}{
		"goCase":          {5, 0, 0, 10, 5},
		"XHeader":         {5, 1, 0, 10, -5},
		"XFooter":         {5, 0, 1, 10, 5},
		"XHeader+XFooter": {5, 1, 1, 10, -5},
	}

	for test, table := range tables {
		getScreen().reset()

		frame, _, _, _ := New(Config{
			test: true,
			Lines:          table.startRows,
			HeaderRows:     table.header,
			FooterRows:     table.footer,
			startRow:       table.startRow,
			PositionPolicy: PolicyFloatOverflow,
		})
		frame.Move(table.moveRows)

		expectedFrameRow := table.startRow + table.moveRows
		actualFrameRow := frame.StartIdx
		if expectedFrameRow != actualFrameRow {
			t.Errorf("Frame.Move(): [case=%s] expected frame to be at %d, but is at %d", test, expectedFrameRow, actualFrameRow)
		}

		if table.header > 0 {

			expectedFrameRow = table.startRow + table.moveRows
			actualFrameRow = frame.Headers[0].row
			if expectedFrameRow != actualFrameRow {
				t.Errorf("Frame.Move(): [case=%s] expected Headers to be at %d, but is at %d", test, expectedFrameRow, actualFrameRow)
			}
		}

		if table.footer > 0 {
			expectedFrameRow = table.startRow + table.startRows + table.header + table.moveRows
			actualFrameRow = frame.Footers[0].row
			if expectedFrameRow != actualFrameRow {
				t.Errorf("Frame.Move(): [case=%s] expected Footers to be at %d, but is at %d", test, expectedFrameRow, actualFrameRow)
			}
		}

		for idx, line := range frame.Lines {
			expectedFrameRow = table.startRow + idx + table.header + table.moveRows
			actualFrameRow = line.row
			if expectedFrameRow != actualFrameRow {
				t.Errorf("Frame.Move(): [case=%s] expected line to be at %d, but is at %d", test, expectedFrameRow, actualFrameRow)
			}
		}

	}
}

// func Test_Frame_Update(t *testing.T) {
// 	terminalHeight = 100
// 	frameRows := 10
//
// 	fn := func(frame *Frame) error {
// 		// if the frame has moved past the bottom of the screen, move it up a bit
// 		if frame.IsPastScreenBottom() {
// 			frameHeight := frame.VisibleHeight()
// 			// offset is how many rows the frame needs to be adjusted to fit on the screen.
// 			// This is the same as how many rows past the edge of the screen this frame currently is.
// 			offset := (frame.StartIdx + frameHeight) - terminalHeight
// 			// offset += 1 // we want to move one line past the frame
// 			frame.Move(-offset)
// 			frame.rowAdvancements += offset
// 		}
//
// 		// if the frame has moved above the top of the screen, move it down a bit
// 		if frame.IsPastScreenTop() {
// 			offset := -1*frame.StartIdx + 1
// 			frame.Move(offset)
// 		}
// 		return nil
// 	}
//
// 	tables := []struct {
// 		startRow int
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
// 		frame := New(XConfig{
// 			XLines:     frameRows,
// 			HeaderRows: false,
// 			FooterRows: false,
// 			startRow:  table.startRow,
// 		})
// 		frame.updateFn = fn
// 		frame.Update()
// 		actualResult := frame.StartIdx
// 		if table.adjustedRow != actualResult {
// 			t.Errorf("Frame.Update(): expected Update row of %d, but is at row %d", table.adjustedRow, actualResult)
// 		}
// 	}
//
// }
