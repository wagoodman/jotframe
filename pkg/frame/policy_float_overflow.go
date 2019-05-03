package frame

type policyFloatOverflow simplePolicy

func newFloatOverflowPolicy(frame *Frame) *policyFloatOverflow {
	return &policyFloatOverflow{
		Frame: frame,
	}
}

// proactive action!
// note: most frame objects don't exist, make changes based on the frame Config
func (policy *policyFloatOverflow) onInit() {
	if policy.Frame.Config.startRow == 0 {
		offset, err := GetCursorRow()
		if err != nil {
			return
		}
		policy.Frame.Config.startRow = offset
		policy.Frame.StartIdx = offset
	}
}

// reactive action!
func (policy *policyFloatOverflow) onResize(adjustment int) {}

func (policy *policyFloatOverflow) onTrail() {
	// write the removed line to the trail log + move the policy down (while advancing the frame)
	if policy.Frame.IsPastScreenBottom() {
		// frame.frame.Move(-1)
		policy.Frame.advance(1)
	} else {
		policy.Frame.move(1)
	}
}

func (policy *policyFloatOverflow) isAllowedMotion(rows int) int {
	return rows
}

func (policy *policyFloatOverflow) isAllowedTrail() bool {
	return true
}
