package frame

type Config struct {
	test           bool
	Lines          int
	startRow       int
	HeaderRows     int
	FooterRows     int
	TrailOnRemove  bool
	PositionPolicy PositionPolicy
	ManualDraw     bool
}

func (config *Config) VisibleHeight() int {
	height := config.Height()
	forwardDrawAreaHeight := terminalHeight - (config.startRow - 1)

	if height > forwardDrawAreaHeight {
		return forwardDrawAreaHeight
	}
	return height
}

func (config *Config) Height() int {
	return config.Lines + config.HeaderRows + config.FooterRows
}
