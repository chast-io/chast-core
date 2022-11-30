package refactoringpipelinemodel

import (
	"path/filepath"

	"github.com/google/uuid"
)

type Pipeline struct {
	OperationLocation      string
	ExecutionGroups        []*ExecutionGroup
	ChangeCaptureLocation  string
	RootFileSystemLocation string
	UUID                   string
}

func NewPipeline(
	operationLocation string,
	changeCaptureLocation string,
	rootFileSystemLocation string,
) *Pipeline {
	absOperationLocation, _ := filepath.Abs(operationLocation)
	absChangeCaptureLocation, _ := filepath.Abs(changeCaptureLocation)
	absRootFileSystemLocation, _ := filepath.Abs(rootFileSystemLocation)

	pipelineUUID := "PIPELINE-" + uuid.New().String()

	return &Pipeline{
		UUID:                   pipelineUUID,
		ExecutionGroups:        make([]*ExecutionGroup, 1),
		OperationLocation:      absOperationLocation,
		ChangeCaptureLocation:  filepath.Join(absChangeCaptureLocation, pipelineUUID),
		RootFileSystemLocation: absRootFileSystemLocation,
	}
}

func (p *Pipeline) GetTempChangeCaptureLocation() string {
	return filepath.Join(p.ChangeCaptureLocation, "tmp")
}

func (p *Pipeline) AddExecutionGroup(executionGroup *ExecutionGroup) {
	executionGroup.withPipeline(p)

	if p.ExecutionGroups[0] == nil {
		p.ExecutionGroups[0] = executionGroup
	} else {
		p.ExecutionGroups = append(p.ExecutionGroups, executionGroup)
	}
}

func (p *Pipeline) GetFinalSteps() []*Step {
	finalSteps := make([]*Step, 0)

	for _, executionGroup := range p.ExecutionGroups {
		for _, step := range executionGroup.Steps {
			if step.IsFinalStep() {
				finalSteps = append(finalSteps, step)
			}
		}
	}

	return finalSteps
}
