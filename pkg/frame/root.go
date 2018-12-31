package frame

var (
	terminalWidth  int
	terminalHeight int
	allFrames      = make([]Frame, 0)
	screenHandlers = make([]ScreenEventHandler, 0)
)

func registerFrame(frame Frame) {
	allFrames = append(allFrames, frame)
}

func addScreenHandler(handler ScreenEventHandler) {
	screenHandlers = append(screenHandlers, handler)
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
		if !frame.IsClosed() {
			frame.Clear()
			frame.Update()
			frame.Draw()
		}
	}
	return nil
}

func Close() error {
	lock := getScreenLock()
	lock.Lock()
	defer lock.Unlock()

	for _, frame := range allFrames {
		err := frame.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
