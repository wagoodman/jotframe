package jotframe

import (
	"sync"

	"github.com/google/uuid"
)

type Line struct {
	id          uuid.UUID
	buffer      []byte
	row         int
	lock        *sync.Mutex
	closeSignal *sync.WaitGroup
	closed      bool
	stale       bool
}

type logicalFrame struct {
	header          *Line
	activeLines     []*Line
	clearRows       []int
	trailRows       []string
	rowAdvancements int
	footer          *Line

	frameStartIdx int
	closeSignal   *sync.WaitGroup
	updateFn      func() error
	closed        bool
	stale         bool
}

type TopFrame struct {
	logicalFrame *logicalFrame
	lock         *sync.Mutex
}

type BottomFrame struct {
	logicalFrame  *logicalFrame
	lock          *sync.Mutex
	trailOnRemove bool
}

type FixedFrame struct {
	logicalFrame  *logicalFrame
	lock          *sync.Mutex
	trailOnRemove bool
}
