package frame

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type PositionPolicy int

const (
	FloatFree    PositionPolicy = iota // allowed to go anyway, even off the screen
	FloatForward                       // similar to free, except once it hits the bottom, it does not go off the screen (it makes more realestate)
	FloatTop                           // top fixed
	FloatBottom                        // bottom fixed
)

func (float PositionPolicy) String() string {
	switch float {
	case FloatFree:
		return "FloatFree"
	case FloatForward:
		return "FloatForward"
	case FloatTop:
		return "FloatTop"
	case FloatBottom:
		return "FloatBottom"
	default:
		return fmt.Sprintf("PositionPolicy=%d?", float)
	}
}

type ScreenEventHandler interface {
	onEvent(*ScreenEvent)
}

type Policy interface {
	// reactive actions
	// onClose()
	onResize(adjustment int)
	// onUpdate()
	onTrail()

	// proactive actions
	onInit()
	allowedMotion(rows int) int
	allowTrail() bool
}

type Config struct {
	Lines          int
	startRow       int
	HasHeader      bool
	HasFooter      bool
	TrailOnRemove  bool
	PositionPolicy PositionPolicy
	ManualDraw     bool
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

type Frame struct {
	config Config
	lock   *sync.Mutex

	header          *Line
	activeLines     []*Line
	clearRows       []int
	trailRows       []string
	rowAdvancements int
	footer          *Line

	policy      Policy
	autoDraw    bool
	topRow      int
	closeSignal *sync.WaitGroup
	closed      bool
	stale       bool
}

type floatTopPolicy struct {
	Frame *Frame
}

type floatBottomPolicy struct {
	Frame *Frame
}

type floatFreePolicy struct {
	Frame *Frame
}

type floatForwardPolicy struct {
	Frame *Frame
}
