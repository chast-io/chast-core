package builder

import (
	"os"
	"path/filepath"

	recipemodel "chast.io/core/internal/recipe/pkg/model"
	runmodel "chast.io/core/internal/run_model/pkg/model"
	"github.com/pkg/errors"
)

var errMissingPrimaryArgument = errors.New("Missing primary argument")

func HandlePrimaryArgument(
	primaryParameter *recipemodel.Parameter,
	variables *runmodel.Variables,
	unparsedArgument string,
) error {
	if primaryParameter == nil {
		return errors.Wrap(
			errMissingPrimaryArgument,
			"No primary parameter defined for this recipe. This is a required field",
		)
	}

	if unparsedArgument == "" {
		if primaryParameter.DefaultValue == "" {
			return errors.New("Missing primary parameter")
		}

		variables.Map[primaryParameter.ID] = primaryParameter.DefaultValue
		variables.TypeDetectionPath = primaryParameter.DefaultValue
		variables.DefaultValueUsed = true

		return nil
	}

	unparsedArgument, absError := filepath.Abs(unparsedArgument)
	if absError != nil {
		return errors.Wrap(absError, "Could not absolutize primary argument path")
	}

	variables.TypeDetectionPath = unparsedArgument
	wordingDir, _ := os.Getwd()

	if err := handleArgument(primaryParameter, unparsedArgument, wordingDir, variables); err != nil {
		return err
	}

	return nil
}