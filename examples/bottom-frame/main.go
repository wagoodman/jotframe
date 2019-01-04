package main

import (
	"fmt"
	"github.com/wagoodman/jotframe/pkg/frame"
	"io"
	"math/rand"
	"time"
)

func renderLine(idx int, line *frame.Line) {
	minMs := 10
	maxMs := 50

	message := fmt.Sprintf("%s %s INITIALIZED", line, time.Now())
	io.WriteString(line, message)
	for idx := 100; idx > 0; idx-- {
		// sleep for a bit...
		randomInterval := rand.Intn(maxMs-minMs) + minMs
		time.Sleep(time.Duration(randomInterval) * time.Millisecond)

		// write a message to this line...
		message := fmt.Sprintf("%s CountDown:%d", line, idx)
		io.WriteString(line, message)

	}
	// write a final message
	message = fmt.Sprintf("%s %s", line, "Closed!")
	io.WriteString(line, message)

	line.Close()
}

func main() {
	rand.Seed(time.Now().Unix())

	config := frame.Config{
		Lines:         5,
		Float:         frame.FloatBottom,
		HasHeader:     true,
		HasFooter:     true,
		TrailOnRemove: true,
	}
	frames := frame.Factory(config)
	fr := frames[0]

	// add a header and footer
	fr.Header().WriteString("This is the best header ever!")
	fr.Header().Close()

	fr.Footer().WriteString("...Followed by the best footer ever...")
	fr.Footer().Close()

	// concurrently write to each line
	for idx := 0; idx < config.Lines; idx++ {
		go renderLine(idx, fr.Lines()[idx])
	}

	// close the frame
	fr.Wait()
	fr.Close()
}
