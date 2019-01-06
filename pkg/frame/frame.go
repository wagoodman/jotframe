package frame

import (
	"fmt"
	"sync"
)

func New(config Config) *Frame {
	frame := &Frame{
		topRow:      config.startRow,
		config:      config,
		closeSignal: &sync.WaitGroup{},
		lock:        getScreenLock(),
		autoDraw:    !config.ManualDraw,
	}

	switch config.PositionPolicy {
	case FloatFree:
		frame.policy = newFloatFreePolicy(frame)
	case FloatForward:
		frame.policy = newFloatForwardPolicy(frame)
	case FloatTop:
		frame.policy = newFloatTopPolicy(frame)
	case FloatBottom:
		frame.policy = newFloatBottomPolicy(frame)
	default:
		panic(fmt.Errorf("unknown policy: %v", config.PositionPolicy))
	}

	frame.policy.onInit()

	var relativeRow int
	if config.HasHeader {
		// todo: should headers have closeSignal waitGroups? or should they be nil?
		frame.header = NewLine(frame.topRow+relativeRow, frame.closeSignal)
		relativeRow++
	}
	for idx := 0; idx < config.Lines; idx++ {
		line := NewLine(frame.topRow+relativeRow, frame.closeSignal)
		frame.activeLines = append(frame.activeLines, line)
		relativeRow++
	}
	if config.HasFooter {
		// todo: should footers have closeSignal waitGroups? or should they be nil?
		frame.footer = NewLine(frame.topRow+relativeRow, frame.closeSignal)
		relativeRow++
	}

	// register frame before drawing to screen
	registerFrame(frame)

	// todo: it's bad that the constructor is writing out to the screen... is it avoidable?
	// adjust the screen such that a known good starting condition is in place
	frame.lock.Lock()
	defer frame.lock.Unlock()

	frame.draw()

	return frame
}

func (frame *Frame) SetAutoDraw(enabled bool){
	frame.autoDraw = enabled
}

func (frame *Frame) Config() Config {
	return frame.config
}

func (frame *Frame) Lines() []*Line {
	return frame.activeLines
}

func (frame *Frame) Header() *Line {
	return frame.header
}

func (frame *Frame) Footer() *Line {
	return frame.footer
}

func (frame *Frame) StartIdx() int {
	return frame.topRow
}


func (frame *Frame) AppendTrail(str string) {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	frame.appendTrail(str)
}

func (frame *Frame) appendTrail(str string) {
	if !frame.policy.allowTrail() {
		return
	}
	frame.trailRows = append(frame.trailRows, str)
	frame.policy.onTrail()

	// todo: what about update/draw here?
}

func (frame *Frame) Height() int {
	height := len(frame.activeLines)
	if frame.header != nil {
		height++
	}
	if frame.footer != nil {
		height++
	}
	return height
}

func (frame *Frame) visibleHeight() int {
	height := frame.Height()
	forwardDrawArea := terminalHeight - (frame.topRow-1)

	if height > forwardDrawArea {
		return forwardDrawArea
	}
	return height

}

func (frame *Frame) IsPastScreenTop() bool {
	// take into account the rows that will be added to the screen realestate
	futureFrameStartIdx := frame.topRow - frame.rowAdvancements

	if futureFrameStartIdx < 1 {
		return true
	}
	return false
}

func (frame *Frame) IsPastScreenBottom() bool {
	height := frame.Height()

	// take into account the rows that will be added to the screen realestate
	futureFrameStartIdx := frame.topRow - frame.rowAdvancements

	// if the frame has moved past the bottom of the screen, move it up a bit
	if futureFrameStartIdx+height > terminalHeight {
		return true
	}
	return false
}

func (frame *Frame) Append() (*Line, error) {
	if frame.IsClosed() {
		return nil, fmt.Errorf("frame is closed")
	}

	frame.lock.Lock()
	defer frame.lock.Unlock()

	var rowIdx int
	if len(frame.activeLines) > 0 {
		rowIdx = frame.activeLines[len(frame.activeLines)-1].row + 1
	} else {
		rowIdx = frame.topRow
		if frame.header != nil {
			rowIdx += 1
		}

	}

	newLine := NewLine(rowIdx, frame.closeSignal)
	frame.activeLines = append(frame.activeLines, newLine)

	if frame.footer != nil {
		frame.footer.move(1)
	}

	frame.policy.onResize(1)

	if frame.autoDraw {
		frame.draw()
	}

	return newLine, nil
}

func (frame *Frame) Prepend() (*Line, error) {
	if frame.IsClosed() {
		return nil, fmt.Errorf("frame is closed")
	}

	frame.lock.Lock()
	defer frame.lock.Unlock()

	rowIdx := frame.topRow
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

	frame.policy.onResize(1)

	if frame.autoDraw {
		frame.draw()
	}

	return newLine, nil
}

func (frame *Frame) Insert(index int) (*Line, error) {
	if frame.IsClosed() {
		return nil, fmt.Errorf("frame is closed")
	}

	frame.lock.Lock()
	defer frame.lock.Unlock()

	if index < 0 || index > len(frame.activeLines) {
		return nil, fmt.Errorf("invalid index given")
	}

	rowIdx := frame.topRow + index
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

	frame.policy.onResize(1)

	if frame.autoDraw {
		frame.draw()
	}

	return newLine, nil
}

func (frame *Frame) indexOf(line *Line) int {
	// find the index of the line object
	matchedIdx := -1
	for idx, item := range frame.activeLines {
		if item == line {
			return idx
		}
	}

	return matchedIdx
}

func (frame *Frame) Remove(line *Line) error {
	if frame.IsClosed() {
		return fmt.Errorf("frame is closed")
	}

	frame.lock.Lock()
	defer frame.lock.Unlock()

	// find the index of the line object
	matchedIdx := frame.indexOf(line)

	if matchedIdx < 0 {
		return fmt.Errorf("could not find line in frame")
	}

	// activeLines that are removed must be closed since any further writes will result in line clashes
	frame.activeLines[matchedIdx].close()
	contents := frame.activeLines[matchedIdx].buffer

	// erase the contents of the last line of the Frame, but persist the line buffer
	if frame.footer != nil {
		frame.clearRows = append(frame.clearRows, frame.footer.row)
	} else {
		frame.clearRows = append(frame.clearRows, frame.activeLines[len(frame.activeLines)-1].row)
	}

	// Remove the line entry from the list
	frame.activeLines = append(frame.activeLines[:matchedIdx], frame.activeLines[matchedIdx+1:]...)

	// move each line index ahead of the deleted element
	for idx := matchedIdx; idx < len(frame.activeLines); idx++ {
		frame.activeLines[idx].move(-1)
	}

	if frame.footer != nil {
		frame.footer.move(-1)
	}

	// apply policies
	if frame.config.TrailOnRemove {
		frame.appendTrail(string(contents))
	} else {
		frame.policy.onResize(-1)
	}

	if frame.autoDraw {
		frame.draw()
	}

	return nil
}

func (frame *Frame) Clear() {
	frame.lock.Lock()
	defer frame.lock.Unlock()

	frame.clear()

	if frame.autoDraw {
		frame.draw()
	}
}

func (frame *Frame) clear() {
	if frame.header != nil {
		frame.clearRows = append(frame.clearRows, frame.header.row)
	}

	for _, line := range frame.activeLines {
		frame.clearRows = append(frame.clearRows, line.row)
	}

	if frame.footer != nil {
		frame.clearRows = append(frame.clearRows, frame.footer.row)
	}
}


func (frame *Frame) Close() {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	frame.close()
	// frame.policy.onClose()
	frame.draw()
}

func (frame *Frame) close() {
	if frame.header != nil {
		frame.header.close()
	}

	for _, line := range frame.activeLines {
		line.close()
	}

	if frame.footer != nil {
		frame.footer.close()
	}

	frame.closed = true
}

// todo: I think this should be decided by the user via a Close() action, not by the indication of closed lines
// since you can always add another line... you don't know when an empty frame should remain open or not
func (frame *Frame) IsClosed() bool {
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

func (frame *Frame) Move(rows int) {
	motion := frame.policy.allowedMotion(rows)
	if motion == 0 {
		return
	}

	frame.lock.Lock()
	defer frame.lock.Unlock()

	frame.move(motion)

	if frame.autoDraw {
		frame.draw()
	}
}

func (frame *Frame) move(rows int) {
	frame.topRow += rows

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

func (frame *Frame) Draw() (errs []error) {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	return frame.draw()
}

func (frame *Frame) draw() (errs []error) {
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
			writeAtRow(frame.trailRows[0], frame.topRow-len(frame.trailRows)+idx)
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
		writeAtRow(message, frame.topRow-len(frame.trailRows)+idx)
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

	if frame.IsClosed() {
		setCursorRow(frame.topRow + frame.Height())
	}

	return errs
}

func (frame *Frame) Wait() {
	frame.closeSignal.Wait()
	// setCursorRow(frame.topRow + frame.height())
}
