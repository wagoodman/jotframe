package frame

func newBottomFrame(config Config) *bottomFrame {
	height := config.Lines
	if config.HasHeader {
		height++
	}
	if config.HasFooter {
		height++
	}

	innerFrame := newLogicalFrame(config)
	frame := &bottomFrame{
		logicalFrame: innerFrame,
		lock:         getScreenLock(),
		config:       config,
	}
	frame.logicalFrame.updateFn = frame.update

	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.logicalFrame.updateAndDraw()

	// make screen realestate if the cursor is already near the bottom row (this preservers the users existing terminal outpu)
	// one assumption is made: the current cursor position is where history (may) start.
	frameHeight := frame.logicalFrame.Height()
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

func (frame *bottomFrame) Config() Config {
	return frame.config
}

func (frame *bottomFrame) StartIdx() int {
	return frame.logicalFrame.StartIdx()
}

func (frame *bottomFrame) Height() int {
	return frame.logicalFrame.Height()
}

func (frame *bottomFrame) Header() *Line {
	return frame.logicalFrame.header
}

func (frame *bottomFrame) Footer() *Line {
	return frame.logicalFrame.footer
}

func (frame *bottomFrame) Lines() []*Line {
	return frame.logicalFrame.activeLines
}

func (frame *bottomFrame) AppendTrail(str string) {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.logicalFrame.updateAndDraw()

	// write the removed line to the trail log + move the frame down (while advancing the frame)
	frame.logicalFrame.appendTrail(str)
	// frame.frame.Move(1)
	frame.logicalFrame.rowAdvancements += 1
}

func (frame *bottomFrame) Append() (*Line, error) {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.logicalFrame.updateAndDraw()

	line, err := frame.logicalFrame.Append()
	if err == nil {
		// appended rows should appear to move upwards on the screen, which means that we should
		// move the entire frame upwards 1 line while making more screen space by 1 line
		frame.logicalFrame.Move(-1)
		frame.logicalFrame.rowAdvancements += 1
	}

	return line, err
}

func (frame *bottomFrame) Prepend() (*Line, error) {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.logicalFrame.updateAndDraw()

	line, err := frame.logicalFrame.Prepend()
	if err == nil {
		// appended rows should appear to move upwards on the screen, which means that we should
		// move the entire frame upwards 1 line while making more screen space by 1 line
		frame.logicalFrame.Move(-1)
		frame.logicalFrame.rowAdvancements += 1
	}

	return line, err
}

func (frame *bottomFrame) Insert(index int) (*Line, error) {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.logicalFrame.updateAndDraw()

	line, err := frame.logicalFrame.Insert(index)
	if err == nil {
		// appended rows should appear to move upwards on the screen, which means that we should
		// move the entire frame upwards 1 line while making more screen space by 1 line
		frame.logicalFrame.Move(-1)
		frame.logicalFrame.rowAdvancements += 1
	}

	return line, err
}

func (frame *bottomFrame) Remove(line *Line) error {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.logicalFrame.updateAndDraw()

	err := frame.logicalFrame.Remove(line)
	if err == nil {
		if frame.config.TrailOnRemove {
			// write the removed line to the trail log + move the frame down
			frame.logicalFrame.appendTrail(string(line.buffer))
		}
		frame.logicalFrame.Move(1)
	}

	return err
}

func (frame *bottomFrame) Close() error {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	// closing the frame moves the cursor, which implies a Update/draw cycle
	defer frame.logicalFrame.updateAndDraw()

	return frame.logicalFrame.Close()
}

func (frame *bottomFrame) Clear() error {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.logicalFrame.updateAndDraw()

	return frame.logicalFrame.Clear()
}

// update any positions based on external data and redraw
func (frame *bottomFrame) update(_ *logicalFrame) error {
	height := frame.logicalFrame.Height()
	targetFrameStartIndex := (terminalHeight - height) + 1
	if frame.logicalFrame.topRow != targetFrameStartIndex {
		// reset the frame and all activeLines to the correct offset. This must be done with new
		// lines since we should not overwrite the trail rows above the frame.
		frame.logicalFrame.rowAdvancements += frame.logicalFrame.topRow - targetFrameStartIndex
	}
	return nil
}

func (frame *bottomFrame) Move(rows int) {
}

func (frame *bottomFrame) Wait() {
	frame.logicalFrame.Wait()
}

func (frame *bottomFrame) IsClosed() bool {
	return frame.logicalFrame.IsClosed()
}

func (frame *bottomFrame) Draw() []error {
	frame.lock.Lock()
	defer frame.lock.Unlock()

	return frame.logicalFrame.Draw()
}

func (frame *bottomFrame) Update() error {
	frame.lock.Lock()
	defer frame.lock.Unlock()

	return frame.logicalFrame.Update()
}
