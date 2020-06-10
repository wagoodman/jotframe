package frame

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/k0kubun/go-ansi"
)

var (
	screenSync sync.Once
	theScr     *screen
)

type screen struct {
	lock     *sync.RWMutex
	events   chan ScreenEvent
	frames   []*Frame
	handlers []EventHandler
	closed   bool
	workers  *sync.WaitGroup
}

func getScreen() *screen {
	screenSync.Do(func() {
		theScr = &screen{
			lock: &sync.RWMutex{},
		}
		theScr.reset()
	})
	return theScr
}

func (scr *screen) reset() {
	scr.lock.Lock()
	defer scr.lock.Unlock()

	theScr.events = make(chan ScreenEvent, 10000000)
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

func (scr *screen) Close() error {
	scr.lock.Lock()
	defer scr.lock.Unlock()

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
	scr.lock.RLock()
	defer scr.lock.RUnlock()

	if !scr.closed {
		scr.events <- ScreenEvent{
			row:   terminalHeight,
			value: []byte(fmt.Sprint(strings.Repeat(lineBreak, rows))),
		}
	}
}

func (scr *screen) writeAtRow(message string, row int) {
	scr.lock.RLock()
	defer scr.lock.RUnlock()

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
			ansi.EraseInLine(2)
			ansi.CursorHorizontalAbsolute(0)

			// write the output
			_, err = os.Stdout.Write(event.value)
			if err != nil {
				fmt.Printf("failed to write payload: %s\n", err)
				scr.closed = true
				return
			}
		}
	}()
}
