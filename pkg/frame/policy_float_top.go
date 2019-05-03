package frame

type policyFloatTop simplePolicy

func newFloatTopPolicy(frame *Frame) *policyFloatTop {
	return &policyFloatTop{
		Frame: frame,
	}
}

// proactive action!
// note: most frame objects don't exist, make changes based on the frame Config
func (policy *policyFloatTop) onInit() {
	policy.Frame.StartIdx = 1
	offset, err := GetCursorRow()
	if err != nil {
		return
	}
	policy.Frame.resetAdvancements()
	policy.Frame.advance(offset - 1)
}

// reactive action!
func (policy *policyFloatTop) onResize(adjustment int) {}

// reactive policy!
func (policy *policyFloatTop) onTrail() {}

// proactive action!
func (policy *policyFloatTop) isAllowedMotion(rows int) int {
	return 0
}

func (policy *policyFloatTop) isAllowedTrail() bool {
	return false
}
