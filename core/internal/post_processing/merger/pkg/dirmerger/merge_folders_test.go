package dirmerger_test

import (
	"chast.io/core/internal/post_processing/merger/pkg/dirmerger"
	"github.com/spf13/afero"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

type mergeFoldersTestCase struct {
	name                  string
	args                  args
	sourceFileStructure   []string
	targetFileStructure   []string
	expectedFileStructure []string
	wantErr               bool
}

type args struct {
	blockOverwrite bool
}

const unionFsHiddenPathSuffix = "_HIDDEN~"

func TestMergeFolders(t *testing.T) {
	tests := []mergeFoldersTestCase{
		{
			name: "Merge two empty folders",
			args: args{
				blockOverwrite: false,
			},
			sourceFileStructure: []string{},
			targetFileStructure: []string{},
			wantErr:             false,
		},
		{
			name: "Merge two folders with one file each [non conflicting]",
			args: args{
				blockOverwrite: false,
			},
			sourceFileStructure:   []string{"/file1"},
			targetFileStructure:   []string{"/file2"},
			expectedFileStructure: []string{"/file1", "/file2"},
			wantErr:               false,
		},
		{
			name: "Merge two folders with one file each [conflicting - blockOverwrite = false]",
			args: args{
				blockOverwrite: false,
			},
			sourceFileStructure:   []string{"/file1"},
			targetFileStructure:   []string{"/file1"},
			expectedFileStructure: []string{"/file1"},
			wantErr:               false,
		},
		{
			name: "Merge two folders with one file each [conflicting - blockOverwrite = true]",
			args: args{
				blockOverwrite: true,
			},
			sourceFileStructure:   []string{"/file1"},
			targetFileStructure:   []string{"/file1"},
			expectedFileStructure: []string{"/file1"},
			wantErr:               true,
		},
		{
			name: "Merge deleted files [blockOverwrite = false]",
			args: args{
				blockOverwrite: false,
			},
			sourceFileStructure:   []string{"/file1" + unionFsHiddenPathSuffix},
			targetFileStructure:   []string{"/file1"},
			expectedFileStructure: []string{"/file1", "/file1" + unionFsHiddenPathSuffix},
			wantErr:               false,
		},
		{
			name: "Merge deleted files [blockOverwrite = true]",
			args: args{
				blockOverwrite: true,
			},
			sourceFileStructure:   []string{"/file1" + unionFsHiddenPathSuffix},
			targetFileStructure:   []string{"/file1"},
			expectedFileStructure: []string{"/file1", "/file1" + unionFsHiddenPathSuffix},
			wantErr:               false,
		},
		{
			name: "Merge deleted folder [blockOverwrite = false]",
			args: args{
				blockOverwrite: false,
			},
			sourceFileStructure:   []string{"/folder1" + unionFsHiddenPathSuffix + "/"},
			targetFileStructure:   []string{"/folder1/"},
			expectedFileStructure: []string{"/folder1/", "/folder1" + unionFsHiddenPathSuffix + "/"},
			wantErr:               false,
		},
		{
			name: "Merge deleted folder [blockOverwrite = true]",
			args: args{
				blockOverwrite: true,
			},
			sourceFileStructure:   []string{"/folder1" + unionFsHiddenPathSuffix + "/"},
			targetFileStructure:   []string{"/folder1/"},
			expectedFileStructure: []string{"/folder1/", "/folder1" + unionFsHiddenPathSuffix + "/"},
			wantErr:               false,
		},

		// TODO: Add test cases.
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			sourceFolder := fileStructureCreator(testCase.sourceFileStructure)
			targetFolder := fileStructureCreator(testCase.targetFileStructure)

			t.Cleanup(func() {
				_ = os.RemoveAll(sourceFolder)
				_ = os.RemoveAll(targetFolder)
			})

			if err := dirmerger.MergeFolders([]string{sourceFolder}, targetFolder, testCase.args.blockOverwrite); (err != nil) != testCase.wantErr {
				t.Errorf("MergeFolders() error = %v, wantErr %v", err, testCase.wantErr)
			}

			checkFolderEquality(t, testCase, targetFolder)
		})
	}
}

func checkFolderEquality(t *testing.T, testCase mergeFoldersTestCase, targetFolder string) {
	t.Helper()

	t.Run("Check file structure equality", func(t *testing.T) {
		t.Parallel()

		actualFileStructure, err := collectPathsInFolder(targetFolder)
		if err != nil {
			t.Fatalf("Could not collect paths in folder %s: %v", targetFolder, err)
		}

		if len(testCase.expectedFileStructure) != len(actualFileStructure) {
			t.Fatalf("MergeFolders() expected %v, got %v", testCase.expectedFileStructure, actualFileStructure)
		}

		sort.Strings(testCase.expectedFileStructure)
		sort.Strings(actualFileStructure)

		for i := range testCase.expectedFileStructure {
			if testCase.expectedFileStructure[i] != actualFileStructure[i] {
				t.Errorf("MergeFolders() expected %v, got %v", testCase.expectedFileStructure, actualFileStructure)
			}
		}
	})
}

func collectPathsInFolder(targetFolder string) ([]string, error) {
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

func fileStructureCreator(filesAndFolders []string) string {
	targetFolder, _ := os.MkdirTemp("", "TestMergeFolders")

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
