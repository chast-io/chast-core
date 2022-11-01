package strategy

type IsolationStrategy = uint8

const (
	UnknownIsolation IsolationStrategy = iota
	OverlayFS        IsolationStrategy = iota
	UnionFS          IsolationStrategy = iota
)
