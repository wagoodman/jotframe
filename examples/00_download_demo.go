package main

import (
	"github.com/wagoodman/jotframe"
	"time"
	"math/rand"
	"fmt"
	"io"
	"bytes"
	"strings"
	"github.com/k0kubun/go-ansi"
	"golang.org/x/sync/semaphore"
	"golang.org/x/net/context"
)

type Resource struct {
	name           string
	totalSize      int
	downloadedSize int
}

func (item Resource) String() string {
	var buffer bytes.Buffer
	const maxBars int = 100

	if item.downloadedSize > item.totalSize {
		item.downloadedSize = item.totalSize
	}

	percent := float32(item.downloadedSize) / float32(item.totalSize)
	if percent > 100 {
		percent = 100.0
	} else if percent < 0 {
		percent = 0.0
	}

	numBars := int(float32(maxBars) * percent)
	numSpaces := maxBars - numBars

	spaces := strings.Repeat(" ", numSpaces)
	bars := strings.Repeat("=", numBars)

	buffer.WriteString(fmt.Sprintf("Downloading %-20s [%s%s] %3.2f%% (%d/%d)", item.name+"...",bars,spaces,percent*100, item.downloadedSize, item.totalSize))

	return buffer.String()
}

func (item Resource) download(line *jotframe.Line) {
	minMs := 10
	maxMs := 100

	message := fmt.Sprintf("Download %s pending...", item.name)
	io.WriteString(line, message)
	for {
		// sleep for a bit...
		randomInterval := rand.Intn(maxMs - minMs) + minMs
		time.Sleep(time.Duration(randomInterval) * time.Millisecond)
		item.downloadedSize += randomInterval

		// write a message to this line...
		io.WriteString(line, item.String())
		if item.downloadedSize >= item.totalSize {
			break
		}

	}
	// write a final message
	message = fmt.Sprintf("Download %s", item.name)
	io.WriteString(line, message)
}


func main() {
	const maxConcurrent = 4
	ansi.CursorHide()
	rand.Seed(time.Now().Unix())

	downloads := []*Resource{
		{"fedora", 18691, 0},
		{"consul", 1580, 0},
		{"bashful", 1720, 0},
		{"mindcraft", 1944, 0},
		{"google-chrome", 2055, 0},
		{"docker", 1699, 0},
		{"docker-compose", 1204, 0},
		{"ubuntu", 15032, 0},
		{"ncdu", 1944, 0},
		{"bridgy", 1403, 0},
		{"kubectl", 2305, 0},
		{"oc", 1898, 0},
		{"cloc", 1733, 0},
		{"ripgrep", 1005, 0},
		{"firefox", 4373, 0},
		{"counter-strike", 4202, 0},
	}
	totalItems := len(downloads)



	/*
	Maybe a frame should have a listener? so you can subscribe to events and act on them (for header and footer updates)

	Allow for a built in worker pool... but its really a frame representing a partial queue?
	 */



	frame := jotframe.NewFixedFrame(0, false, true, true)

	frame.Footer().WriteString(fmt.Sprintf("Completed Downloads: [0/%d]", totalItems))

	// worker pool
	ctx := context.TODO()
	sem := semaphore.NewWeighted(maxConcurrent)

	completedItems := 0
	for _, item := range downloads {
		sem.Acquire(ctx, 1)
		line, _ := frame.Append()
		go func(fun func(line *jotframe.Line), line *jotframe.Line) {
			defer sem.Release(1)

			fun(line)
			completedItems++
			frame.Remove(line)

			if completedItems == totalItems {
				frame.Footer().WriteString("All Downloads Complete!")
				frame.Footer().Close()
			} else {
				frame.Footer().WriteString(fmt.Sprintf("Downloads: [%d/%d]", completedItems, totalItems))
			}

			if len(frame.Lines()) == 0 {
				frame.Close()
			}

		}(item.download, line)
	}

	frame.Wait()

	ansi.CursorShow()
}

