package refactoringpipelinebuilder

import (
	"strconv"

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

	for i, runModelsInStage := range isolatedExecutionOrder {
		stage := refactoringpipelinemodel.NewStage(strconv.Itoa(i + 1))

		for _, runModel := range runModelsInStage {
			step := refactoringpipelinemodel.NewStep(runModel)
			stage.AddStep(step)
		}

		pipeline.AddStage(stage)
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
