package main

import (
	"github.com/wagoodman/jotframe/pkg/frame"
	"io"
	"math/rand"
	"time"
)

func renderLine(idx int, line *frame.Line, fr *frame.Frame) {

	for i := 0; i < 50; i++ {
		io.WriteString(line, line.String())
		time.Sleep(time.Duration(50) * time.Millisecond)
	}

	// line.Close()
	fr.Remove(line)
}

func main() {
	rand.Seed(time.Now().Unix())

	totalLines := 20
	config := frame.Config{
		Lines:          0, //totalLines,
		HasHeader:      true,
		HasFooter:      true,
		TrailOnRemove:  true,
		PositionPolicy: frame.PolicyFloatOverflow,
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
		time.Sleep(time.Duration(120) * time.Millisecond)
		line, _ := fr.Append()
		// go renderLine(idx, fr.Lines[idx], fr)
		go renderLine(idx, line, fr)
	}
	time.Sleep(time.Duration(10 * time.Second))

	// close the frame
	fr.Wait()
	fr.Close()
	frame.Close()
}
