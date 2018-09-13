package jotframe

import (
	"fmt"
)

func newLogicalFrameAt(rows int, hasHeader, hasFooter bool, destinationRow int) *logicalFrame {
	frame := &logicalFrame{}
	frame.frameStartIdx = destinationRow

	var relativeRow int
	if hasHeader {
		frame.header = NewLine(frame.frameStartIdx + relativeRow)
		relativeRow++
	}
	for idx := 0; idx < rows; idx++ {
		frame.append()
	}
	if hasFooter {
		frame.footer = NewLine(frame.frameStartIdx + len(frame.activeLines) + relativeRow)
		relativeRow++
	}

	registerFrame(frame)

	return frame
}

func (frame *logicalFrame) appendTrail(str string) {
	frame.trailRows = append(frame.trailRows, str)
}

func (frame *logicalFrame) height() int {
	height := len(frame.activeLines)
	if frame.header != nil {
		height++
	}
	if frame.footer != nil {
		height++
	}
	return height
}

func (frame *logicalFrame) append() (*Line, error) {
	if frame.closed {
		return nil, fmt.Errorf("frame is closed")
	}

	var rowIdx int
	if len(frame.activeLines) > 0 {
		rowIdx = frame.activeLines[len(frame.activeLines)-1].row + 1
	} else {
		rowIdx = frame.frameStartIdx
		if frame.header != nil {
			rowIdx += 1
		}

	}

	newLine := NewLine(rowIdx)
	frame.activeLines = append(frame.activeLines, newLine)

	if frame.footer != nil {
		frame.footer.move(1)
	}

	return newLine, nil
}

func (frame *logicalFrame) prepend() (*Line, error) {
	if frame.closed {
		return nil, fmt.Errorf("frame is closed")
	}

	rowIdx := frame.frameStartIdx
	if frame.header != nil {
		rowIdx += 1
	}

	newLine := NewLine(rowIdx)
	for _, line := range frame.activeLines {
		line.move(1)
	}
	frame.activeLines = append([]*Line{newLine}, frame.activeLines...)

	if frame.footer != nil {
		frame.footer.move(1)
	}

	return newLine, nil
}

func (frame *logicalFrame) insert(index int) (*Line, error) {
	if frame.closed {
		return nil, fmt.Errorf("frame is closed")
	}

	if index < 0 || index > len(frame.activeLines) {
		return nil, fmt.Errorf("invalid index given")
	}

	rowIdx := frame.frameStartIdx +index
	if frame.header != nil {
		rowIdx += 1
	}

	newLine := NewLine(rowIdx)

	frame.activeLines = append(frame.activeLines, nil)
	copy(frame.activeLines[index+1:], frame.activeLines[index:])
	frame.activeLines[index] = newLine

	// bump the indexes for other rows
	for idx := index+1; idx < len(frame.activeLines); idx++ {
		frame.activeLines[idx].move(1)
	}

	if frame.footer != nil {
		frame.footer.move(1)
	}

	return newLine, nil
}

func (frame *logicalFrame) remove(line *Line) error {
	if frame.closed {
		return fmt.Errorf("frame is closed")
	}

	// find the index of the line object
	matchedIdx := -1
	for idx, item := range frame.activeLines {
		if item == line {
			matchedIdx = idx
			break
		}
	}

	if matchedIdx < 0 {
		return fmt.Errorf("could not find line in frame")
	}

	// activeLines that are removed must be closed since any further writes will result in line clashes
	frame.activeLines[matchedIdx].close()

	// erase the contents of the last line of the logicalFrame, but persist the line buffer
	if frame.footer != nil {
		frame.clearRows = append(frame.clearRows, frame.footer.row)
	} else {
		frame.clearRows = append(frame.clearRows, frame.activeLines[len(frame.activeLines)-1].row)
	}

	// remove the line entry from the list
	frame.activeLines = append(frame.activeLines[:matchedIdx], frame.activeLines[matchedIdx+1:]...)

	// move each line index ahead of the deleted element
	for idx := matchedIdx+1; idx < len(frame.activeLines); idx++ {
		frame.activeLines[idx].move(-1)
	}

	if frame.footer != nil {
		frame.footer.move(-1)
	}

	return nil
}

func (frame *logicalFrame) clear() error {

	if frame.header != nil {
		frame.clearRows = append(frame.clearRows, frame.header.row)
	}

	for _, line := range frame.activeLines {
		frame.clearRows = append(frame.clearRows, line.row)
	}

	if frame.footer != nil {
		frame.clearRows = append(frame.clearRows, frame.footer.row)
	}
	return nil
}

func (frame *logicalFrame) close() error {

	if frame.header != nil {
		err := frame.header.close()
		if err != nil {
			return err
		}
	}

	for _, line := range frame.activeLines {
		err := line.close()
		if err != nil {
			return err
		}
	}

	if frame.footer != nil {
		err := frame.footer.close()
		if err != nil {
			return err
		}
	}
	frame.closed = true
	return nil
}


func (frame *logicalFrame) move(rows int) error {
	// todo: check real screen dimensions for moving past bottom of screen
	if frame.frameStartIdx+ rows < 0 {
		return fmt.Errorf("unable to move past screen dimensions")
	}
	frame.frameStartIdx += rows

	// todo: instead of clearing all frame lines, only clear the ones affected
	frame.clear()

	// bump rows and redraw entire frame
	if frame.header != nil {
		frame.header.move(rows)
	}
	for _, line := range frame.activeLines {
		line.move(rows)
	}
	if frame.footer != nil {
		frame.footer.move(rows)
	}

	return nil
}


// ensure that the frame is within the bounds of the terminal
func (frame *logicalFrame) update() error {
	height := frame.height()

	// take into account the rows that will be added to the screen realestate
	futureFrameStartIdx := frame.frameStartIdx - frame.rowPreAdvancements

	// if the frame has moved past the bottom of the screen, move it up a bit
	if futureFrameStartIdx + height > terminalHeight {
		offset := ((terminalHeight - height)+1) - futureFrameStartIdx
		return frame.move(offset)
	}

	// if the frame has moved past the bottom of the screen, move it down a bit
	if futureFrameStartIdx < 1 {
		offset := 1 - futureFrameStartIdx
		return frame.move(offset)
	}

	return nil
}


func (frame *logicalFrame) updateAndDraw() {
	if frame.updateFn != nil {
		frame.updateFn()
	}

	// don't allow any update function to draw outside of the screen dimensions
	frame.update()

	frame.draw()
}

func (frame *logicalFrame) draw() error {

	// clear any marked lines (preserving the buffer) while these indexes still exist
	for _, row := range frame.clearRows {
		err := clearRow(row)
		if err != nil {
			return err
		}
	}
	frame.clearRows = make([]int, 0)

	// advance the screen while adding any trail lines
	for idx := 0; idx < frame.rowPreAdvancements; idx++ {
		advanceScreen(1)
		if idx < len(frame.trailRows) {
			writeAtRow(frame.trailRows[0], frame.frameStartIdx - len(frame.trailRows) + idx)
			if len(frame.trailRows) >= 1 {
				frame.trailRows = frame.trailRows[1:]
			} else {
				frame.trailRows = make([]string, 0)
			}
		}
	}
	frame.rowPreAdvancements = 0

	// append any remaining trail rows
	for idx, message := range frame.trailRows {
		writeAtRow(message, frame.frameStartIdx - len(frame.trailRows) + idx)
	}
	frame.trailRows = make([]string, 0)

	// paint all stale lines to the screen
	if frame.header != nil {
		if frame.header.stale || frame.stale {
			_, err := frame.header.write(frame.header.buffer)
			if err != nil {
				return err
			}
		}
	}

	for _, line := range frame.activeLines {
		if line.stale || frame.stale {
			_, err := line.write(line.buffer)
			if err != nil {
				return err
			}
		}
	}

	if frame.footer != nil {
		if frame.footer.stale || frame.stale {
			_, err := frame.footer.write(frame.footer.buffer)
			if err != nil {
				return err
			}
		}
	}

	return nil
}