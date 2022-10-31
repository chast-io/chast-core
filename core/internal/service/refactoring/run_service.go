package refactoringService

import (
	refactoringPipelineBuilder "chast.io/core/internal/core/pipeline_builder/refactoring"
	refactoringPipelineModel "chast.io/core/internal/model/pipeline/refactoring"
	"chast.io/core/internal/model/run_models/refactoring"
	"chast.io/core/internal/recipe/run_model_builder"
	"chast.io/core/internal/runner/local"
	util "chast.io/core/pkg/util"
	"github.com/pkg/errors"
)
import (
	"chast.io/core/internal/recipe/parser"
)

func Run(recipeFile *util.File, args ...string) error {
	parsedRecipe, err := parser.ParseRecipe(recipeFile)
	if err != nil {
		panic(err)
	}

	runModel, err := run_model_builder.BuildRunModel(parsedRecipe, args, recipeFile.ParentDirectory)
	if err != nil {
		return err
	}

	var pipeline *refactoringPipelineModel.Pipeline
	switch m := (*runModel).(type) {
	case refactoring.RunModel:
		pipeline = refactoringPipelineBuilder.BuildRunPipeline(&m)
	default:
		return errors.Errorf("Provided recipe is not a refactoring recipe")
	}

	if err := local.NewRunner(true, false).Run(pipeline); err != nil {
		return err
	}

	return nil
}
