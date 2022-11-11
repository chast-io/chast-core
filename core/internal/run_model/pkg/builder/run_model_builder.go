package builder

import (
	recipemodel "chast.io/core/internal/recipe/pkg/model"
	"chast.io/core/internal/run_model/internal/builder"
	refactoringrunmodelbuilder "chast.io/core/internal/run_model/pkg/builder/refactoring"
	runmodel "chast.io/core/internal/run_model/pkg/model"
	"github.com/pkg/errors"
)

type RunModelBuilder interface {
	BuildRunModel(
		*recipemodel.Recipe,
		*runmodel.Variables,
		[]string,
		[]runmodel.UnparsedFlag,
	) (*runmodel.RunModel, error)
}

func BuildRunModel(
	parsedRecipe *recipemodel.Recipe,
	arguments []string,
	flags []runmodel.UnparsedFlag,
	recipeDirectory string,
) (*runmodel.RunModel, error) {
	runModelBuilder, baseRecipe, builderError := getBuilder(parsedRecipe)
	if builderError != nil {
		return nil, builderError
	}

	variables := runmodel.NewVariables(recipeDirectory)

	if err := builder.HandleFlags(baseRecipe, variables, flags); err != nil {
		return nil, err
	}

	runModel, runModelBuildError := runModelBuilder.BuildRunModel(parsedRecipe, variables, arguments, flags)

	return runModel, errors.Wrap(runModelBuildError, "Failed to build run model")
}

func getBuilder(parsedRecipe *recipemodel.Recipe) (RunModelBuilder, *recipemodel.BaseRecipe, error) { //nolint:lll,ireturn // This is a factory function that returns a builder for generating a run model
	switch concreteRecipe := (*parsedRecipe).(type) {
	case *recipemodel.RefactoringRecipe:
		var runModelBuilder RunModelBuilder = refactoringrunmodelbuilder.NewRunModelBuilder()

		return runModelBuilder, &concreteRecipe.BaseRecipe, nil
	default:
		return nil, nil, errors.Errorf("No run model builder for recipe of type %T", concreteRecipe.GetRecipeType())
	}
}
