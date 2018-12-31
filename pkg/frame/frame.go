package frame

import (
	"fmt"
)

func (float FloatRule) String() string {
	switch float {
	case FloatFree:
		return "FloatFree"
	case FloatTop:
		return "FloatTop"
	case FloatBottom:
		return "FloatBottom"
	default:
		return fmt.Sprintf("FloatRule=%d?", float)
	}
}

func Factory(configs ...Config) []Frame {
	frames := make([]Frame, 0)

	// todo: check config continuity/validity

	for _, config := range configs {
		switch config.Float {
		case FloatFree:
			startRow, err := GetCursorRow()
			if err != nil {
				panic(err)
			}
			config.startRow = startRow
			frames = append(frames, newFreeFrame(config))
		case FloatTop:
			config.startRow = 1
			frames = append(frames, newTopFrame(config))
		case FloatBottom:
			height := config.Lines
			if config.HasHeader {
				height++
			}
			if config.HasFooter {
				height++
			}
			// the screen index starts at 1 (not 0), hence the +1
			config.startRow = (terminalHeight - height) + 1
			frames = append(frames, newBottomFrame(config))
		default:
			panic(fmt.Errorf("unknown FloatRule: %v", config.Float))
		}
	}
	return frames
}
