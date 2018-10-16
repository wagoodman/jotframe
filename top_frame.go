package jotframe

import "fmt"

func NewTopFrame(rows int, hasHeader, hasFooter bool) *TopFrame {
	logicalFrame := newLogicalFrameAt(rows, hasHeader, hasFooter, 0)
	frame := &TopFrame{
		logicalFrame: logicalFrame,
		lock:         getScreenLock(),
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

func (frame *TopFrame) Header() *Line {
	return frame.logicalFrame.header
}

func (frame *TopFrame) Footer() *Line {
	return frame.logicalFrame.footer
}

func (frame *TopFrame) Lines() []*Line {
	return frame.logicalFrame.activeLines
}

func (frame *TopFrame) Append() (*Line, error) {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.logicalFrame.updateAndDraw()

	return frame.logicalFrame.append()
}

func (frame *TopFrame) Prepend() (*Line, error) {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.logicalFrame.updateAndDraw()

	return frame.logicalFrame.prepend()
}

func (frame *TopFrame) Insert(index int) (*Line, error) {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.logicalFrame.updateAndDraw()

	return frame.logicalFrame.insert(index)
}

func (frame *TopFrame) Remove(line *Line) error {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.logicalFrame.updateAndDraw()

	return frame.logicalFrame.remove(line)
}

func (frame *TopFrame) Close() error {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	// closing the frame moves the cursor, which implies a update/draw cycle
	defer frame.logicalFrame.updateAndDraw()

	return frame.logicalFrame.close()
}

func (frame *TopFrame) Clear() error {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.logicalFrame.updateAndDraw()

	return frame.logicalFrame.clear()
}

func (frame *TopFrame) ClearAndClose() error {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.logicalFrame.updateAndDraw()

	err := frame.logicalFrame.clear()
	if err != nil {
		return err
	}

	return frame.logicalFrame.close()
}

func (frame *TopFrame) Wait() {
	frame.logicalFrame.wait()
}

// // update any positions based on external data and redraw
// func (frame *TopFrame) update() error {
// 	if frame.frame.frameStartIdx != 1 {
// 		offset := 1 - frame.frame.frameStartIdx
// 		return frame.frame.move(offset)
// 	}
// 	return nil
// }
