package jotframe

type BottomFrame struct {
	rawFrame *FixedFrame
}

func NewBottomFrame(rows int, hasHeader, hasFooter bool) (*BottomFrame, error) {
	height := rows
	if hasHeader {
		height++
	}
	if hasFooter {
		height++
	}

	frameTopRow := (terminalHeight - height)+1

	innerFrame, err := NewFixedFrameAt(rows, hasHeader, hasFooter, frameTopRow)
	frame := &BottomFrame{
		rawFrame: innerFrame,
	}
	frame.rawFrame.updateFn = frame.update

	return frame, err
}

func (frame *BottomFrame) Header() *Line {
	return frame.rawFrame.header
}

func (frame *BottomFrame) Footer() *Line {
	return frame.rawFrame.footer
}

func (frame *BottomFrame) Lines() []*Line {
	return frame.rawFrame.lines
}

func (frame *BottomFrame) Append() (*Line, error) {
	return frame.rawFrame.Append()
}

func (frame *BottomFrame) Prepend() (*Line, error) {
	return frame.rawFrame.Prepend()
}

func (frame *BottomFrame) Insert(index int) (*Line, error) {
	return frame.rawFrame.Insert(index)
}

func (frame *BottomFrame) Remove(line *Line) error {
	return frame.rawFrame.Remove(line)
}

func (frame *BottomFrame) Advance(rows int) error {
	return frame.rawFrame.Advance(rows)
}

func (frame *BottomFrame) Retreat(rows int) error {
	return frame.rawFrame.Retreat(rows)
}

func (frame *BottomFrame) Draw() error {
	return frame.rawFrame.Draw()
}

func (frame *BottomFrame) Close() error {
	return frame.rawFrame.Close()
}

func (frame *BottomFrame) Clear() error {
	return frame.rawFrame.Clear()
}

func (frame *BottomFrame) ClearAndClose() error {
	return frame.rawFrame.ClearAndClose()
}

// func (frame *BottomFrame) Info() string {
// 	height := len(frame.rawFrame.lines)
// 	if frame.rawFrame.header != nil {
// 		height++
// 	}
// 	if frame.rawFrame.footer != nil {
// 		height++
// 	}
//
// 	targetFrameStartIndex := terminalHeight - height
// 	offset := targetFrameStartIndex - frame.rawFrame.startScreenIndex
// 	return fmt.Sprintf("Offset: %d (screenHeight:%d  targetIdx:%d  actualIdx:%d  frameHeight:%d)", offset, terminalHeight, targetFrameStartIndex, frame.rawFrame.startScreenIndex, height)
// }

func (frame *BottomFrame) update() bool {
	height := len(frame.rawFrame.lines)
	if frame.rawFrame.header != nil {
		height++
	}
	if frame.rawFrame.footer != nil {
		height++
	}

	targetFrameStartIndex := (terminalHeight - height)+1
	if frame.rawFrame.startScreenIndex != targetFrameStartIndex {
		// reset the frame and all lines to the correct offset
		offset := targetFrameStartIndex - frame.rawFrame.startScreenIndex

		frame.rawFrame.startScreenIndex += offset

		// erase any affected rows so no artifacts are left behind
		frame.rawFrame.clear(true)

		// bump rows and redraw entire FixedFrame
		if frame.rawFrame.header != nil {
			frame.rawFrame.header.row += offset
		}
		for _, line := range frame.rawFrame.lines {
			line.row += offset
		}
		if frame.rawFrame.footer != nil {
			frame.rawFrame.footer.row += offset
		}

		// redraw the entire frame
		frame.rawFrame.draw()

		return true
	}
	return false
}