package frame

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

var (
	screenSync sync.Once
	theScr     *screen
)

type screen struct {
	lock      *sync.RWMutex
	closeLock *sync.RWMutex
	events    chan ScreenEvent
	frames    []*Frame
	handlers  []EventHandler
	closed    bool
	workers   *sync.WaitGroup
	output    *os.File
}

func getScreen() *screen {
	screenSync.Do(func() {
		theScr = &screen{
			lock:      &sync.RWMutex{},
			closeLock: &sync.RWMutex{},
			output:    os.Stdout,
		}
		theScr.reset()
	})
	return theScr
}

func (scr *screen) setWriter(writer *os.File) {
	scr.output = writer
	// now there is a different fd which to ask for screen dimensions from
	updateScreenDimensions()
}

func (scr *screen) reset() {
	scr.lock.Lock()
	defer scr.lock.Unlock()

	theScr.events = make(chan ScreenEvent, 100000)
	theScr.frames = make([]*Frame, 0)
	theScr.handlers = make([]EventHandler, 0)
	theScr.workers = &sync.WaitGroup{}
}

func (scr *screen) register(frame *Frame) error {
	if len(scr.frames) > 0 {
		return fmt.Errorf("only one frame is allowed on the screen")
	}
	scr.frames = append(scr.frames, frame)
	return nil
}

func (scr *screen) addScreenHandler(handler EventHandler) {
	scr.handlers = append(scr.handlers, handler)
}

func (scr *screen) refresh() error {
	for _, frame := range scr.frames {
		if !frame.IsClosed() {
			frame.clear()
			frame.draw()
		}
	}
	return nil
}

func Close() error {
	return getScreen().Close()
}

func (scr *screen) Close() error {
	scr.closeLock.Lock()
	defer scr.closeLock.Unlock()

	for _, frame := range scr.frames {
		frame.close()
	}
	// allow the frames to exist as a trail now. advance the screen to allow room for the cursor.
	row, _ := GetCursorRow()
	if row == terminalHeight {
		scr.events <- ScreenEvent{
			row:   terminalHeight,
			value: []byte(lineBreak),
		}
	}

	scr.closed = true
	close(scr.events)

	scr.workers.Wait()

	return nil
}

func (scr *screen) advance(rows int) {
	scr.closeLock.RLock()
	defer scr.closeLock.RUnlock()

	if !scr.closed {
		scr.events <- ScreenEvent{
			row:   terminalHeight,
			value: []byte(fmt.Sprint(strings.Repeat(lineBreak, rows))),
		}
	}
}

func (scr *screen) writeAtRow(message string, row int) {
	scr.closeLock.RLock()
	defer scr.closeLock.RUnlock()

	if !scr.closed {
		scr.events <- ScreenEvent{
			row:   row,
			value: []byte(strings.Replace(message, lineBreak, "", -1)),
		}
	}
}

// TODO: this should be written as a frame handler
func (scr *screen) Run() {
	scr.workers.Add(1)

	go func() {
		defer scr.workers.Done()

		for event := range scr.events {
			// clear the row
			err := setCursorRow(event.row)
			if err != nil {
				fmt.Printf("failed to set cursor row: %s\n", err)
				scr.closed = true
				return
			}
			// erase line (set mode=2)
			_, err = fmt.Fprintf(scr.output, "\x1b[%dK", 2)
			if err != nil {
				fmt.Printf("failed to erase line: %s\n", err)
				scr.closed = true
				return
			}
			// set cursor horizontal absolute position to 0
			_, err = fmt.Fprintf(scr.output, "\x1b[%dG", 0)
			if err != nil {
				fmt.Printf("failed to set horizontal position: %s\n", err)
				scr.closed = true
				return
			}

			// write the output
			_, err = scr.output.Write(event.value)
			if err != nil {
				fmt.Printf("failed to write payload: %s\n", err)
				scr.closed = true
				return
			}
		}
	}()
}
