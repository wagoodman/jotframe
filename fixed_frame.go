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
		frame: innerFrame,
		lock: getScreenLock(),
		trailOnRemove: includeTrailOnRemove,
	}
	frame.frame.updateFn = nil

	return frame
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

func (frame *FixedFrame) AppendTrail(str string) {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.frame.updateAndDraw()

	// write the removed line to the trail log + move the frame down (while advancing the frame)
	frame.frame.appendTrail(str)
	// frame.frame.move(1)
	frame.frame.rowPreAdvancements += 1
}

func (frame *FixedFrame) Append() (*Line, error) {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.frame.updateAndDraw()

	return frame.frame.append()
}

func (frame *FixedFrame) Prepend() (*Line, error) {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.frame.updateAndDraw()

	return frame.frame.prepend()
}

func (frame *FixedFrame) Insert(index int) (*Line, error) {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.frame.updateAndDraw()

	return frame.frame.insert(index)
}

func (frame *FixedFrame) Remove(line *Line) error {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.frame.updateAndDraw()

	if frame.trailOnRemove {
		// write the removed line to the trail log + move the frame down
		frame.frame.appendTrail(string(line.buffer))
	}
	frame.frame.move(1)

	return frame.frame.remove(line)
}

func (frame *FixedFrame) Move(rows int) error {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.frame.updateAndDraw()

	return frame.frame.move(rows)
}

func (frame *FixedFrame) Close() error {
	frame.lock.Lock()
	defer frame.lock.Unlock()

	return frame.frame.close()
}

func (frame *FixedFrame) Clear() error {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.frame.updateAndDraw()

	return frame.frame.clear()
}

func (frame *FixedFrame) ClearAndClose() error {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.frame.updateAndDraw()

	err := frame.frame.clear()
	if err != nil {
		return err
	}
	return frame.frame.close()
}
