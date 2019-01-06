package frame

func newFloatFreePolicy(frame *Frame) *floatFreePolicy {
	return &floatFreePolicy{
		Frame: frame,
	}
}

// proactive action!
// note: most frame objects don't exist, make changes based on the frame config
func (policy *floatFreePolicy) onInit() {
	if policy.Frame.config.startRow == 0 {
		offset, err := GetCursorRow()
		if err != nil {
			return
		}
		policy.Frame.config.startRow = offset
		policy.Frame.topRow = offset
	}
}

// reactive action!
func (policy *floatFreePolicy) onResize(adjustment int) {
	if policy.Frame.IsPastScreenBottom() {
		// make more screen realestate
		policy.Frame.move(-adjustment)
		policy.Frame.rowAdvancements += adjustment
	}
}

func (policy *floatFreePolicy) onTrail() {
	// write the removed line to the trail log + move the policy down (while advancing the frame)
	if policy.Frame.IsPastScreenBottom() {
		// frame.frame.Move(-1)
		policy.Frame.rowAdvancements += 1
	} else {
		policy.Frame.move(1)
	}
}

// func (policy *floatFreePolicy) onUpdate() {}

// func (policy *floatFreePolicy) onClose() {}

func (policy *floatFreePolicy) allowedMotion(rows int) int {
	return rows
}


func (policy *floatFreePolicy) allowTrail() bool {
	return true
}
