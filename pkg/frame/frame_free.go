package frame

func newFreeFrame(config Config) *floatingFrame {

	innerFrame := newLogicalFrame(config)
	frame := &floatingFrame{
		logicalFrame: innerFrame,
		lock:         getScreenLock(),
		config:       config,
	}
	frame.logicalFrame.updateFn = nil

	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.logicalFrame.updateAndDraw()

	return frame
}

func (frame *floatingFrame) Config() Config {
	return frame.config
}

func (frame *floatingFrame) StartIdx() int {
	return frame.logicalFrame.StartIdx()
}

func (frame *floatingFrame) Height() int {
	return frame.logicalFrame.Height()
}

func (frame *floatingFrame) Header() *Line {
	return frame.logicalFrame.header
}

func (frame *floatingFrame) Footer() *Line {
	return frame.logicalFrame.footer
}

func (frame *floatingFrame) Lines() []*Line {
	return frame.logicalFrame.activeLines
}

// func (frame *floatingFrame) appendTrail(str string) {
// 	frame.lock.Lock()
// 	defer frame.lock.Unlock()
// 	defer frame.logicalFrame.updateAndDraw()
//
// 	// write the removed line to the trail log + move the frame down (while advancing the frame)
// 	frame.logicalFrame.appendTrail(str)
// 	if frame.logicalFrame.IsPastScreenBottom() {
// 		// frame.frame.Move(-1)
// 		frame.logicalFrame.rowAdvancements += 1
// 	} else {
// 		frame.logicalFrame.Move(1)
// 	}
// }

func (frame *floatingFrame) Append() (*Line, error) {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.logicalFrame.updateAndDraw()

	line, err := frame.logicalFrame.Append()
	if err == nil {
		if frame.logicalFrame.IsPastScreenBottom() {
			// make more screen realestate
			frame.logicalFrame.Move(-1)
			frame.logicalFrame.rowAdvancements += 1
		}
	}

	return line, err
}

func (frame *floatingFrame) Prepend() (*Line, error) {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.logicalFrame.updateAndDraw()

	line, err := frame.logicalFrame.Prepend()
	if err == nil {
		if frame.logicalFrame.IsPastScreenBottom() {
			// make more screen realestate
			frame.logicalFrame.Move(-1)
			frame.logicalFrame.rowAdvancements += 1
		}
	}

	return line, err
}

func (frame *floatingFrame) Insert(index int) (*Line, error) {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.logicalFrame.updateAndDraw()

	line, err := frame.logicalFrame.Insert(index)
	if err == nil {
		if frame.logicalFrame.IsPastScreenBottom() {
			// make more screen realestate
			frame.logicalFrame.Move(-1)
			frame.logicalFrame.rowAdvancements += 1
		}
	}

	return line, err
}

func (frame *floatingFrame) Remove(line *Line) error {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.logicalFrame.updateAndDraw()

	err := frame.logicalFrame.Remove(line)

	if err == nil && frame.config.TrailOnRemove {
		// write the removed line to the trail log + move the frame down
		frame.logicalFrame.appendTrail(string(line.buffer))
		frame.logicalFrame.Move(1)
	}

	return err
}

func (frame *floatingFrame) Move(rows int) {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.logicalFrame.updateAndDraw()

	frame.logicalFrame.Move(rows)
}

func (frame *floatingFrame) Close() error {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	// closing the frame moves the cursor, which implies a update/draw cycle
	defer frame.logicalFrame.updateAndDraw()

	return frame.logicalFrame.Close()
}

func (frame *floatingFrame) Clear() error {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.logicalFrame.updateAndDraw()

	return frame.logicalFrame.Clear()
}

func (frame *floatingFrame) Wait() {
	frame.logicalFrame.Wait()
}

func (frame *floatingFrame) IsClosed() bool {
	return frame.logicalFrame.IsClosed()
}

func (frame *floatingFrame) Draw() []error {
	frame.lock.Lock()
	defer frame.lock.Unlock()

	return frame.logicalFrame.Draw()
}

func (frame *floatingFrame) Update() error {
	frame.lock.Lock()
	defer frame.lock.Unlock()

	return frame.logicalFrame.Update()
}
