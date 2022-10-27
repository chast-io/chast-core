package parser

import (
	"chast.io/core/internal/model/recipe"
	util "chast.io/core/pkg/util"
	"errors"
	"gopkg.in/yaml.v3"
	"log"
	"strings"
)

type RecipeParser interface {
	ParseRecipe(data *[]byte) (*recipe.Recipe, error)
}

func ParseRecipe(file util.FileReader) (*recipe.Recipe, error) {
	fileData := file.Read()
	parser, err := getParser(fileData)
	if err != nil {
		return nil, err
	}
	return parser.ParseRecipe(fileData)
}

func getParser(fileData *[]byte) (RecipeParser, error) {
	recipeType, err := getRecipeType(fileData)
	if err != nil {
		return nil, err
	}
	switch recipeType {
	case recipe.Refactoring:
		return &RefactoringParser{}, nil
	default:
		return nil, errors.New("unknown config type - available types: refactoring")
	}
}

func getRecipeType(data *[]byte) (recipe.ChastOperationType, error) {
	var plainConfigRoot recipe.RecipeInfo
	err := yaml.Unmarshal(*data, &plainConfigRoot)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	switch strings.ToLower(plainConfigRoot.Type) {
	case "refactoring":
		if plainConfigRoot.Version == "1" || plainConfigRoot.Version == "1.0" {
			return recipe.Refactoring, nil
		}
		return recipe.Refactoring, errors.New("unknown refactoring version - only version 1.0 is supported")
	default:
		return recipe.Unknown, nil
	}
}
