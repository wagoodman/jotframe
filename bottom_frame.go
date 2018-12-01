package jotframe

func NewBottomFrame(rows int, hasHeader, hasFooter bool, includeTrailOnRemove bool) *BottomFrame {
	height := rows
	if hasHeader {
		height++
	}
	if hasFooter {
		height++
	}

	// the screen index starts at 1 (not 0), hence the +1
	frameTopRow := (terminalHeight - height) + 1

	innerFrame := newLogicalFrameAt(rows, hasHeader, hasFooter, frameTopRow)
	frame := &BottomFrame{
		logicalFrame:  innerFrame,
		lock:          getScreenLock(),
		trailOnRemove: includeTrailOnRemove,
	}
	frame.logicalFrame.updateFn = frame.update

	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.logicalFrame.updateAndDraw()

	// make screen realestate if the cursor is already near the bottom row (this preservers the users existing terminal outpu)
	// one assumption is made: the current cursor position is where history (may) start.
	frameHeight := frame.logicalFrame.height()
	currentRow, err := GetCursorRow()
	if err != nil {
		panic(err)
	}

	// if we start drawing now, we'll be past the bottom of the screen, preserve the current terminal history
	if currentRow+frameHeight > terminalHeight {
		offset := currentRow - ((terminalHeight - height) + 1)
		frame.logicalFrame.rowAdvancements += offset
	}

	return frame
}

func (frame *BottomFrame) Header() *Line {
	return frame.logicalFrame.header
}

func (frame *BottomFrame) Footer() *Line {
	return frame.logicalFrame.footer
}

func (frame *BottomFrame) Lines() []*Line {
	return frame.logicalFrame.activeLines
}

func (frame *BottomFrame) AppendTrail(str string) {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.logicalFrame.updateAndDraw()

	// write the removed line to the trail log + move the frame down (while advancing the frame)
	frame.logicalFrame.appendTrail(str)
	// frame.frame.move(1)
	frame.logicalFrame.rowAdvancements += 1
}

func (frame *BottomFrame) Append() (*Line, error) {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.logicalFrame.updateAndDraw()

	line, err := frame.logicalFrame.append()
	if err == nil {
		// appended rows should appear to move upwards on the screen, which means that we should
		// move the entire frame upwards 1 line while making more screen space by 1 line
		frame.logicalFrame.move(-1)
		frame.logicalFrame.rowAdvancements += 1
	}

	return line, err
}

func (frame *BottomFrame) Prepend() (*Line, error) {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.logicalFrame.updateAndDraw()

	line, err := frame.logicalFrame.prepend()
	if err == nil {
		// appended rows should appear to move upwards on the screen, which means that we should
		// move the entire frame upwards 1 line while making more screen space by 1 line
		frame.logicalFrame.move(-1)
		frame.logicalFrame.rowAdvancements += 1
	}

	return line, err
}

func (frame *BottomFrame) Insert(index int) (*Line, error) {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.logicalFrame.updateAndDraw()

	line, err := frame.logicalFrame.insert(index)
	if err == nil {
		// appended rows should appear to move upwards on the screen, which means that we should
		// move the entire frame upwards 1 line while making more screen space by 1 line
		frame.logicalFrame.move(-1)
		frame.logicalFrame.rowAdvancements += 1
	}

	return line, err
}

func (frame *BottomFrame) Remove(line *Line) error {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.logicalFrame.updateAndDraw()

	err := frame.logicalFrame.remove(line)
	if err == nil {
		if frame.trailOnRemove {
			// write the removed line to the trail log + move the frame down
			frame.logicalFrame.appendTrail(string(line.buffer))
		}
		frame.logicalFrame.move(1)
	}

	return err
}

func (frame *BottomFrame) Close() error {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	// closing the frame moves the cursor, which implies a update/draw cycle
	defer frame.logicalFrame.updateAndDraw()

	return frame.logicalFrame.close()
}

func (frame *BottomFrame) Clear() error {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.logicalFrame.updateAndDraw()

	return frame.logicalFrame.clear()
}

func (frame *BottomFrame) ClearAndClose() error {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.logicalFrame.updateAndDraw()

	err := frame.logicalFrame.clear()
	if err != nil {
		return err
	}
	return frame.logicalFrame.close()
}

// update any positions based on external data and redraw
func (frame *BottomFrame) update() error {
	height := frame.logicalFrame.height()
	targetFrameStartIndex := (terminalHeight - height) + 1
	if frame.logicalFrame.frameStartIdx != targetFrameStartIndex {
		// reset the frame and all activeLines to the correct offset. This must be done with new
		// lines since we should not overwrite the trail rows above the frame.
		frame.logicalFrame.rowAdvancements += frame.logicalFrame.frameStartIdx - targetFrameStartIndex
	}
	return nil
}

func (frame *BottomFrame) Wait() {
	frame.logicalFrame.wait()
}
