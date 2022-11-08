package refactoringpipelinemodel

import (
	"path/filepath"

	"chast.io/core/internal/run_model/pkg/model/refactoring"
)

type Step struct {
	// TODO replace single run model with generic run model containing only the necessary information
	RunModel              *refactoring.SingleRunModel
	UUID                  string
	ChangeCaptureLocation string
	OperationLocation     string
}

func NewStep(runModel *refactoring.SingleRunModel) *Step {
	runUUID := runModel.Run.GetUUID()

	return &Step{ //nolint:exhaustruct // rest initialized in withStage
		RunModel: runModel,
		UUID:     runUUID,
	}
}

func (s *Step) withStage(stage *Stage) {
	s.ChangeCaptureLocation = filepath.Join(stage.ChangeCaptureLocation, s.UUID)
	s.OperationLocation = filepath.Join(stage.OperationLocation, s.UUID)
}
