package run_model_builder

import (
	"chast.io/core/internal/model/recipe"
	"chast.io/core/internal/model/run_models"
	"chast.io/core/internal/recipe/run_model_builder/refactoring"
	"chast.io/core/pkg/util/collection"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type RunModelBuilder interface {
	BuildRunModel(*recipe.Recipe, *run_models.ParsedArguments) (*run_models.RunModel, error)
}

func BuildRunModel(parsedRecipe *recipe.Recipe, arguments []string, workingDirectory string) (*run_models.RunModel, error) {
	runModelBuilder, baseRecipe, err := getBuilder(parsedRecipe)
	if err != nil {
		return nil, err
	}

	parsedArguments, err := handleArguments(baseRecipe, arguments)
	if err != nil {
		return nil, err
	}

	parsedArguments.WorkingDirectory = workingDirectory
	return runModelBuilder.BuildRunModel(parsedRecipe, parsedArguments)
}

func getBuilder(parsedRecipe *recipe.Recipe) (RunModelBuilder, *recipe.BaseRecipe, error) {
	switch m := (*parsedRecipe).(type) {
	case *recipe.RefactoringRecipe:
		var runModelBuilder RunModelBuilder = refactoring.NewRunModelBuilder()
		return runModelBuilder, &m.BaseRecipe, nil
	default:
		return nil, nil, errors.Errorf("No run model builder for recipe of type %T", m.GetRecipeType())
	}
}

func handleArguments(baseRecipe *recipe.BaseRecipe, args []string) (*run_models.ParsedArguments, error) {
	requiredArgsCount := collection.Count(baseRecipe.Arguments, func(argument recipe.Argument) bool { return argument.Required })
	if len(args) < requiredArgsCount {
		return nil, errors.Errorf("Not enough arguments passed. Expected %d, got %d", len(baseRecipe.Arguments), len(args))
	}
	argsMap := make(map[string]string)
	wordingDir, _ := os.Getwd()

	for i, argDecl := range baseRecipe.Arguments {
		if i < len(args) {
			if err := verifyArgument(&argDecl, args[i]); err != nil {
				return nil, err
			}
			if argDecl.Type == "Path" && !strings.HasPrefix(args[i], "/") {
				argsMap[argDecl.ID], _ = filepath.Abs(filepath.Join(wordingDir, args[i]))
			} else {
				argsMap[argDecl.ID] = args[i]
			}
		}
	}
	return &run_models.ParsedArguments{Arguments: argsMap, UnmappedArguments: args[len(argsMap):]}, nil
}

func verifyArgument(argument *recipe.Argument, value string) error {
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
