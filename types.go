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
	header      *Line
	activeLines []*Line
	clearLines  []*Line
	footer      *Line

	lock          *sync.Mutex
	frameStartIdx int
	updateFn      func() error
	closed        bool
	stale         bool
}

type TopFrame struct {
	frame *logicalFrame
}

type BottomFrame struct {
	frame *logicalFrame
}

type FixedFrame struct {
	frame *logicalFrame
}