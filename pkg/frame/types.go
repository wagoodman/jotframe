package frame

import (
	"sync"

	"github.com/google/uuid"
)

type FloatRule int

const (
	FloatFree FloatRule = iota
	FloatTop
	FloatBottom
)

type ScreenEventHandler interface {
	onEvent(*ScreenEvent)
}

type Frame interface {
	StartIdx() int
	Config() Config
	Height() int

	Header() *Line
	Footer() *Line
	Lines() []*Line
	Append() (*Line, error)
	Prepend() (*Line, error)
	Insert(index int) (*Line, error)
	Remove(line *Line) error
	Move(rows int)

	Clear() error
	Close() error
	IsClosed() bool

	Update() error
	Draw() []error

	Wait()
}

type Config struct {
	Lines         int
	startRow      int
	HasHeader     bool
	HasFooter     bool
	TrailOnRemove bool
	Float         FloatRule
}

type ScreenEvent struct {
	value []byte
	row   int
}

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
	config Config

	header          *Line
	activeLines     []*Line
	clearRows       []int
	trailRows       []string
	rowAdvancements int
	footer          *Line

	topRow      int
	closeSignal *sync.WaitGroup
	updateFn    func(*logicalFrame) error
	closed      bool
	stale       bool
}

type topFrame struct {
	logicalFrame *logicalFrame
	lock         *sync.Mutex
	config       Config
}

type bottomFrame struct {
	logicalFrame *logicalFrame
	lock         *sync.Mutex
	config       Config
}

type floatingFrame struct {
	logicalFrame *logicalFrame
	lock         *sync.Mutex
	config       Config
}
