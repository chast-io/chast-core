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
		WorkingDirectory:  workingDirectory,
		Map:               make(map[string]string),
		DefaultValueUsed:  false,
		TypeDetectionPath: "",
	}
}

type UnparsedFlag struct {
	Name  string
	Value string
}
