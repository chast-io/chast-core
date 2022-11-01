package model

type ParsedArguments struct {
	Arguments         map[string]string
	UnmappedArguments []string
	WorkingDirectory  string
}
