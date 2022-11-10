package refactoringpipelinebuilder_test

import (
	"testing"

	uut "chast.io/core/internal/pipeline/pkg/builder/refactoring"
	"chast.io/core/internal/run_model/pkg/model/refactoring"
)

// region Helpers
func builderDummyRunModelWithSingleRun() *refactoring.RunModel {
	return &refactoring.RunModel{
		Run: []*refactoring.Run{
			{
				ID:                 "run1",
				Dependencies:       []*refactoring.Run{},
				SupportedLanguages: []string{"java"},
				Docker: refactoring.Docker{
					DockerImage: "dockerImage1",
				},
				Local: refactoring.Local{
					RequiredTools: []refactoring.RequiredTool{
						{
							Description: "description1",
							CheckCmd:    "checkCmd1",
						},
					},
				},
				Command: refactoring.Command{
					Cmds: [][]string{
						{"cmd1"},
					},
					WorkingDirectory: "workingDirectory1",
				},
			},
		},
	}
}

// endregion

// region BuildRunPipeline [SingleRun]
func TestBuildRunPipeline_SingleRun(t *testing.T) {
	t.Parallel()

	runModel1 := builderDummyRunModelWithSingleRun()

	actualPipeline, _ := uut.BuildRunPipeline(runModel1)

	t.Run("should set UUID", func(t *testing.T) {
		t.Parallel()
		if actualPipeline.UUID == "" {
			t.Error("Expected pipeline UUID to be set, but was empty")
		}
	})

	t.Run("should set operation location", func(t *testing.T) {
		t.Parallel()
		if actualPipeline.OperationLocation != "/tmp/chast" {
			t.Errorf("expected operation location to be %s but was %s", "/tmp/chast", actualPipeline.OperationLocation)
		}
	})

	t.Run("should set change capture folder", func(t *testing.T) {
		t.Parallel()
		if actualPipeline.ChangeCaptureLocation != "/tmp/chast-changes/"+actualPipeline.UUID {
			t.Errorf("Expected pipeline ChangeCaptureLocation to be '/tmp/chast-changes/%s', but was '%s'", actualPipeline.UUID, actualPipeline.ChangeCaptureLocation)
		}
	})

	t.Run("should set root file system location", func(t *testing.T) {
		t.Parallel()
		if actualPipeline.RootFileSystemLocation != "/" {
			t.Errorf("Expected pipeline RootFileSystemLocation to be '/', but was '%s'", actualPipeline.RootFileSystemLocation)
		}
	})

	t.Run("should set stages", func(t *testing.T) {
		t.Parallel()
		if len(actualPipeline.Stages) != 1 {
			t.Errorf("Expected pipeline to have 1 stage, but had %d", len(actualPipeline.Stages))
		}
	})
}

// endregion

// region BuildRunPipeline [MultipleRuns]

func TestBuildRunPipeline_MultipleRuns_WithStages(t *testing.T) {
	t.Parallel()

	run1 := &refactoring.Run{
		ID:                 "run1",
		Dependencies:       []*refactoring.Run{},
		SupportedLanguages: []string{},
		Docker:             refactoring.Docker{},  //nolint:exhaustruct // not required for test
		Local:              refactoring.Local{},   //nolint:exhaustruct // not required for test
		Command:            refactoring.Command{}, //nolint:exhaustruct // not required for test
	}

	run2 := &refactoring.Run{
		ID:                 "run2",
		Dependencies:       []*refactoring.Run{},
		SupportedLanguages: []string{},
		Docker:             refactoring.Docker{},  //nolint:exhaustruct // not required for test
		Local:              refactoring.Local{},   //nolint:exhaustruct // not required for test
		Command:            refactoring.Command{}, //nolint:exhaustruct // not required for test
	}

	run3deps1 := &refactoring.Run{
		ID:                 "run3deps1",
		Dependencies:       []*refactoring.Run{run1},
		SupportedLanguages: []string{},
		Docker:             refactoring.Docker{},  //nolint:exhaustruct // not required for test
		Local:              refactoring.Local{},   //nolint:exhaustruct // not required for test
		Command:            refactoring.Command{}, //nolint:exhaustruct // not required for test
	}

	run4deps1and2and3 := &refactoring.Run{
		ID:                 "run4deps1and2and3",
		Dependencies:       []*refactoring.Run{run1, run2, run3deps1},
		SupportedLanguages: []string{},
		Docker:             refactoring.Docker{},  //nolint:exhaustruct // not required for test
		Local:              refactoring.Local{},   //nolint:exhaustruct // not required for test
		Command:            refactoring.Command{}, //nolint:exhaustruct // not required for test
	}

	run5 := &refactoring.Run{
		ID:                 "run5",
		Dependencies:       []*refactoring.Run{},
		SupportedLanguages: []string{},
		Docker:             refactoring.Docker{},  //nolint:exhaustruct // not required for test
		Local:              refactoring.Local{},   //nolint:exhaustruct // not required for test
		Command:            refactoring.Command{}, //nolint:exhaustruct // not required for test
	}

	runModel := &refactoring.RunModel{
		Run: []*refactoring.Run{
			run1,
			run2,
			run3deps1,
			run4deps1and2and3,
			run5,
		},
	}

	actualPipeline, _ := uut.BuildRunPipeline(runModel)

	t.Run("should set stages", func(t *testing.T) {
		t.Parallel()
		if len(actualPipeline.Stages) != 3 {
			t.Errorf("Expected pipeline to have 3 stages, but had %d", len(actualPipeline.Stages))
		}
	})

	t.Run("should set stage 1", func(t *testing.T) {
		t.Parallel()
		if len(actualPipeline.Stages[0].Steps) != 3 {
			t.Errorf("Expected stage 1 to have 3 steps, but had %d", len(actualPipeline.Stages[0].Steps))
		}
	})

	t.Run("should set stage 2", func(t *testing.T) {
		t.Parallel()
		if len(actualPipeline.Stages[1].Steps) != 1 {
			t.Errorf("Expected stage 2 to have 1 step, but had %d", len(actualPipeline.Stages[1].Steps))
		}
	})

	t.Run("should set stage 3", func(t *testing.T) {
		t.Parallel()
		if len(actualPipeline.Stages[2].Steps) != 1 {
			t.Errorf("Expected stage 3 to have 1 step, but had %d", len(actualPipeline.Stages[2].Steps))
		}
	})

	t.Run("should set stage 1 step 1", func(t *testing.T) {
		t.Parallel()
		if actualPipeline.Stages[0].Steps[0].RunModel.Run != run1 {
			t.Errorf("Expected stage 1 step 1 to be run1, but was %s", actualPipeline.Stages[0].Steps[0].RunModel.Run.ID)
		}
	})

	t.Run("should set stage 1 step 2", func(t *testing.T) {
		t.Parallel()
		if actualPipeline.Stages[0].Steps[1].RunModel.Run != run2 {
			t.Errorf("Expected stage 1 step 2 to be run2, but was %s", actualPipeline.Stages[0].Steps[1].RunModel.Run.ID)
		}
	})

	t.Run("should set stage 1 step 3", func(t *testing.T) {
		t.Parallel()
		if actualPipeline.Stages[0].Steps[2].RunModel.Run != run5 {
			t.Errorf("Expected stage 1 step 3 to be run5, but was %s", actualPipeline.Stages[2].Steps[0].RunModel.Run.ID)
		}
	})

	t.Run("should set stage 2 step 1", func(t *testing.T) {
		t.Parallel()
		if actualPipeline.Stages[1].Steps[0].RunModel.Run != run3deps1 {
			t.Errorf("Expected stage 2 step 1 to be run3deps1, but was %s", actualPipeline.Stages[1].Steps[0].RunModel.Run.ID)
		}
	})

	t.Run("should set stage 3 step 1", func(t *testing.T) {
		t.Parallel()
		if actualPipeline.Stages[2].Steps[0].RunModel.Run != run4deps1and2and3 {
			t.Errorf("Expected stage 3 step 1 to be run4depsAll, but was %s", actualPipeline.Stages[2].Steps[0].RunModel.Run.ID)
		}
	})
}

// endregion

// region Cyclic Dependency Detection

func TestBuildExecutionOrder_CyclicDependencyDetection(t *testing.T) {
	t.Parallel()

	run1 := &refactoring.Run{
		ID:                 "run1",
		Dependencies:       []*refactoring.Run{},
		SupportedLanguages: []string{},
		Docker:             refactoring.Docker{},  //nolint:exhaustruct // not required for test
		Local:              refactoring.Local{},   //nolint:exhaustruct // not required for test
		Command:            refactoring.Command{}, //nolint:exhaustruct // not required for test
	}

	run2 := &refactoring.Run{
		ID:                 "run2",
		Dependencies:       []*refactoring.Run{run1},
		SupportedLanguages: []string{},
		Docker:             refactoring.Docker{},  //nolint:exhaustruct // not required for test
		Local:              refactoring.Local{},   //nolint:exhaustruct // not required for test
		Command:            refactoring.Command{}, //nolint:exhaustruct // not required for test
	}

	run3 := &refactoring.Run{
		ID:                 "run3",
		Dependencies:       []*refactoring.Run{run2},
		SupportedLanguages: []string{},
		Docker:             refactoring.Docker{},  //nolint:exhaustruct // not required for test
		Local:              refactoring.Local{},   //nolint:exhaustruct // not required for test
		Command:            refactoring.Command{}, //nolint:exhaustruct // not required for test
	}

	run4 := &refactoring.Run{
		ID:                 "run4",
		Dependencies:       []*refactoring.Run{run3},
		SupportedLanguages: []string{},
		Docker:             refactoring.Docker{},  //nolint:exhaustruct // not required for test
		Local:              refactoring.Local{},   //nolint:exhaustruct // not required for test
		Command:            refactoring.Command{}, //nolint:exhaustruct // not required for test
	}

	run5 := &refactoring.Run{
		ID:                 "run5",
		Dependencies:       []*refactoring.Run{},
		SupportedLanguages: []string{},
		Docker:             refactoring.Docker{},  //nolint:exhaustruct // not required for test
		Local:              refactoring.Local{},   //nolint:exhaustruct // not required for test
		Command:            refactoring.Command{}, //nolint:exhaustruct // not required for test
	}

	run2.Dependencies = []*refactoring.Run{run1, run3} // introduce cyclic dependency

	runModel := &refactoring.RunModel{
		Run: []*refactoring.Run{
			run1,
			run2,
			run3,
			run4,
			run5,
		},
	}

	_, err := uut.BuildRunPipeline(runModel)

	if err == nil {
		t.Error("expected error to be returned but was nil")
	}
}

// endregion
