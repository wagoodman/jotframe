package main

import (
	"fmt"
	"io"
	"math/rand"
	"time"

	"github.com/wagoodman/jotframe"
)

func main() {
	rand.Seed(time.Now().Unix())

	renderLine := func(idx int, line *jotframe.Line, frame *jotframe.FixedFrame) {
		minMs := 10
		maxMs := 50

		message := fmt.Sprintf("%s %s INITIALIZED", line, time.Now())
		io.WriteString(line, message)
		for idx := 10; idx > 0; idx-- {
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

	// create 5 lines within a frame
	lines := 50
	frame := jotframe.NewFixedFrame(lines, false, false, true)

	// concurrently write to each line
	for idx := 0; idx < lines; idx++ {
		go renderLine(idx, frame.Lines()[idx], frame)
	}

	// close the frame
	frame.Wait()
	frame.Close()
}
