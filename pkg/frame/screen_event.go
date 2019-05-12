package frame

type EventHandler interface {
	onEvent(*ScreenEvent)
}

type ScreenEvent struct {
	value []byte
	row   int
}

func newScreenEvent(line *Line) *ScreenEvent {
	e := &ScreenEvent{
		row:   line.row,
		value: make([]byte, len(line.buffer)),
	}
	copy(e.value, line.buffer)
	return e
}
