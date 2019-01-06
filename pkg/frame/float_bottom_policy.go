package frame

func newFloatBottomPolicy(frame *Frame) *floatBottomPolicy {
	return &floatBottomPolicy{
		Frame: frame,
	}
}

// proactive action!
// note: most frame objects don't exist, make changes based on the frame config
func (policy *floatBottomPolicy) onInit() {
	height := policy.Frame.config.Lines
	if policy.Frame.config.HasHeader {
		height++
	}
	if policy.Frame.config.HasFooter {
		height++
	}
	// the screen index starts at 1 (not 0), hence the +1
	policy.Frame.topRow = (terminalHeight - height) + 1
}

// reactive action!
func (policy *floatBottomPolicy) onResize(adjustment int) {
	if adjustment > 0 {
		// Grow in size:
		// appended rows should appear to move upwards on the screen, which means that we should
		// move the entire frame upwards 1 line while making more screen space by 1 line
		policy.Frame.move(-adjustment)
		policy.Frame.rowAdvancements += adjustment
	} else if adjustment < 0 {
		// Shrink in size:
		policy.Frame.move(adjustment)
	}

}

// reactive policy!
func (policy *floatBottomPolicy) onTrail() {
	// write the removed line to the trail log + move the policy down (while advancing the frame)
	// policy.frame.move(1)
	policy.Frame.rowAdvancements += 1
}

// reactive action!
// update any positions based on external data and redraw
// func (policy *floatBottomPolicy) onUpdate() {
// 	height := policy.Frame.Height()
// 	targetFrameStartRow := (terminalHeight - height) + 1
// 	if policy.Frame.topRow != targetFrameStartRow {
// 		// reset the policy and all activeLines to the correct offset. This must be done with new
// 		// lines since we should not overwrite the trail rows above the policy.
// 		policy.Frame.rowAdvancements += policy.Frame.topRow - targetFrameStartRow
// 	}
// }

// proactive policy!
// func (policy *floatBottomPolicy) onClose() {
// 	// allow new real estate to be created for the cursor to be placed after the frame at the bottom of the screen
// 	// policy.frame.rowAdvancements += 1
// 	// advanceScreen(1)
//
// 	// no no: it is possible for a bottom frame to exist without the cursor at the bottom of the screen
// 	// do nothing!
// }

// proactive action!
func (policy *floatBottomPolicy) allowedMotion(rows int) int {
	return 0
}

func (policy *floatBottomPolicy) allowTrail() bool {
	return true
}
