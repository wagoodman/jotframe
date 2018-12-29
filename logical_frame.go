package jotframe

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"sync"
)

func newLogicalFrameAt(rows int, hasHeader, hasFooter bool, firstRow int) *logicalFrame {
	frame := &logicalFrame{}
	frame.frameStartIdx = firstRow
	frame.closeSignal = &sync.WaitGroup{}

	var relativeRow int
	if hasHeader {
		// todo: should headers have closeSignal waitGroups? or should they be nil?
		frame.header = NewLine(frame.frameStartIdx+relativeRow, frame.closeSignal)
		relativeRow++
	}
	for idx := 0; idx < rows; idx++ {
		frame.append()
	}
	if hasFooter {
		// todo: should footers have closeSignal waitGroups? or should they be nil?
		frame.footer = NewLine(frame.frameStartIdx+len(frame.activeLines)+relativeRow, frame.closeSignal)
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

func (frame *logicalFrame) visibleHeight() int {
	height := frame.height()

	if height > terminalHeight {
		return terminalHeight
	}
	return height
}

func (frame *logicalFrame) isAtOrPastScreenTop() bool {
	// take into account the rows that will be added to the screen realestate
	futureFrameStartIdx := frame.frameStartIdx - frame.rowAdvancements

	if futureFrameStartIdx < 1 {
		return true
	}
	return false
}

func (frame *logicalFrame) isAtOrPastScreenBottom() bool {
	height := frame.height()

	// take into account the rows that will be added to the screen realestate
	futureFrameStartIdx := frame.frameStartIdx - frame.rowAdvancements

	logrus.Debug(futureFrameStartIdx, "=", frame.frameStartIdx, "-", frame.rowAdvancements)
	// if the frame has moved past the bottom of the screen, move it up a bit
	if futureFrameStartIdx + height > terminalHeight {
		return true
	}
	return false
}

func (frame *logicalFrame) append() (*Line, error) {
	if frame.isClosed() {
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

	newLine := NewLine(rowIdx, frame.closeSignal)
	frame.activeLines = append(frame.activeLines, newLine)

	if frame.footer != nil {
		frame.footer.move(1)
	}

	return newLine, nil
}

func (frame *logicalFrame) prepend() (*Line, error) {
	if frame.isClosed() {
		return nil, fmt.Errorf("frame is closed")
	}

	rowIdx := frame.frameStartIdx
	if frame.header != nil {
		rowIdx += 1
	}

	newLine := NewLine(rowIdx, frame.closeSignal)

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
	if frame.isClosed() {
		return nil, fmt.Errorf("frame is closed")
	}

	if index < 0 || index > len(frame.activeLines) {
		return nil, fmt.Errorf("invalid index given")
	}

	rowIdx := frame.frameStartIdx + index
	if frame.header != nil {
		rowIdx += 1
	}

	newLine := NewLine(rowIdx, frame.closeSignal)

	frame.activeLines = append(frame.activeLines, nil)
	copy(frame.activeLines[index+1:], frame.activeLines[index:])
	frame.activeLines[index] = newLine

	// bump the indexes for other rows
	for idx := index + 1; idx < len(frame.activeLines); idx++ {
		frame.activeLines[idx].move(1)
	}

	if frame.footer != nil {
		frame.footer.move(1)
	}

	return newLine, nil
}

func (frame *logicalFrame) indexOf(line *Line) int {
	// find the index of the line object
	matchedIdx := -1
	for idx, item := range frame.activeLines {
		if item == line {
			return idx
		}
	}

	return matchedIdx
}

func (frame *logicalFrame) remove(line *Line) error {
	if frame.isClosed() {
		return fmt.Errorf("frame is closed")
	}

	// find the index of the line object
	matchedIdx := frame.indexOf(line)

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
	for idx := matchedIdx; idx < len(frame.activeLines); idx++ {
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

	// make screen realestate if the cursor is already near the bottom row (this preservers the users existing terminal output)
	// if frame.isAtOrPastScreenBottom() {
	// 	height := frame.visibleHeight()
	// 	offset := frame.frameStartIdx - ((terminalHeight - height) + 1)
	// 	offset += 1 // we want to move one line past the frame
	// 	frame.move(-offset)
	// 	frame.rowAdvancements += offset
	// }

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

// todo: I think this should be decided by the user via a close() aciton, not by the indication of closed lines
// since you can always add another line... you don't know when an empty frame should remain open or not
func (frame *logicalFrame) isClosed() bool {
	// if frame.header != nil {
	// 	if !frame.header.closed {
	// 		return false
	// 	}
	// }
	//
	// for _, line := range frame.activeLines {
	// 	if !line.closed {
	// 		return false
	// 	}
	// }
	//
	// if frame.footer != nil {
	// 	if !frame.footer.closed {
	// 		return false
	// 	}
	// }
	// return true
	return frame.closed
}

func (frame *logicalFrame) move(rows int) {
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
}

// ensure that the frame is within the bounds of the terminal
func (frame *logicalFrame) update() error {
	if frame.updateFn != nil {
		err := frame.updateFn(frame)
		if err != nil {
			return err
		}
	}

	return nil
}

func (frame *logicalFrame) updateAndDraw() (errs []error)  {
	errs = make([]error, 0)

	err := frame.update()
	if err != nil {
		errs = append(errs, err)
	}

	return append(errs, frame.draw()...)
}

func (frame *logicalFrame) draw() (errs []error) {
	errs = make([]error, 0)

	// clear any marked lines (preserving the buffer) while these indexes still exist
	for _, row := range frame.clearRows {
		err := clearRow(row)
		if err != nil {
			errs = append(errs, err)
		}
	}
	frame.clearRows = make([]int, 0)

	// advance the screen while adding any trail lines
	for idx := 0; idx < frame.rowAdvancements; idx++ {
		advanceScreen(1)
		if idx < len(frame.trailRows) {
			writeAtRow(frame.trailRows[0], frame.frameStartIdx-len(frame.trailRows)+idx)
			if len(frame.trailRows) >= 1 {
				frame.trailRows = frame.trailRows[1:]
			} else {
				frame.trailRows = make([]string, 0)
			}
		}
	}
	frame.rowAdvancements = 0

	// append any remaining trail rows
	for idx, message := range frame.trailRows {
		writeAtRow(message, frame.frameStartIdx-len(frame.trailRows)+idx)
	}
	frame.trailRows = make([]string, 0)

	// paint all stale lines to the screen
	if frame.header != nil {
		if frame.header.stale || frame.stale {
			_, err := frame.header.write(frame.header.buffer)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}

	for _, line := range frame.activeLines {
		if line.stale || frame.stale {
			_, err := line.write(line.buffer)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}

	if frame.footer != nil {
		if frame.footer.stale || frame.stale {
			_, err := frame.footer.write(frame.footer.buffer)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}

	if frame.isClosed() {
		setCursorRow(frame.frameStartIdx + frame.height())
	}

	return errs
}

func (frame *logicalFrame) wait() {
	frame.closeSignal.Wait()
	// setCursorRow(frame.frameStartIdx + frame.height())
}
