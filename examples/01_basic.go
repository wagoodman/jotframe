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
		// write a message to this line...
		message := fmt.Sprintf("%s LineIdx:%d", line, idx)
		io.WriteString(line, message)

		line.Close()
	}

	// create 5 lines within a frame
	lines := 15
	frame := jotframe.NewFixedFrame(lines, false, false, true)

	// concurrently write to each line
	for idx := 0; idx < lines; idx++ {
		renderLine(idx, frame.Lines()[idx], frame)
	}

	// close the frame
	frame.Wait()
	frame.Close()
}
