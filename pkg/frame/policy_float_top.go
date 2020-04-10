package frame

type floatTopPolicy struct {
	Frame *Frame
}


func newFloatTopPolicy(frame *Frame) *floatTopPolicy {
	return &floatTopPolicy{
		Frame: frame,
	}
}

// proactive action!
// note: most frame objects don't exist, make changes based on the frame config
func (policy *floatTopPolicy) onInit() {
	policy.Frame.Config.startRow = 1
	policy.Frame.startIdx = 1
	offset, err := GetCursorRow()
	if err != nil {
		return
	}
	policy.Frame.rowAdvancements = offset - 1
}

// reactive action!
func (policy *floatTopPolicy) onResize(adjustment int) {}

// reactive policy!
func (policy *floatTopPolicy) onTrail() {}

// reactive action!
// update any positions based on external data and redraw
// func (policy *floatTopPolicy) onUpdate() {
// 	if policy.Frame.topRow != 1 {
// 		policy.Frame.move(policy.Frame.topRow - 1)
// 	}
// }

// proactive policy!
// func (policy *floatTopPolicy) onClose() {}

// proactive action!
func (policy *floatTopPolicy) allowedMotion(rows int) int {
	return 0
}

func (policy *floatTopPolicy) isAllowedTrail() bool {
	return false
}
