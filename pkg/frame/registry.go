package frame

var (
	allFrames      = make([]*Frame, 0)
	screenHandlers = make([]ScreenEventHandler, 0)
)

func registerFrame(frame *Frame) {
	allFrames = append(allFrames, frame)
}

func addScreenHandler(handler ScreenEventHandler) {
	screenHandlers = append(screenHandlers, handler)
}

func Refresh() error {
	lock := getScreenLock()
	lock.Lock()
	defer lock.Unlock()

	return refresh()
}

func refresh() error {
	for _, frame := range allFrames {
		if !frame.IsClosed() {
			frame.clear()
			// frame.update()
			frame.draw()
		}
	}
	return nil
}

func Close() error {
	lock := getScreenLock()
	lock.Lock()
	defer lock.Unlock()

	for _, frame := range allFrames {
		frame.close()
	}
	// allow the frames to exist as a trail now. advance the screen to allow room for the cursor.
	row, _ := GetCursorRow()
	if row == terminalHeight {
		advanceScreen(1)
	}

	return nil
}
