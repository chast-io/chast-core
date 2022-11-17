package builder

import (
	"os"
	"strconv"
	"strings"

	"chast.io/core/internal/internal_util/collection"
	recipemodel "chast.io/core/internal/recipe/pkg/model"
	runmodel "chast.io/core/internal/run_model/pkg/model"
	"github.com/joomcode/errorx"
)

type handleFlagsMapper interface {
	GetFlags() []recipemodel.Flag
	GetFlagsMap() map[string]*recipemodel.Flag
}

func HandleFlags(
	flagsMapper handleFlagsMapper,
	variables *runmodel.Variables,
	unparsedFlags []runmodel.UnparsedFlag,
) error {
	requiredFlagsCount := collection.Count(
		flagsMapper.GetFlags(),
		func(flag recipemodel.Flag) bool { return flag.Required && flag.DefaultValue == "" },
	)

	flagsDefinitionMap := flagsMapper.GetFlagsMap()
	coveredRequiredFlags := 0
	wordingDir, _ := os.Getwd()

	for _, flag := range unparsedFlags {
		flagDefinition := flagsDefinitionMap[flag.Name]
		if flagDefinition == nil {
			return errorx.IllegalArgument.New("Unknown flag %s", flag.Name)
		}

		if flagDefinition.Required && flagDefinition.DefaultValue == "" {
			coveredRequiredFlags++
		}

		value := flag.Value

		if err := verifyFlagValue(flagDefinition, value); err != nil {
			return err
		}

		value, absolutizePathFlagError := absolutizePath(value, flagDefinition.TypeExtension, wordingDir)
		if absolutizePathFlagError != nil {
			return absolutizePathFlagError
		}

		variables.Map[flagDefinition.Name] = value
	}

	if coveredRequiredFlags != requiredFlagsCount {
		return errorx.IllegalFormat.New("Not all required flags are set")
	}

	return nil
}

func verifyFlagValue(flagDefinition *recipemodel.Flag, value string) error {
	if flagDefinition.Type == "" {
		return nil
	}

	switch flagDefinition.Type {
	case "string":
		break
	case "bool":
		if !(value == "true" || value == "yes" || value == "false" || value == "no") {
			return errorx.IllegalArgument.New("Flag %s is not a boolean. Passed parameter: %s", flagDefinition.Name, value)
		}
	case "int":
		if _, err := strconv.Atoi(value); err != nil {
			return errorx.IllegalArgument.New("Flag %s is not an integer. Passed parameter: %s", flagDefinition.Name, value)
		}
	default:
		if strings.HasSuffix(flagDefinition.Type, "Path") && flagDefinition.Extensions != nil {
			if err := verifyPathExtension(value, flagDefinition.Extensions); err != nil {
				return err
			}

			break
		}

		return errorx.IllegalArgument.New("Unknown flag type %v", flagDefinition.Type)
	}

	return nil
}
