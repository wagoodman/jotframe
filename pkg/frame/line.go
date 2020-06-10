package frame

import (
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/google/uuid"
)

type Line struct {
	id      uuid.UUID
	buffer  []byte
	height  int
	frame   *Frame
	row     int
	lock    *sync.RWMutex
	visible bool
	closed  bool
	stale   bool
	events  chan ScreenEvent
}

func NewLine(row int, events chan ScreenEvent) *Line {
	return &Line{
		id:      uuid.New(),
		row:     row,
		lock:    getScreen().lock,
		stale:   true,
		events:  events,
		height:  1,
		visible: true,
	}
}

// todo: the line is blocking on write for all handlers, this should not be the case
func (line *Line) notify() error {
	scr := getScreen()

	event := newScreenEvent(line)
	line.events <- *event
	if len(scr.handlers) == 0 {
		return nil
	}
	for _, handler := range scr.handlers {
		handler.onEvent(event)
	}
	return nil
}

func (line *Line) Id() uuid.UUID {
	return line.id
}

func (line *Line) Row() int {
	return line.row
}

func (line *Line) Remove() error {
	line.lock.Lock()
	defer line.lock.Unlock()

	if line.frame != nil {
		err := line.frame.remove(line, false)
		if err != nil {
			return err
		}
	}
	return nil
}

func (line *Line) Hide() error {
	line.lock.Lock()
	defer line.lock.Unlock()

	line.visible = false
	line.stale = true
	line.height = 0

	if line.frame != nil {
		err := line.frame.remove(line, true)
		if err != nil {
			return err
		}
	}
	return nil
}

func (line *Line) Show() error {
	line.lock.Lock()
	defer line.lock.Unlock()

	line.visible = true
	line.stale = true
	line.height = 1

	if line.frame != nil {
		_, idx := line.frame.indexOf(line)
		_, err := line.frame.insert(idx, true)
		if err != nil {
			return err
		}
	}
	return nil
}

func (line *Line) IsClosed() bool {
	return line.closed
}

func (line Line) String() string {
	return fmt.Sprintf("<Line row:%d buff:%d id:%v>", line.row, len(line.buffer), line.id)
}

func (line *Line) move(rows int) {
	line.row += rows
	line.stale = true
}

func (line *Line) Clear() error {
	if line.closed {
		return fmt.Errorf("line is closed")
	}

	line.lock.Lock()
	defer line.lock.Unlock()

	return line.clear(false)
}

func (line *Line) clear(preserveBuffer bool) error {
	if !preserveBuffer {
		line.buffer = []byte("")
	}

	return line.notify()
}

func (line *Line) Read(buff []byte) (int, error) {
	line.lock.Lock()
	defer line.lock.Unlock()
	return line.read(buff)
}

func (line *Line) read(buff []byte) (int, error) {
	buff = line.buffer
	return len(buff), nil
}

func (line *Line) WriteString(str string) error {
	// WriteString already uses Draw() which will implicitly lock
	_, err := io.WriteString(line, str)

	return err
}

func (line *Line) Write(buff []byte) (int, error) {
	line.lock.Lock()
	defer line.lock.Unlock()

	return line.write(buff)
}

func (line *Line) write(buff []byte) (int, error) {
	if line.closed {
		return -1, fmt.Errorf("line is closed")
	}

	line.buffer = []byte(strings.Split(string(buff), lineBreak)[0])

	if line.row < 0 || line.row > terminalHeight {
		return -1, fmt.Errorf("line is out of bounds (row=%d)", line.row)
	}

	return len(line.buffer), line.notify()
}

func (line *Line) WriteStringAndClose(str string) (int, error) {
	// WriteString already uses Draw() which will implicitly lock
	numBytes, err := io.WriteString(line, str)
	if err != nil {
		return -1, err
	}
	return numBytes, line.close()
}

func (line *Line) WriteAndClose(buff []byte) (int, error) {
	line.lock.Lock()
	defer line.lock.Unlock()

	numBytes, err := line.write(buff)
	if err != nil {
		return -1, err
	}

	return numBytes, line.close()
}

func (line *Line) ClearAndClose() error {
	line.lock.Lock()
	defer line.lock.Unlock()

	err := line.clear(false)
	if err != nil {
		return err
	}
	return line.close()
}

func (line *Line) Open() error {
	line.lock.Lock()
	defer line.lock.Unlock()

	return line.open()
}

func (line *Line) open() error {
	line.closed = false

	return nil
}

func (line *Line) Close() error {
	line.lock.Lock()
	defer line.lock.Unlock()

	return line.close()
}

func (line *Line) close() error {
	line.closed = true

	return nil
}
