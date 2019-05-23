package frame

import (
	"context"
	"fmt"
	"sync"
)

type frameSection int

const (
	sectionUnknown frameSection = iota
	sectionHeader
	sectionBody
	sectionFooter
)

var sections = []frameSection{sectionHeader, sectionBody, sectionFooter}

type Frame struct {
	Config Config
	lock   *sync.Mutex

	startIdx    int
	HeaderLines []*Line
	BodyLines   []*Line
	FooterLines []*Line

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
		startIdx: config.startRow,
		Config:   config,
		lock:     getScreen().lock,
		autoDraw: !config.ManualDraw,
		events:   getScreen().events,
	}

	switch config.PositionPolicy {
	case PolicyOverflow:
		frame.policy = newOverflowPolicy(frame)
	default:
		panic(fmt.Errorf("unknown policy: %v", config.PositionPolicy))
	}

	// set the frame start row
	frame.policy.onInit()

	for idx := 0; idx < config.HeaderRows; idx++ {
		// todo: should headers have closeSignal waitGroups? or should they be nil?
		line := NewLine(frame.startIdx+idx, frame.events)
		frame.HeaderLines = append(frame.HeaderLines, line)
	}
	for idx := 0; idx < config.Lines; idx++ {
		line := NewLine(frame.startIdx+config.HeaderRows+idx, frame.events)
		frame.BodyLines = append(frame.BodyLines, line)
	}
	for idx := 0; idx < config.FooterRows; idx++ {
		// todo: should footers have closeSignal waitGroups? or should they be nil?
		line := NewLine(frame.startIdx+config.HeaderRows+config.Lines+idx, frame.events)
		frame.FooterLines = append(frame.FooterLines, line)
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
	newLine.frame = frame
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

func (frame *Frame) visibleBodyLines() int {
	height := 0
	for _, line := range frame.BodyLines {
		height += line.height
	}
	return height
}

func (frame *Frame) visibleHeaderLines() int {
	height := 0
	for _, line := range frame.HeaderLines {
		height += line.height
	}
	return height
}

func (frame *Frame) visibleFooterLines() int {
	height := 0
	for _, line := range frame.FooterLines {
		height += line.height
	}
	return height
}


func (frame *Frame) Height() int {
	return frame.visibleBodyLines() + frame.visibleFooterLines() + frame.visibleHeaderLines()
}

func (frame *Frame) VisibleHeight() int {
	height := frame.Height()
	forwardDrawArea := terminalHeight - (frame.startIdx - 1)

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
	futureFrameStartIdx := frame.startIdx - frame.rowAdvancements

	if futureFrameStartIdx < 1 {
		return true
	}
	return false
}

func (frame *Frame) IsPastScreenBottom() bool {
	height := frame.Height()

	// take into account the rows that will be added to the screen realestate
	futureFrameStartIdx := frame.startIdx - frame.rowAdvancements

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
	headerLines := frame.visibleHeaderLines()
	if headerLines > 0 {
		rowIdx = frame.HeaderLines[headerLines-1].row + 1
	} else {
		rowIdx = frame.startIdx
	}

	newLine := frame.newLine(rowIdx)
	frame.HeaderLines = append(frame.HeaderLines, newLine)

	for _, line := range frame.BodyLines {
		line.move(1)
	}

	for _, footer := range frame.FooterLines {
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
	footerLines := frame.visibleFooterLines()
	bodyLines := frame.visibleBodyLines()
	headerLines := frame.visibleHeaderLines()
	if footerLines > 0 {
		rowIdx = frame.FooterLines[footerLines-1].row + 1
	} else {
		rowIdx = bodyLines + headerLines
	}

	newLine := frame.newLine(rowIdx)
	frame.FooterLines = append(frame.FooterLines, newLine)

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
	bodyLines := frame.visibleBodyLines()
	if bodyLines > 0 {
		rowIdx = frame.BodyLines[bodyLines-1].row + 1
	} else {
		rowIdx = frame.startIdx
		if frame.HeaderLines != nil {
			rowIdx += 1
		}

	}

	newLine := frame.newLine(rowIdx)
	frame.BodyLines = append(frame.BodyLines, newLine)

	for _, footer := range frame.FooterLines {
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

	rowIdx := frame.startIdx

	newLine := frame.newLine(rowIdx)

	for _, header := range frame.HeaderLines {
		header.move(1)
	}

	frame.HeaderLines = append([]*Line{newLine}, frame.HeaderLines...)

	for _, line := range frame.BodyLines {
		line.move(1)
	}

	for _, footer := range frame.FooterLines {
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

	rowIdx := frame.startIdx + frame.visibleBodyLines() + frame.visibleHeaderLines()

	newLine := frame.newLine(rowIdx)

	for _, header := range frame.FooterLines {
		header.move(1)
	}

	frame.FooterLines = append([]*Line{newLine}, frame.FooterLines...)


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

	rowIdx := frame.startIdx
	if frame.HeaderLines != nil {
		rowIdx += 1
	}

	newLine := frame.newLine(rowIdx)

	for _, line := range frame.BodyLines {
		line.move(1)
	}
	frame.BodyLines = append([]*Line{newLine}, frame.BodyLines...)

	for _, footer := range frame.FooterLines {
		footer.move(1)
	}

	frame.policy.onResize(1)

	if frame.autoDraw {
		frame.draw()
	}

	return newLine, nil
}

func (frame *Frame) Insert(index int) (*Line, error) {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	return frame.insert(index, false)
}

func (frame *Frame) insert(index int, show bool) (*Line, error) {
	if frame.IsClosed() {
		return nil, fmt.Errorf("frame is closed")
	}

	if index < 0 || index > frame.visibleBodyLines() {
		return nil, fmt.Errorf("invalid index given")
	}

	rowIdx := frame.startIdx + index + frame.visibleHeaderLines()

	var newLine *Line
	if show {
		newLine = frame.BodyLines[index]
	} else {
		newLine = frame.newLine(rowIdx)
		frame.BodyLines = append(frame.BodyLines, nil)
		copy(frame.BodyLines[index+1:], frame.BodyLines[index:])
		frame.BodyLines[index] = newLine
	}

	// bump the indexes for other rows
	for idx := index + 1; idx < len(frame.BodyLines); idx++ {
		frame.BodyLines[idx].move(1)
	}

	for _, footer := range frame.FooterLines {
		footer.move(1)
	}

	frame.policy.onResize(1)

	if frame.autoDraw {
		frame.draw()
	}

	return newLine, nil
}

func (frame *Frame) section(section frameSection) *[]*Line {
	switch section {
	case sectionHeader:
		return &frame.HeaderLines
	case sectionBody:
		return &frame.BodyLines
	case sectionFooter:
		return &frame.FooterLines
	default:
		return nil
	}
}

func (frame *Frame) fetch(section frameSection, index int) *Line {
	if index < 0 {
		return nil
	}

	var source = frame.section(section)

	if source == nil || index >= len(*source) {
		return nil
	}

	return (*source)[index]
}

func (frame *Frame) indexOf(line *Line) (frameSection, int) {
	// find the index of the line object
	for idx, item := range frame.BodyLines {
		if item == line {
			return sectionBody, idx
		}
	}

	for idx, item := range frame.HeaderLines {
		if item == line {
			return sectionHeader, idx
		}
	}

	for idx, item := range frame.FooterLines {
		if item == line {
			return sectionFooter, idx
		}
	}

	return sectionUnknown, -1
}

func (frame *Frame) lastVisibleLineIdx() (frameSection, int) {
	for secIdx := len(sections)-1; secIdx > 0; secIdx-- {
		secName := sections[secIdx]
		section := frame.section(secName)
		for idx := len(*section) - 1; idx >= 0; idx-- {
			if (*section)[idx].visible {
				return secName, idx
			}
		}
	}

	return sectionUnknown, -1
}

func (frame *Frame) moveAfter(moveAdj int, section frameSection, index int) error {
	for _, iterSection := range sections {
		if iterSection < section {
			continue
		}
		source := frame.section(iterSection)
		var startIdx int
		if section == iterSection {
			startIdx = index
		}
		for idx := startIdx; idx < len(*source); idx++ {
			(*source)[idx].move(moveAdj)
		}
	}
	return nil
}

func (frame *Frame) Remove(line *Line) error {

	frame.lock.Lock()
	defer frame.lock.Unlock()

	return frame.remove(line, false)
}

func (frame *Frame) remove(line *Line, hide bool) error {
	if frame.IsClosed() {
		return fmt.Errorf("frame is closed")
	}

	// find the index of the line object
	section, matchedIdx := frame.indexOf(line)
	if matchedIdx < 0 || section == sectionUnknown {
		return fmt.Errorf("could not find line in frame")
	}
	source := frame.section(section)

	// lines that are removed must be closed since any further writes will result in line clashes
	if !hide {
		line.close()
	}
	contents := line.buffer

	// erase the contents of the last line of the Frame, but persist the line buffer
	if (line.visible && !hide) || hide {
		lastVisibleLineSection, lastVisibleLineIdx := frame.lastVisibleLineIdx()
		if lastVisibleLineSection == sectionUnknown || lastVisibleLineIdx < 0 {
			return fmt.Errorf("frame is empty")
		}
		clearRowsSource := frame.section(lastVisibleLineSection)

		frame.clearRows = append(frame.clearRows, (*clearRowsSource)[lastVisibleLineIdx].row)
	}

	// Remove the line entry from the list
	if !hide {
		*source = append((*source)[:matchedIdx], (*source)[matchedIdx+1:]...)
	} else {
		matchedIdx++
	}

	// no need to adjust any lines if the line was hidden already
	if !line.visible && !hide {
		return nil
	}

	frame.moveAfter(-1, section, matchedIdx)

	// apply policies
	if !hide && frame.Config.TrailOnRemove {
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
	for _, header := range frame.HeaderLines {
		frame.clearRows = append(frame.clearRows, header.row)
	}

	for _, line := range frame.BodyLines {
		frame.clearRows = append(frame.clearRows, line.row)
	}

	for _, footer := range frame.FooterLines {
		frame.clearRows = append(frame.clearRows, footer.row)
	}
}

func (frame *Frame) Close() {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	frame.close()
	// frame.policy.onClose()

	// todo: make this a screen write, starting at the last line of the frame and inserting a
	// newline char with the screen loop (should we kill the screen loop afterwards?)
	fmt.Println()


	frame.draw()
}

func (frame *Frame) close() {
	for _, header := range frame.HeaderLines {
		header.close()
	}

	for _, line := range frame.BodyLines {
		line.close()
	}

	for _, footer := range frame.FooterLines {
		footer.close()
	}

	frame.closed = true
}

// todo: I think this should be decided by the user via a Close() action, not by the indication of closed lines
// since you can always add another line... you don't know when an empty frame should remain open or not
func (frame *Frame) IsClosed() bool {
	// if frame.HeaderLines != nil {
	// 	if !frame.HeaderLines.closed {
	// 		return false
	// 	}
	// }
	//
	// for _, line := range frame.BodyLines {
	// 	if !line.closed {
	// 		return false
	// 	}
	// }
	//
	// if frame.FooterLines != nil {
	// 	if !frame.FooterLines.closed {
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
	frame.startIdx += rows

	// todo: instead of clearing all frame lines, only clear the ones affected
	frame.clear()

	// bump rows and redraw entire frame
	for _, header := range frame.HeaderLines {
		header.move(rows)
	}
	for _, line := range frame.BodyLines {
		line.move(rows)
	}
	for _, footer := range frame.FooterLines {
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
			getScreen().writeAtRow(frame.trailRows[0], frame.startIdx-len(frame.trailRows)+idx)
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
		getScreen().writeAtRow(message, frame.startIdx-len(frame.trailRows)+idx)
	}
	frame.trailRows = make([]string, 0)

	// paint all stale lines to the screen
	for _, header := range frame.HeaderLines {
		if header.visible && (header.stale || frame.stale) {
			_, err := header.write(header.buffer)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}

	for _, line := range frame.BodyLines {
		if line.visible && (line.stale || frame.stale) {
			_, err := line.write(line.buffer)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}

	for _, footer := range frame.FooterLines {
		if footer.visible && (footer.stale || frame.stale) {
			_, err := footer.write(footer.buffer)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}

	// if frame.IsClosed() {
	// 	setCursorRow(frame.startIdx + frame.Height())
	// }

	return errs
}

