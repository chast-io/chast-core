package recipe

type ChastOperationType int8

const (
	Unknown     ChastOperationType = iota
	Refactoring ChastOperationType = iota
)