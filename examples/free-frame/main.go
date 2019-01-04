package main

import (
	"fmt"
	"github.com/wagoodman/jotframe/pkg/frame"
	"io"
	"math/rand"
	"time"
)

func renderLine(idx int, line *frame.Line) {
	// write a message to this line...
	message := fmt.Sprintf("%s --------------- LineIdx:%d", line, idx)
	io.WriteString(line, message)

	line.Close()
}

func main() {
	rand.Seed(time.Now().Unix())

	config := frame.Config{
		Lines: 15,
		Float: frame.FloatFree,
	}

	frames := frame.Factory(config)
	fr := frames[0]

	// concurrently write to each line
	for idx := 0; idx < config.Lines; idx++ {
		renderLine(idx, fr.Lines()[idx])
	}

	// close the frame
	fr.Wait()
	fr.Close()
}
