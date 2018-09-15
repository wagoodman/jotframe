package main

import (
"github.com/wagoodman/jotframe"
"time"
"math/rand"
"fmt"
"io"
)

func main() {
	rand.Seed(time.Now().Unix())

	renderLine := func(idx int, line *jotframe.Line, frame *jotframe.FixedFrame) {
		minMs := 10
		maxMs := 50

		message := fmt.Sprintf("%s %s INITIALIZED", line, time.Now())
		io.WriteString(line, message)
		for idx := 100 ; idx > 0 ; idx-- {
			// sleep for a bit...
			randomInterval := rand.Intn(maxMs - minMs) + minMs
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
	lines := 5
	frame := jotframe.NewFixedFrame(lines, true, true, true)

	// add a header and footer
	frame.Header().WriteString("This is the best header ever!")
	frame.Header().Close()

	frame.Footer().WriteString("...Followed by the best footer ever...")
	frame.Footer().Close()

	// concurrently write to each line
	for idx := 0; idx < lines; idx++ {
		go renderLine(idx, frame.Lines()[idx], frame)
	}

	// close the frame
	frame.Wait()
	frame.Close()
}

