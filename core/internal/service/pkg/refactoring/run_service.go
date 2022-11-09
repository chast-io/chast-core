package refactoringservice

import (
	refactoringPipelineBuilder "chast.io/core/internal/pipeline/pkg/builder/refactoring"
	refactoringpipelinemodel "chast.io/core/internal/pipeline/pkg/model/refactoring"
	"chast.io/core/internal/post_processing/pipelinereport"
	"chast.io/core/internal/recipe/pkg/parser"
	"chast.io/core/internal/run_model/pkg/builder"
	"chast.io/core/internal/run_model/pkg/model/refactoring"
	"chast.io/core/internal/runner/pkg/local"
	util "chast.io/core/pkg/util/fs/file"
	"github.com/pkg/errors"
)

func Run(recipeFile *util.File, args ...string) error {
	parsedRecipe, recipeParseError := parser.ParseRecipe(recipeFile)
	if recipeParseError != nil {
		panic(recipeParseError)
	}

	runModel, runModelBuildError := builder.BuildRunModel(parsedRecipe, args, recipeFile.ParentDirectory)
	if runModelBuildError != nil {
		return errors.Wrap(runModelBuildError, "Failed to build run model")
	}

	var pipeline *refactoringpipelinemodel.Pipeline
	var pipelineBuildError error
	switch m := (*runModel).(type) {
	case refactoring.RunModel:
		pipeline, pipelineBuildError = refactoringPipelineBuilder.BuildRunPipeline(&m)
	default:
		return errors.Errorf("Provided recipe is not a refactoring recipe")
	}

	if pipelineBuildError != nil {
		return errors.Wrap(pipelineBuildError, "Failed to build pipeline")
	}

	if err := local.NewRunner(true, false).Run(pipeline); err != nil {
		return errors.Wrap(err, "Failed to run pipeline")
	}

	report, reportError := pipelinereport.BuildReport(pipeline)
	if reportError != nil {
		return errors.Wrap(reportError, "Failed to generate report")
	}

	report.PrintFileTree(true)
	report.PrintChanges(true)

	//changedFilesRelative, recipeParseError := report.ChangedFilesRelative()
	//if recipeParseError != nil {
	//	return errors.Wrap(recipeParseError, "Failed to generate report")
	//}
	//
	//for _, line := range changedFilesRelative {
	//	println(line)
	//}

	return nil
}
