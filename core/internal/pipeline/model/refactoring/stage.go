package refactoringPipelineModel

import (
	"github.com/google/uuid"
)

type Stage struct {
	Name string
	UUID string

	Steps []*Step

	prev *Stage
}

func NewStage(name string) *Stage {
	UUID := "STAGE-"
	if name != "" {
		UUID = name + "-"
	}
	name = name + uuid.New().String()

	return &Stage{
		Name: name,
		UUID: UUID,
	}
}

func (s *Stage) AddStep(step *Step) {
	s.Steps = append(s.Steps, step)
}

func (s *Stage) WithPipeline(targetPipeline *Pipeline) {
	for _, step := range s.Steps {
		step.WithPipeline(targetPipeline)
	}
}

func (s *Stage) GetPrevChangeCaptureFolders() []string {
	prevCaptureFolders := make([]string, 0)
	for prev := s.prev; prev != nil; prev = prev.prev {
		for _, step := range prev.Steps {
			prevCaptureFolders = append(prevCaptureFolders, step.ChangeCaptureFolder)
		}
	}
	return prevCaptureFolders
}
