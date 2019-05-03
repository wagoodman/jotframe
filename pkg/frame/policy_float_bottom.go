package frame

type policyFloatBottom simplePolicy

func newFloatBottomPolicy(frame *Frame) *policyFloatBottom {
	return &policyFloatBottom{
		Frame: frame,
	}
}

// proactive action!
// note: most frame objects don't exist, make changes based on the frame Config
func (policy *policyFloatBottom) onInit() {
	height := policy.Frame.Config.Lines
	if policy.Frame.Config.HasHeader {
		height++
	}
	if policy.Frame.Config.HasFooter {
		height++
	}
	// the screen index starts at 1 (not 0), hence the +1
	policy.Frame.StartIdx = (terminalHeight - height) + 1
}

// reactive action!
func (policy *policyFloatBottom) onResize(adjustment int) {
	if adjustment > 0 {
		// Grow in size:
		// appended rows should appear to move upwards on the screen, which means that we should
		// move the entire frame upwards 1 line while making more screen space by 1 line
		policy.Frame.move(-adjustment)
		policy.Frame.advance(adjustment)
	} else if adjustment < 0 {
		// Shrink in size:
		policy.Frame.move(adjustment)
	}

}

// reactive policy!
func (policy *policyFloatBottom) onTrail() {
	// write the removed line to the trail log + move the policy down (while advancing the frame)
	// policy.frame.move(1)
	policy.Frame.advance(1)
}

// proactive action!
func (policy *policyFloatBottom) isAllowedMotion(rows int) int {
	return 0
}

func (policy *policyFloatBottom) isAllowedTrail() bool {
	return true
}
