package jotframe

import (
	"testing"
	"strconv"
)

func Test_NewLogicalFrame(t *testing.T) {

	tables := []struct {
		rows int
		hasHeader bool
		hasFooter bool
		destinationRow int
		expectedHeaderRow int
		expectedFooterRow int
		expectedLineRows []int
	}{
		{5, false, false, 10, -1, -1, []int{10, 11, 12, 13, 14} },
		{5, true, false, 10, 10, -1, []int{11, 12, 13, 14, 15} },
		{5, false, true, 10, -1, 15, []int{10, 11, 12, 13, 14} },
		{5, true, true, 10, 10, 16, []int{11, 12, 13, 14, 15} },
	}

	for _, table := range tables {
		frame := newLogicalFrameAt(table.rows, table.hasHeader, table.hasFooter, table.destinationRow)

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
		rows int
		hasHeader bool
		hasFooter bool
		expectedHeight int
	}{
		{5, false, false, 5},
		{5, true, false, 6 },
		{5, false, true, 6},
		{5, true, true, 7},
	}

	for _, table := range tables {
		frame := newLogicalFrameAt(table.rows, table.hasHeader, table.hasFooter, 10)
		actualHeight := frame.height()
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
		isAtBottom bool
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
		frame := newLogicalFrameAt(frameRows, false, false, table.destinationRow)
		actualResult := frame.isAtOrPastScreenBottom()
		if table.isAtBottom != actualResult {
			t.Errorf("LogicalFrame.isAtOrPastScreenBottom(): expected result of %v, but found %v (startIdx:%d)", table.isAtBottom, actualResult, frame.frameStartIdx)
		}
	}

}


func Test_LogicalFrame_Append(t *testing.T) {

	tables := []struct {
		appendRows int
		hasHeader bool
		hasFooter bool
		destinationRow int
		expectedHeaderRow int
		expectedFooterRow int
		expectedLineRows []int
	}{
		{1, false, false, 10, -1, -1, []int{10} },
		{1, true, false, 10, 10, -1, []int{11} },
		{1, false, true, 10, -1, 11, []int{10} },
		{1, true, true, 10, 10, 12, []int{11} },
		{5, false, false, 10, -1, -1, []int{10, 11, 12, 13, 14} },
		{5, true, false, 10, 10, -1, []int{11, 12, 13, 14, 15} },
		{5, false, true, 10, -1, 15, []int{10, 11, 12, 13, 14} },
		{5, true, true, 10, 10, 16, []int{11, 12, 13, 14, 15} },
	}

	for _, table := range tables {
		frame := newLogicalFrameAt(0, table.hasHeader, table.hasFooter, table.destinationRow)

		// append rows...
		for idx := 0; idx < table.appendRows; idx++ {
			line, err := frame.append()
			if err != nil {
				t.Errorf("LogicalFrame.append(): expected no error on append(), got %v", err)
			}
			// write out each index while appending, later check the order
			line.buffer = []byte(strconv.Itoa(idx))
		}

		// check if the number of rows matches
		if len(frame.activeLines) != table.appendRows {
			t.Errorf("LogicalFrame.append(): expected %d number of lines, got %d", table.appendRows, len(frame.activeLines))
		}

		// check the contents of each line
		for idx := 0; idx < table.appendRows; idx++ {
			line := frame.activeLines[idx]
			actualNum, err := strconv.Atoi(string(line.buffer))
			if err != nil {
				t.Errorf("LogicalFrame.append(): expected no error on line read (%v), got %v", actualNum, err)
			}
			if actualNum != idx {
				t.Errorf("LogicalFrame.append(): expected %d, got %d", idx, actualNum)
			}
		}

		// ensure the screen row values are correct relative to the given starting row
		var actualRow int

		for idx, expectedRow := range table.expectedLineRows {
			actualRow = frame.activeLines[idx].row
			if expectedRow != actualRow {
				t.Errorf("LogicalFrame.append(): expected line row to start at %d, but starts at %d", expectedRow, actualRow)
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
		{1, false, false, 10, -1, -1, []int{10} },
		{1, true, false, 10, 10, -1, []int{11} },
		{1, false, true, 10, -1, 11, []int{10} },
		{1, true, true, 10, 10, 12, []int{11} },
		{5, false, false, 10, -1, -1, []int{10, 11, 12, 13, 14} },
		{5, true, false, 10, 10, -1, []int{11, 12, 13, 14, 15} },
		{5, false, true, 10, -1, 15, []int{10, 11, 12, 13, 14} },
		{5, true, true, 10, 10, 16, []int{11, 12, 13, 14, 15} },
	}

	for _, table := range tables {
		frame := newLogicalFrameAt(0, table.hasHeader, table.hasFooter, table.destinationRow)

		// append rows...
		for idx := 0; idx < table.prependRows; idx++ {
			line, err := frame.prepend()
			if err != nil {
				t.Errorf("LogicalFrame.prepend(): expected no error on prepend(), got %v", err)
			}
			// write out each index while appending, later check the order
			line.buffer = []byte(strconv.Itoa(idx))
		}

		// check if the number of rows matches
		if len(frame.activeLines) != table.prependRows {
			t.Errorf("LogicalFrame.prepend(): expected %d number of lines, got %d", table.prependRows, len(frame.activeLines))
		}

		// check the contents of each line
		for idx := 0; idx < table.prependRows; idx++ {
			// note: indexes should be in reverse
			line := frame.activeLines[table.prependRows-idx-1]
			actualNum, err := strconv.Atoi(string(line.buffer))
			if err != nil {
				t.Errorf("LogicalFrame.prepend(): expected no error on line read (%v), got %v", actualNum, err)
			}
			if actualNum != idx {
				t.Errorf("LogicalFrame.prepend(): expected %d, got %d", idx, actualNum)
			}
		}

		// ensure the screen row values are correct relative to the given starting row
		var actualRow int

		for idx, expectedRow := range table.expectedLineRows {
			actualRow = frame.activeLines[idx].row
			if expectedRow != actualRow {
				t.Errorf("LogicalFrame.prepend(): expected line row to start at %d, but starts at %d", expectedRow, actualRow)
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
		{false, false, 10, -1, -1, []int{10, 11, 12, 13, 14} },
		{true, false, 10, 10, -1, []int{11, 12, 13, 14, 15} },
		{false, true, 10, -1, 15, []int{10, 11, 12, 13, 14} },
		{true, true, 10, 10, 16, []int{11, 12, 13, 14, 15} },
	}

	for _, table := range tables {
		frame := newLogicalFrameAt(4, table.hasHeader, table.hasFooter, table.destinationRow)
		finalRows := 5
		insertIdx := 2
		// append rows...

		line, err := frame.insert(2)
		if err != nil {
			t.Errorf("LogicalFrame.insert(): expected no error on insert(), got %v", err)
		}
		// write out each index while appending, later check the order
		line.buffer = []byte(strconv.Itoa(insertIdx))


		// check if the number of rows matches
		if len(frame.activeLines) != finalRows {
			t.Errorf("LogicalFrame.insert(): expected %d number of lines, got %d", finalRows, len(frame.activeLines))
		}

		// check the contents of the inserted line

		fetchedLine := frame.activeLines[insertIdx]
		actualNum, err := strconv.Atoi(string(fetchedLine.buffer))
		if err != nil {
			t.Errorf("LogicalFrame.insert(): expected no error on line read (%v), got %v", actualNum, err)
		}
		if actualNum != insertIdx {
			t.Errorf("LogicalFrame.insert(): expected %d, got %d", insertIdx, actualNum)
		}


		// ensure the screen row values are correct relative to the given starting row
		var actualRow int

		for idx, expectedRow := range table.expectedLineRows {
			actualRow = frame.activeLines[idx].row
			if expectedRow != actualRow {
				t.Errorf("LogicalFrame.prepend(): expected line row to start at %d, but starts at %d", expectedRow, actualRow)
			}
		}
	}
}



// func Test_LogicalFrame_Remove(t *testing.T) {
//
// 	tables := []struct {
// 		startRows         int
// 		hasHeader         bool
// 		hasFooter         bool
// 		destinationRow    int
// 		expectedHeaderRow int
// 		expectedFooterRow int
// 		expectedLineRows  []int
// 	}{
// 		{6, false, false, 10, -1, -1, []int{10, 11, 12, 13, 14} },
// 		{6, true, false, 10, 10, -1, []int{11, 12, 13, 14, 15} },
// 		{6, false, true, 10, -1, 15, []int{10, 11, 12, 13, 14} },
// 		{6, true, true, 10, 10, 16, []int{11, 12, 13, 14, 15} },
// 	}
//
// 	for _, table := range tables {
// 		frame := newLogicalFrameAt(table.startRows, table.hasHeader, table.hasFooter, table.destinationRow)
// 		rmIdx := 2
//
// 		// add content to rows
// 		for idx, line := range frame.activeLines {
// 			line.buffer = []byte(strconv.Itoa(idx))
// 		}
//
// 		// check if the number of rows matches
// 		if len(frame.activeLines) != table.startRows {
// 			t.Errorf("LogicalFrame.append(): expected %d number of lines, got %d", table.startRows, len(frame.activeLines))
// 		}
//
// 		// check the contents of each line
// 		for idx := 0; idx < table.startRows; idx++ {
// 			line := frame.activeLines[idx]
// 			actualNum, err := strconv.Atoi(string(line.buffer))
// 			if err != nil {
// 				t.Errorf("LogicalFrame.append(): expected no error on line read (%v), got %v", actualNum, err)
// 			}
//
// 			if actualNum != idx {
// 				t.Errorf("LogicalFrame.append(): expected %d, got %d", idx, actualNum)
// 			}
// 		}
//
// 		// ensure the screen row values are correct relative to the given starting row
// 		var actualRow int
//
// 		for idx, expectedRow := range table.expectedLineRows {
// 			actualRow = frame.activeLines[idx].row
// 			if expectedRow != actualRow {
// 				t.Errorf("LogicalFrame.append(): expected line row to start at %d, but starts at %d", expectedRow, actualRow)
// 			}
// 		}
// 	}
// }