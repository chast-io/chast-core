package parser

import (
	"fmt"
	"strings"

	"chast.io/core/internal/internal_util/collection"
	chastlog "chast.io/core/internal/logger"
	refactroingdependencygraph "chast.io/core/internal/recipe/internal/refactoring/dependency_graph"
	recipemodel "chast.io/core/internal/recipe/pkg/model"
	"github.com/joomcode/errorx"
	"gopkg.in/yaml.v3"
)

type RefactoringParser struct{}

func (parser *RefactoringParser) ParseRecipe(data *[]byte) (*recipemodel.Recipe, error) {
	var refactoringRecipe *recipemodel.RefactoringRecipe

	decoder := yaml.NewDecoder(strings.NewReader(string(*data)))
	decoder.KnownFields(true)

	if err := decoder.Decode(&refactoringRecipe); err != nil {
		return nil, errorx.Decorate(err, "Error parsing refactoring recipe")
	}

	if err := validateRecipe(refactoringRecipe); err != nil {
		return nil, errorx.Decorate(err, "Error validating refactoring recipe")
	}

	var recipe recipemodel.Recipe = refactoringRecipe

	return &recipe, nil
}

func validateRecipe(recipe *recipemodel.RefactoringRecipe) error {
	if err := validateRuns(recipe.Runs); err != nil {
		return errorx.Decorate(err, "Error validating primary parameter")
	}

	supportedExtensionsOfRuns := collection.Reduce(recipe.Runs, func(run recipemodel.Run, acc []string) []string {
		return append(acc, run.SupportedExtensions...)
	}, make([]string, 0))

	if err := validatePrimaryParameter(recipe.PrimaryParameter, supportedExtensionsOfRuns); err != nil {
		return errorx.Decorate(err, "Error validating primary parameter")
	}

	return nil
}

func validateRuns(runs []recipemodel.Run) error {
	if len(runs) == 0 {
		return errorx.IllegalFormat.New("At least one run is required")
	}

	presentRunIds := make(map[string]bool)
	for runIndex := range runs {
		if err := validateID(&runs[runIndex], presentRunIds); err != nil {
			return errorx.Decorate(err, "Error validating run ID")
		}

		if err := validateRun(&runs[runIndex]); err != nil {
			return errorx.Decorate(err, "Error validating run")
		}
	}

	if err := validateDependencies(runs, presentRunIds); err != nil {
		return errorx.Decorate(err, "Error validating run dependencies")
	}

	dependencyGraph := refactroingdependencygraph.BuildDependencyGraph(runs)
	if dependencyGraph.HasCycles() {
		return errorx.IllegalArgument.New("Recipe dependencies contains a cycle")
	}

	return nil
}

func validateID(run *recipemodel.Run, presentRunIds map[string]bool) error {
	if presentRunIds[run.ID] {
		return errorx.IllegalArgument.New(fmt.Sprintf("Duplicate run ID '%s'", run.ID))
	}

	presentRunIds[run.ID] = true

	return nil
}

func validateDependencies(runs []recipemodel.Run, presentRunIds map[string]bool) error {
	for _, run := range runs {
		for _, dependency := range run.Dependencies {
			if !presentRunIds[dependency] {
				return errorx.IllegalArgument.New(fmt.Sprintf("Run '%s' depends on unknown run '%s'", run.ID, dependency))
			}

			if dependency == run.ID {
				return errorx.IllegalArgument.New(fmt.Sprintf("Run '%s' depends on itself", run.ID))
			}
		}
	}

	return nil
}

func validateRun(run *recipemodel.Run) error {
	if run == nil {
		return errorx.IllegalFormat.New("Run cannot be nil")
	}

	// VALIDATE FLAGS

	if run.Script == nil || len(run.Script) == 0 {
		return errorx.IllegalFormat.New("Run script is required")
	}

	// TODO add change locations

	return nil
}

func validatePrimaryParameter(parameter *recipemodel.Parameter, supportedExtensions []string) error {
	if parameter == nil {
		return errorx.IllegalFormat.New("Primary parameter is required")
	}

	if parameter.ID == "" {
		parameter.ID = "primaryParameter" // TODO make this configurable
		chastlog.Log.Printf("Primary parameter ID is not set and falls back to '%s'", parameter.ID)
	}

	if !parameter.RequiredExtension.Required {
		// TODO show message if it was explicitly set to false
		parameter.RequiredExtension.Required = true
	}

	// TODO make this configurable
	options := []string{"filePath", "folderPath", "wildcardPath", "string", "int", "boolean"}
	if parameter.TypeExtension.Type == "" {
		return errorx.IllegalFormat.New(fmt.Sprintf("Primary parameter type is required. Options: %s", options))
	}

	if !collection.Include(options, parameter.TypeExtension.Type) {
		return errorx.IllegalFormat.New(fmt.Sprintf("Must be of type %s", options))
	}

	if parameter.TypeExtension.Extensions == nil || len(parameter.TypeExtension.Extensions) == 0 {
		parameter.TypeExtension.Extensions = supportedExtensions
	} else {
		return errorx.IllegalFormat.New(
			"Primary parameter can not contain extensions as they are defined by the supported extensions of the runs",
		)
	}

	if parameter.DescriptionExtension.Description == "" {
		chastlog.Log.Println("It is advised to provide a description for the primary parameter.")
	}

	return nil
}
