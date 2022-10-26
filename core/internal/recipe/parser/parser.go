package parser

import (
	"chast.io/core/internal/recipe/model"
	util "chast.io/core/pkg/util"
	"errors"
	"gopkg.in/yaml.v3"
	"log"
	"strings"
)

import (
	"chast.io/core/internal/model/run_models"
)

type RecipeParser interface {
	ParseRecipe(data *[]byte)
	VerifyRecipeAndBuildModel() (*run_models.RunModel, error)
}

func ParseRecipe(file util.FileReader) (*run_models.RunModel, error) {
	fileData := file.Read()
	parser, err := getParser(fileData)
	if err != nil {
		return nil, err
	}
	parser.ParseRecipe(fileData)
	return parser.VerifyRecipeAndBuildModel()
}

func getParser(fileData *[]byte) (RecipeParser, error) {
	recipeType, err := getRecipeType(fileData)
	if err != nil {
		return nil, err
	}
	switch recipeType {
	case model.Refactoring:
		return &RefactoringParser{}, nil
	default:
		return nil, errors.New("unknown config type - available types: refactoring")
	}
}

func getRecipeType(data *[]byte) (model.ChastOperationType, error) {
	var plainConfigRoot model.RecipeInfo
	err := yaml.Unmarshal(*data, &plainConfigRoot)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	switch strings.ToLower(plainConfigRoot.Type) {
	case "refactoring":
		if plainConfigRoot.Version == "1" || plainConfigRoot.Version == "1.0" {
			return model.Refactoring, nil
		}
		return model.Refactoring, errors.New("unknown refactoring version - only version 1.0 is supported")
	default:
		return model.Unknown, nil
	}
}
