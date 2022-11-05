package refactoringpipelinemodel

import (
	"path/filepath"

	"github.com/google/uuid"
)

type Stage struct {
	UUID                string
	ChangeCaptureFolder string
	OperationLocation   string

	Steps []*Step

	prev *Stage
}

func NewStage(name string) *Stage {
	extendedUUID := "STAGE-"
	if name != "" {
		extendedUUID += name + "-"
	}

	extendedUUID += uuid.New().String()

	return &Stage{ //nolint:exhaustruct // rest initialized in WithPipeline
		UUID:  extendedUUID,
		Steps: make([]*Step, 0),
		prev:  nil,
	}
}

func (s *Stage) AddStep(step *Step) {
	s.Steps = append(s.Steps, step)
}

func (s *Stage) WithPipeline(targetPipeline *Pipeline) {
	for _, step := range s.Steps {
		s.ChangeCaptureFolder = filepath.Join(targetPipeline.ChangeCaptureFolder, "tmp", s.UUID)
		s.OperationLocation = filepath.Join(targetPipeline.OperationLocation, s.UUID)
		step.WithStage(s)
	}
}

func (s *Stage) GetPrevChangeCaptureFolders() []string {
	if s.prev == nil {
		return []string{}
	}

	return append(s.prev.GetPrevChangeCaptureFolders(), s.prev.ChangeCaptureFolder)
}
