package main

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"strings"
	"time"

	"github.com/k0kubun/go-ansi"
	"github.com/wagoodman/jotframe"
)

// TODO: The worker wrapper object should take the following:
// 1. A slice of objects that implements a Worker interface.
//      these objects will do the actual work to be done...
// 2. An (optional) header function that satisfies a Worker interface.
// 3. An (optional) footer function that satisfies a Worker interface.
// 4. A number of maximum items to work on concurrently.

type Resource struct {
	name           string
	totalSize      int
	downloadedSize int
	// TODO: this should take a pointer to a channel (for the footer)
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

	buffer.WriteString(fmt.Sprintf("Downloading %-20s [%s%s] %3.2f%% (%d/%d)", item.name+"...", bars, spaces, percent*100, item.downloadedSize, item.totalSize))

	return buffer.String()
}

// TODO: use this to actually update the footer
func footerProcessor(line *jotframe.Line) {

	// This function reads from a channel... items come in when worker tasks have been completed
	// We then update the footer
}

func (item Resource) Work(line *jotframe.Line) {
	minMs := 10
	maxMs := 100

	message := fmt.Sprintf("Download %s pending...", item.name)
	io.WriteString(line, message)
	for {
		// sleep for a bit...
		randomInterval := rand.Intn(maxMs-minMs) + minMs
		time.Sleep(time.Duration(randomInterval) * time.Millisecond)
		item.downloadedSize += randomInterval

		// write a message to this line...
		io.WriteString(line, item.String())
		if item.downloadedSize >= item.totalSize {
			break
		}

	}
	// write a final message
	message = fmt.Sprintf("Downloaded %s", item.name)
	io.WriteString(line, message)

	// TODO: when this task completes, we should write out to the footer channel
}

func main() {
	const maxConcurrent = 4
	ansi.CursorHide()
	rand.Seed(time.Now().Unix())

	downloads := []Resource{
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

	wq := jotframe.NewWorkQueue(maxConcurrent)
	for _, item := range downloads {
		wq.AddWork(item)
	}
	wq.Work()

}
