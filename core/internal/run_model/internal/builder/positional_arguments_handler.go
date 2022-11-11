package builder

import (
	"os"

	"chast.io/core/internal/internal_util/collection"
	recipemodel "chast.io/core/internal/recipe/pkg/model"
	runmodel "chast.io/core/internal/run_model/pkg/model"
	"github.com/pkg/errors"
)

type handlePositionalArgumentsMapper interface {
	GetPositionalParameters() []recipemodel.Parameter
}

func HandlePositionalArguments(
	baseRecipe handlePositionalArgumentsMapper,
	variables *runmodel.Variables,
	arguments []string,
) error {
	positionalParameters := baseRecipe.GetPositionalParameters()

	requiredArgsCount := collection.Count(
		positionalParameters,
		func(argument recipemodel.Parameter) bool { return argument.Required && argument.DefaultValue == "" },
	)
	if len(arguments) < requiredArgsCount {
		return errors.Errorf("Not enough positional arguments passed. Expected %d, got %d", len(positionalParameters), len(arguments))
	}

	wordingDir, _ := os.Getwd()

	for index, parameter := range positionalParameters {
		if index < len(arguments) {
			if variables.DefaultValueUsed {
				return errors.New("After using a default value, no more required positional arguments are allowed")
			}

			argument := arguments[index]

			if err := handleArgument(&positionalParameters[index], argument, wordingDir, variables); err != nil {
				return err
			}
		} else if parameter.Required {
			if parameter.DefaultValue == "" {
				return errors.Errorf("Missing positional argument %s", parameter.ID)
			}

			variables.Map[parameter.ID] = parameter.DefaultValue
			variables.DefaultValueUsed = true
		}
	}

	return nil
}
