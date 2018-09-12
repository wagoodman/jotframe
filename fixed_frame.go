package jotframe


func NewFixedFrame(rows int, hasHeader, hasFooter bool) (*FixedFrame, error) {
	currentRow, err := GetCursorRow()
	if err != nil {
		return nil, err
	}

	return NewFixedFrameAt(rows, hasHeader, hasFooter, currentRow)
}

func NewFixedFrameAt(rows int, hasHeader, hasFooter bool, destinationRow int) (*FixedFrame, error) {
	innerFrame, err := newLogicalFrameAt(rows, hasHeader, hasFooter, destinationRow)
	frame := &FixedFrame{
		frame: innerFrame,
	}
	frame.frame.updateFn = frame.update

	return frame, err
}

func (frame *FixedFrame) Header() *Line {
	return frame.frame.header
}

func (frame *FixedFrame) Footer() *Line {
	return frame.frame.footer
}

func (frame *FixedFrame) Lines() []*Line {
	return frame.frame.activeLines
}

func (frame *FixedFrame) Append() (*Line, error) {
	defer frame.frame.updateAndDraw()
	return frame.frame.Append()
}

func (frame *FixedFrame) Prepend() (*Line, error) {
	defer frame.frame.updateAndDraw()
	return frame.frame.Prepend()
}

func (frame *FixedFrame) Insert(index int) (*Line, error) {
	defer frame.frame.updateAndDraw()
	return frame.frame.Insert(index)
}

func (frame *FixedFrame) Remove(line *Line) error {
	defer frame.frame.updateAndDraw()
	return frame.frame.Remove(line)
}

func (frame *FixedFrame) Close() error {
	return frame.frame.Close()
}

func (frame *FixedFrame) Clear() error {
	frame.frame.lock.Lock()
	defer frame.frame.lock.Unlock()
	defer frame.frame.updateAndDraw()

	return frame.frame.clear()
}

func (frame *FixedFrame) ClearAndClose() error {
	frame.frame.lock.Lock()
	defer frame.frame.lock.Unlock()
	defer frame.frame.updateAndDraw()

	err := frame.frame.clear()
	if err != nil {
		return err
	}
	return frame.frame.close()
}

// update any positions based on external data and redraw
func (frame *FixedFrame) update() error {
	// height := len(frame.frame.activeLines)
	// if frame.frame.header != nil {
	// 	height++
	// }
	// if frame.frame.footer != nil {
	// 	height++
	// }
	//
	// targetFrameStartIndex := (terminalHeight - height)+1
	// if frame.frame.frameStartIdx != targetFrameStartIndex {
	// 	// reset the frame and all activeLines to the correct offset
	// 	offset := targetFrameStartIndex - frame.frame.frameStartIdx
	//
	// 	return frame.frame.move(offset)
	// }
	return nil
}
