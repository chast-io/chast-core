package refactoringPipelineModel

import (
	"chast.io/core/internal/run_model/model/refactoring"
	"path/filepath"
)

type Step struct {
	// TODO replace single run model with generic run model containing only the necessary information
	RunModel            *refactoring.SingleRunModel
	UUID                string
	ChangeCaptureFolder string
	OperationLocation   string
}

func NewStep(runModel *refactoring.SingleRunModel) *Step {
	UUID := runModel.Run.GetUUID()
	return &Step{
		RunModel: runModel,
		UUID:     UUID,
	}
}

func (s *Step) WithPipeline(targetPipeline *Pipeline) {
	s.ChangeCaptureFolder = filepath.Join(targetPipeline.ChangeCaptureFolder, "tmp", s.UUID)
	s.OperationLocation = filepath.Join(targetPipeline.OperationLocation, s.UUID)
}
