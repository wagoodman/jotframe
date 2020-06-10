package frame

type floatForwardPolicy struct {
	Frame *Frame
}

func newFloatForwardPolicy(frame *Frame) *floatForwardPolicy {
	return &floatForwardPolicy{
		Frame: frame,
	}
}

// proactive action!
// note: most frame objects don't exist, make changes based on the frame config
func (policy *floatForwardPolicy) onInit() {
	if policy.Frame.Config.startRow == 0 {
		offset, err := GetCursorRow()
		if err != nil {
			return
		}
		policy.Frame.Config.startRow = offset
		policy.Frame.startIdx = offset
	}

	// we may be starting near the bottom of the screen, make room on the screen if necessary
	adjustment := policy.Frame.Config.Height() - policy.Frame.Config.VisibleHeight()

	// lets not pass the top of the screen by default
	if policy.Frame.startIdx-adjustment-1 < 0 {
		adjustment += policy.Frame.startIdx - adjustment - 1
	}

	policy.Frame.startIdx -= adjustment
	policy.Frame.rowAdvancements += adjustment
}

// reactive action!
func (policy *floatForwardPolicy) onResize(adjustment int) {
	if policy.Frame.IsPastScreenBottom() {
		policy.Frame.move(-adjustment)
		policy.Frame.rowAdvancements += adjustment
	}
}

func (policy *floatForwardPolicy) onTrail() {
	// write the removed line to the trail log + move the policy down (while advancing the frame)
	if policy.Frame.IsPastScreenBottom() {
		// frame.frame.Move(-1)
		policy.Frame.rowAdvancements += 1
	} else {
		policy.Frame.move(1)
	}
}

// func (policy *floatForwardPolicy) onUpdate() {}

// func (policy *floatForwardPolicy) onClose() {}

func (policy *floatForwardPolicy) allowedMotion(rows int) int {
	return rows
}

func (policy *floatForwardPolicy) isAllowedTrail() bool {
	return true
}
