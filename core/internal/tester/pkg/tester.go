package tester

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"chast.io/core/internal/internal_util/collection"
	chastlog "chast.io/core/internal/logger"
	refactoringpipelinemodel "chast.io/core/internal/pipeline/pkg/model/refactoring"
	recipemodel "chast.io/core/internal/recipe/pkg/model"
	"chast.io/core/internal/recipe/pkg/parser"
	refactoringservice "chast.io/core/internal/service/pkg/refactoring"
	util "chast.io/core/pkg/util/fs/file"
	"github.com/joomcode/errorx"
	"github.com/spf13/afero"
)

func Test(recipeFile *util.File) {
	parsedRecipe, recipeParseError := parser.ParseRecipe(recipeFile)
	if recipeParseError != nil {
		panic(recipeParseError)
	}

	switch concreteRecipe := (*parsedRecipe).(type) {
	case *recipemodel.RefactoringRecipe:
		workingDir := recipeFile.ParentDirectory

		if len(concreteRecipe.Tests) == 0 {
			chastlog.Log.Infof("No tests found for recipe %s", recipeFile.AbsolutePath)

			return
		}

		for index, test := range concreteRecipe.Tests {
			primaryArgument := filepath.Join(workingDir, "tests", test.ID, "input")
			pipeline, recipeRunError := refactoringservice.Run(
				recipeFile,
				append([]string{primaryArgument}, test.Args...),
				convertFlags(test.Flags),
			)

			if recipeRunError != nil {
				panic(recipeParseError)
			}

			compareResults(&concreteRecipe.Tests[index], pipeline, workingDir)
		}
	default:
		panic(errorx.UnsupportedOperation.New("No run model builder for recipe of type %T", concreteRecipe.GetRecipeType()))
	}
}

func convertFlags(flags []string) []refactoringservice.FlagParameter {
	return collection.Map(flags, func(flag string) refactoringservice.FlagParameter {
		split := strings.Split(flag, "=")

		return refactoringservice.FlagParameter{
			Name:  split[0],
			Value: split[1],
		}
	})
}

func compareResults(test *recipemodel.Test, pipeline *refactoringpipelinemodel.Pipeline, workingDir string) {
	expectedOutputFolder := filepath.Join(workingDir, "tests", test.ID, "expected")

	if !checkFolderEquality(expectedOutputFolder, pipeline.GetFinalChangeCaptureLocation()) {
		chastlog.Log.Errorf("Test %s failed", test.ID)
	} else {
		chastlog.Log.Infof("Test %s passed", test.ID)
	}
}

func checkFolderEquality(checkFolder string, expectedOutputFolder string) bool {
	expectedFileStructure, expectedPathCollectionError := collectPathsInFolder(expectedOutputFolder)
	if expectedPathCollectionError != nil {
		chastlog.Log.Errorf("Could not collect paths in folder %s: %v", expectedOutputFolder, expectedPathCollectionError)

		return false
	}

	actualFileStructure, actualPathCollectionError := collectPathsInFolder(checkFolder)
	if actualPathCollectionError != nil {
		chastlog.Log.Errorf("Could not collect paths in folder %s: %v", checkFolder, actualPathCollectionError)

		return false
	}

	if len(expectedFileStructure) != len(actualFileStructure) {
		chastlog.Log.Errorf("MergeFolders() expected %v, got %v", expectedFileStructure, actualFileStructure)

		return false
	}

	sort.Strings(expectedFileStructure)
	sort.Strings(actualFileStructure)

	for i := range expectedFileStructure {
		if expectedFileStructure[i] != actualFileStructure[i] {
			chastlog.Log.Errorf("MergeFolders() expected %v, got %v", expectedFileStructure, actualFileStructure)

			return false
		}

		if !compareFiles(
			filepath.Join(expectedOutputFolder, expectedFileStructure[i]),
			filepath.Join(checkFolder, actualFileStructure[i])) {
			return false
		}

	}

	return true
}

func collectPathsInFolder(targetFolder string) ([]string, error) {
	actualFileStructure := make([]string, 0)

	if err := filepath.Walk(targetFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			actualFileStructure = append(actualFileStructure, strings.TrimPrefix(path, targetFolder))
		} else {
			if empty, _ := afero.IsEmpty(afero.NewOsFs(), path); empty {
				folder := strings.TrimPrefix(path, targetFolder) + "/"
				if folder != "/" {
					actualFileStructure = append(actualFileStructure, folder)
				}
			}
		}

		return nil
	}); err != nil {
		return nil, errorx.ExternalError.Wrap(err, "Could not walk folder %s", targetFolder)
	}

	return actualFileStructure, nil
}

func compareFiles(actualFilePath string, expectedFilePath string) bool {
	actualFile, actualFileOpenError := os.Open(actualFilePath)
	if actualFileOpenError != nil {
		chastlog.Log.Errorf("Could not open file %s: %v", actualFilePath, actualFileOpenError)

		return false
	}

	expectedFile, expectedFileOpenError := os.Open(expectedFilePath)
	if expectedFileOpenError != nil {
		chastlog.Log.Errorf("Could not open file %s: %v", expectedFilePath, expectedFileOpenError)

		return false
	}

	actualFileContent, actualFileReadError := io.ReadAll(actualFile)
	if actualFileReadError != nil {
		chastlog.Log.Errorf("Could not read file %s: %v", actualFilePath, actualFileReadError)

		return false
	}

	expectedFileContent, expectedFileReadError := io.ReadAll(expectedFile)
	if expectedFileReadError != nil {
		chastlog.Log.Errorf("Could not read file %s: %v", expectedFilePath, expectedFileReadError)

		return false
	}

	return bytes.Equal(actualFileContent, expectedFileContent)
}
