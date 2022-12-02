package refactoringpipelinemodel

import (
	"path/filepath"

	"chast.io/core/internal/internal_util/collection"
	"chast.io/core/internal/run_model/pkg/model/refactoring"
)

type Step struct {
	// TODO replace single run model with generic run model containing only the necessary information
	RunModel              *refactoring.SingleRunModel
	UUID                  string
	ChangeCaptureLocation string
	OperationLocation     string

	Pipeline     *Pipeline
	Dependencies []*Step
	Dependents   []*Step
}

func NewStep(runModel *refactoring.SingleRunModel) *Step {
	runUUID := runModel.Run.GetUUID()

	return &Step{ //nolint:exhaustruct // rest initialized in withStage
		RunModel: runModel,
		UUID:     runUUID,
	}
}

func (s *Step) withPipeline(pipeline *Pipeline) {
	s.ChangeCaptureLocation = filepath.Join(pipeline.GetTempChangeCaptureLocation(), s.UUID)
	s.OperationLocation = filepath.Join(pipeline.OperationLocation, s.UUID)
	s.Pipeline = pipeline
}

func (s *Step) AddDependency(dependency *Step) {
	s.Dependencies = append(s.Dependencies, dependency)
	dependency.Dependents = append(dependency.Dependents, s)
}

func (s *Step) IsFinalStep() bool {
	return len(s.Dependents) == 0
}

func (s *Step) GetFinalChangesLocation() string {
	return s.ChangeCaptureLocation + "-final"
}

func (s *Step) GetChangesStagingLocation() string {
	return s.ChangeCaptureLocation + "-staging"
}

func (s *Step) GetMergedPreviousChangesLocation() string {
	return s.ChangeCaptureLocation + "-prev"
}

func (s *Step) GetPreviousChangeCaptureLocations() []string {
	locations := make([]string, 0)
	steps := make([]*Step, 0)

	steps = append(steps, s.Dependencies...)

	i := 0
	for i < len(steps) {
		step := steps[i]
		locations = collection.Prepend(locations, step.ChangeCaptureLocation)
		steps = append(steps, step.Dependencies...)
		i++
	}

	return locations
}

func (s *Step) ChangeFilteringLocations() *refactoring.ChangeLocations {
	return s.RunModel.Run.ChangeLocations
}
