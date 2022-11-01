package local

import (
	"chast.io/core/internal/change_isolator"
	"chast.io/core/internal/pipeline/model/refactoring"
	"github.com/pkg/errors"
)

type Runner struct {
	isolated bool
	parallel bool
}

func NewRunner(isolated bool, parallel bool) *Runner {
	return &Runner{
		isolated: isolated,
		parallel: parallel,
	}
}

func (r *Runner) Run(pipeline *refactoringPipelineModel.Pipeline) error {
	if r.isolated && !r.parallel {
		return sequentialRun(pipeline)
	}
	return errors.Errorf("Unisolated and parallel execution is not yet implemented")
}

func sequentialRun(pipeline *refactoringPipelineModel.Pipeline) error {
	for _, stage := range pipeline.Stages {
		for _, step := range stage.Steps {
			err := runIsolated(step, stage, pipeline)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func runIsolated(
	step *refactoringPipelineModel.Step,
	stage *refactoringPipelineModel.Stage,
	pipeline *refactoringPipelineModel.Pipeline) error {

	var nsContext = change_isolator.NewNamespaceContext(
		pipeline.RootFileSystemLocation,
		stage.GetPrevChangeCaptureFolders(),
		step.ChangeCaptureFolder,
		step.OperationLocation,
		step.RunModel.Run.Command.WorkingDirectory,
		step.RunModel.Run.Command.Cmds,
	)

	if err := change_isolator.RunCommandInIsolatedEnvironment(nsContext); err != nil {
		return errors.Errorf("Error running command in isolated environment - %s", err)
	}
	return nil
}
