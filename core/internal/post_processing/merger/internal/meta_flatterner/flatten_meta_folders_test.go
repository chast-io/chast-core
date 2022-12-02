package metaflatterner_test

import (
	wildcardstring "chast.io/core/internal/internal_util/wildcard_string"
	"os"
	"testing"

	metaflatterner "chast.io/core/internal/post_processing/merger/internal/meta_flatterner"
	testhelper "chast.io/core/internal/post_processing/merger/internal/test_helpers"
	"chast.io/core/internal/post_processing/merger/pkg/mergeoptions"
)

// region Test Helpers
type flattenMetaFolderTestCase struct {
	name                  string
	args                  mergeFoldersArgs
	sourceFileStructure   []string
	expectedFileStructure []string
	wantErr               bool
}

type mergeFoldersArgs struct {
	getMergeOptions func() *mergeoptions.MergeOptions
}

const unionFsMetaFolder = ".unionfs-fuse"
const unionFsHiddenPathSuffix = "_HIDDEN~"

// endregion
func TestFlattenMetaFolder(t *testing.T) {
	t.Parallel()

	tests := []flattenMetaFolderTestCase{
		{
			name: "should flatten meta folder [no conflict]",
			args: mergeFoldersArgs{
				getMergeOptions: mergeoptions.NewMergeOptions,
			},
			sourceFileStructure: []string{
				"/folder1/file1.txt",

				"/" + unionFsMetaFolder + "/folder2/file2.txt" + unionFsHiddenPathSuffix,
			},
			expectedFileStructure: []string{
				"/folder1/file1.txt",
				"/folder2/file2.txt" + unionFsHiddenPathSuffix,
			},
			wantErr: false,
		},
		{
			name: "should flatten meta folder [conflicting folder - blockOverwrite = true]",
			args: mergeFoldersArgs{
				getMergeOptions: func() *mergeoptions.MergeOptions {
					options := mergeoptions.NewMergeOptions()
					options.BlockOverwrite = true

					return options
				},
			},
			sourceFileStructure: []string{
				"/folder1/file1.txt",

				"/" + unionFsMetaFolder + "/folder1" + unionFsHiddenPathSuffix + "/",
				"/" + unionFsMetaFolder + "/folder1/file1.txt" + unionFsHiddenPathSuffix,
			},
			expectedFileStructure: []string{
				"/folder1/file1.txt",

				"/folder1" + unionFsHiddenPathSuffix + "/file1.txt" + unionFsHiddenPathSuffix,
			},
			wantErr: true,
		},
		{
			name: "should flatten meta folder [conflicting file no error - blockOverwrite = true]",
			args: mergeFoldersArgs{
				getMergeOptions: func() *mergeoptions.MergeOptions {
					options := mergeoptions.NewMergeOptions()
					options.BlockOverwrite = true

					return options
				},
			},
			sourceFileStructure: []string{
				"/folder1/file1.txt",

				"/" + unionFsMetaFolder + "/folder1/file1.txt" + unionFsHiddenPathSuffix,
			},
			expectedFileStructure: []string{
				"/folder1/file1.txt",

				"/folder1/file1.txt" + unionFsHiddenPathSuffix,
			},
			wantErr: false,
		},
		{
			name: "should flatten meta folder [conflicting file no error - blockOverwrite = false]",
			args: mergeFoldersArgs{
				getMergeOptions: func() *mergeoptions.MergeOptions {
					options := mergeoptions.NewMergeOptions()
					options.BlockOverwrite = false

					return options
				},
			},
			sourceFileStructure: []string{
				"/folder1/file1.txt",

				"/" + unionFsMetaFolder + "/folder1/file1.txt" + unionFsHiddenPathSuffix,
			},
			expectedFileStructure: []string{
				"/folder1/file1.txt",

				"/folder1/file1.txt" + unionFsHiddenPathSuffix,
			},
			wantErr: false,
		},
		{
			name: "should flatten meta folder [no conflict - dry run]",
			args: mergeFoldersArgs{
				getMergeOptions: func() *mergeoptions.MergeOptions {
					options := mergeoptions.NewMergeOptions()
					options.DryRun = true

					return options
				},
			},
			sourceFileStructure: []string{
				"/folder1/file1.txt",

				"/" + unionFsMetaFolder + "/folder2/file2.txt" + unionFsHiddenPathSuffix,
			},
			expectedFileStructure: []string{
				"/folder1/file1.txt",

				"/" + unionFsMetaFolder + "/folder2/file2.txt" + unionFsHiddenPathSuffix,
			},
			wantErr: false,
		},
		{
			name: "should flatten meta folder [conflicting folder - blockOverwrite = true - dry run]",
			args: mergeFoldersArgs{
				getMergeOptions: func() *mergeoptions.MergeOptions {
					options := mergeoptions.NewMergeOptions()
					options.DryRun = true
					options.BlockOverwrite = true

					return options
				},
			},
			sourceFileStructure: []string{
				"/folder1/file1.txt",

				"/" + unionFsMetaFolder + "/folder1" + unionFsHiddenPathSuffix + "/",
				"/" + unionFsMetaFolder + "/folder1/file1.txt" + unionFsHiddenPathSuffix,
			},
			expectedFileStructure: []string{
				"/folder1/file1.txt",

				"/" + unionFsMetaFolder + "/folder1" + unionFsHiddenPathSuffix + "/",
				"/" + unionFsMetaFolder + "/folder1/file1.txt" + unionFsHiddenPathSuffix,
			},
			wantErr: true,
		},
		{
			name: "should flatten meta folder [conflicting file no error - blockOverwrite = true - dry run]",
			args: mergeFoldersArgs{
				getMergeOptions: func() *mergeoptions.MergeOptions {
					options := mergeoptions.NewMergeOptions()
					options.DryRun = true
					options.BlockOverwrite = true

					return options
				},
			},
			sourceFileStructure: []string{
				"/folder1/file1.txt",

				"/" + unionFsMetaFolder + "/folder1/file1.txt" + unionFsHiddenPathSuffix,
			},
			expectedFileStructure: []string{
				"/folder1/file1.txt",

				"/" + unionFsMetaFolder + "/folder1/file1.txt" + unionFsHiddenPathSuffix,
			},
			wantErr: false,
		},
		{
			name: "should flatten meta folder [conflicting file no error - blockOverwrite = false - dry run]",
			args: mergeFoldersArgs{
				getMergeOptions: func() *mergeoptions.MergeOptions {
					options := mergeoptions.NewMergeOptions()
					options.DryRun = true
					options.BlockOverwrite = false

					return options
				},
			},
			sourceFileStructure: []string{
				"/folder1/file1.txt",

				"/" + unionFsMetaFolder + "/folder1/file1.txt" + unionFsHiddenPathSuffix,
			},
			expectedFileStructure: []string{
				"/folder1/file1.txt",

				"/" + unionFsMetaFolder + "/folder1/file1.txt" + unionFsHiddenPathSuffix,
			},
			wantErr: false,
		},
		{
			name: "should filter paths",
			args: mergeFoldersArgs{
				getMergeOptions: func() *mergeoptions.MergeOptions {
					options := mergeoptions.NewMergeOptions()
					options.Inclusions = []*wildcardstring.WildcardString{wildcardstring.NewWildcardString("/folder1")}
					options.Exclusions = []*wildcardstring.WildcardString{wildcardstring.NewWildcardString("/folder2")}

					return options
				},
			},
			sourceFileStructure: []string{
				"/" + unionFsMetaFolder + "/folder1" + unionFsHiddenPathSuffix + "/",
				"/" + unionFsMetaFolder + "/folder2" + unionFsHiddenPathSuffix + "/",
				"/" + unionFsMetaFolder + "/folder3" + unionFsHiddenPathSuffix + "/",
			},
			expectedFileStructure: []string{
				"/folder1" + unionFsHiddenPathSuffix + "/",
			},
			wantErr: false,
		},
	}

	for i := range tests {
		testCase := tests[i]

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			sourceFolder := testhelper.FileStructureCreator(testCase.sourceFileStructure, testCase.name)

			t.Cleanup(func() {
				_ = os.RemoveAll(sourceFolder)
			})

			if err := metaflatterner.FlattenMetaFolder(sourceFolder, testCase.args.getMergeOptions()); (err != nil) != testCase.wantErr {
				t.Errorf("FlattenMetaFolder() error = %v, wantErr %v", err, testCase.wantErr)
			}

			testhelper.CheckFolderEquality(t, testCase.expectedFileStructure, sourceFolder)
		})
	}
}
