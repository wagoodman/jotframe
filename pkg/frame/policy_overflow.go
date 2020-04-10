package frame

type policyOverflow struct {
	Frame *Frame
}

func newOverflowPolicy(frame *Frame) *policyOverflow {
	return &policyOverflow{
		Frame: frame,
	}
}

// proactive action!
// note: most frame objects don't exist, make changes based on the frame Config
func (policy *policyOverflow) onInit() {
	if policy.Frame.Config.startRow == 0 {
		offset, err := GetCursorRow()
		if err != nil {
			return
		}
		policy.Frame.Config.startRow = offset
		policy.Frame.startIdx = offset
	}
}

// reactive action!
func (policy *policyOverflow) onResize(adjustment int) {}

func (policy *policyOverflow) onTrail() {
	// write the removed line to the trail log + move the policy down (while advancing the frame)
	if policy.Frame.IsPastScreenBottom() {
		policy.Frame.advance(1)
	} else {
		policy.Frame.move(1)
	}
}

func (policy *policyOverflow) allowedMotion(rows int) int {
	return rows
}

func (policy *policyOverflow) isAllowedTrail() bool {
	return true
}
