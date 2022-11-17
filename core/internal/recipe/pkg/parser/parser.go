package parser

import (
	"strings"

	chastlog "chast.io/core/internal/logger"
	recipemodel "chast.io/core/internal/recipe/pkg/model"
	"github.com/joomcode/errorx"
	"gopkg.in/yaml.v3"
)

type RecipeParser interface {
	ParseRecipe(data *[]byte) (*recipemodel.Recipe, error)
}

type fileReader interface {
	Read() *[]byte
}

func ParseRecipe(file fileReader) (*recipemodel.Recipe, error) {
	fileData := file.Read()

	parser, err := getParser(fileData)
	if err != nil {
		return nil, errorx.InternalError.Wrap(err, "Failed to get parser")
	}

	recipe, parseRecipeErr := parser.ParseRecipe(fileData)
	if parseRecipeErr != nil {
		return nil, errorx.InternalError.Wrap(parseRecipeErr, "Failed to parse recipe")
	}

	return recipe, nil
}

func getParser(fileData *[]byte) (RecipeParser, error) { //nolint:ireturn // Factory function
	recipeType, err := getRecipeType(fileData)
	if err != nil {
		return nil, err
	}

	switch recipeType { //nolint:exhaustive // Others are handled by default case
	case recipemodel.Refactoring:
		return &RefactoringParser{}, nil
	default:
		return nil, errorx.UnsupportedOperation.New("Unknown config type. Available types: refactoring")
	}
}

func getRecipeType(data *[]byte) (recipemodel.ChastOperationType, error) {
	var plainConfigRoot recipemodel.RecipeInfo
	err := yaml.Unmarshal(*data, &plainConfigRoot)

	if err != nil {
		chastlog.Log.Fatalf("error: %v", err)
	}

	switch strings.ToLower(plainConfigRoot.Type) {
	case "refactoring":
		if plainConfigRoot.Version == "1" || plainConfigRoot.Version == "1.0" {
			return recipemodel.Refactoring, nil
		}

		return recipemodel.Refactoring,
			errorx.UnsupportedVersion.New("Unknown refactoring version. Only version 1.0 is supported")

	default:
		return recipemodel.Unknown, nil
	}
}
