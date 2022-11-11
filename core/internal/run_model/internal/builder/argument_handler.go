package builder

import (
	"strconv"
	"strings"

	recipemodel "chast.io/core/internal/recipe/pkg/model"
	runmodel "chast.io/core/internal/run_model/pkg/model"
	"github.com/pkg/errors"
)

func handleArgument(
	parameter *recipemodel.Parameter,
	argument string,
	wordingDir string,
	variables *runmodel.Variables,
) error {
	if err := verifyArgument(parameter, argument); err != nil {
		return err
	}

	argument, absolutizePathFlagError := absolutizePath(argument, parameter.TypeExtension, wordingDir)
	if absolutizePathFlagError != nil {
		return absolutizePathFlagError
	}

	variables.Map[parameter.ID] = argument

	return nil
}

func verifyArgument(parameter *recipemodel.Parameter, value string) error {
	if parameter.Type == "" {
		return nil
	}

	switch parameter.Type {
	case "string":
		break
	case "bool":
		if !(value == "true" || value == "yes" || value == "false" || value == "no") {
			return errors.Errorf("Parameter %s is not a boolean. Passed parameter: %s", parameter.ID, value)
		}
	case "int":
		if _, err := strconv.Atoi(value); err != nil {
			return errors.Errorf("Parameter %s is not an integer. Passed parameter: %s", parameter.ID, value)
		}
	default:
		if strings.HasSuffix(parameter.Type, "Path") {
			if err := verifyPath(parameter, value); err != nil {
				return err
			}
			return nil
		}

		return errors.Errorf("Unknown parameter type %v", parameter.Type)
	}

	return nil
}

func verifyPath(parameter *recipemodel.Parameter, value string) error {
	if parameter.Extensions == nil {
		return nil
	}

	if err := verifyPathExtension(value, parameter.Extensions); err != nil {
		return err
	}

	return nil
}
