package strategy

type IsolationStrategy = uint8

const (
	_         IsolationStrategy = iota
	OverlayFS IsolationStrategy = iota
	UnionFS   IsolationStrategy = iota
)
