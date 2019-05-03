package frame

type policyFloatForwardWindow simplePolicy

func newFloatForwardWindowPolicy(frame *Frame) *policyFloatForwardWindow {
	return &policyFloatForwardWindow{
		Frame: frame,
	}
}

// proactive action!
// note: most frame objects don't exist, make changes based on the frame Config
func (policy *policyFloatForwardWindow) onInit() {
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
func (policy *policyFloatForwardWindow) onResize(adjustment int) {
	if policy.Frame.IsPastScreenBottom() {
		policy.Frame.move(-adjustment)
		policy.Frame.advance(adjustment)
	}
}

func (policy *policyFloatForwardWindow) onTrail() {
}

func (policy *policyFloatForwardWindow) isAllowedMotion(rows int) int {
	return rows
}

func (policy *policyFloatForwardWindow) isAllowedTrail() bool {
	return false
}
