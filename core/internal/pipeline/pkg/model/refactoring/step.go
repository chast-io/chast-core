package refactoringpipelinemodel

import (
	"path/filepath"

	"chast.io/core/internal/run_model/pkg/model/refactoring"
)

type Step struct {
	// TODO replace single run model with generic run model containing only the necessary information
	RunModel            *refactoring.SingleRunModel
	UUID                string
	ChangeCaptureFolder string
	OperationLocation   string
}

func NewStep(runModel *refactoring.SingleRunModel) *Step {
	runUUID := runModel.Run.GetUUID()

	return &Step{ //nolint:exhaustruct // rest initialized in WithStage
		RunModel: runModel,
		UUID:     runUUID,
	}
}

func (s *Step) WithStage(stage *Stage) {
	s.ChangeCaptureFolder = filepath.Join(stage.ChangeCaptureFolder, s.UUID)
	s.OperationLocation = filepath.Join(stage.OperationLocation, s.UUID)
}
