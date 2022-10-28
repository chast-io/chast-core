package refactoringPipelineModel

import (
	"chast.io/core/internal/model/run_models/refactoring"
	"path/filepath"
)

type Step struct {
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

func (s Step) WithPipeline(targetPipeline *Pipeline) Step {
	s.ChangeCaptureFolder = filepath.Join(targetPipeline.ChangeCaptureFolder, "tmp", s.UUID)
	s.OperationLocation = filepath.Join(targetPipeline.OperationLocation, s.UUID)
	return s
}
