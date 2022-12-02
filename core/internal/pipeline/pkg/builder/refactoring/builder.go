package refactoringpipelinebuilder

import (
	"chast.io/core/internal/internal_util/collection"
	dependencygraph "chast.io/core/internal/pipeline/internal/dependency_graph"
	refactoringpipelinemodel "chast.io/core/internal/pipeline/pkg/model/refactoring"
	refactoringRunModelIsolator "chast.io/core/internal/run_model/pkg/isolator/refactoring"
	"chast.io/core/internal/run_model/pkg/model/refactoring"
	"github.com/joomcode/errorx"
)

func BuildRunPipeline(runModel *refactoring.RunModel) (*refactoringpipelinemodel.Pipeline, error) {
	// TODO verify id uniqueness
	isolatedExecutionOrder, isolatedExecutionOrderBuildError := buildIsolatedExecutionOrder(runModel)
	if isolatedExecutionOrderBuildError != nil {
		return nil, errorx.InternalError.Wrap(isolatedExecutionOrderBuildError, "failed to build isolated execution order")
	}

	// TODO make configurable
	pipeline := refactoringpipelinemodel.NewPipeline("/tmp/chast/", "/tmp/chast-changes/", "/")

	stepsLookup := make(map[*refactoring.Run]*refactoringpipelinemodel.Step)

	for _, runModelsInStage := range isolatedExecutionOrder {
		executionGroup := refactoringpipelinemodel.NewExecutionGroup()

		for _, runModel := range runModelsInStage {
			step := refactoringpipelinemodel.NewStep(runModel)
			stepsLookup[runModel.Run] = step

			for _, dependency := range runModel.Run.Dependencies {
				step.AddDependency(stepsLookup[dependency])
			}

			executionGroup.AddStep(step)
		}

		pipeline.AddExecutionGroup(executionGroup)
	}

	return pipeline, nil
}

func buildIsolatedExecutionOrder(
	runModel *refactoring.RunModel,
) ([][]*refactoring.SingleRunModel, error) {
	executionOrder, executionOrderBuildError := dependencygraph.BuildExecutionOrder(runModel)
	if executionOrderBuildError != nil {
		return nil, errorx.InternalError.Wrap(executionOrderBuildError, "failed to build execution order")
	}

	return collection.Map(executionOrder, func(run []*refactoring.Run) []*refactoring.SingleRunModel {
		return collection.Map(run, func(run *refactoring.Run) *refactoring.SingleRunModel {
			return refactoringRunModelIsolator.Isolate(runModel, run)
		})
	}), nil
}
