package refactoringService

import (
	refactoringPipelineBuilder "chast.io/core/internal/pipeline/builder/refactoring"
	"chast.io/core/internal/pipeline/model/refactoring"
	"chast.io/core/internal/run_model/builder"
	"chast.io/core/internal/run_model/model/refactoring"
	"chast.io/core/internal/runner/local"
	util "chast.io/core/pkg/util/fs"
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

	runModel, err := builder.BuildRunModel(parsedRecipe, args, recipeFile.ParentDirectory)
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
