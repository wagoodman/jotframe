package jotframe

type TopFrame struct {
	rawFrame *FixedFrame
}

func NewTopFrame(rows int, hasHeader, hasFooter bool) (*TopFrame, error) {
	innerFrame, err := NewFixedFrameAt(rows, hasHeader, hasFooter, 0)
	frame := &TopFrame{
		rawFrame: innerFrame,
	}
	frame.rawFrame.updateFn = frame.update

	return frame, err
}

func (frame *TopFrame) Header() *Line {
	return frame.rawFrame.header
}

func (frame *TopFrame) Footer() *Line {
	return frame.rawFrame.footer
}

func (frame *TopFrame) Lines() []*Line {
	return frame.rawFrame.lines
}

func (frame *TopFrame) Append() (*Line, error) {
	return frame.rawFrame.Append()
}

func (frame *TopFrame) Prepend() (*Line, error) {
	return frame.rawFrame.Prepend()
}

func (frame *TopFrame) Remove(line *Line) error {
	return frame.rawFrame.Remove(line)
}

func (frame *TopFrame) Insert(index int) (*Line, error) {
	return frame.rawFrame.Insert(index)
}

func (frame *TopFrame) Advance(rows int) error {
	return frame.rawFrame.Advance(rows)
}

func (frame *TopFrame) Retreat(rows int) error {
	return frame.rawFrame.Retreat(rows)
}

func (frame *TopFrame) Draw() error {
	return frame.rawFrame.Draw()
}

func (frame *TopFrame) Close() error {
	return frame.rawFrame.Close()
}

func (frame *TopFrame) Clear() error {
	return frame.rawFrame.Clear()
}

func (frame *TopFrame) ClearAndClose() error {
	return frame.rawFrame.ClearAndClose()
}

func (frame *TopFrame) update() bool {
	return false
}