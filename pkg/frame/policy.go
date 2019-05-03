package frame

import (
	"fmt"
)

type PositionPolicy int

const (
	PolicyFloatOverflow     PositionPolicy = iota // allowed to go anywhere, even off the screen
	PolicyFloatForwardTrail                       // similar to free, except once it hits the bottom, it does not go off the screen (it makes more realestate). If the frame is too large for the screen, overflow (including headers) occurs at the top of the screen.
	PolicyFloatForwardWindow                      // similar to forward-trail, except once it hits the bottom, it does not go off the screen... instead it will act like a bottom-frame with the header fixed to the top of the screen. (it does NOT make more realestate, but instead buffers the unseen output and flushes it all to the screen at the end.). The header and footer stays on the screen while content is overflowed.
	PolicyFloatTop                                // top fixed
	PolicyFloatBottom                             // bottom fixed
)

func (float PositionPolicy) String() string {
	switch float {
	case PolicyFloatOverflow:
		return "PolicyFloatOverflow"
	case PolicyFloatForwardTrail:
		return "policyFloatForwardTrail"
	case PolicyFloatForwardWindow:
		return "policyFloatForwardWindow"
	case PolicyFloatTop:
		return "PolicyFloatTop"
	case PolicyFloatBottom:
		return "PolicyFloatBottom"
	default:
		return fmt.Sprintf("PositionPolicy=%d?", float)
	}
}

type Policy interface {
	// reactive actions
	// onClose()
	onResize(adjustment int)
	// onUpdate()
	onTrail()

	// proactive actions (admission)
	onInit()
	isAllowedMotion(rows int) int
	isAllowedTrail() bool
}

type simplePolicy struct {
	Frame *Frame
}
