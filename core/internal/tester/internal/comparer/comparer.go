package comparer

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	chastlog "chast.io/core/internal/logger"
	refactoringpipelinemodel "chast.io/core/internal/pipeline/pkg/model/refactoring"
	recipemodel "chast.io/core/internal/recipe/pkg/model"
	"github.com/joomcode/errorx"
	"github.com/spf13/afero"
)

func CompareResults(test *recipemodel.Test, pipeline *refactoringpipelinemodel.Pipeline, workingDir string) {
	expectedOutputFolderPath, _ := filepath.Abs(filepath.Join(workingDir, "tests", test.ID, "expected"))
	inputFolderPath, _ := filepath.Abs(filepath.Join(workingDir, "tests", test.ID, "input"))

	if !checkFolderEquality(pipeline.GetFinalChangeCaptureLocation(), expectedOutputFolderPath, inputFolderPath) {
		chastlog.Log.Errorf("Test %s failed", test.ID)
	} else {
		chastlog.Log.Infof("Test %s passed", test.ID)
	}
}

func checkFolderEquality(checkFolder string, expectedOutputFolder string, inputFolderPath string) bool {
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
		chastlog.Log.Errorf("Expected %v files, got %v", len(expectedFileStructure), len(actualFileStructure))

		return false
	}

	sort.Strings(expectedFileStructure)
	sort.Strings(actualFileStructure)

	for index := range expectedFileStructure {
		if expectedFileStructure[index] != strings.TrimPrefix(actualFileStructure[index], inputFolderPath) {
			chastlog.Log.Errorf("Expected %v, got %v", expectedFileStructure[index], actualFileStructure[index])

			return false
		}

		if !compareFiles(
			filepath.Join(checkFolder, actualFileStructure[index]),
			filepath.Join(expectedOutputFolder, expectedFileStructure[index]),
		) {
			chastlog.Log.Errorf("File does not match the expected.")
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

	actualFileContent = bytes.TrimSpace(actualFileContent)

	expectedFileContent, expectedFileReadError := io.ReadAll(expectedFile)
	if expectedFileReadError != nil {
		chastlog.Log.Errorf("Could not read file %s: %v", expectedFilePath, expectedFileReadError)

		return false
	}

	expectedFileContent = bytes.TrimSpace(expectedFileContent)

	isEqual := bytes.Equal(actualFileContent, expectedFileContent)

	if !isEqual {
		chastlog.Log.Debugf("Expected: \n%s\n\nGot: \n%s\n", expectedFileContent, actualFileContent)
	}

	return isEqual
}
