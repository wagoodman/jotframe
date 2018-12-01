package jotframe

func NewFixedFrame(rows int, hasHeader, hasFooter, includeTrailOnRemove bool) *FixedFrame {
	currentRow, err := GetCursorRow()
	if err != nil {
		panic(err)
	}

	return NewFixedFrameAt(rows, hasHeader, hasFooter, includeTrailOnRemove, currentRow)
}

func NewFixedFrameAt(rows int, hasHeader, hasFooter, includeTrailOnRemove bool, destinationRow int) *FixedFrame {
	innerFrame := newLogicalFrameAt(rows, hasHeader, hasFooter, destinationRow)
	frame := &FixedFrame{
		logicalFrame:  innerFrame,
		lock:          getScreenLock(),
		trailOnRemove: includeTrailOnRemove,
	}
	frame.logicalFrame.updateFn = nil

	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.logicalFrame.updateAndDraw()

	return frame
}

func (frame *FixedFrame) Header() *Line {
	return frame.logicalFrame.header
}

func (frame *FixedFrame) Footer() *Line {
	return frame.logicalFrame.footer
}

func (frame *FixedFrame) Lines() []*Line {
	return frame.logicalFrame.activeLines
}

func (frame *FixedFrame) AppendTrail(str string) {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.logicalFrame.updateAndDraw()

	// write the removed line to the trail log + move the frame down (while advancing the frame)
	frame.logicalFrame.appendTrail(str)
	if frame.logicalFrame.isAtOrPastScreenBottom() {
		// frame.frame.move(-1)
		frame.logicalFrame.rowAdvancements += 1
	} else {
		frame.logicalFrame.move(1)
	}
}

func (frame *FixedFrame) Append() (*Line, error) {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.logicalFrame.updateAndDraw()

	line, err := frame.logicalFrame.append()
	if err == nil {
		if frame.logicalFrame.isAtOrPastScreenBottom() {
			// make more screen realestate
			frame.logicalFrame.move(-1)
			frame.logicalFrame.rowAdvancements += 1
		}
	}

	return line, err
}

func (frame *FixedFrame) Prepend() (*Line, error) {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.logicalFrame.updateAndDraw()

	line, err := frame.logicalFrame.prepend()
	if err == nil {
		if frame.logicalFrame.isAtOrPastScreenBottom() {
			// make more screen realestate
			frame.logicalFrame.move(-1)
			frame.logicalFrame.rowAdvancements += 1
		}
	}

	return line, err
}

func (frame *FixedFrame) Insert(index int) (*Line, error) {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.logicalFrame.updateAndDraw()

	line, err := frame.logicalFrame.insert(index)
	if err == nil {
		if frame.logicalFrame.isAtOrPastScreenBottom() {
			// make more screen realestate
			frame.logicalFrame.move(-1)
			frame.logicalFrame.rowAdvancements += 1
		}
	}

	return line, err
}

func (frame *FixedFrame) Remove(line *Line) error {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.logicalFrame.updateAndDraw()

	err := frame.logicalFrame.remove(line)

	if err == nil && frame.trailOnRemove {
		// write the removed line to the trail log + move the frame down
		frame.logicalFrame.appendTrail(string(line.buffer))
		frame.logicalFrame.move(1)
	}

	return err
}

func (frame *FixedFrame) Move(rows int) error {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.logicalFrame.updateAndDraw()

	return frame.logicalFrame.move(rows)
}

func (frame *FixedFrame) Close() error {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	// closing the frame moves the cursor, which implies a update/draw cycle
	defer frame.logicalFrame.updateAndDraw()

	return frame.logicalFrame.close()
}

func (frame *FixedFrame) Clear() error {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.logicalFrame.updateAndDraw()

	return frame.logicalFrame.clear()
}

func (frame *FixedFrame) ClearAndClose() error {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.logicalFrame.updateAndDraw()

	err := frame.logicalFrame.clear()
	if err != nil {
		return err
	}
	return frame.logicalFrame.close()
}

func (frame *FixedFrame) Wait() {
	frame.logicalFrame.wait()
}
