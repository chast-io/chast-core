package parser_test

import (
	"os"
	"testing"

	recipemodel "chast.io/core/internal/recipe/pkg/model"
	"chast.io/core/internal/recipe/pkg/parser"
)

func TestParseRecipe_Refactoring(t *testing.T) {
	t.Parallel()

	t.Run("CompleteValid", func(t *testing.T) {
		t.Parallel()
		testParseRecipeRefactoringCompleteValid(t)
	})

	t.Run("DuplicateIds", func(t *testing.T) {
		t.Parallel()
		testDuplicateIds(t)
	})

	t.Run("Invalid Dependencies", func(t *testing.T) {
		t.Parallel()
		testInvalidDependencies(t)
	})

	t.Run("Cyclic Dependencies", func(t *testing.T) {
		t.Parallel()
		testCyclicDependencies(t)
	})

	t.Run("Self Referencing Dependencies", func(t *testing.T) {
		t.Parallel()
		testSelfReferencingDependencies(t)
	})
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

			if len(recipe.Runs) != 3 {
				t.Errorf("Expected recipe runs to have 3 entries, but had %d", len(recipe.Runs))
			}

			testRun(t, &recipe.Runs[0], recipemodel.Run{
				ID:                  "rearrange_class_members_java",
				Dependencies:        make([]string, 0),
				SupportedExtensions: []string{"java"},
				Flags:               make([]recipemodel.Flag, 0),
				Docker:              nil,
				Local:               nil,
				Script:              []string{"java -jar ./file.jar $inputFile $configFile > ${inputFile}.out"},
				ChangeLocations:     make([]string, 0),
			})

			testRun(t, &recipe.Runs[1], recipemodel.Run{
				ID:                  "mv_files",
				Dependencies:        []string{"rearrange_class_members_java"},
				SupportedExtensions: make([]string, 0),
				Flags:               make([]recipemodel.Flag, 0),
				Docker:              nil,
				Local:               nil,
				Script:              []string{"mv ${inputFile}.out $inputFile"},
				ChangeLocations:     make([]string, 0),
			})

			testRun(t, &recipe.Runs[2], recipemodel.Run{
				ID:                  "rearrange_class_members_cs",
				Dependencies:        make([]string, 0),
				SupportedExtensions: []string{"cs"},
				Flags:               make([]recipemodel.Flag, 0),
				Docker:              nil,
				Local:               nil,
				Script:              []string{"java -jar ./file.jar $inputFile $configFile > ${inputFile}.out", "mv ${inputFile}.out $inputFile"},
				ChangeLocations:     make([]string, 0),
			})
		})

		// TODO: test tests section
	})
}

func testDuplicateIds(t *testing.T) {
	t.Helper()

	fileData, err := os.ReadFile("testdata/refactoring_parser/duplicate_ids_recipe.yml")
	if err != nil {
		t.Fatalf("Error reading test recipe: %v", err)
	}

	refactoringParser := &parser.RefactoringParser{}
	_, parseError := refactoringParser.ParseRecipe(&fileData)

	if parseError == nil {
		t.Fatal("Expected error, but was nil")
	}
}

func testInvalidDependencies(t *testing.T) {
	t.Helper()

	fileData, err := os.ReadFile("testdata/refactoring_parser/invalid_dependencies_recipe.yml")
	if err != nil {
		t.Fatalf("Error reading test recipe: %v", err)
	}

	refactoringParser := &parser.RefactoringParser{}
	_, parseError := refactoringParser.ParseRecipe(&fileData)

	if parseError == nil {
		t.Fatal("Expected error, but was nil")
	}
}

func testCyclicDependencies(t *testing.T) {
	t.Helper()

	t.Run("Cyclic Dependency at Start", func(t *testing.T) {
		t.Parallel()

		fileData, err := os.ReadFile("testdata/refactoring_parser/cyclic_dependencies_recipe_start.yml")
		if err != nil {
			t.Fatalf("Error reading test recipe: %v", err)
		}

		refactoringParser := &parser.RefactoringParser{}
		_, parseError := refactoringParser.ParseRecipe(&fileData)

		if parseError == nil {
			t.Fatal("Expected error, but was nil")
		}
	})

	t.Run("Cyclic Dependency in Middle", func(t *testing.T) {
		t.Parallel()

		fileData, err := os.ReadFile("testdata/refactoring_parser/cyclic_dependencies_recipe_middle.yml")
		if err != nil {
			t.Fatalf("Error reading test recipe: %v", err)
		}

		refactoringParser := &parser.RefactoringParser{}
		_, parseError := refactoringParser.ParseRecipe(&fileData)

		if parseError == nil {
			t.Fatal("Expected error, but was nil")
		}
	})

	t.Run("Cyclic Dependency at End", func(t *testing.T) {
		t.Parallel()

		fileData, err := os.ReadFile("testdata/refactoring_parser/cyclic_dependencies_recipe_end.yml")
		if err != nil {
			t.Fatalf("Error reading test recipe: %v", err)
		}

		refactoringParser := &parser.RefactoringParser{}
		_, parseError := refactoringParser.ParseRecipe(&fileData)

		if parseError == nil {
			t.Fatal("Expected error, but was nil")
		}
	})
}

func testSelfReferencingDependencies(t *testing.T) {
	t.Helper()

	fileData, err := os.ReadFile("testdata/refactoring_parser/self_referencing_dependencies.yml")
	if err != nil {
		t.Fatalf("Error reading test recipe: %v", err)
	}

	refactoringParser := &parser.RefactoringParser{}
	_, parseError := refactoringParser.ParseRecipe(&fileData)

	if parseError == nil {
		t.Fatal("Expected error, but was nil")
	}
}
