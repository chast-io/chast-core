package refactoringPipelineModel

type Pipeline struct {
	OperationLocation      string
	Steps                  []*Step
	ChangeCaptureFolder    string
	RootFileSystemLocation string
}

func NewPipeline(
	operationLocation string,
	changeCaptureFolder string,
	rootFileSystemLocation string,
) *Pipeline {
	return &Pipeline{
		OperationLocation:      operationLocation,
		Steps:                  make([]*Step, 2),
		ChangeCaptureFolder:    changeCaptureFolder,
		RootFileSystemLocation: rootFileSystemLocation,
	}
}

func (p *Pipeline) AddStep(step *Step) {
	pipelinedStep := step.WithPipeline(p)
	p.Steps = append(p.Steps, &pipelinedStep)
}
