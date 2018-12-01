package jotframe

import (
	"github.com/sirupsen/logrus"
	"os"
)

var (
	terminalWidth     int
	terminalHeight    int
	allFrames = make([]*logicalFrame, 0)
)

func registerFrame(frame *logicalFrame) {
	allFrames = append(allFrames, frame)
}

func init() {
	// fetch initial values
	terminalWidth, terminalHeight = getTerminalSize()

	go pollSignals()

	initLogging()
}

func initLogging() {

	logFileObj, err := os.OpenFile("debug.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	Formatter := new(logrus.TextFormatter)
	Formatter.DisableTimestamp = true
	logrus.SetFormatter(Formatter)
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetOutput(logFileObj)

	logrus.Debug("Starting...")
}

func Refresh() error {
	lock := getScreenLock()
	lock.Lock()
	defer lock.Unlock()

	return refresh()
}

func refresh() error {
	for _, frame := range allFrames {
		if !frame.isClosed() {
			frame.clear()
			frame.updateAndDraw()
		}
	}
	return nil
}

func Close() error {
	lock := getScreenLock()
	lock.Lock()
	defer lock.Unlock()

	for _, frame := range allFrames {
		err := frame.close()
		if err != nil {
			return err
		}
	}
	return nil
}
