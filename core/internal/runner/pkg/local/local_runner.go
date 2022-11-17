package local

import (
	changeisolator "chast.io/core/internal/changeisolator/pkg"
	"chast.io/core/internal/changeisolator/pkg/namespace"
	"chast.io/core/internal/changeisolator/pkg/strategy"
	chastlog "chast.io/core/internal/logger"
	refactoringPipelineModel "chast.io/core/internal/pipeline/pkg/model/refactoring"
	refactoringpipelinecleanup "chast.io/core/internal/post_processing/cleanup/refactoring"
	"github.com/joomcode/errorx"
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
	chastlog.Log.Printf("Running pipeline %s", pipeline.UUID)

	if r.isolated && !r.parallel {
		return sequentialRun(pipeline)
	}

	return errorx.NotImplemented.New("Unisolated and parallel execution is not yet implemented")
}

func sequentialRun(pipeline *refactoringPipelineModel.Pipeline) error {
	for _, stage := range pipeline.Stages {
		for _, step := range stage.Steps {
			if err := runIsolated(step, stage, pipeline); err != nil {
				return errorx.InternalError.Wrap(err, "Error running isolated")
			}
		}

		if err := refactoringpipelinecleanup.CleanupStage(stage); err != nil {
			return errorx.InternalError.Wrap(err, "Error cleaning up stage")
		}
	}

	if err := refactoringpipelinecleanup.CleanupPipeline(pipeline); err != nil {
		return errorx.InternalError.Wrap(err, "Failed to cleanup pipeline")
	}

	return nil
}

func runIsolated(
	step *refactoringPipelineModel.Step,
	stage *refactoringPipelineModel.Stage,
	pipeline *refactoringPipelineModel.Pipeline,
) error {
	var nsContext = namespace.NewContext(
		pipeline.RootFileSystemLocation,
		stage.GetPrevChangeCaptureLocations(),
		step.ChangeCaptureLocation,
		step.OperationLocation,
		step.RunModel.Run.Command.WorkingDirectory,
		step.RunModel.Run.Command.Cmds,
		strategy.UnionFS,
	)

	if err := changeisolator.RunCommandInIsolatedEnvironment(nsContext); err != nil {
		return errorx.InternalError.New("Error running command in isolated environment - %s", err)
	}

	return nil
}
