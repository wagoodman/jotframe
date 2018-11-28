package jotframe

import (
	"context"

	"github.com/k0kubun/go-ansi"
	"golang.org/x/sync/semaphore"
)

// Worker is working
type Worker interface {
	Work(*Line)
}

type WorkQueue struct {
	maxConcurrent int64
	queue         []interface{}
}

func NewWorkQueue(maxConcurrent int64) *WorkQueue {
	return &WorkQueue{
		maxConcurrent: maxConcurrent,
	}
}

func (wq *WorkQueue) AddWork(work interface{}) {
	wq.queue = append(wq.queue, work)
}

func (wq *WorkQueue) Work() {
	frame := NewFixedFrame(0, false, false, true)
	// worker pool
	ctx := context.TODO()
	sem := semaphore.NewWeighted(wq.maxConcurrent)

	for _, item := range wq.queue {
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
