package builder

import (
	"os"
	"path/filepath"

	recipemodel "chast.io/core/internal/recipe/pkg/model"
	runmodel "chast.io/core/internal/run_model/pkg/model"
	"github.com/joomcode/errorx"
)

func HandlePrimaryArgument(
	primaryParameter *recipemodel.Parameter,
	variables *runmodel.Variables,
	unparsedArgument string,
) error {
	if primaryParameter == nil {
		return errorx.IllegalArgument.New(
			"Missing primary argument. No primary parameter defined for this recipe. This is a required field",
		)
	}

	if unparsedArgument == "" {
		if primaryParameter.DefaultValue == "" {
			return errorx.IllegalArgument.New("Missing primary parameter")
		}

		absoluteDefaultValue, absError := filepath.Abs(primaryParameter.DefaultValue)
		if absError != nil {
			return errorx.ExternalError.Wrap(absError, "Could not absolutize default primary argument path")
		}

		variables.Map[primaryParameter.ID] = absoluteDefaultValue
		variables.TypeDetectionPath = absoluteDefaultValue
		variables.DefaultValueUsed = true

		return nil
	}

	unparsedArgument, absError := filepath.Abs(unparsedArgument)
	if absError != nil {
		return errorx.ExternalError.Wrap(absError, "Could not absolutize primary argument path")
	}

	variables.TypeDetectionPath = unparsedArgument
	wordingDir, _ := os.Getwd()

	if err := handleArgument(primaryParameter, unparsedArgument, wordingDir, variables); err != nil {
		return err
	}

	return nil
}
