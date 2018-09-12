package main

import (
	"math/rand"
	"github.com/wagoodman/jotframe"
	"time"
	"fmt"
	"io"
	"github.com/k0kubun/go-ansi"
)

func randomInt(min, max int) int {
	// this is a good test for race conditions ;)
	// rand.Seed(time.Now().Unix())
	return rand.Intn(max - min) + min
}

func randomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(65 + rand.Intn(25))  //A=65 and Z = 65+25
	}
	return string(bytes)
}

func renderLine(idx int, line *jotframe.Line, frame *jotframe.BottomFrame) {
	for {
		randomTime := randomInt(100, 500)
		time.Sleep(time.Duration(randomTime) * time.Millisecond)
		// message := fmt.Sprintf("%d: %+v --> Message:'%v'",  idx, time.Now(), randomString(randomInt(5, 50)))
		message := fmt.Sprintf("%s %s %d", line, time.Now(), idx)
		_, err := io.WriteString(line, message)
		if err != nil {
			break
		}
	}
}

// func renderBar(bar *jotframe.FixedFrame, frame *jotframe.FixedFrame) {
func renderBar(bar *jotframe.FixedFrame, frame *jotframe.BottomFrame) {
	for {
		time.Sleep(time.Duration(200) * time.Millisecond)
		// message := fmt.Sprintf("%d: %+v --> Message:'%v'",  idx, time.Now(), randomString(randomInt(5, 50)))
		// message := fmt.Sprintf("Width: %d  Height: %d", jotframe.terminalWidth, jotframe.terminalHeight)
		w, h := jotframe.GetTerminalSize()
		message := fmt.Sprintf("Width: %d  Height: %d", w, h)
		_, err := io.WriteString(bar.Lines()[0], message)
		if err != nil {
			break
		}
	}
}


func main() {
	ansi.CursorHide()
	lines := 5
	frame, err := jotframe.NewBottomFrame(lines, true, true)
	if err != nil {
		panic(err)
	}

	bar, err := jotframe.NewFixedFrameAt(1, false, false, 50)
	if err != nil {
		panic(err)
	}
	go renderBar(bar, frame)

	rand.Seed(time.Now().Unix())

	frame.Header().WriteString("header!")
	frame.Footer().WriteString("footer!")
	for idx := 0; idx < lines; idx++ {
		go renderLine(idx, frame.Lines()[idx], frame)
	}

	for idx := 0; idx < lines; idx++ {
		time.Sleep(time.Duration(1000) * time.Millisecond)
		line, err := frame.Append()
		if err != nil {
			panic(err)
		}
		go renderLine(lines + idx, line, frame)
	}


	for idx := len(frame.Lines())-1; idx > 0; idx-- {
		time.Sleep(time.Duration(1000) * time.Millisecond)
		frame.Lines()[idx].ClearAndClose()
		// frame.Remove(frame.lines[idx])
		// frame.Advance(2)
	}

	frame.ClearAndClose()
	bar.ClearAndClose()

}




//
// func main() {
// 	lines := 5
// 	frame, err := jotframe.NewFixedFrame(lines, true, true)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	bar, err := jotframe.NewFixedFrameAt(1, false, false, 50)
// 	if err != nil {
// 		panic(err)
// 	}
// 	go renderBar(bar, frame)
//
// 	rand.Seed(time.Now().Unix())
//
// 	frame.header.WriteString("header!")
// 	frame.footer.WriteString("footer!")
// 	for idx := 0; idx < lines; idx++ {
// 		go renderLine(idx, frame.lines[idx], frame)
// 	}
//
// 	for idx := 0; idx < lines; idx++ {
// 		time.Sleep(time.Duration(500) * time.Millisecond)
// 		line, err := frame.Append()
// 		if err != nil {
// 			panic(err)
// 		}
// 		go renderLine(lines + idx, line, frame)
// 	}
//
//
// 	for idx := len(frame.lines)-1; idx > 0; idx-- {
// 		time.Sleep(time.Duration(500) * time.Millisecond)
// 		frame.lines[idx].ClearAndClose()
// 		// frame.Remove(frame.lines[idx])
// 		// frame.Advance(2)
// 	}
//
// 	frame.ClearAndClose()
// 	bar.ClearAndClose()
//
// }
