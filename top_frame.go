package jotframe


func NewTopFrame(rows int, hasHeader, hasFooter bool) (*TopFrame, error) {
	innerFrame, err := newLogicalFrameAt(rows, hasHeader, hasFooter, 0)
	frame := &TopFrame{
		frame: innerFrame,
	}
	frame.frame.updateFn = frame.update

	return frame, err
}

func (frame *TopFrame) Header() *Line {
	return frame.frame.header
}

func (frame *TopFrame) Footer() *Line {
	return frame.frame.footer
}

func (frame *TopFrame) Lines() []*Line {
	return frame.frame.activeLines
}

func (frame *TopFrame) Append() (*Line, error) {
	defer frame.frame.updateAndDraw()
	return frame.frame.Append()
}

func (frame *TopFrame) Prepend() (*Line, error) {
	defer frame.frame.updateAndDraw()
	return frame.frame.Prepend()
}

func (frame *TopFrame) Insert(index int) (*Line, error) {
	defer frame.frame.updateAndDraw()
	return frame.frame.Insert(index)
}

func (frame *TopFrame) Remove(line *Line) error {
	defer frame.frame.updateAndDraw()
	return frame.frame.Remove(line)
}

func (frame *TopFrame) Close() error {
	return frame.frame.Close()
}

func (frame *TopFrame) Clear() error {
	frame.frame.lock.Lock()
	defer frame.frame.lock.Unlock()
	defer frame.frame.updateAndDraw()

	return frame.frame.clear()
}

func (frame *TopFrame) ClearAndClose() error {
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
func (frame *TopFrame) update() error {
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
