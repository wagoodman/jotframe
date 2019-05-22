package frame

import (
	"context"
	"fmt"
	"github.com/k0kubun/go-ansi"
	"os"
	"strings"
	"sync"
)

var (
	screenSync sync.Once
	theScr     *screen
)

type screen struct {
	lock   *sync.Mutex
	events chan ScreenEvent
	frames  []*Frame
	handlers []EventHandler
}

func getScreen() *screen {
	screenSync.Do(func() {
		theScr = &screen{
			lock: &sync.Mutex{},
			events: make(chan ScreenEvent, 10000000),
			frames: make([]*Frame, 0),
			handlers: make([]EventHandler, 0),
		}
	})
	return theScr
}

func (scr *screen) reset() {
	theScr = &screen{
		lock: &sync.Mutex{},
		events: make(chan ScreenEvent, 10000000),
		frames: make([]*Frame, 0),
		handlers: make([]EventHandler, 0),
	}
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

// func Refresh() error {
// 	lock := getScreenLock()
// 	lock.Lock()
// 	defer lock.Unlock()
//
// 	return refresh()
// }

func (scr *screen) refresh() error {
	for _, frame := range scr.frames {
		if !frame.IsClosed() {
			frame.clear()
			// frame.update()
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
		scr.advance(1)
	}

	return nil
}



// TODO: this should be a ScreenEvent!
func (scr *screen) advance(rows int) {
	scr.events <- ScreenEvent{
		row: terminalHeight,
		value: []byte(fmt.Sprint(strings.Repeat(lineBreak, rows))),
	}
}

func (scr *screen) writeAtRow(message string, row int) {
	scr.events <- ScreenEvent{
		row: row,
		value: []byte(strings.Replace(message, lineBreak, "", -1)),
	}
}

// // todo: assumes VT100, not cross platform
// func clearScreen() {
// 	fmt.Print("\x1b[2J")
// }



// TODO: this should be written as a frame handler
func (scr *screen) Run(ctx context.Context) error {
	// f, err := os.OpenFile("event.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	// if err != nil {
	// 	log.Fatalf("error opening file: %v", err)
	// }
	// defer f.Close()
	//
	// log.SetOutput(f)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case event := <-scr.events:
			// log.Println(fmt.Sprintf("%d '%s'", event.row, string(event.value)))
			// clear the row
			err := setCursorRow(event.row)
			if err != nil {
				return err
			}
			ansi.EraseInLine(2)
			ansi.CursorHorizontalAbsolute(0)

			// write the output
			_, err = os.Stdout.Write(event.value)
			if err != nil {
				fmt.Println(err)
				return err
			}
		}
	}
	// setCursorRow(frame.startIdx + frame.height())
	return nil
}
