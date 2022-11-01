package parser

import (
	. "chast.io/core/internal/recipe/model"
	"gopkg.in/yaml.v3"
	"strings"
)

type RefactoringParser struct {
}

func (parser *RefactoringParser) ParseRecipe(data *[]byte) (*Recipe, error) {
	var refactoringRecipe *RefactoringRecipe

	decoder := yaml.NewDecoder(strings.NewReader(string(*data)))
	decoder.KnownFields(true)

	err := decoder.Decode(&refactoringRecipe)
	if err != nil {
		return nil, err
	}

	var recipe Recipe = refactoringRecipe
	return &recipe, nil
}
