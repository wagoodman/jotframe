package recipe

import (
	"context"

	"github.com/wagoodman/jotframe/pkg/frame"

	"github.com/k0kubun/go-ansi"
	"golang.org/x/sync/semaphore"
)

// Worker is working
type Worker interface {
	Work(*frame.Line)
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
	fr, _ := frame.New(frame.Config{
		Lines:          0,
		HeaderRows:     0,
		FooterRows:     0,
		TrailOnRemove:  true,
		PositionPolicy: frame.PolicyFloatForward,
		ManualDraw:     false,
	})

	// worker pool
	ctx := context.TODO()
	sem := semaphore.NewWeighted(wq.maxConcurrent)

	for _, item := range wq.queue {
		worker, _ := item.(Worker)
		sem.Acquire(ctx, 1)
		line, _ := fr.Append()
		jotFunc := func(userFunc func(line *frame.Line), line *frame.Line) {
			defer sem.Release(1)
			userFunc(line)
			fr.Remove(line)
		}
		go jotFunc(worker.Work, line)
	}

	fr.Close()

	ansi.CursorShow()
}
