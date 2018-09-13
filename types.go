package jotframe

import (
	"sync"
	"github.com/satori/go.uuid"
)

type Line struct {
	id         uuid.UUID
	buffer     []byte
	row        int
	lock       *sync.Mutex
	closed     bool
	stale      bool
}

type logicalFrame struct  {
	header             *Line
	activeLines        []*Line
	clearRows          []int
	trailRows          []string
	rowPreAdvancements int
	footer             *Line

	frameStartIdx int
	updateFn      func() error
	closed        bool
	stale         bool
}

type TopFrame struct {
	frame *logicalFrame
	lock          *sync.Mutex
}

type BottomFrame struct {
	frame *logicalFrame
	lock          *sync.Mutex
}

type FixedFrame struct {
	frame *logicalFrame
	lock          *sync.Mutex
}