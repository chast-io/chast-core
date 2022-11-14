package refactoringpipelinemodel

import (
	"path/filepath"

	"github.com/google/uuid"
)

type Stage struct {
	UUID                  string
	ChangeCaptureLocation string
	OperationLocation     string

	Steps []*Step

	prev     *Stage
	pipeline *Pipeline
}

func NewStage(name string) *Stage {
	extendedUUID := "STAGE-"
	if name != "" {
		extendedUUID += name + "-"
	}

	extendedUUID += uuid.New().String()

	return &Stage{ //nolint:exhaustruct // rest initialized in withPipeline
		UUID:  extendedUUID,
		Steps: make([]*Step, 0),
		prev:  nil,
	}
}

func (s *Stage) AddStep(step *Step) {
	if step == nil {
		return
	}

	s.Steps = append(s.Steps, step)
	if s.pipeline != nil {
		s.setStepDetails(step)
	}
}

func (s *Stage) withPipeline(targetPipeline *Pipeline) {
	s.pipeline = targetPipeline
	s.setStageDetails()

	for _, step := range s.Steps {
		s.setStepDetails(step)
	}
}

func (s *Stage) setStageDetails() {
	s.ChangeCaptureLocation = filepath.Join(s.pipeline.ChangeCaptureLocation, "tmp", s.UUID)
	s.OperationLocation = filepath.Join(s.pipeline.OperationLocation, s.UUID)
}

func (s *Stage) setStepDetails(step *Step) {
	step.withStage(s)
}

func (s *Stage) GetPrevChangeCaptureLocations() []string {
	if s.prev == nil {
		return []string{}
	}

	return append(s.prev.GetPrevChangeCaptureLocations(), s.prev.ChangeCaptureLocation)
}
