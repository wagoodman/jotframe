package jotframe

func NewBottomFrame(rows int, hasHeader, hasFooter bool) *BottomFrame {
	height := rows
	if hasHeader {
		height++
	}
	if hasFooter {
		height++
	}

	// todo: why plus 1?
	frameTopRow := (terminalHeight - height) + 1

	innerFrame := newLogicalFrameAt(rows, hasHeader, hasFooter, frameTopRow)
	frame := &BottomFrame{
		frame: innerFrame,
		lock: getScreenLock(),
	}
	frame.frame.updateFn = frame.update

	return frame
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

func (frame *BottomFrame) AppendTrail(str string) {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.frame.updateAndDraw()

	// write the removed line to the trail log + move the frame down (while advancing the frame)
	frame.frame.appendTrail(str)
	// frame.frame.move(1)
	frame.frame.rowPreAdvancements += 1
}

func (frame *BottomFrame) Append() (*Line, error) {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.frame.updateAndDraw()

	// appended rows should appear to move upwards on the screen, which means that we should
	// move the entire frame upwards 1 line while making more screen space by 1 line
	frame.frame.move(-1)
	frame.frame.rowPreAdvancements += 1

	return frame.frame.append()
}

func (frame *BottomFrame) Prepend() (*Line, error) {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.frame.updateAndDraw()

	return frame.frame.prepend()
}

func (frame *BottomFrame) Insert(index int) (*Line, error) {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.frame.updateAndDraw()

	return frame.frame.insert(index)
}

func (frame *BottomFrame) Remove(line *Line) error {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.frame.updateAndDraw()

	// write the removed line to the trail log + move the frame down
	frame.frame.appendTrail(string(line.buffer))
	frame.frame.move(1)

	return frame.frame.remove(line)
}

func (frame *BottomFrame) Close() error {
	frame.lock.Lock()
	defer frame.lock.Unlock()

	return frame.frame.close()
}

func (frame *BottomFrame) Clear() error {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.frame.updateAndDraw()

	return frame.frame.clear()
}

func (frame *BottomFrame) ClearAndClose() error {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.frame.updateAndDraw()

	err := frame.frame.clear()
	if err != nil {
		return err
	}
	return frame.frame.close()
}

// update any positions based on external data and redraw
func (frame *BottomFrame) update() error {
	height := frame.frame.height()
	targetFrameStartIndex := (terminalHeight - height)+1
	if frame.frame.frameStartIdx != targetFrameStartIndex {
		// reset the frame and all activeLines to the correct offset. This must be done with new
		// lines since we should not overwrite the trail rows above the frame.
		frame.frame.rowPreAdvancements += frame.frame.frameStartIdx - targetFrameStartIndex
	}
	return nil
}
