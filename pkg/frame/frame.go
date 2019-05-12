package frame

import (
	"context"
	"fmt"
	"sync"
)

type Frame struct {
	Config Config
	lock   *sync.Mutex

	StartIdx int
	Headers  []*Line
	Lines    []*Line
	Footers  []*Line

	clearRows       []int
	trailRows       []string
	rowAdvancements int

	events      chan ScreenEvent
	policy      Policy
	autoDraw    bool
	closed      bool
	stale       bool
}

func New(config Config) (*Frame, error, context.Context, context.CancelFunc) {
	frame := &Frame{
		StartIdx:    config.startRow,
		Config:      config,
		lock:        getScreen().lock,
		autoDraw:    !config.ManualDraw,
		events:      getScreen().events,
	}

	switch config.PositionPolicy {
	case PolicyFloatOverflow:
		frame.policy = newFloatOverflowPolicy(frame)
	default:
		panic(fmt.Errorf("unknown policy: %v", config.PositionPolicy))
	}

	frame.policy.onInit()

	for idx := 0; idx < config.HeaderRows; idx++ {
		// todo: should headers have closeSignal waitGroups? or should they be nil?
		line := NewLine(frame.StartIdx+idx, frame.events)
		frame.Headers = append(frame.Headers, line)
	}
	for idx := 0; idx < config.Lines; idx++ {
		line := NewLine(frame.StartIdx+config.HeaderRows+idx, frame.events)
		frame.Lines = append(frame.Lines, line)
	}
	for idx := 0; idx < config.FooterRows; idx++ {
		// todo: should footers have closeSignal waitGroups? or should they be nil?
		line := NewLine(frame.StartIdx+config.HeaderRows+config.Lines+idx, frame.events)
		frame.Footers = append(frame.Footers, line)
	}

	// register frame before drawing to screen
	if !config.test {
		err := getScreen().register(frame)
		if err != nil {
			return nil, err, nil, nil
		}
	}

	// todo: it's bad that the constructor is writing out to the screen... is it avoidable?
	// adjust the screen such that a known good starting condition is in place
	frame.lock.Lock()
	defer frame.lock.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	if !config.test {
		go getScreen().Run(ctx)
	}
	frame.draw()

	return frame, nil, ctx, cancel
}

func (frame *Frame) newLine(rowIdx int) *Line {
	newLine := NewLine(rowIdx, frame.events)
	return newLine
}

func (frame *Frame) SetAutoDraw(enabled bool) {
	frame.autoDraw = enabled
}

func (frame *Frame) AppendTrail(str string) {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	frame.appendTrail(str)
}

func (frame *Frame) appendTrail(str string) {
	if !frame.policy.isAllowedTrail() {
		return
	}
	frame.trailRows = append(frame.trailRows, str)
	frame.policy.onTrail()

	// todo: what about update/draw here?
}

func (frame *Frame) Height() int {
	height := len(frame.Lines)
	if frame.Headers != nil {
		height++
	}
	if frame.Footers != nil {
		height++
	}
	return height
}

func (frame *Frame) VisibleHeight() int {
	height := frame.Height()
	forwardDrawArea := terminalHeight - (frame.StartIdx - 1)

	if height > forwardDrawArea {
		return forwardDrawArea
	}
	return height

}

func (frame *Frame) resetAdvancements() {
	frame.rowAdvancements = 0
}

func (frame *Frame) advance(rows int) {
	frame.rowAdvancements += rows
}

func (frame *Frame) IsPastScreenTop() bool {
	// take into account the rows that will be added to the screen realestate
	futureFrameStartIdx := frame.StartIdx - frame.rowAdvancements

	if futureFrameStartIdx < 1 {
		return true
	}
	return false
}

func (frame *Frame) IsPastScreenBottom() bool {
	height := frame.Height()

	// take into account the rows that will be added to the screen realestate
	futureFrameStartIdx := frame.StartIdx - frame.rowAdvancements

	// if the frame has moved past the bottom of the screen, move it up a bit
	if futureFrameStartIdx+height > terminalHeight {
		return true
	}
	return false
}

func (frame *Frame) AppendHeader() (*Line, error) {
	if frame.IsClosed() {
		return nil, fmt.Errorf("frame is closed")
	}

	frame.lock.Lock()
	defer frame.lock.Unlock()

	var rowIdx int
	if len(frame.Headers) > 0 {
		rowIdx = frame.Headers[len(frame.Headers)-1].row + 1
	} else {
		rowIdx = frame.StartIdx
	}

	newLine := frame.newLine(rowIdx)
	frame.Headers = append(frame.Headers, newLine)

	for _, line := range frame.Lines{
		line.move(1)
	}

	for _, footer := range frame.Footers{
		footer.move(1)
	}

	frame.policy.onResize(1)

	if frame.autoDraw {
		frame.draw()
	}

	return newLine, nil
}

func (frame *Frame) AppendFooter() (*Line, error) {
	if frame.IsClosed() {
		return nil, fmt.Errorf("frame is closed")
	}

	frame.lock.Lock()
	defer frame.lock.Unlock()

	var rowIdx int
	if len(frame.Footers) > 0 {
		rowIdx = frame.Footers[len(frame.Footers)-1].row + 1
	} else {
		rowIdx = len(frame.Lines) + len(frame.Headers)
	}

	newLine := frame.newLine(rowIdx)
	frame.Footers = append(frame.Footers, newLine)

	frame.policy.onResize(1)

	if frame.autoDraw {
		frame.draw()
	}

	return newLine, nil
}


func (frame *Frame) Append() (*Line, error) {
	if frame.IsClosed() {
		return nil, fmt.Errorf("frame is closed")
	}

	frame.lock.Lock()
	defer frame.lock.Unlock()

	var rowIdx int
	if len(frame.Lines) > 0 {
		rowIdx = frame.Lines[len(frame.Lines)-1].row + 1
	} else {
		rowIdx = frame.StartIdx
		if frame.Headers != nil {
			rowIdx += 1
		}

	}

	newLine := frame.newLine(rowIdx)
	frame.Lines = append(frame.Lines, newLine)

	for _, footer := range frame.Footers{
		footer.move(1)
	}

	frame.policy.onResize(1)

	if frame.autoDraw {
		frame.draw()
	}

	return newLine, nil
}

func (frame *Frame) PrependHeader() (*Line, error) {
	if frame.IsClosed() {
		return nil, fmt.Errorf("frame is closed")
	}

	frame.lock.Lock()
	defer frame.lock.Unlock()

	rowIdx := frame.StartIdx

	newLine := frame.newLine(rowIdx)

	for _, header := range frame.Headers {
		header.move(1)
	}

	frame.Headers = append([]*Line{newLine}, frame.Headers...)

	for _, line := range frame.Lines {
		line.move(1)
	}

	for _, footer := range frame.Footers{
		footer.move(1)
	}

	frame.policy.onResize(1)

	if frame.autoDraw {
		frame.draw()
	}

	return newLine, nil
}

func (frame *Frame) PrependFooter() (*Line, error) {
	if frame.IsClosed() {
		return nil, fmt.Errorf("frame is closed")
	}

	frame.lock.Lock()
	defer frame.lock.Unlock()

	rowIdx := frame.StartIdx + len(frame.Lines) + len(frame.Headers)

	newLine := frame.newLine(rowIdx)

	for _, header := range frame.Footers {
		header.move(1)
	}

	frame.Footers = append([]*Line{newLine}, frame.Footers...)


	frame.policy.onResize(1)

	if frame.autoDraw {
		frame.draw()
	}

	return newLine, nil
}

func (frame *Frame) Prepend() (*Line, error) {
	if frame.IsClosed() {
		return nil, fmt.Errorf("frame is closed")
	}

	frame.lock.Lock()
	defer frame.lock.Unlock()

	rowIdx := frame.StartIdx
	if frame.Headers != nil {
		rowIdx += 1
	}

	newLine := frame.newLine(rowIdx)

	for _, line := range frame.Lines {
		line.move(1)
	}
	frame.Lines = append([]*Line{newLine}, frame.Lines...)

	for _, footer := range frame.Footers{
		footer.move(1)
	}

	frame.policy.onResize(1)

	if frame.autoDraw {
		frame.draw()
	}

	return newLine, nil
}

func (frame *Frame) Insert(index int) (*Line, error) {
	if frame.IsClosed() {
		return nil, fmt.Errorf("frame is closed")
	}

	frame.lock.Lock()
	defer frame.lock.Unlock()

	if index < 0 || index > len(frame.Lines) {
		return nil, fmt.Errorf("invalid index given")
	}

	rowIdx := frame.StartIdx + index
	if frame.Headers != nil {
		rowIdx += 1
	}

	newLine := frame.newLine(rowIdx)

	frame.Lines = append(frame.Lines, nil)
	copy(frame.Lines[index+1:], frame.Lines[index:])
	frame.Lines[index] = newLine

	// bump the indexes for other rows
	for idx := index + 1; idx < len(frame.Lines); idx++ {
		frame.Lines[idx].move(1)
	}

	for _, footer := range frame.Footers{
		footer.move(1)
	}

	frame.policy.onResize(1)

	if frame.autoDraw {
		frame.draw()
	}

	return newLine, nil
}

func (frame *Frame) indexOf(line *Line) int {
	// find the index of the line object
	matchedIdx := -1
	for idx, item := range frame.Lines {
		if item == line {
			return idx
		}
	}

	return matchedIdx
}

func (frame *Frame) Remove(line *Line) error {
	if frame.IsClosed() {
		return fmt.Errorf("frame is closed")
	}

	frame.lock.Lock()
	defer frame.lock.Unlock()

	// find the index of the line object
	matchedIdx := frame.indexOf(line)

	if matchedIdx < 0 {
		return fmt.Errorf("could not find line in frame")
	}

	// Lines that are removed must be closed since any further writes will result in line clashes
	frame.Lines[matchedIdx].close()
	contents := frame.Lines[matchedIdx].buffer

	// erase the contents of the last line of the Frame, but persist the line buffer
	if len(frame.Footers) >0 {
		frame.clearRows = append(frame.clearRows, frame.Footers[len(frame.Footers)-1].row)
	} else {
		frame.clearRows = append(frame.clearRows, frame.Lines[len(frame.Lines)-1].row)
	}

	// Remove the line entry from the list
	frame.Lines = append(frame.Lines[:matchedIdx], frame.Lines[matchedIdx+1:]...)

	// move each line index ahead of the deleted element
	for idx := matchedIdx; idx < len(frame.Lines); idx++ {
		frame.Lines[idx].move(-1)
	}

	for _, footer := range frame.Footers{
		footer.move(-1)
	}

	// apply policies
	if frame.Config.TrailOnRemove {
		frame.appendTrail(string(contents))
	} else {
		frame.policy.onResize(-1)
	}

	if frame.autoDraw {
		frame.draw()
	}

	return nil
}

func (frame *Frame) Clear() {
	frame.lock.Lock()
	defer frame.lock.Unlock()

	frame.clear()

	if frame.autoDraw {
		frame.draw()
	}
}

func (frame *Frame) clear() {
	for _, header := range frame.Headers {
		frame.clearRows = append(frame.clearRows, header.row)
	}

	for _, line := range frame.Lines {
		frame.clearRows = append(frame.clearRows, line.row)
	}

	for _, footer := range frame.Footers {
		frame.clearRows = append(frame.clearRows, footer.row)
	}
}

func (frame *Frame) Close() {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	frame.close()
	// frame.policy.onClose()
	frame.draw()
}

func (frame *Frame) close() {
	for _, header := range frame.Headers{
		header.close()
	}

	for _, line := range frame.Lines {
		line.close()
	}

	for _, footer := range frame.Footers{
		footer.close()
	}

	frame.closed = true
}

// todo: I think this should be decided by the user via a Close() action, not by the indication of closed lines
// since you can always add another line... you don't know when an empty frame should remain open or not
func (frame *Frame) IsClosed() bool {
	// if frame.Headers != nil {
	// 	if !frame.Headers.closed {
	// 		return false
	// 	}
	// }
	//
	// for _, line := range frame.Lines {
	// 	if !line.closed {
	// 		return false
	// 	}
	// }
	//
	// if frame.Footers != nil {
	// 	if !frame.Footers.closed {
	// 		return false
	// 	}
	// }
	// return true
	return frame.closed
}

func (frame *Frame) Move(rows int) {
	motion := frame.policy.isAllowedMotion(rows)
	if motion == 0 {
		return
	}

	frame.lock.Lock()
	defer frame.lock.Unlock()

	frame.move(motion)

	if frame.autoDraw {
		frame.draw()
	}
}

func (frame *Frame) move(rows int) {
	frame.StartIdx += rows

	// todo: instead of clearing all frame lines, only clear the ones affected
	frame.clear()

	// bump rows and redraw entire frame
	for _, header := range frame.Headers{
		header.move(rows)
	}
	for _, line := range frame.Lines {
		line.move(rows)
	}
	for _, footer := range frame.Footers{
		footer.move(rows)
	}
}

func (frame *Frame) Draw() (errs []error) {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	return frame.draw()
}

func (frame *Frame) draw() (errs []error) {
	errs = make([]error, 0)

	// clear any marked lines (preserving the buffer) while these indexes still exist
	for _, row := range frame.clearRows {
		frame.events <- ScreenEvent{
			row:   row,
			value: []byte{},
		}
	}
	frame.clearRows = make([]int, 0)

	// advance the screen while adding any trail lines
	for idx := 0; idx < frame.rowAdvancements; idx++ {
		getScreen().advance(1)
		if idx < len(frame.trailRows) {
			getScreen().writeAtRow(frame.trailRows[0], frame.StartIdx-len(frame.trailRows)+idx)
			if len(frame.trailRows) >= 1 {
				frame.trailRows = frame.trailRows[1:]
			} else {
				frame.trailRows = make([]string, 0)
			}
		}
	}
	frame.rowAdvancements = 0

	// append any remaining trail rows
	for idx, message := range frame.trailRows {
		getScreen().writeAtRow(message, frame.StartIdx-len(frame.trailRows)+idx)
	}
	frame.trailRows = make([]string, 0)

	// paint all stale lines to the screen
	for _, header := range frame.Headers {
		if header.stale || frame.stale {
			_, err := header.write(header.buffer)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}

	for _, line := range frame.Lines {
		if line.stale || frame.stale {
			_, err := line.write(line.buffer)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}

	for _, footer := range frame.Footers {
		if footer.stale || frame.stale {
			_, err := footer.write(footer.buffer)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}

	// if frame.IsClosed() {
	// 	setCursorRow(frame.StartIdx + frame.Height())
	// }

	return errs
}

