package frame

import "fmt"

func newTopFrame(config Config) *topFrame {
	logicalFrame := newLogicalFrame(config)
	frame := &topFrame{
		logicalFrame: logicalFrame,
		lock:         getScreenLock(),
		config:       config,
	}

	// TODO: Add this later
	frame.logicalFrame.updateFn = nil

	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.logicalFrame.updateAndDraw()

	// Clear screen and set cursor back to top of window
	fmt.Print("\033[2J")

	return frame
}

func (frame *topFrame) Config() Config {
	return frame.config
}

func (frame *topFrame) Move(rows int) {
}

func (frame *topFrame) StartIdx() int {
	return frame.logicalFrame.StartIdx()
}

func (frame *topFrame) Height() int {
	return frame.logicalFrame.Height()
}

func (frame *topFrame) Header() *Line {
	return frame.logicalFrame.header
}

func (frame *topFrame) Footer() *Line {
	return frame.logicalFrame.footer
}

func (frame *topFrame) Lines() []*Line {
	return frame.logicalFrame.activeLines
}

func (frame *topFrame) Append() (*Line, error) {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.logicalFrame.updateAndDraw()

	return frame.logicalFrame.Append()
}

func (frame *topFrame) Prepend() (*Line, error) {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.logicalFrame.updateAndDraw()

	return frame.logicalFrame.Prepend()
}

func (frame *topFrame) Insert(index int) (*Line, error) {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.logicalFrame.updateAndDraw()

	return frame.logicalFrame.Insert(index)
}

func (frame *topFrame) Remove(line *Line) error {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.logicalFrame.updateAndDraw()

	return frame.logicalFrame.Remove(line)
}

func (frame *topFrame) Close() error {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	// closing the frame moves the cursor, which implies a update/draw cycle
	defer frame.logicalFrame.updateAndDraw()

	return frame.logicalFrame.Close()
}

func (frame *topFrame) Clear() error {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.logicalFrame.updateAndDraw()

	return frame.logicalFrame.Clear()
}

func (frame *topFrame) IsClosed() bool {
	return frame.logicalFrame.IsClosed()
}

func (frame *topFrame) Wait() {
	frame.logicalFrame.Wait()
}

func (frame *topFrame) Draw() []error {
	frame.lock.Lock()
	defer frame.lock.Unlock()

	return frame.logicalFrame.Draw()
}

func (frame *topFrame) Update() error {
	frame.lock.Lock()
	defer frame.lock.Unlock()

	return frame.logicalFrame.Update()
}

// // update any positions based on external data and redraw
// func (frame *topFrame) Update() error {
// 	if frame.frame.topRow != 1 {
// 		offset := 1 - frame.frame.topRow
// 		return frame.frame.Move(offset)
// 	}
// 	return nil
// }
