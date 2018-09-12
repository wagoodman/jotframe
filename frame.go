package jotframe

type Frame interface {
	Header() *Line
	Footer() *Line
	Lines() []*Line
	Append() (*Line, error)
	Prepend() (*Line, error)
	Insert(index int) (*Line, error)
	Remove(line *Line) error
	Advance(rows int) error
	Retreat(rows int) error
	Draw() error
	Close() error
	Clear() error
	ClearAndClose() error
	update() error
}