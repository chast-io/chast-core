package parser_test

import (
	"os"
	"testing"

	recipemodel "chast.io/core/internal/recipe/pkg/model"
	"chast.io/core/internal/recipe/pkg/parser"
)

func TestParseRecipe_Refactoring(t *testing.T) {
	t.Parallel()

	testParseRecipeRefactoringCompleteValid(t)
}

func testParseRecipeRefactoringCompleteValid(t *testing.T) {
	t.Helper()

	fileData, err := os.ReadFile("testdata/refactoring_parser/complete_valid_recipe.yml")
	if err != nil {
		t.Fatalf("Error reading test recipe: %v", err)
	}

	refactoringParser := &parser.RefactoringParser{}
	genericRecipe, parseError := refactoringParser.ParseRecipe(&fileData)

	if parseError != nil {
		t.Fatalf("Expected no error, but was '%v'", parseError)
	}

	var recipe *recipemodel.RefactoringRecipe

	switch actualRecipe := (*genericRecipe).(type) {
	case *recipemodel.RefactoringRecipe:
		recipe = actualRecipe
	default:
		t.Fatalf("Expected recipe to be of type RefactoringRecipe, but was %T", genericRecipe)
	}

	t.Run("Recipe", func(t *testing.T) {
		t.Parallel()

		if recipe == nil {
			t.Error("Expected recipe to be set, but was nil")
		}

		t.Run("PrimaryParameter", func(t *testing.T) {
			t.Parallel()

			if recipe.PrimaryParameter == nil {
				t.Error("Expected recipe primary parameter to be set, but was nil")
			}

			testParameter(t, recipe.PrimaryParameter, recipemodel.Parameter{
				ID: "inputFile",
				RequiredExtension: recipemodel.RequiredExtension{
					Required:     true,
					DefaultValue: "./src/main/java",
				},
				TypeExtension: recipemodel.TypeExtension{
					Type:       "filePath",
					Extensions: make([]string, 0),
				},
				DescriptionExtension: recipemodel.DescriptionExtension{
					Description:     "The file to be refactored.",
					LongDescription: "Long description of the primary parameter.",
				},
			})
		})

		t.Run("Runs", func(t *testing.T) {
			t.Parallel()

			if recipe.Runs == nil || len(recipe.Runs) == 0 {
				t.Error("Expected recipe runs to be set, but was nil")
			}
		})

		//t.Run("Tests", func(t *testing.T) {
		//	t.Skip("TODO")
		//})
	})
}
