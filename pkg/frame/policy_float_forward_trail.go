package frame

type policyFloatForwardTrail simplePolicy

func newFloatForwardTrailPolicy(frame *Frame) *policyFloatForwardTrail {
	return &policyFloatForwardTrail{
		Frame: frame,
	}
}

// proactive action!
// note: most frame objects don't exist, make changes based on the frame Config
func (policy *policyFloatForwardTrail) onInit() {
	if policy.Frame.Config.startRow == 0 {
		offset, err := GetCursorRow()
		if err != nil {
			return
		}
		policy.Frame.Config.startRow = offset
		policy.Frame.StartIdx = offset
	}

	// we may be starting near the bottom of the screen, make room on the screen if necessary
	adjustment := policy.Frame.Config.Height() - policy.Frame.Config.VisibleHeight()

	// lets not pass the top of the screen by default
	if policy.Frame.StartIdx-adjustment-1 < 0 {
		adjustment += policy.Frame.StartIdx - adjustment - 1
	}

	policy.Frame.StartIdx -= adjustment
	policy.Frame.advance(adjustment)
}

// reactive action!
func (policy *policyFloatForwardTrail) onResize(adjustment int) {
	if policy.Frame.IsPastScreenBottom() {
		policy.Frame.move(-adjustment)
		policy.Frame.advance(adjustment)
	}
}

func (policy *policyFloatForwardTrail) onTrail() {
	// write the removed line to the trail log + move the policy down (while advancing the frame)
	if policy.Frame.IsPastScreenBottom() {
		// frame.frame.Move(-1)
		policy.Frame.advance(1)
	} else {
		policy.Frame.move(1)
	}
}

func (policy *policyFloatForwardTrail) isAllowedMotion(rows int) int {
	return rows
}

func (policy *policyFloatForwardTrail) isAllowedTrail() bool {
	return true
}
