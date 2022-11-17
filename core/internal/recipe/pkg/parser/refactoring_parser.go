package parser

import (
	"fmt"
	"strings"

	"chast.io/core/internal/internal_util/collection"
	refactroingdependencygraph "chast.io/core/internal/recipe/internal/refactoring/dependency_graph"
	recipemodel "chast.io/core/internal/recipe/pkg/model"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type RefactoringParser struct{}

func (parser *RefactoringParser) ParseRecipe(data *[]byte) (*recipemodel.Recipe, error) {
	var refactoringRecipe *recipemodel.RefactoringRecipe

	decoder := yaml.NewDecoder(strings.NewReader(string(*data)))
	decoder.KnownFields(true)

	if err := decoder.Decode(&refactoringRecipe); err != nil {
		return nil, errors.Wrap(err, "Error parsing refactoring recipe")
	}

	if err := validateRecipe(refactoringRecipe); err != nil {
		return nil, errors.Wrap(err, "Error validating refactoring recipe")
	}

	var recipe recipemodel.Recipe = refactoringRecipe

	return &recipe, nil
}

func validateRecipe(recipe *recipemodel.RefactoringRecipe) error {
	if err := validateRuns(recipe.Runs); err != nil {
		return errors.Wrap(err, "Error validating primary parameter")
	}

	supportedExtensionsOfRuns := collection.Reduce(recipe.Runs, func(run recipemodel.Run, acc []string) []string {
		return append(acc, run.SupportedExtensions...)
	}, make([]string, 0))

	if err := validatePrimaryParameter(recipe.PrimaryParameter, supportedExtensionsOfRuns); err != nil {
		return errors.Wrap(err, "Error validating primary parameter")
	}

	return nil
}

// TODO: check dependencies
//   - dependencies must exist
//   - dependencies must not be self-referencing
//
// TODO: check for duplicate IDs
func validateRuns(runs []recipemodel.Run) error {
	if len(runs) == 0 {
		return errors.New("At least one run is required")
	}

	presentRunIds := make(map[string]bool)
	for runIndex := range runs {
		if err := validateID(&runs[runIndex], presentRunIds); err != nil {
			return errors.Wrap(err, "Error validating run ID")
		}

		if err := validateRun(&runs[runIndex]); err != nil {
			return errors.Wrap(err, "Error validating run")
		}
	}

	if err := validateDependencies(runs, presentRunIds); err != nil {
		return errors.Wrap(err, "Error validating run dependencies")
	}

	dependencyGraph := refactroingdependencygraph.BuildDependencyGraph(runs)
	if dependencyGraph.HasCycles() {
		return errors.New("Recipe dependencies contains a cycle")
	}

	return nil
}

func validateID(run *recipemodel.Run, presentRunIds map[string]bool) error {
	if presentRunIds[run.ID] {
		return errors.New(fmt.Sprintf("Duplicate run ID '%s'", run.ID)) //nolint:revive,lll // TODO: use new error strategy #18
	}

	presentRunIds[run.ID] = true

	return nil
}

func validateDependencies(runs []recipemodel.Run, presentRunIds map[string]bool) error {
	for _, run := range runs {
		for _, dependency := range run.Dependencies {
			if !presentRunIds[dependency] {
				return errors.New(fmt.Sprintf("Run '%s' depends on unknown run '%s'", run.ID, dependency)) //nolint:revive,lll // TODO: use new error strategy #18
			}

			if dependency == run.ID {
				return errors.New(fmt.Sprintf("Run '%s' depends on itself", run.ID)) //nolint:revive,lll // TODO: use new error strategy #18
			}
		}
	}

	return nil
}

func validateRun(run *recipemodel.Run) error {
	if run == nil {
		return errors.New("Run cannot be nil")
	}

	// VALIDATE FLAGS

	if run.Script == nil || len(run.Script) == 0 {
		return errors.New("Run script is required")
	}

	// TODO add change locations

	return nil
}

var errInvalidPrimaryParameter = errors.New("Invalid primary parameter error: ")

func validatePrimaryParameter(parameter *recipemodel.Parameter, supportedExtensions []string) error {
	if parameter == nil {
		return errors.New("Primary parameter is required")
	}

	if parameter.ID == "" {
		parameter.ID = "primaryParameter" // TODO make this configurable
		log.Printf("Primary parameter ID is not set and falls back to '%s'\n", parameter.ID)
	}

	if !parameter.RequiredExtension.Required {
		// TODO show message if it was explicitly set to false
		parameter.RequiredExtension.Required = true
	}

	// TODO make this configurable
	options := []string{"filePath", "folderPath", "wildcardPath", "string", "int", "boolean"}
	if parameter.TypeExtension.Type == "" {
		return errors.Wrap(
			errInvalidPrimaryParameter,
			fmt.Sprintf("Primary parameter type is required. Options: %s", options),
		)
	}

	if !collection.Include(options, parameter.TypeExtension.Type) {
		return errors.Wrap(errInvalidPrimaryParameter, fmt.Sprintf("Must be of type %s", options))
	}

	if parameter.TypeExtension.Extensions == nil || len(parameter.TypeExtension.Extensions) == 0 {
		parameter.TypeExtension.Extensions = supportedExtensions
	} else {
		return errors.Wrap(
			errInvalidPrimaryParameter,
			"Primary parameter can not contain extensions as they are defined by the supported extensions of the runs",
		)
	}

	if parameter.DescriptionExtension.Description == "" {
		log.Println("It is advised to provide a description for the primary parameter.")
	}

	return nil
}
