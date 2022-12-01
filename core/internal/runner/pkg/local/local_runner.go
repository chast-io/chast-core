package local

import (
	pipelinepostprocessor "chast.io/core/internal/post_processing/pipeline_post_processor/pkg/refactoring"
	steppostprocessor "chast.io/core/internal/post_processing/step_post_processor/pkg/refactoring"
	"os"

	changeisolator "chast.io/core/internal/changeisolator/pkg"
	"chast.io/core/internal/changeisolator/pkg/namespace"
	"chast.io/core/internal/changeisolator/pkg/strategy"
	chastlog "chast.io/core/internal/logger"
	refactoringPipelineModel "chast.io/core/internal/pipeline/pkg/model/refactoring"
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
	for _, stage := range pipeline.ExecutionGroups {
		for _, step := range stage.Steps {
			chastlog.Log.Printf("Running step %s", step.UUID)

			if err := runIsolated(step); err != nil {
				return errorx.InternalError.Wrap(err, "Error running isolated")
			}
		}
	}

	if err := pipelinepostprocessor.Process(pipeline); err != nil {
		return errorx.InternalError.Wrap(err, "Error running post processing")
	}

	return nil
}

func runIsolated(
	step *refactoringPipelineModel.Step,
) error {
	if err := os.MkdirAll(step.GetMergedPreviousChangesLocation(), os.ModePerm); err != nil {
		return errorx.ExternalError.Wrap(err, "Failed to create previous changes directory")
	}

	var nsContext = namespace.NewContext(
		step.Pipeline.RootFileSystemLocation,
		step.GetPreviousChangesLocations(),
		step.ChangeCaptureLocation,
		step.OperationLocation,
		step.RunModel.Run.Command.WorkingDirectory,
		step.RunModel.Run.Command.Cmds,
		strategy.UnionFS,
	)

	if err := changeisolator.RunCommandInIsolatedEnvironment(nsContext); err != nil {
		return errorx.InternalError.New("Error running command in isolated environment - %s", err)
	}

	if err := steppostprocessor.Process(step); err != nil {
		return errorx.InternalError.Wrap(err, "Error running post processing")
	}

	return nil
}
