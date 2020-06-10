package frame

type policyFloatBottom struct {
	Frame       *Frame
	largestSize int
	currentSize int
}

func newFloatBottomPolicy(frame *Frame) *policyFloatBottom {
	return &policyFloatBottom{
		Frame: frame,
	}
}

// proactive action!
// note: most frame objects don't exist, make changes based on the frame config
func (policy *policyFloatBottom) onInit() {
	height := policy.Frame.Config.Lines
	height += policy.Frame.Config.HeaderRows
	height += policy.Frame.Config.FooterRows

	// the screen index starts at 1 (not 0), hence the +1
	offset := (terminalHeight - height) + 1
	policy.Frame.startIdx = offset
	policy.Frame.Config.startRow = offset
	policy.largestSize = height
	policy.currentSize = height
}

// reactive action!
func (policy *policyFloatBottom) onResize(adjustment int) {
	policy.currentSize += adjustment
	if adjustment > 0 {
		// Grow in size:
		// appended rows should appear to move upwards on the screen, which means that we should
		// move the entire frame upwards 1 line while making more screen space by 1 line
		policy.Frame.move(-adjustment)
		if policy.currentSize > policy.largestSize {
			policy.Frame.rowAdvancements += adjustment
		}
	} else if adjustment < 0 {
		// Shrink in size:
		policy.Frame.move(-adjustment)
	}
	if policy.currentSize > policy.largestSize {
		policy.largestSize = policy.currentSize
	}
}

// reactive policy!
func (policy *policyFloatBottom) onTrail() {
	// write the removed line to the trail log + move the policy down (while advancing the frame)
	policy.Frame.rowAdvancements += 1
}

// proactive action!
func (policy *policyFloatBottom) allowedMotion(rows int) int {
	return 0
}

func (policy *policyFloatBottom) isAllowedTrail() bool {
	return true
}
