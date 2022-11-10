package refactoringpipelinemodel

import (
	"path/filepath"

	"github.com/google/uuid"
)

type Pipeline struct {
	OperationLocation      string
	Stages                 []*Stage
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
		Stages:                 make([]*Stage, 1),
		OperationLocation:      absOperationLocation,
		ChangeCaptureLocation:  filepath.Join(absChangeCaptureLocation, pipelineUUID),
		RootFileSystemLocation: absRootFileSystemLocation,
	}
}

func (p *Pipeline) AddStage(stage *Stage) {
	stage.withPipeline(p)

	if p.Stages[0] == nil {
		p.Stages[0] = stage
	} else {
		stage.prev = p.Stages[len(p.Stages)-1]
		p.Stages = append(p.Stages, stage)
	}
}
