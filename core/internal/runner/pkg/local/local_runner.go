package local

import (
	"chast.io/core/internal/changeisolator/pkg"
	"chast.io/core/internal/changeisolator/pkg/namespace"
	"chast.io/core/internal/changeisolator/pkg/strategy"
	refactoringPipelineModel "chast.io/core/internal/pipeline/pkg/model/refactoring"
	refactoringpipelinecleanup "chast.io/core/internal/post_processing/cleanup/refactoring"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
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
	log.Printf("Running pipeline %s", pipeline.UUID)

	if r.isolated && !r.parallel {
		return sequentialRun(pipeline)
	}

	return errors.Errorf("Unisolated and parallel execution is not yet implemented")
}

func sequentialRun(pipeline *refactoringPipelineModel.Pipeline) error {
	for _, stage := range pipeline.Stages {
		for _, step := range stage.Steps {
			if err := runIsolated(step, stage, pipeline); err != nil {
				return errors.Wrap(err, "Error running isolated")
			}
		}

		if err := refactoringpipelinecleanup.CleanupStage(stage); err != nil {
			return errors.Wrap(err, "Error cleaning up stage")
		}
	}

	if err := refactoringpipelinecleanup.CleanupPipeline(pipeline); err != nil {
		return errors.Wrap(err, "Failed to cleanup pipeline")
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

	if err := pkg.RunCommandInIsolatedEnvironment(nsContext); err != nil {
		return errors.Errorf("Error running command in isolated environment - %s", err)
	}

	return nil
}
