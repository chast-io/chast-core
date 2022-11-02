package refactoringpipelinemodel

import (
	"path/filepath"

	"github.com/google/uuid"
)

type Pipeline struct {
	OperationLocation      string
	Stages                 []*Stage
	ChangeCaptureFolder    string
	RootFileSystemLocation string
	UUID                   string
}

func NewPipeline(
	operationLocation string,
	changeCaptureFolder string,
	rootFileSystemLocation string,
) *Pipeline {
	absOperationLocation, _ := filepath.Abs(operationLocation)
	absChangeCaptureFolder, _ := filepath.Abs(changeCaptureFolder)
	absRootFileSystemLocation, _ := filepath.Abs(rootFileSystemLocation)

	pipelineUUID := "PIPELINE-" + uuid.New().String()

	return &Pipeline{
		UUID:                   pipelineUUID,
		Stages:                 make([]*Stage, 1),
		OperationLocation:      absOperationLocation,
		ChangeCaptureFolder:    filepath.Join(absChangeCaptureFolder, pipelineUUID),
		RootFileSystemLocation: absRootFileSystemLocation,
	}
}

func (p *Pipeline) AddStage(stage *Stage) {
	stage.WithPipeline(p)

	if p.Stages[0] == nil {
		p.Stages[0] = stage
	} else {
		stage.prev = p.Stages[len(p.Stages)-1]
		p.Stages = append(p.Stages, stage)
	}
}
