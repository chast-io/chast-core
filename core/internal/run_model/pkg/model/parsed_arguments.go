package runmodel

type ParsedArguments struct {
	Arguments         map[string]string
	UnmappedArguments []string
	WorkingDirectory  string
}
