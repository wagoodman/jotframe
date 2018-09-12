package jotframe

var (
	terminalWidth  int
	terminalHeight int
	allFrames      = make([]*logicalFrame, 0)
)

func registerFrame(frame *logicalFrame) {
	allFrames = append(allFrames, frame)
}

func init() {
	// fetch initial values
	terminalWidth, terminalHeight = getTerminalSize()

	go pollSignals()
}

func Refresh() error {
	lock := getScreenLock()
	lock.Lock()
	defer lock.Unlock()

	return refresh()
}

func refresh() error {
	for _, frame := range allFrames {
		if !frame.closed {
			frame.clear()
			frame.updateAndDraw()
		}
	}
	return nil
}

func Close() error {
	// each logicalFrame.Close() will implicitly lock
	for _, frame := range allFrames {
		err := frame.Close()
		if err != nil {
			return err
		}
	}
	return nil
}