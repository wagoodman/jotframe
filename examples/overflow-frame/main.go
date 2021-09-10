package main

import (
	"io"
	"math/rand"
	"os"
	"time"

	"github.com/wagoodman/jotframe/pkg/frame"
)

func renderLine(idx int, line *frame.Line, fr *frame.Frame) {

	io.WriteString(line, line.String())
	time.Sleep(time.Duration((idx+1)*100)*time.Millisecond + 1*time.Second)

	line.Hide()

	// line.Close()

	time.Sleep(time.Duration((idx+1)*100)*time.Millisecond + 1*time.Second)

	// fr.Remove(line)

	line.Show()

}

func main() {
	rand.Seed(time.Now().Unix())

	totalLines := 10
	config := frame.Config{
		Lines:          0,
		HeaderRows:     1,
		FooterRows:     1,
		TrailOnRemove:  false,
		PositionPolicy: frame.PolicyFloatTop,
		ManualDraw:     false,
		Output:         os.Stderr,
	}
	fr, err := frame.New(config)
	if err != nil {
		panic(err)
	}

	// add a header and footer
	fr.HeaderLines[0].WriteString("This is the best header ever!")
	// fr.HeaderLines.Close()

	fr.FooterLines[0].WriteString("...Followed by the best footer ever...")
	// fr.FooterLines.Close()

	// concurrently write to each line
	time.Sleep(200 * time.Millisecond)
	for idx := 0; idx < totalLines; idx++ {
		time.Sleep(time.Duration(90) * time.Millisecond)
		line, _ := fr.Append()
		// go renderLine(idx, fr.BodyLines[idx], fr)
		go renderLine(idx, line, fr)
	}

	// time.Sleep(time.Duration(700 * time.Millisecond))
	// header2, err := fr.AppendHeader()
	// if err != nil{
	// 	panic(err)
	// }
	// header2.WriteString("append header")
	// time.Sleep(time.Duration(700 * time.Millisecond))
	// header3, err := fr.PrependHeader()
	// if err != nil{
	// 	panic(err)
	// }
	// header3.WriteString("prepend header")
	//
	// time.Sleep(time.Duration(700 * time.Millisecond))
	//
	//
	// footer2, err := fr.AppendFooter()
	// if err != nil{
	// 	panic(err)
	// }
	// footer2.WriteString("append footer")
	//
	// time.Sleep(time.Duration(700 * time.Millisecond))
	//
	// footer3, err := fr.PrependFooter()
	// if err != nil{
	// 	panic(err)
	// }
	// footer3.WriteString("prepend footer")

	time.Sleep(time.Duration(5 * time.Second))

	// close the frame
	fr.Close()
}
