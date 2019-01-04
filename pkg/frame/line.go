package frame

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/k0kubun/go-ansi"
)

func NewLine(row int, closeSignal *sync.WaitGroup) *Line {
	if closeSignal != nil {
		closeSignal.Add(1)
	}
	return &Line{
		id:          uuid.New(),
		row:         row,
		lock:        getScreenLock(),
		stale:       true,
		closeSignal: closeSignal,
	}
}

func (line *Line) notify() {
	if len(screenHandlers) == 0 {
		return
	}
	event := newScreenEvent(line)
	for _, handler := range screenHandlers {
		handler.onEvent(event)
	}
}

func (line *Line) Id() uuid.UUID {
	return line.id
}

func (line *Line) Row() int {
	return line.row
}

func (line *Line) IsClosed() bool {
	return line.closed
}

func (line Line) String() string {
	return fmt.Sprintf("<Line row:%d closed:%v bufferLen:%d>", line.row, line.closed, len(line.buffer))
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

	err := clearRow(line.row)
	if err != nil {
		return err
	}

	if !preserveBuffer {
		line.buffer = []byte("")
	}
	line.notify()
	return nil
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
	if line.closed {
		return -1, fmt.Errorf("line is closed")
	}
	line.lock.Lock()
	defer line.lock.Unlock()

	return line.write(buff)
}

func (line *Line) write(buff []byte) (int, error) {
	line.buffer = []byte(strings.Split(string(buff), "\n")[0])

	if line.row < 0 || line.row > terminalHeight {
		return -1, fmt.Errorf("line is out of bounds (row=%d)", line.row)
	}

	err := clearRow(line.row)
	if err != nil {
		return -1, err
	}
	ansi.CursorHorizontalAbsolute(0)

	line.notify()
	return os.Stdout.Write(line.buffer)
}

func (line *Line) WriteStringAndClose(str string) (int, error) {
	// WriteString already uses Draw() which will implicitly lock
	defer line.close()
	return io.WriteString(line, str)
}

func (line *Line) WriteAndClose(buff []byte) (int, error) {
	line.lock.Lock()
	defer line.lock.Unlock()

	defer line.close()
	return line.write(buff)
}

func (line *Line) ClearAndClose() error {
	line.lock.Lock()
	defer line.lock.Unlock()

	defer line.close()
	return line.clear(false)
}

func (line *Line) Open() error {
	line.lock.Lock()
	defer line.lock.Unlock()

	return line.open()
}

func (line *Line) open() error {
	if line.closed {
		line.closed = false
		if line.closeSignal != nil {
			line.closeSignal.Add(1)
		}
	}

	return nil
}

func (line *Line) Close() error {
	line.lock.Lock()
	defer line.lock.Unlock()

	return line.close()
}

func (line *Line) close() error {
	if !line.closed {
		line.closed = true
		if line.closeSignal != nil {
			line.closeSignal.Done()
		}
	}

	return nil
}
