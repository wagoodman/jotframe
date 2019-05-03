package main

import (
"fmt"
"github.com/wagoodman/jotframe/pkg/frame"
"io"
"math/rand"
"time"
)

func renderLine(idx int, line *frame.Line, fr *frame.Frame) {

	message := fmt.Sprintf("%s %s INITIALIZED", line, time.Now())
	io.WriteString(line, message)

	time.Sleep(time.Duration(100*idx) * time.Millisecond)

	// line.Close()
	fr.Remove(line)
}

func main() {
	rand.Seed(time.Now().Unix())

	totalLines := 10
	config := frame.Config{
		Lines:          totalLines,
		HasHeader:      true,
		HasFooter:      true,
		TrailOnRemove:  true,
		PositionPolicy: frame.PolicyFloatForwardWindow,
		ManualDraw:     false,
	}
	fr := frame.New(config)

	// add a header and footer
	fr.Header.WriteString("This is the best header ever!")
	fr.Header.Close()

	fr.Footer.WriteString("...Followed by the best footer ever...")
	fr.Footer.Close()

	// concurrently write to each line
	for idx := 0; idx < totalLines; idx++ {
		// line, _ := fr.Append()
		line := fr.Lines[idx]
		// go renderLine(idx, fr.Lines[idx], fr)
		go renderLine(idx, line, fr)
	}
	time.Sleep(time.Duration(10 * time.Second))

	// close the frame
	fr.Wait()
	fr.Close()
	frame.Close()
}
