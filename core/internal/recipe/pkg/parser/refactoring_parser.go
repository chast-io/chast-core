package parser

import (
	"strings"

	recipemodel "chast.io/core/internal/recipe/pkg/model"
	"github.com/pkg/errors"
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

	var recipe recipemodel.Recipe = refactoringRecipe

	return &recipe, nil
}
