package refactoringservice

import (
	"chast.io/core/internal/internal_util/collection"
	refactoringPipelineBuilder "chast.io/core/internal/pipeline/pkg/builder/refactoring"
	refactoringpipelinemodel "chast.io/core/internal/pipeline/pkg/model/refactoring"
	"chast.io/core/internal/post_processing/pipelinereport"
	"chast.io/core/internal/recipe/pkg/parser"
	"chast.io/core/internal/run_model/pkg/builder"
	runmodel "chast.io/core/internal/run_model/pkg/model"
	"chast.io/core/internal/run_model/pkg/model/refactoring"
	"chast.io/core/internal/runner/pkg/local"
	util "chast.io/core/pkg/util/fs/file"
	"github.com/joomcode/errorx"
)

func Run(
	recipeFile *util.File,
	args []string,
	flags []struct {
		name  string
		value string
	},
) error {
	parsedRecipe, recipeParseError := parser.ParseRecipe(recipeFile)
	if recipeParseError != nil {
		panic(recipeParseError)
	}

	runModel, runModelBuildError := builder.BuildRunModel(parsedRecipe, args, mapFlags(flags), recipeFile.ParentDirectory)
	if runModelBuildError != nil {
		return errorx.InternalError.Wrap(runModelBuildError, "Failed to build run model")
	}

	var pipeline *refactoringpipelinemodel.Pipeline

	var pipelineBuildError error

	switch m := (*runModel).(type) {
	case refactoring.RunModel:
		pipeline, pipelineBuildError = refactoringPipelineBuilder.BuildRunPipeline(&m)
	default:
		return errorx.InternalError.New("Provided recipe is not a refactoring recipe")
	}

	if pipelineBuildError != nil {
		return errorx.InternalError.Wrap(pipelineBuildError, "Failed to build pipeline")
	}

	if err := local.NewRunner(true, false).Run(pipeline); err != nil {
		return errorx.InternalError.Wrap(err, "Failed to run pipeline")
	}

	report, reportError := pipelinereport.BuildReport(pipeline)
	if reportError != nil {
		return errorx.InternalError.Wrap(reportError, "Failed to generate report")
	}

	report.PrintFileTree(true)
	report.PrintChanges(true)

	return nil
}

func mapFlags(flags []struct {
	name  string
	value string
}) []runmodel.UnparsedFlag {
	return collection.Map(flags, func(flag struct {
		name  string
		value string
	}) runmodel.UnparsedFlag {
		return runmodel.UnparsedFlag{
			Name:  flag.name,
			Value: flag.value,
		}
	})
}
