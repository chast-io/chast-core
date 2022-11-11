package parser

import (
	"chast.io/core/internal/internal_util/collection"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"

	recipemodel "chast.io/core/internal/recipe/pkg/model"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type RefactoringParser struct{}

// TODO validate recipe! After optional positional arguments, no more required positional arguments are allowed
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
	// TODO: check dependencies (dependencies must exist, no cycles, etc.)

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

func validateRuns(runs []recipemodel.Run) error {
	if runs == nil || len(runs) == 0 {
		return errors.New("At least one run is required")
	}

	for _, run := range runs {
		if err := validateRun(&run); err != nil {
			return errors.Wrap(err, "Error validating run")
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

func validatePrimaryParameter(parameter *recipemodel.Parameter, supportedExtensions []string) error {
	if parameter == nil {
		return errors.New("Primary parameter is required")
	}

	if parameter.ID == "" {
		parameter.ID = "primaryParameter" // TODO make this configurable
		log.Printf("Primary parameter ID is not set and falls back to '%s'\n", parameter.ID)
	}

	if parameter.RequiredExtension.Required == false {
		// TODO show message if it was explicitly set to false
		parameter.RequiredExtension.Required = true
	}

	options := []string{"filePath", "folderPath", "wildcardPath", "string", "int", "boolean"} // TODO make this configurable
	if parameter.TypeExtension.Type == "" {
		return errors.New(fmt.Sprintf("Primary parameter type is required. Options: %s", options))
	}

	if !collection.Include(options, parameter.TypeExtension.Type) {
		return errors.New(fmt.Sprintf("Invalid primary parameter type. Must be of type %s", options))
	}

	if parameter.TypeExtension.Extensions == nil || len(parameter.TypeExtension.Extensions) == 0 {
		parameter.TypeExtension.Extensions = supportedExtensions
	} else {
		return errors.New("Primary parameter can not contain extensions as they are defined by the supported extensions of the runs.")
	}

	if parameter.DescriptionExtension.Description == "" {
		log.Println("It is advised to provide a description for the primary parameter.")
	}

	return nil
}
