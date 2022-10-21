package refactoring

import "chast.io/core/internal/model/run_models"

type RunModel struct {
	run_models.RunModel
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
	Cmd         string
}

type Command struct {
	Cmd              []string
	WorkingDirectory string
}
