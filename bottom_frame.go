package jotframe

func NewBottomFrame(rows int, hasHeader, hasFooter bool) (*BottomFrame, error) {
	height := rows
	if hasHeader {
		height++
	}
	if hasFooter {
		height++
	}

	// todo: why plus 1?
	frameTopRow := (terminalHeight - height) + 1

	innerFrame, err := newLogicalFrameAt(rows, hasHeader, hasFooter, frameTopRow)
	frame := &BottomFrame{
		frame: innerFrame,
	}
	frame.frame.updateFn = frame.update

	return frame, err
}

func (frame *BottomFrame) Header() *Line {
	return frame.frame.header
}

func (frame *BottomFrame) Footer() *Line {
	return frame.frame.footer
}

func (frame *BottomFrame) Lines() []*Line {
	return frame.frame.activeLines
}

func (frame *BottomFrame) Append() (*Line, error) {
	defer frame.frame.updateAndDraw()
	return frame.frame.Append()
}

func (frame *BottomFrame) Prepend() (*Line, error) {
	defer frame.frame.updateAndDraw()
	return frame.frame.Prepend()
}

func (frame *BottomFrame) Insert(index int) (*Line, error) {
	defer frame.frame.updateAndDraw()
	return frame.frame.Insert(index)
}

func (frame *BottomFrame) Remove(line *Line) error {
	defer frame.frame.updateAndDraw()
	return frame.frame.Remove(line)
}

func (frame *BottomFrame) Close() error {
	return frame.frame.Close()
}

func (frame *BottomFrame) Clear() error {
	frame.frame.lock.Lock()
	defer frame.frame.lock.Unlock()
	defer frame.frame.updateAndDraw()

	return frame.frame.clear()
}

func (frame *BottomFrame) ClearAndClose() error {
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
func (frame *BottomFrame) update() error {
	height := len(frame.frame.activeLines)
	if frame.frame.header != nil {
		height++
	}
	if frame.frame.footer != nil {
		height++
	}

	targetFrameStartIndex := (terminalHeight - height)+1
	if frame.frame.frameStartIdx != targetFrameStartIndex {
		// reset the frame and all activeLines to the correct offset
		offset := targetFrameStartIndex - frame.frame.frameStartIdx

		return frame.frame.move(offset)
	}
	return nil
}
