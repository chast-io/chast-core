package runmodel

// TODO add function => accept interface and return struct
type RunModel interface{}

type Variables struct {
	WorkingDirectory  string
	Map               map[string]string
	DefaultValueUsed  bool
	TypeDetectionPath string
}

func NewVariables(workingDirectory string) *Variables {
	return &Variables{
		Map:              make(map[string]string),
		WorkingDirectory: workingDirectory,
		DefaultValueUsed: false,
	}
}

type UnparsedFlag struct {
	Name  string
	Value string
}
