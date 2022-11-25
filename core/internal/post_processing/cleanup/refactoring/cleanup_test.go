package refactoringpipelinecleanup_test

import (
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"

	chastlog "chast.io/core/internal/logger"
	refactoringpipelinemodel "chast.io/core/internal/pipeline/pkg/model/refactoring"
	uut "chast.io/core/internal/post_processing/cleanup/refactoring"
	"chast.io/core/internal/run_model/pkg/model/refactoring"
	"github.com/spf13/afero"
)

// region Helpers
func dummyPipeline(stages []*refactoringpipelinemodel.Stage) *refactoringpipelinemodel.Pipeline {
	temp, err := os.MkdirTemp("", "cleanup_test_dummyPipeline")
	if err != nil {
		panic(err)
	}

	pipeline := refactoringpipelinemodel.NewPipeline(
		filepath.Join(temp, "operationLocation"),
		filepath.Join(temp, "changeCaptureLocation"),
		filepath.Join(temp, "rootFileSystemLocation"))

	for _, stage := range stages {
		pipeline.AddStage(stage)
	}

	return pipeline
}

func dummyStage(nr int, steps []*refactoringpipelinemodel.Step) *refactoringpipelinemodel.Stage {
	stage := refactoringpipelinemodel.NewStage("stage_name" + strconv.Itoa(nr))

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
			Docker:             &refactoring.Docker{},  //nolint:exhaustruct // not required for test
			Local:              &refactoring.Local{},   //nolint:exhaustruct // not required for test
			Command:            &refactoring.Command{}, //nolint:exhaustruct // not required for test
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

	for _, stage := range pipeline.Stages {
		for _, step := range stage.Steps {
			_ = os.MkdirAll(step.OperationLocation, os.ModePerm)
			_ = os.MkdirAll(step.ChangeCaptureLocation, os.ModePerm)
		}
	}
}

func createSubPaths(pipeline *refactoringpipelinemodel.Pipeline, paths [][]string) {
	stageIndex := 0
	stepIndex := 0

	for _, stageStep := range paths {
		for _, path := range stageStep {
			location := filepath.Join(pipeline.Stages[stageIndex].Steps[stepIndex].ChangeCaptureLocation, path)
			if strings.HasSuffix(path, "/") {
				_ = os.MkdirAll(location, os.ModePerm)
			} else {
				_ = os.MkdirAll(filepath.Dir(location), os.ModePerm)
				_, _ = os.Create(location)
			}
		}

		stepIndex++
		if stepIndex >= len(pipeline.Stages[stageIndex].Steps) {
			stepIndex = 0
			stageIndex++
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

func TestCleanupPipeline(t *testing.T) {
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
				getPipeline: nil,
			},
			wantErr:               true,
			expectedFileStructure: nil,
		},
		{
			name: "should merge stages and steps folders and files",
			args: args{
				getPipeline: func() *refactoringpipelinemodel.Pipeline {
					return dummyPipeline(
						[]*refactoringpipelinemodel.Stage{
							dummyStage(1,
								[]*refactoringpipelinemodel.Step{
									dummyStep(1),
									dummyStep(2),
								},
							),
							dummyStage(2,
								[]*refactoringpipelinemodel.Step{
									dummyStep(3),
									dummyStep(4),
								},
							),
							dummyStage(3,
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
			name: "should merge stages and steps and delete empty folders",
			args: args{
				getPipeline: func() *refactoringpipelinemodel.Pipeline {
					return dummyPipeline(
						[]*refactoringpipelinemodel.Stage{
							dummyStage(1,
								[]*refactoringpipelinemodel.Step{
									dummyStep(1),
									dummyStep(2),
								},
							),
							dummyStage(2,
								[]*refactoringpipelinemodel.Step{
									dummyStep(3),
									dummyStep(4),
								},
							),
							dummyStage(3,
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

				{"/folder3/"},
				{"/folder4/"},

				{"/folder5/file.go"},
				{"/folder6/file.go"},
			},
			expectedFileStructure: []string{
				"/folder1/file.go",
				"/folder2/file.go",

				"/folder5/file.go",
				"/folder6/file.go",
			},
		},
		{
			name: "should merge stages and steps and keep marked as deleted paths",
			args: args{
				getPipeline: func() *refactoringpipelinemodel.Pipeline {
					return dummyPipeline(
						[]*refactoringpipelinemodel.Stage{
							dummyStage(1,
								[]*refactoringpipelinemodel.Step{
									dummyStep(1),
									dummyStep(2),
								},
							),
							dummyStage(2,
								[]*refactoringpipelinemodel.Step{
									dummyStep(3),
									dummyStep(4),
								},
							),
							dummyStage(3,
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
						[]*refactoringpipelinemodel.Stage{
							dummyStage(1,
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
						[]*refactoringpipelinemodel.Stage{
							dummyStage(1,
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
	}

	for i := range tests {
		testCase := tests[i]

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			var pipeline *refactoringpipelinemodel.Pipeline

			if testCase.args.getPipeline != nil {
				pipeline = testCase.args.getPipeline()

				createFolders(pipeline)
				if testCase.changedPaths != nil {
					createSubPaths(pipeline, testCase.changedPaths)
				}

				t.Cleanup(func() { cleanupCleanup(pipeline) })
			}

			if err := uut.CleanupPipeline(pipeline); (err != nil) != testCase.wantErr {
				t.Errorf("CleanupPipeline() error = %v, wantErr %v", err, testCase.wantErr)
			} else if err != nil {
				chastlog.Log.Debugf("Reported error: %v", err)
			}

			if pipeline == nil {
				return
			}

			if testCase.expectedFileStructure != nil {
				checkFolderEquality(t, testCase.expectedFileStructure, pipeline.ChangeCaptureLocation)
			}

			operationLocationsShouldBeDeleted(t, pipeline)
			rootFileSystemShouldBeUnaltered(t, pipeline)
		})
	}
}

// endregion

// region CleanupStage

func TestCleanupStage(t *testing.T) {
	t.Parallel()

	type args struct {
		stage *refactoringpipelinemodel.Stage
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}

	for i := range tests {
		testCase := tests[i]

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			if err := uut.CleanupStage(testCase.args.stage); (err != nil) != testCase.wantErr {
				t.Errorf("CleanupStage() error = %v, wantErr %v", err, testCase.wantErr)
			}
		})
	}
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

func operationLocationsShouldBeDeleted(t *testing.T, pipeline *refactoringpipelinemodel.Pipeline) {
	t.Helper()

	t.Run("Operation locations should be empty", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewOsFs()

		exists, err := afero.Exists(fs, pipeline.OperationLocation)
		if err != nil {
			t.Fatalf("Could not check if operation location exists: %v", err)
		}

		if exists {
			t.Errorf("Operation locations should be deleted")
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
