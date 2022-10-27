package refactoring

type RunModel struct {
	SupportedLanguages []string
	Run                []Run
}

type Run struct {
	Docker  Docker
	Local   Local
	Command Command
}

type Docker struct {
	DockerImage string
}

type Local struct {
	RequiredTools []RequiredTool
}

type RequiredTool struct {
	Description string
	CheckCmd    string
}

type Command struct {
	Cmds             [][]string
	WorkingDirectory string
}
