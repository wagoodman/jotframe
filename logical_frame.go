package jotframe

import (
	"fmt"
)

func newLogicalFrameAt(rows int, hasHeader, hasFooter bool, destinationRow int) (*logicalFrame, error) {
	// todo: check real screen dimensions for moving past bottom of screen
	if destinationRow < 0 {
		return nil, fmt.Errorf("unable to move past screen dimensions")
	}

	frame := &logicalFrame{}
	frame.lock = getScreenLock()
	frame.frameStartIdx = destinationRow

	var relativeRow int
	if hasHeader {
		frame.header = NewLine(frame.frameStartIdx + relativeRow)
		relativeRow++
	}
	for idx := 0; idx < rows; idx++ {
		frame.Append()
	}
	if hasFooter {
		frame.footer = NewLine(frame.frameStartIdx + len(frame.activeLines) + relativeRow)
		relativeRow++
	}

	registerFrame(frame)

	return frame, nil
}

func (frame *logicalFrame) Append() (*Line, error) {
	if frame.closed {
		return nil, fmt.Errorf("frame is closed")
	}

	frame.lock.Lock()
	defer frame.lock.Unlock()

	var rowIdx int
	if len(frame.activeLines) > 0 {
		rowIdx = frame.activeLines[len(frame.activeLines)-1].row + 1
	} else {
		rowIdx = frame.frameStartIdx +1
	}

	newLine := NewLine(rowIdx)
	frame.activeLines = append(frame.activeLines, newLine)

	if frame.footer != nil {
		frame.footer.move(1)
	}

	return newLine, nil
}

func (frame *logicalFrame) Prepend() (*Line, error) {
	if frame.closed {
		return nil, fmt.Errorf("frame is closed")
	}

	frame.lock.Lock()
	defer frame.lock.Unlock()

	newLine := NewLine(frame.frameStartIdx +1)
	for _, line := range frame.activeLines {
		line.move(1)
	}
	frame.activeLines = append([]*Line{newLine}, frame.activeLines...)

	if frame.footer != nil {
		frame.footer.move(1)
	}

	return newLine, nil
}

func (frame *logicalFrame) Insert(index int) (*Line, error) {
	if frame.closed {
		return nil, fmt.Errorf("frame is closed")
	}

	if index < 0 || index > len(frame.activeLines) {
		return nil, fmt.Errorf("invalid index given")
	}

	frame.lock.Lock()
	defer frame.lock.Unlock()

	newLine := NewLine(frame.frameStartIdx +index)

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

func (frame *logicalFrame) Remove(line *Line) error {
	if frame.closed {
		return fmt.Errorf("frame is closed")
	}

	frame.lock.Lock()
	defer frame.lock.Unlock()

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
		frame.clearLines = append(frame.clearLines, frame.footer)
	} else {
		frame.clearLines = append(frame.clearLines, frame.activeLines[len(frame.activeLines)-1])
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

func (frame *logicalFrame) Retreat(rows int) error {
	if frame.closed {
		return fmt.Errorf("frame is closed")
	}

	if rows <= 0 {
		return fmt.Errorf("invalid row retreat amount given")
	}
	frame.lock.Lock()
	defer frame.lock.Unlock()

	return frame.move(rows*-1)
}

func (frame *logicalFrame) Advance(rows int) error {
	if frame.closed {
		return fmt.Errorf("frame is closed")
	}

	if rows <= 0 {
		return fmt.Errorf("invalid row advancement given")
	}
	frame.lock.Lock()
	defer frame.lock.Unlock()

	return frame.move(rows)
}

func (frame *logicalFrame) Close() error {
	if frame.closed {
		return fmt.Errorf("frame is closed")
	}

	frame.lock.Lock()
	defer frame.lock.Unlock()

	return frame.close()
}


func (frame *logicalFrame) clear() error {

	if frame.header != nil {
		frame.clearLines = append(frame.clearLines, frame.header)
	}

	for _, line := range frame.activeLines {
		frame.clearLines = append(frame.clearLines, line)
	}

	if frame.footer != nil {
		frame.clearLines = append(frame.clearLines, frame.footer)
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

	// todo: make this clear only the affected lines (instead of whole frame)
	// erase any affected rows
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


func (frame *logicalFrame) updateAndDraw() {
	if frame.updateFn != nil {
		frame.updateFn()
	}
	frame.draw()
}

func (frame *logicalFrame) draw() error {

	// clear any marked lines (preserving the buffer)
	for _, line := range frame.clearLines {
		line.clear(true)
	}
	frame.clearLines = make([]*Line, 0)

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