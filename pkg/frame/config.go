package frame

func (config *Config) VisibleHeight() int {
	height := config.Height()
	forwardDrawAreaHeight := terminalHeight - (config.startRow - 1)

	if height > forwardDrawAreaHeight {
		return forwardDrawAreaHeight
	}
	return height
}

func (config *Config) Height() int {
	height := config.Lines
	if config.HasHeader {
		height++
	}
	if config.HasFooter {
		height++
	}
	return height
}
