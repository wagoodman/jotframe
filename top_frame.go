package jotframe

import "fmt"

func NewTopFrame(rows int, hasHeader, hasFooter bool) *TopFrame {
	logicalFrame := newLogicalFrameAt(rows, hasHeader, hasFooter, 0)
	frame := &TopFrame{
		frame: logicalFrame,
		lock:  getScreenLock(),
	}

	// TODO: Add this later
	frame.frame.updateFn = nil

	frame.lock.Lock()
	defer frame.lock.Unlock()
	defer frame.frame.updateAndDraw()

	// Clear screen and set cursor back to top of window
	fmt.Print("\033[2J")

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
	// closing the frame moves the cursor, which implies a update/draw cycle
	defer frame.frame.updateAndDraw()

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

func (frame *TopFrame) Wait() {
	frame.frame.wait()
}

// // update any positions based on external data and redraw
// func (frame *TopFrame) update() error {
// 	if frame.frame.frameStartIdx != 1 {
// 		offset := 1 - frame.frame.frameStartIdx
// 		return frame.frame.move(offset)
// 	}
// 	return nil
// }
