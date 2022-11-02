package builder

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"chast.io/core/internal/internal_util/collection"
	recipemodel "chast.io/core/internal/recipe/pkg/model"
	refactoringrunmodelbuilder "chast.io/core/internal/run_model/pkg/builder/refactoring"
	runmodel "chast.io/core/internal/run_model/pkg/model"
	"github.com/pkg/errors"
)

type RunModelBuilder interface {
	BuildRunModel(*recipemodel.Recipe, *runmodel.ParsedArguments) (*runmodel.RunModel, error)
}

func BuildRunModel(
	parsedRecipe *recipemodel.Recipe,
	arguments []string,
	workingDirectory string,
) (*runmodel.RunModel, error) {
	runModelBuilder, baseRecipe, builderError := getBuilder(parsedRecipe)
	if builderError != nil {
		return nil, builderError
	}

	parsedArguments, argumentsHandleError := handleArguments(baseRecipe, arguments)
	if argumentsHandleError != nil {
		return nil, argumentsHandleError
	}

	parsedArguments.WorkingDirectory = workingDirectory

	runModel, runModelBuildError := runModelBuilder.BuildRunModel(parsedRecipe, parsedArguments)

	return runModel, errors.Wrap(runModelBuildError, "Failed to build run model")
}

func getBuilder(parsedRecipe *recipemodel.Recipe) (RunModelBuilder, *recipemodel.BaseRecipe, error) { //nolint:lll,ireturn // This is a factory function that returns a builder for generating a run model
	switch concreateRecipe := (*parsedRecipe).(type) {
	case *recipemodel.RefactoringRecipe:
		var runModelBuilder RunModelBuilder = refactoringrunmodelbuilder.NewRunModelBuilder()

		return runModelBuilder, &concreateRecipe.BaseRecipe, nil
	default:
		return nil, nil, errors.Errorf("No run model builder for recipe of type %T", concreateRecipe.GetRecipeType())
	}
}

func handleArguments(baseRecipe *recipemodel.BaseRecipe, args []string) (*runmodel.ParsedArguments, error) {
	requiredArgsCount := collection.Count(
		baseRecipe.Arguments,
		func(argument recipemodel.Argument) bool { return argument.Required },
	)
	if len(args) < requiredArgsCount {
		return nil, errors.Errorf("Not enough arguments passed. Expected %d, got %d", len(baseRecipe.Arguments), len(args))
	}

	argsMap := make(map[string]string)
	wordingDir, _ := os.Getwd()

	for argIndex, argDecl := range baseRecipe.Arguments {
		if argIndex < len(args) {
			if err := verifyArgument(&baseRecipe.Arguments[argIndex], args[argIndex]); err != nil {
				return nil, err
			}

			if argDecl.Type == "Path" && !strings.HasPrefix(args[argIndex], "/") {
				argsMap[argDecl.ID], _ = filepath.Abs(filepath.Join(wordingDir, args[argIndex]))
			} else {
				argsMap[argDecl.ID] = args[argIndex]
			}
		}
	}

	return &runmodel.ParsedArguments{
		Arguments:         argsMap,
		UnmappedArguments: args[len(argsMap):],
		WorkingDirectory:  wordingDir,
	}, nil
}

func verifyArgument(argument *recipemodel.Argument, value string) error {
	if argument.Type != "" {
		switch argument.Type {
		case "string":
			break
		case "Bool":
			if !(value == "true" || value == "yes" || value == "false" || value == "no") {
				return errors.Errorf("Argument %s is not a boolean. Passed argument: %s", argument.ID, value)
			}
		case "Int":
			if _, err := strconv.Atoi(value); err != nil {
				return errors.Errorf("Argument %s is not an integer. Passed argument: %s", argument.ID, value)
			}
		case "Path":
			break
		default:
			return errors.Errorf("Unknown argument type %v", argument.Type)
		}
	}

	return nil
}
