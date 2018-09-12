package jotframe

var (
	terminalWidth  int
	terminalHeight int
	allFrames      = make([]*FixedFrame, 0)
)

func registerFrame(frame *FixedFrame) {
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
			frame.update()
			frame.clear(true)
			frame.draw()
		}
	}
	return nil
}

func Close() error {
	// each FixedFrame.Close() will implicitly lock
	for _, frame := range allFrames {
		err := frame.Close()
		if err != nil {
			return err
		}
	}
	return nil
}