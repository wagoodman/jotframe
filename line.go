package jotframe

import (
	"github.com/k0kubun/go-ansi"
	"os"
	"fmt"
	"strings"
	"github.com/satori/go.uuid"
	"io"
)

func NewLine(row int) *Line {
	line := &Line{}
	line.id = uuid.Must(uuid.NewV4())
	line.row = row
	line.lock = getScreenLock()
	line.stale = true
	return line
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
	return fmt.Sprintf("<Line id:%s idx:%d closed:%v buffer:%d>", line.id, line.row, line.closed, len(line.buffer))
}

func (line *Line) move(rows int) error {
	line.row += rows
	line.stale = true
	return nil
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

	err := setCursorRow(line.row)
	if err != nil {
		return err
	}
	ansi.EraseInLine(2)

	if !preserveBuffer {
		line.buffer = []byte("")
	}
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

	err := setCursorRow(line.row)
	if err != nil {
		return -1, err
	}
	ansi.EraseInLine(2)
	ansi.CursorHorizontalAbsolute(0)

	line.buffer = []byte(strings.Split(string(buff), "\n")[0])

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


func (line *Line) Close() error {
	line.lock.Lock()
	defer line.lock.Unlock()

	return line.close()
}

func (line *Line) close() error {
	line.closed = true
	return nil
}