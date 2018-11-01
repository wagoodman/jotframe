package jotframe

import (
	"context"

	ansi "github.com/k0kubun/go-ansi"
	"golang.org/x/sync/semaphore"
)

// Worker is working
type Worker interface {
	Work(*Line)
}

func WorkQueue(maxConcurrent int64, workerArray []interface{}) {
	frame := NewFixedFrame(0, false, false, true)
	// worker pool
	ctx := context.TODO()
	sem := semaphore.NewWeighted(maxConcurrent)

	for _, item := range workerArray {
		worker, _ := item.(Worker)
		sem.Acquire(ctx, 1)
		line, _ := frame.Append()
		jotFunc := func(userFunc func(line *Line), line *Line) {
			defer sem.Release(1)

			userFunc(line)
			frame.Remove(line)
		}
		go jotFunc(worker.Work, line)
	}

	frame.Wait()

	ansi.CursorShow()
}
