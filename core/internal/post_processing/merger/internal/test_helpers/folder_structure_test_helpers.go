package testhelper

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/spf13/afero"
)

func CheckFolderEquality(t *testing.T, expectedFileStructure []string, checkFolder string) {
	t.Helper()

	t.Run("Check file structure equality", func(t *testing.T) {
		t.Parallel()

		actualFileStructure, err := CollectPathsInFolder(checkFolder)
		if err != nil {
			t.Fatalf("Could not collect paths in folder %s: %v", checkFolder, err)
		}

		if len(expectedFileStructure) != len(actualFileStructure) {
			t.Fatalf("MergeFolders() expected %v, got %v", expectedFileStructure, actualFileStructure)
		}

		sort.Strings(expectedFileStructure)
		sort.Strings(actualFileStructure)

		for i := range expectedFileStructure {
			if expectedFileStructure[i] != actualFileStructure[i] {
				t.Errorf("MergeFolders() expected %v, got %v", expectedFileStructure, actualFileStructure)
			}
		}
	})
}

func CollectPathsInFolder(targetFolder string) ([]string, error) {
	actualFileStructure := make([]string, 0)
	err := filepath.Walk(targetFolder, func(path string, info os.FileInfo, err error) error {
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
	})

	return actualFileStructure, err
}

func FileStructureCreator(filesAndFolders []string, name string) string {
	cleanedName := regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(name, "")
	cleanedName = strings.ReplaceAll(cleanedName, " ", "_")
	targetFolder, _ := os.MkdirTemp("", cleanedName)

	for _, fileOrFolder := range filesAndFolders {
		_ = os.MkdirAll(filepath.Join(targetFolder, filepath.Dir(fileOrFolder)), os.ModePerm)
		if strings.HasSuffix(fileOrFolder, "/") {
			_ = os.MkdirAll(filepath.Join(targetFolder, fileOrFolder), os.ModePerm)
		} else {
			_, _ = os.Create(filepath.Join(targetFolder, fileOrFolder))
		}
	}

	return targetFolder
}
