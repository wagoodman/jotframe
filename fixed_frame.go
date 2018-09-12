package jotframe

import (
	"sync"
	"fmt"
)


type FixedFrame struct  {
	header *Line
	lines  []*Line
	footer *Line

	lock               *sync.Mutex
	updateFn           func() bool
	startScreenIndex   int
	closed             bool
}


func NewFixedFrameAt(rows int, hasHeader, hasFooter bool, destinationRow int) (*FixedFrame, error) {
	// todo: check real screen dimensions for moving past bottom of screen
	if destinationRow < 0 {
		return nil, fmt.Errorf("unable to move FixedFrame past screen dimensions")
	}

	setCursorRow(destinationRow)
	return NewFixedFrame(rows, hasHeader, hasFooter)
}

func NewFixedFrame(rows int, hasHeader, hasFooter bool) (*FixedFrame, error) {
	frame := &FixedFrame{}
	frame.lock = getScreenLock()
	currentRow, err := getCursorRow()
	if err != nil {
		return nil, err
	}
	frame.startScreenIndex = currentRow

	var relativeRow int
	if hasHeader {
		frame.header = NewLine(frame.startScreenIndex+relativeRow)
		relativeRow++
	}
	for idx := 0; idx < rows; idx++ {
		frame.Append()
	}
	if hasFooter {
		frame.footer = NewLine(frame.startScreenIndex+len(frame.lines)+relativeRow)
		relativeRow++
	}

	registerFrame(frame)

	return frame, nil
}

func (frame *FixedFrame) Header() *Line {
	return frame.header
}

func (frame *FixedFrame) Footer() *Line {
	return frame.footer
}

func (frame *FixedFrame) Lines() []*Line {
	return frame.lines
}

func (frame *FixedFrame) Append() (*Line, error) {
	if frame.closed {
		return nil, fmt.Errorf("FixedFrame is closed")
	}

	frame.lock.Lock()
	defer frame.lock.Unlock()

	var rowIdx int
	if len(frame.lines) > 0 {
		rowIdx = frame.lines[len(frame.lines)-1].row + 1
	} else {
		rowIdx = frame.startScreenIndex+1
	}

	newLine := NewLine(rowIdx)
	frame.lines = append(frame.lines, newLine)

	// only the new line should be updated on the screen and the footer
	newLine.clear(true)

	if frame.footer != nil {
		frame.footer.row++
	}

	// the frame placement may need to change based on external factors (and the new frame shape)
	didRedraw := frame.update()

	// if the frame placement did not move, then we still need to redraw the affected lines
	if !didRedraw {
		newLine.clear(true)

		if frame.footer != nil {
			frame.footer.write(frame.footer.buffer)
		}
	}


	return newLine, nil
}

func (frame *FixedFrame) Prepend() (*Line, error) {
	if frame.closed {
		return nil, fmt.Errorf("FixedFrame is closed")
	}

	frame.lock.Lock()
	defer frame.lock.Unlock()

	newLine := NewLine(frame.startScreenIndex+1)
	for _, line := range frame.lines {
		line.row++
	}
	frame.lines = append([]*Line{newLine}, frame.lines...)

	if frame.footer != nil {
		frame.footer.row++
	}

	// the frame placement may need to change based on external factors (and the new frame shape)
	didRedraw := frame.update()

	// if the frame placement did not move, then we still need to redraw the affected lines
	if !didRedraw {
		for _, line := range frame.lines {
			line.write(line.buffer)
		}

		if frame.footer != nil {
			frame.footer.write(frame.footer.buffer)
		}
	}

	return newLine, nil
}

func (frame *FixedFrame) Insert(index int) (*Line, error) {
	if frame.closed {
		return nil, fmt.Errorf("FixedFrame is closed")
	}

	if index < 0 || index > len(frame.lines) {
		return nil, fmt.Errorf("invalid index given")
	}

	frame.lock.Lock()
	defer frame.lock.Unlock()

	newLine := NewLine(frame.startScreenIndex+index)

	frame.lines = append(frame.lines, nil)
	copy(frame.lines[index+1:], frame.lines[index:])
	frame.lines[index] = newLine

	// bump the indexes for other rows
	for idx := index+1; idx < len(frame.lines); idx++ {
		frame.lines[idx].row++
	}

	if frame.footer != nil {
		frame.footer.row++
	}

	// the frame placement may need to change based on external factors (and the new frame shape)
	didRedraw := frame.update()

	// if the frame placement did not move, then we still need to redraw the affected lines
	if !didRedraw {
		for idx := index+1; idx < len(frame.lines); idx++ {
			frame.lines[idx].write(frame.lines[idx].buffer)
		}

		newLine.clear(true)

		if frame.footer != nil {
			frame.footer.write(frame.footer.buffer)
		}
	}

	return newLine, nil
}

func (frame *FixedFrame) Remove(line *Line) error {
	if frame.closed {
		return fmt.Errorf("FixedFrame is closed")
	}

	frame.lock.Lock()
	defer frame.lock.Unlock()

	// find the index of the line object
	matchedIdx := -1
	for idx, item := range frame.lines {
		if item == line {
			matchedIdx = idx
			break
		}
	}

	if matchedIdx < 0 {
		return fmt.Errorf("could not find line in FixedFrame")
	}

	// lines that are removed must be closed since any further writes will result in line clashes
	frame.lines[matchedIdx].close()

	// erase the contents of the last line of the FixedFrame, but persist the line buffer
	if frame.footer != nil {
		frame.footer.clear(true)
	} else {
		frame.lines[len(frame.lines)-1].clear(true)
	}

	// remove the line entry from the list
	frame.lines = append(frame.lines[:matchedIdx], frame.lines[matchedIdx+1:]...)

	// move each line index ahead of the deleted element
	for idx := matchedIdx+1; idx < len(frame.lines); idx++ {
		frame.lines[idx].row--
		frame.lines[idx].write(frame.lines[idx].buffer)
	}

	if frame.footer != nil {
		frame.footer.row--
		frame.footer.write(frame.footer.buffer)
	}

	// the frame placement may need to change based on external factors (and the new frame shape)
	didRedraw := frame.update()

	// if the frame placement did not move, then we still need to redraw the affected lines
	if !didRedraw {
		for idx := matchedIdx+1; idx < len(frame.lines); idx++ {
			frame.lines[idx].write(frame.lines[idx].buffer)
		}

		if frame.footer != nil {
			frame.footer.write(frame.footer.buffer)
		}
	}

	return nil
}

func (frame *FixedFrame) Retreat(rows int) error {
	if frame.closed {
		return fmt.Errorf("FixedFrame is closed")
	}

	if rows <= 0 {
		return fmt.Errorf("invalid row retreat amount given")
	}
	frame.lock.Lock()
	defer frame.lock.Unlock()

	return frame.move(rows*-1)
}

func (frame *FixedFrame) Advance(rows int) error {
	if frame.closed {
		return fmt.Errorf("FixedFrame is closed")
	}

	if rows <= 0 {
		return fmt.Errorf("invalid row advancement given")
	}
	frame.lock.Lock()
	defer frame.lock.Unlock()

	return frame.move(rows)
}

func (frame *FixedFrame) Clear() error {
	if frame.closed {
		return fmt.Errorf("FixedFrame is closed")
	}

	frame.lock.Lock()
	defer frame.lock.Unlock()

	return frame.clear(false)
}

func (frame *FixedFrame) ClearAndClose() error {
	if frame.closed {
		return fmt.Errorf("FixedFrame is closed")
	}

	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.close()

	return frame.clear(false)
}

func (frame *FixedFrame) Close() error {
	if frame.closed {
		return fmt.Errorf("FixedFrame is closed")
	}

	frame.lock.Lock()
	defer frame.lock.Unlock()

	return frame.close()
}

func (frame *FixedFrame) Draw() error {
	if frame.closed {
		return fmt.Errorf("FixedFrame is closed")
	}
	frame.lock.Lock()
	defer frame.lock.Unlock()

	return frame.draw()
}



func (frame *FixedFrame) clear(preserveBuffer bool) error {

	if frame.header != nil {
		frame.header.clear(preserveBuffer)
	}

	for _, line := range frame.lines {
		line.clear(preserveBuffer)
	}

	if frame.footer != nil {
		frame.footer.clear(preserveBuffer)
	}
	return nil
}

func (frame *FixedFrame) close() error {

	if frame.header != nil {
		err := frame.header.close()
		if err != nil {
			return err
		}
	}

	for _, line := range frame.lines {
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

func (frame *FixedFrame) draw() error {

	if frame.header != nil {
		_, err := frame.header.write(frame.header.buffer)
		if err != nil {
			return err
		}
	}

	for _, line := range frame.lines {
		_, err := line.write(line.buffer)
		if err != nil {
			return err
		}
	}

	if frame.footer != nil {
		_, err := frame.footer.write(frame.footer.buffer)
		if err != nil {
			return err
		}
	}
	return nil
}

func (frame *FixedFrame) move(rows int) error {
	// todo: check real screen dimensions for moving past bottom of screen
	if frame.startScreenIndex + rows < 0 {
		return fmt.Errorf("unable to move FixedFrame past screen dimensions")
	}
	frame.startScreenIndex += rows

	// erase any affected rows
	frame.clear(true)

	// bump rows and redraw entire FixedFrame
	if frame.header != nil {
		frame.header.row += rows
	}
	for _, line := range frame.lines {
		line.row += rows
	}
	if frame.footer != nil {
		frame.footer.row += rows
	}

	// the frame placement may need to change based on external factors
	didRedraw := frame.update()
	if !didRedraw {
		return frame.draw()
	}

	return nil
}

// check the terminal size against the current FixedFrame strategy
func (frame *FixedFrame) update() bool {
	if frame.updateFn != nil {
		return frame.updateFn()
	}
	return false
}