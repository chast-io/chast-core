package dirmerger_test

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"chast.io/core/internal/post_processing/merger/pkg/dirmerger"
	"github.com/spf13/afero"
)

type mergeFoldersTestCase struct {
	name                  string
	args                  mergeFoldersArgs
	sourceFileStructure   []string
	targetFileStructure   []string
	expectedFileStructure []string
	wantErr               bool
}

type mergeFoldersArgs struct {
	getMergeOptions func() *dirmerger.MergeOptions
}

const unionFsHiddenPathSuffix = "_HIDDEN~"

func TestMergeFolders(t *testing.T) { //nolint:maintidx // Test function
	tests := []mergeFoldersTestCase{
		{
			name: "Merge two empty folders",
			args: mergeFoldersArgs{
				getMergeOptions: func() *dirmerger.MergeOptions {
					options := dirmerger.NewMergeOptions()

					return options
				},
			},
			sourceFileStructure: []string{},
			targetFileStructure: []string{},
			wantErr:             false,
		},
		{
			name: "Merge two folders with one file each [non conflicting]",
			args: mergeFoldersArgs{
				getMergeOptions: func() *dirmerger.MergeOptions {
					options := dirmerger.NewMergeOptions()

					return options
				},
			},
			sourceFileStructure:   []string{"/file1"},
			targetFileStructure:   []string{"/file2"},
			expectedFileStructure: []string{"/file1", "/file2"},
			wantErr:               false,
		},
		{
			name: "Merge two folders with one file each [conflicting - blockOverwrite = false]",
			args: mergeFoldersArgs{
				getMergeOptions: func() *dirmerger.MergeOptions {
					options := dirmerger.NewMergeOptions()
					options.BlockOverwrite = false

					return options
				},
			},
			sourceFileStructure:   []string{"/file1"},
			targetFileStructure:   []string{"/file1"},
			expectedFileStructure: []string{"/file1"},
			wantErr:               false,
		},
		{
			name: "Merge two folders with one file each [conflicting - blockOverwrite = true]",
			args: mergeFoldersArgs{
				getMergeOptions: func() *dirmerger.MergeOptions {
					options := dirmerger.NewMergeOptions()
					options.BlockOverwrite = true

					return options
				},
			},
			sourceFileStructure:   []string{"/file1"},
			targetFileStructure:   []string{"/file1"},
			expectedFileStructure: []string{"/file1"},
			wantErr:               true,
		},
		{
			name: "Merge deleted file with file [blockOverwrite = false]",
			args: mergeFoldersArgs{
				getMergeOptions: func() *dirmerger.MergeOptions {
					options := dirmerger.NewMergeOptions()
					options.BlockOverwrite = false

					return options
				},
			},
			sourceFileStructure:   []string{"/file1" + unionFsHiddenPathSuffix},
			targetFileStructure:   []string{"/file1"},
			expectedFileStructure: []string{"/file1" + unionFsHiddenPathSuffix},
			wantErr:               false,
		},
		{
			name: "Merge file with deleted file [blockOverwrite = false]",
			args: mergeFoldersArgs{
				getMergeOptions: func() *dirmerger.MergeOptions {
					options := dirmerger.NewMergeOptions()
					options.BlockOverwrite = false

					return options
				},
			},
			sourceFileStructure:   []string{"/file1"},
			targetFileStructure:   []string{"/file1" + unionFsHiddenPathSuffix},
			expectedFileStructure: []string{"/file1"},
			wantErr:               false,
		},
		{
			name: "Merge deleted file with file [blockOverwrite = true]",
			args: mergeFoldersArgs{
				getMergeOptions: func() *dirmerger.MergeOptions {
					options := dirmerger.NewMergeOptions()
					options.BlockOverwrite = true

					return options
				},
			},
			sourceFileStructure:   []string{"/file1" + unionFsHiddenPathSuffix},
			targetFileStructure:   []string{"/file1"},
			expectedFileStructure: []string{"/file1"},
			wantErr:               true,
		},
		{
			name: "Merge file with deleted file [blockOverwrite = true]",
			args: mergeFoldersArgs{
				getMergeOptions: func() *dirmerger.MergeOptions {
					options := dirmerger.NewMergeOptions()
					options.BlockOverwrite = true

					return options
				},
			},
			sourceFileStructure:   []string{"/file1"},
			targetFileStructure:   []string{"/file1" + unionFsHiddenPathSuffix},
			expectedFileStructure: []string{"/file1" + unionFsHiddenPathSuffix},
			wantErr:               true,
		},
		{
			name: "Merge deleted folder with folder [blockOverwrite = false]",
			args: mergeFoldersArgs{
				getMergeOptions: func() *dirmerger.MergeOptions {
					options := dirmerger.NewMergeOptions()
					options.BlockOverwrite = false

					return options
				},
			},
			sourceFileStructure:   []string{"/folder1" + unionFsHiddenPathSuffix + "/"},
			targetFileStructure:   []string{"/folder1/"},
			expectedFileStructure: []string{"/folder1" + unionFsHiddenPathSuffix + "/"},
			wantErr:               false,
		},
		{
			name: "Merge folder with deleted folder [blockOverwrite = false]",
			args: mergeFoldersArgs{
				getMergeOptions: func() *dirmerger.MergeOptions {
					options := dirmerger.NewMergeOptions()
					options.BlockOverwrite = false

					return options
				},
			},
			sourceFileStructure:   []string{"/folder1/"},
			targetFileStructure:   []string{"/folder1" + unionFsHiddenPathSuffix + "/"},
			expectedFileStructure: []string{"/folder1/"},
			wantErr:               false,
		},
		{
			name: "Merge deleted folder with folder [blockOverwrite = true]",
			args: mergeFoldersArgs{
				getMergeOptions: func() *dirmerger.MergeOptions {
					options := dirmerger.NewMergeOptions()
					options.BlockOverwrite = true

					return options
				},
			},
			sourceFileStructure:   []string{"/folder1" + unionFsHiddenPathSuffix + "/"},
			targetFileStructure:   []string{"/folder1/"},
			expectedFileStructure: []string{"/folder1/"},
			wantErr:               true,
		},
		{
			name: "Merge folder with deleted folder [blockOverwrite = true]",
			args: mergeFoldersArgs{
				getMergeOptions: func() *dirmerger.MergeOptions {
					options := dirmerger.NewMergeOptions()
					options.BlockOverwrite = true

					return options
				},
			},
			sourceFileStructure:   []string{"/folder1/"},
			targetFileStructure:   []string{"/folder1" + unionFsHiddenPathSuffix + "/"},
			expectedFileStructure: []string{"/folder1" + unionFsHiddenPathSuffix + "/"},
			wantErr:               true,
		},
		{
			name: "Delete folders that are marked as deleted after merge",
			args: mergeFoldersArgs{
				getMergeOptions: func() *dirmerger.MergeOptions {
					options := dirmerger.NewMergeOptions()
					options.DeleteMarkedAsDeletedPaths = true

					return options
				},
			},
			sourceFileStructure:   []string{"/folder1" + unionFsHiddenPathSuffix + "/"},
			targetFileStructure:   []string{"/folder1/"},
			expectedFileStructure: []string{},
			wantErr:               false,
		},
		{
			name: "Delete files that are marked as deleted after merge",
			args: mergeFoldersArgs{
				getMergeOptions: func() *dirmerger.MergeOptions {
					options := dirmerger.NewMergeOptions()
					options.DeleteMarkedAsDeletedPaths = true

					return options
				},
			},
			sourceFileStructure:   []string{"/folder1/file1" + unionFsHiddenPathSuffix},
			targetFileStructure:   []string{"/folder1/file1"},
			expectedFileStructure: []string{"/folder1/"},
			wantErr:               false,
		},
		{
			name: "Delete empty folders",
			args: mergeFoldersArgs{
				getMergeOptions: func() *dirmerger.MergeOptions {
					options := dirmerger.NewMergeOptions()
					options.DeleteEmptyFolders = true

					return options
				},
			},
			sourceFileStructure: []string{
				"/folder1/folder1/",
				"/folder1/folder2/file1",
			},
			targetFileStructure: []string{"/folder2/file1"},
			expectedFileStructure: []string{
				"/folder1/folder2/file1",
				"/folder2/file1",
			},
			wantErr: false,
		},
		{
			name: "Delete empty folders except root",
			args: mergeFoldersArgs{
				getMergeOptions: func() *dirmerger.MergeOptions {
					options := dirmerger.NewMergeOptions()
					options.DeleteEmptyFolders = true

					return options
				},
			},
			sourceFileStructure: []string{
				"/folder1/folder1/",
			},
			targetFileStructure: []string{
				"/folder1/folder2/",
			},
			expectedFileStructure: make([]string, 0),
			wantErr:               false,
		},
		{
			name: "Delete folders that are marked as deleted after merge and delete empty folders",
			args: mergeFoldersArgs{
				getMergeOptions: func() *dirmerger.MergeOptions {
					options := dirmerger.NewMergeOptions()
					options.DeleteMarkedAsDeletedPaths = true
					options.DeleteEmptyFolders = true

					return options
				},
			},
			sourceFileStructure: []string{
				"/folder1/folder1/",
				"/folder1/folder2" + unionFsHiddenPathSuffix + "/",
				"/folder1/folder3/file1",
			},
			targetFileStructure: []string{
				"/folder1/folder2/file1",
			},
			expectedFileStructure: []string{
				"/folder1/folder3/file1",
			},
			wantErr: false,
		},

		// TODO: Skip locations
	}

	t.Parallel()

	for i := range tests {
		testCase := tests[i]
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			sourceFolder := fileStructureCreator(testCase.sourceFileStructure)
			targetFolder := fileStructureCreator(testCase.targetFileStructure)

			t.Cleanup(func() {
				_ = os.RemoveAll(sourceFolder)
				_ = os.RemoveAll(targetFolder)
			})

			mergeEntities := dirmerger.NewMergeEntity(sourceFolder, nil)

			if err := dirmerger.MergeFolders([]dirmerger.MergeEntity{mergeEntities}, targetFolder, testCase.args.getMergeOptions()); (err != nil) != testCase.wantErr {
				t.Errorf("MergeFolders() error = %v, wantErr %v", err, testCase.wantErr)
			}

			checkFolderEquality(t, testCase.expectedFileStructure, targetFolder)
		})
	}

	t.Run("Copy mode", func(t *testing.T) {
		t.Parallel()

		t.Run("Should copy folders", func(t *testing.T) {
			t.Parallel()
			options := dirmerger.NewMergeOptions()
			options.CopyMode = true

			sourceFileStructure := []string{
				"/folder1/folder1/",
			}
			targetFileStructure := make([]string, 0)
			expectedFileStructure := []string{
				"/folder1/folder1/",
			}
			expectedSourceFileStructure := []string{
				"/folder1/folder1/",
			}

			sourceFolder := fileStructureCreator(sourceFileStructure)
			targetFolder := fileStructureCreator(targetFileStructure)

			t.Cleanup(func() {
				_ = os.RemoveAll(sourceFolder)
				_ = os.RemoveAll(targetFolder)
			})

			mergeEntities := dirmerger.NewMergeEntity(sourceFolder, nil)

			if err := dirmerger.MergeFolders([]dirmerger.MergeEntity{mergeEntities}, targetFolder, options); err != nil {
				t.Errorf("MergeFolders() error = %v, wantErr %v", err, false)
			}

			checkFolderEquality(t, expectedSourceFileStructure, sourceFolder)
			checkFolderEquality(t, expectedFileStructure, targetFolder)
		})
		t.Run("Should copy files", func(t *testing.T) {
			t.Parallel()
			options := dirmerger.NewMergeOptions()
			options.CopyMode = true

			sourceFileStructure := []string{
				"/folder1/folder1/file",
			}
			targetFileStructure := make([]string, 0)
			expectedFileStructure := []string{
				"/folder1/folder1/file",
			}
			expectedSourceFileStructure := []string{
				"/folder1/folder1/file",
			}

			sourceFolder := fileStructureCreator(sourceFileStructure)
			targetFolder := fileStructureCreator(targetFileStructure)

			t.Cleanup(func() {
				_ = os.RemoveAll(sourceFolder)
				_ = os.RemoveAll(targetFolder)
			})

			mergeEntities := dirmerger.NewMergeEntity(sourceFolder, nil)

			if err := dirmerger.MergeFolders([]dirmerger.MergeEntity{mergeEntities}, targetFolder, options); err != nil {
				t.Errorf("MergeFolders() error = %v, wantErr %v", err, false)
			}

			checkFolderEquality(t, expectedSourceFileStructure, sourceFolder)
			checkFolderEquality(t, expectedFileStructure, targetFolder)
		})
	})
}

func checkFolderEquality(t *testing.T, expectedFileStructure []string, checkFolder string) {
	t.Helper()

	t.Run("Check file structure equality", func(t *testing.T) {
		t.Parallel()

		actualFileStructure, err := collectPathsInFolder(checkFolder)
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

type areMergeableTestCase struct {
	name                string
	args                areMergeableArgs
	sourceFileStructure []string
	targetFileStructure []string
	expectedMergeable   bool
	wantErr             bool
}
type areMergeableArgs struct {
	getMergeOptions func() *dirmerger.MergeOptions
}

func TestAreMergeable(t *testing.T) {
	tests := []areMergeableTestCase{
		{
			name:                "Merge two non conflicting files [blockOverwrite = false]",
			sourceFileStructure: []string{"folder1/file1"},
			targetFileStructure: []string{"folder1/file2"},
			expectedMergeable:   true,
			wantErr:             false,
			args: areMergeableArgs{
				getMergeOptions: func() *dirmerger.MergeOptions {
					options := dirmerger.NewMergeOptions()
					options.BlockOverwrite = false

					return options
				},
			},
		},
		{
			name:                "Merge two conflicting files [blockOverwrite = false]",
			sourceFileStructure: []string{"folder1/file1"},
			targetFileStructure: []string{"folder1/file1"},
			expectedMergeable:   true,
			wantErr:             false,
			args: areMergeableArgs{
				getMergeOptions: func() *dirmerger.MergeOptions {
					options := dirmerger.NewMergeOptions()
					options.BlockOverwrite = false

					return options
				},
			},
		},
		{
			name:                "Merge deleted file with file [blockOverwrite = false]",
			sourceFileStructure: []string{"folder1/file1" + unionFsHiddenPathSuffix},
			targetFileStructure: []string{"folder1/file1"},
			expectedMergeable:   true,
			wantErr:             false,
			args: areMergeableArgs{
				getMergeOptions: func() *dirmerger.MergeOptions {
					options := dirmerger.NewMergeOptions()
					options.BlockOverwrite = false

					return options
				},
			},
		},
		{
			name:                "Merge file with deleted file [blockOverwrite = false]",
			sourceFileStructure: []string{"folder1/file1"},
			targetFileStructure: []string{"folder1/file1" + unionFsHiddenPathSuffix},
			expectedMergeable:   true,
			wantErr:             false,
			args: areMergeableArgs{
				getMergeOptions: func() *dirmerger.MergeOptions {
					options := dirmerger.NewMergeOptions()
					options.BlockOverwrite = false

					return options
				},
			},
		},
		// === blockOverwrite = true ===
		{
			name:                "Merge two non conflicting files [blockOverwrite = true]",
			sourceFileStructure: []string{"folder1/file1"},
			targetFileStructure: []string{"folder1/file2"},
			expectedMergeable:   true,
			wantErr:             false,
			args: areMergeableArgs{
				getMergeOptions: func() *dirmerger.MergeOptions {
					options := dirmerger.NewMergeOptions()
					options.BlockOverwrite = true

					return options
				},
			},
		},
		{
			name:                "Merge two conflicting files [blockOverwrite = true]",
			sourceFileStructure: []string{"folder1/file1"},
			targetFileStructure: []string{"folder1/file1"},
			expectedMergeable:   false,
			wantErr:             false,
			args: areMergeableArgs{
				getMergeOptions: func() *dirmerger.MergeOptions {
					options := dirmerger.NewMergeOptions()
					options.BlockOverwrite = true

					return options
				},
			},
		},
		{
			name:                "Merge deleted file with file [blockOverwrite = true]",
			sourceFileStructure: []string{"folder1/file1" + unionFsHiddenPathSuffix},
			targetFileStructure: []string{"folder1/file1"},
			expectedMergeable:   false,
			wantErr:             false,
			args: areMergeableArgs{
				getMergeOptions: func() *dirmerger.MergeOptions {
					options := dirmerger.NewMergeOptions()
					options.BlockOverwrite = true

					return options
				},
			},
		},
		{
			name:                "Merge file with deleted file [blockOverwrite = true]",
			sourceFileStructure: []string{"folder1/file1"},
			targetFileStructure: []string{"folder1/file1" + unionFsHiddenPathSuffix},
			expectedMergeable:   false,
			wantErr:             false,
			args: areMergeableArgs{
				getMergeOptions: func() *dirmerger.MergeOptions {
					options := dirmerger.NewMergeOptions()
					options.BlockOverwrite = true

					return options
				},
			},
		},
		{
			name:                "Overwrites dry run option",
			sourceFileStructure: []string{"folder1/file1"},
			targetFileStructure: []string{"folder1/file2"},
			expectedMergeable:   true,
			wantErr:             false,
			args: areMergeableArgs{
				getMergeOptions: func() *dirmerger.MergeOptions {
					options := dirmerger.NewMergeOptions()
					options.DryRun = false

					return options
				},
			},
		},
	}

	t.Parallel()

	for i := range tests {
		testCase := tests[i]
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			sourceFolder := fileStructureCreator(testCase.sourceFileStructure)
			targetFolder := fileStructureCreator(testCase.targetFileStructure)

			sourceFileStructure, _ := collectPathsInFolder(sourceFolder)
			targetFileStructure, _ := collectPathsInFolder(targetFolder)

			t.Cleanup(func() {
				_ = os.RemoveAll(sourceFolder)
				_ = os.RemoveAll(targetFolder)
			})

			mergeEntities := dirmerger.NewMergeEntity(sourceFolder, nil)

			mergeable, err := dirmerger.AreMergeable([]dirmerger.MergeEntity{mergeEntities}, targetFolder, testCase.args.getMergeOptions())
			if (err != nil) != testCase.wantErr {
				t.Errorf("MergeFolders() error = %v, wantErr %v", err, testCase.wantErr)
			}

			if mergeable != testCase.expectedMergeable {
				t.Errorf("MergeFolders() expected %v, got %v", testCase.expectedMergeable, mergeable)
			}

			checkFolderEquality(t, sourceFileStructure, sourceFolder)
			checkFolderEquality(t, targetFileStructure, targetFolder)
		})
	}

	t.Run("Check if options are not modified", func(t *testing.T) {
		t.Parallel()

		sourceFolder := fileStructureCreator([]string{"folder1/file1"})
		targetFolder := fileStructureCreator([]string{"folder1/file2"})

		t.Cleanup(func() {
			_ = os.RemoveAll(sourceFolder)
			_ = os.RemoveAll(targetFolder)
		})

		options := dirmerger.NewMergeOptions()
		options.DryRun = false
		options.BlockOverwrite = true
		options.DeleteMarkedAsDeletedPaths = true

		mergeEntities := dirmerger.NewMergeEntity(sourceFolder, nil)

		_, _ = dirmerger.AreMergeable([]dirmerger.MergeEntity{mergeEntities}, targetFolder, options)

		if options.DryRun != false || options.BlockOverwrite != true || options.DeleteMarkedAsDeletedPaths != true {
			t.Errorf("MergeFolders() options were modified")
		}
	})
}
