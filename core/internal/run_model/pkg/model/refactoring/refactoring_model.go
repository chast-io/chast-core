package refactoring

import "github.com/google/uuid"

type RunModel struct {
	Run []*Run
}

type SingleRunModel struct {
	Run *Run
}

type Run struct {
	ID                 string
	uuid               string
	Dependencies       []*Run
	SupportedLanguages []string
	Docker             *Docker
	Local              *Local
	Command            *Command
	ChangeLocations    *ChangeLocations
}

type ChangeLocations struct {
	Include []string
	Exclude []string
}

func (run *Run) GetUUID() string {
	if run.uuid == "" {
		id := run.ID
		if id != "" {
			id += "-"
		}

		run.uuid = id + uuid.New().String()
	}

	return run.uuid
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
