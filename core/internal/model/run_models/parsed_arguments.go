package run_models

type ParsedArguments struct {
	Arguments         map[string]string
	UnmappedArguments []string
	WorkingDirectory  string
}
