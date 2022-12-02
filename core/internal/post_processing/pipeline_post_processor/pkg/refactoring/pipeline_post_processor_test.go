package pipelinepostprocessor_test

import (
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"

	chastlog "chast.io/core/internal/logger"
	refactoringpipelinemodel "chast.io/core/internal/pipeline/pkg/model/refactoring"
	uut "chast.io/core/internal/post_processing/pipeline_post_processor/pkg/refactoring"
	steppostprocessor "chast.io/core/internal/post_processing/step_post_processor/pkg/refactoring"
	"chast.io/core/internal/run_model/pkg/model/refactoring"
	"github.com/spf13/afero"
)

// region Helpers
func dummyPipeline(executionGroups []*refactoringpipelinemodel.ExecutionGroup) *refactoringpipelinemodel.Pipeline {
	temp, err := os.MkdirTemp("", "pipeline_post_processor_test_dummyPipeline")
	if err != nil {
		panic(err)
	}

	pipeline := refactoringpipelinemodel.NewPipeline(
		filepath.Join(temp, "operationLocation"),
		filepath.Join(temp, "changeCaptureLocation"),
		filepath.Join(temp, "rootFileSystemLocation"))

	for _, executionGroup := range executionGroups {
		pipeline.AddExecutionGroup(executionGroup)
	}

	return pipeline
}

func dummyExecutionGroup(steps []*refactoringpipelinemodel.Step) *refactoringpipelinemodel.ExecutionGroup {
	stage := refactoringpipelinemodel.NewExecutionGroup()

	for _, step := range steps {
		stage.AddStep(step)
	}

	return stage
}

func dummyStep(nr int) *refactoringpipelinemodel.Step {
	runModel := &refactoring.SingleRunModel{
		Run: &refactoring.Run{
			ID:                 "runId" + strconv.Itoa(nr),
			Dependencies:       make([]*refactoring.Run, 0),
			SupportedLanguages: []string{"java"},
			Docker:             &refactoring.Docker{},          //nolint:exhaustruct // not required for test
			Local:              &refactoring.Local{},           //nolint:exhaustruct // not required for test
			Command:            &refactoring.Command{},         //nolint:exhaustruct // not required for test
			ChangeLocations:    &refactoring.ChangeLocations{}, //nolint:exhaustruct // not required for test
		},
	}

	return refactoringpipelinemodel.NewStep(runModel)
}

func createFolders(pipeline *refactoringpipelinemodel.Pipeline) {
	_ = os.MkdirAll(filepath.Join(pipeline.RootFileSystemLocation, "boot"), os.ModePerm)
	_ = os.MkdirAll(filepath.Join(pipeline.RootFileSystemLocation, "dev"), os.ModePerm)
	_ = os.MkdirAll(filepath.Join(pipeline.RootFileSystemLocation, "etc"), os.ModePerm)
	_ = os.MkdirAll(filepath.Join(pipeline.RootFileSystemLocation, "home"), os.ModePerm)
	_ = os.MkdirAll(filepath.Join(pipeline.RootFileSystemLocation, "tmp"), os.ModePerm)
	_ = os.MkdirAll(filepath.Join(pipeline.RootFileSystemLocation, "var"), os.ModePerm)

	for _, stage := range pipeline.ExecutionGroups {
		for _, step := range stage.Steps {
			_ = os.MkdirAll(step.OperationLocation, os.ModePerm)
			_ = os.MkdirAll(step.ChangeCaptureLocation, os.ModePerm)
		}
	}
}

func createSubPaths(pipeline *refactoringpipelinemodel.Pipeline, paths [][]string) {
	executionGroupIndex := 0
	stepIndex := 0

	for _, executionGroup := range paths {
		for _, path := range executionGroup {
			location := filepath.Join(pipeline.ExecutionGroups[executionGroupIndex].Steps[stepIndex].ChangeCaptureLocation, path)
			if strings.HasSuffix(path, "/") {
				_ = os.MkdirAll(location, os.ModePerm)
			} else {
				_ = os.MkdirAll(filepath.Dir(location), os.ModePerm)
				_, _ = os.Create(location)
			}
		}

		stepIndex++
		if stepIndex >= len(pipeline.ExecutionGroups[executionGroupIndex].Steps) {
			stepIndex = 0
			executionGroupIndex++
		}
	}
}

func cleanupCleanup(pipeline *refactoringpipelinemodel.Pipeline) {
	parent := filepath.Dir(pipeline.RootFileSystemLocation)
	if filepath.Dir(parent) != os.TempDir() {
		panic("cleanupCleanup() should only be used in tests and only on temporary folders")
	}

	_ = os.RemoveAll(parent)
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

// endregion

// region CleanupPipeline

//nolint:gocognit
func TestProcess(t *testing.T) {
	t.Parallel()

	type args struct {
		getPipeline func() *refactoringpipelinemodel.Pipeline
	}

	tests := []struct {
		name                  string
		args                  args
		wantErr               bool
		changedPaths          [][]string
		expectedFileStructure []string
	}{
		{
			name: "should fail if pipeline is nil",
			args: args{
				getPipeline: func() *refactoringpipelinemodel.Pipeline {
					return nil
				},
			},
			wantErr:               true,
			expectedFileStructure: nil,
		},
		{
			name: "should merge steps folders and files",
			args: args{
				getPipeline: func() *refactoringpipelinemodel.Pipeline {
					return dummyPipeline(
						[]*refactoringpipelinemodel.ExecutionGroup{
							dummyExecutionGroup(
								[]*refactoringpipelinemodel.Step{
									dummyStep(1),
									dummyStep(2),
								},
							),
							dummyExecutionGroup(
								[]*refactoringpipelinemodel.Step{
									dummyStep(3),
									dummyStep(4),
								},
							),
							dummyExecutionGroup(
								[]*refactoringpipelinemodel.Step{
									dummyStep(5),
									dummyStep(6),
								},
							),
						},
					)
				},
			},
			wantErr: false,
			changedPaths: [][]string{
				{"/folder1/file.go"},
				{"/folder2/file.go"},

				{"/folder3/file.go"},
				{"/folder4/file.go"},

				{"/folder5/file.go"},
				{"/folder6/file.go"},
			},
			expectedFileStructure: []string{
				"/folder1/file.go",
				"/folder2/file.go",

				"/folder3/file.go",
				"/folder4/file.go",

				"/folder5/file.go",
				"/folder6/file.go",
			},
		},
		{
			name: "should merge steps and keep marked as deleted paths",
			args: args{
				getPipeline: func() *refactoringpipelinemodel.Pipeline {
					step1 := dummyStep(1)
					step2 := dummyStep(2)
					step3 := dummyStep(3)
					step4 := dummyStep(4)
					step5 := dummyStep(5)
					step6 := dummyStep(6)

					step3.AddDependency(step2)
					step4.AddDependency(step1)

					step5.AddDependency(step4)
					step6.AddDependency(step3)

					return dummyPipeline(
						[]*refactoringpipelinemodel.ExecutionGroup{
							dummyExecutionGroup(
								[]*refactoringpipelinemodel.Step{step1, step2},
							),
							dummyExecutionGroup(
								[]*refactoringpipelinemodel.Step{step3, step4},
							),
							dummyExecutionGroup(
								[]*refactoringpipelinemodel.Step{step5, step6},
							),
						},
					)
				},
			},
			wantErr: false,
			changedPaths: [][]string{
				{"/folder1/file.go"},
				{"/folder2/file.go"},

				{"/folder3/file.go"},
				{"/.unionfs-fuse/folder1_HIDDEN~/"},

				{"/folder5/file.go"},
				{"/.unionfs-fuse/folder3_HIDDEN~/"},
			},
			expectedFileStructure: []string{
				"/folder1_HIDDEN~/",
				"/folder2/file.go",

				"/folder3_HIDDEN~/",

				"/folder5/file.go",
			},
		},
		{
			name: "should return error if the same file is edited in same stage",
			args: args{
				getPipeline: func() *refactoringpipelinemodel.Pipeline {
					return dummyPipeline(
						[]*refactoringpipelinemodel.ExecutionGroup{
							dummyExecutionGroup(
								[]*refactoringpipelinemodel.Step{
									dummyStep(1),
									dummyStep(2),
								},
							),
						},
					)
				},
			},
			wantErr: true,
			changedPaths: [][]string{
				{"/folder1/file.go"},
				{"/folder1/file.go"},
			},
			expectedFileStructure: nil,
		},
		{
			name: "should return error if the same file is edited and deleted in same stage",
			args: args{
				getPipeline: func() *refactoringpipelinemodel.Pipeline {
					return dummyPipeline(
						[]*refactoringpipelinemodel.ExecutionGroup{
							dummyExecutionGroup(
								[]*refactoringpipelinemodel.Step{
									dummyStep(1),
									dummyStep(2),
								},
							),
						},
					)
				},
			},
			wantErr: true,
			changedPaths: [][]string{
				{"/folder1/file.go"},
				{"/folder1/file.go_HIDDEN~"},
			},
			expectedFileStructure: nil,
		},
		{
			name: "should be able to create deleted path",
			args: args{
				getPipeline: func() *refactoringpipelinemodel.Pipeline {
					step1 := dummyStep(1)
					step2 := dummyStep(2)
					step3 := dummyStep(3)

					step2.AddDependency(step1)
					step3.AddDependency(step2)

					return dummyPipeline(
						[]*refactoringpipelinemodel.ExecutionGroup{
							dummyExecutionGroup(
								[]*refactoringpipelinemodel.Step{step1},
							),
							dummyExecutionGroup(
								[]*refactoringpipelinemodel.Step{step2},
							),
							dummyExecutionGroup(
								[]*refactoringpipelinemodel.Step{step3},
							),
						},
					)
				},
			},
			wantErr: false,
			changedPaths: [][]string{
				{"/folder1/file.go"},

				{"/.unionfs-fuse/folder1_HIDDEN~/"},

				{"/folder1/"},
			},
			expectedFileStructure: []string{
				"/folder1/",
			},
		},
	}

	for i := range tests {
		testCase := tests[i]

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			pipeline := testCase.args.getPipeline()
			if pipeline != nil {
				createFolders(pipeline)
				if testCase.changedPaths != nil {
					createSubPaths(pipeline, testCase.changedPaths)
				}

				t.Cleanup(func() { cleanupCleanup(pipeline) })

				for _, executionGroup := range pipeline.ExecutionGroups {
					for _, step := range executionGroup.Steps {
						if err := steppostprocessor.Process(step); err != nil {
							t.Fatalf("unexpected error: %v", err)
						}
					}
				}
			}

			if err := uut.Process(pipeline); (err != nil) != testCase.wantErr {
				t.Errorf("CleanupPipeline() error = %v, wantErr %v", err, testCase.wantErr)
			} else if err != nil {
				chastlog.Log.Debugf("Reported error: %v", err)
			}

			if pipeline == nil {
				return
			}

			if testCase.expectedFileStructure != nil {
				checkFolderEquality(t, testCase.expectedFileStructure, pipeline.GetFinalChangeCaptureLocation())
			}

			operationLocationsShouldBeEmpty(t, pipeline)
			rootFileSystemShouldBeUnaltered(t, pipeline)
		})
	}
}

// endregion

// region CleanupStep

func TestUnionFSCleanupStep(t *testing.T) {
	t.Parallel()
	// TODO cases where meta folders are not initially merged
}

// endregion

// region Test Helpers

func checkFolderEquality(t *testing.T, expectedFileStructure []string, checkFolder string) {
	t.Helper()

	t.Run("Check file structure equality", func(t *testing.T) {
		t.Parallel()

		actualFileStructure, err := collectPathsInFolder(checkFolder)
		if err != nil {
			t.Fatalf("Could not collect paths in folder %s: %v", checkFolder, err)
		}

		if len(expectedFileStructure) != len(actualFileStructure) {
			t.Fatalf("Expected %v paths, got %v:\n%v\n%v",
				len(expectedFileStructure), len(actualFileStructure),
				expectedFileStructure, actualFileStructure,
			)
		}

		sort.Strings(expectedFileStructure)
		sort.Strings(actualFileStructure)

		for i := range expectedFileStructure {
			if expectedFileStructure[i] != actualFileStructure[i] {
				t.Errorf("Expected \n%v\n, got\n%v", expectedFileStructure, actualFileStructure)
			}
		}
	})
}

func operationLocationsShouldBeEmpty(t *testing.T, pipeline *refactoringpipelinemodel.Pipeline) {
	t.Helper()

	t.Run("Operation locations should be empty", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewOsFs()

		empty, err := afero.IsEmpty(fs, pipeline.OperationLocation)
		if err != nil {
			t.Fatalf("Could not check if operation location is empty: %v", err)
		}

		if !empty {
			t.Errorf("Operation locations should be empty")
		}
	})
}

func rootFileSystemShouldBeUnaltered(t *testing.T, pipeline *refactoringpipelinemodel.Pipeline) {
	t.Helper()

	t.Run("Root file system should be unaltered", func(t *testing.T) {
		t.Parallel()

		checkFolderEquality(t, []string{
			"/boot/",
			"/dev/",
			"/etc/",
			"/home/",
			"/tmp/",
			"/var/",
		}, pipeline.RootFileSystemLocation)
	})
}

// endregion
