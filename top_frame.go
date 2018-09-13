package jotframe


func NewTopFrame(rows int, hasHeader, hasFooter bool) *TopFrame {
	innerFrame := newLogicalFrameAt(rows, hasHeader, hasFooter, 1)
	frame := &TopFrame{
		frame: innerFrame,
		lock: getScreenLock(),
	}
	frame.frame.updateFn = frame.update

	return frame
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
	frame.lock.Lock()
	defer frame.lock.Unlock()

	defer frame.frame.updateAndDraw()
	return frame.frame.append()
}

func (frame *TopFrame) Prepend() (*Line, error) {
	frame.lock.Lock()
	defer frame.lock.Unlock()

	defer frame.frame.updateAndDraw()
	return frame.frame.prepend()
}

func (frame *TopFrame) Insert(index int) (*Line, error) {
	frame.lock.Lock()
	defer frame.lock.Unlock()

	defer frame.frame.updateAndDraw()
	return frame.frame.insert(index)
}

func (frame *TopFrame) Remove(line *Line) error {
	frame.lock.Lock()
	defer frame.lock.Unlock()

	defer frame.frame.updateAndDraw()
	return frame.frame.remove(line)
}

func (frame *TopFrame) Close() error {
	frame.lock.Lock()
	defer frame.lock.Unlock()

	return frame.frame.close()
}

func (frame *TopFrame) Clear() error {
	frame.lock.Lock()
	defer frame.lock.Unlock()

	defer frame.frame.updateAndDraw()

	return frame.frame.clear()
}

func (frame *TopFrame) ClearAndClose() error {
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
func (frame *TopFrame) update() error {
	if frame.frame.frameStartIdx != 1 {
		offset := 1 - frame.frame.frameStartIdx
		return frame.frame.move(offset)
	}
	return nil
}
