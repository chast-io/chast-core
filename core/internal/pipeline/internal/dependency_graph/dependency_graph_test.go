package dependencygraph_test

import (
	"testing"

	uut "chast.io/core/internal/pipeline/internal/dependency_graph"
	"chast.io/core/internal/run_model/pkg/model/refactoring"
)

// region Helpers
func dependencyGraphDummyRunModelWithSingleRun() *refactoring.RunModel {
	return &refactoring.RunModel{
		Run: []*refactoring.Run{
			{
				ID:                 "run1",
				Dependencies:       []*refactoring.Run{},
				SupportedLanguages: []string{"java"},
				Docker: &refactoring.Docker{
					DockerImage: "dockerImage1",
				},
				Local: &refactoring.Local{
					RequiredTools: []refactoring.RequiredTool{
						{
							Description: "description1",
							CheckCmd:    "checkCmd1",
						},
					},
				},
				Command: &refactoring.Command{
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
func TestBuildExecutionOrder_SingleRun(t *testing.T) {
	t.Parallel()

	runModel1 := dependencyGraphDummyRunModelWithSingleRun()

	executionOrder, _ := uut.BuildExecutionOrder(runModel1)

	t.Run("should return a single stage", func(t *testing.T) {
		t.Parallel()

		if len(executionOrder) != 1 {
			t.Errorf("expected execution order to contain 1 stage but was %d", len(executionOrder))
		}
	})

	t.Run("should return a single run", func(t *testing.T) {
		t.Parallel()

		if len(executionOrder[0]) != 1 {
			t.Errorf("expected execution order to contain 1 run but was %d", len(executionOrder[0]))
		}
	})
}

// endregion

// region BuildRunPipeline [MultipleRuns]

func TestBuildExecutionOrder_MultipleRuns_WithStages(t *testing.T) {
	t.Parallel()

	run1 := &refactoring.Run{
		ID:                 "run1",
		Dependencies:       []*refactoring.Run{},
		SupportedLanguages: []string{},
		Docker:             &refactoring.Docker{},  //nolint:exhaustruct // not required for test
		Local:              &refactoring.Local{},   //nolint:exhaustruct // not required for test
		Command:            &refactoring.Command{}, //nolint:exhaustruct // not required for test
	}

	run2 := &refactoring.Run{
		ID:                 "run2",
		Dependencies:       []*refactoring.Run{},
		SupportedLanguages: []string{},
		Docker:             &refactoring.Docker{},  //nolint:exhaustruct // not required for test
		Local:              &refactoring.Local{},   //nolint:exhaustruct // not required for test
		Command:            &refactoring.Command{}, //nolint:exhaustruct // not required for test
	}

	run3deps1 := &refactoring.Run{
		ID:                 "run3deps1",
		Dependencies:       []*refactoring.Run{run1},
		SupportedLanguages: []string{},
		Docker:             &refactoring.Docker{},  //nolint:exhaustruct // not required for test
		Local:              &refactoring.Local{},   //nolint:exhaustruct // not required for test
		Command:            &refactoring.Command{}, //nolint:exhaustruct // not required for test
	}

	run4deps1and2and3 := &refactoring.Run{
		ID:                 "run4deps1and2and3",
		Dependencies:       []*refactoring.Run{run1, run2, run3deps1},
		SupportedLanguages: []string{},
		Docker:             &refactoring.Docker{},  //nolint:exhaustruct // not required for test
		Local:              &refactoring.Local{},   //nolint:exhaustruct // not required for test
		Command:            &refactoring.Command{}, //nolint:exhaustruct // not required for test
	}

	run5 := &refactoring.Run{
		ID:                 "run5",
		Dependencies:       []*refactoring.Run{},
		SupportedLanguages: []string{},
		Docker:             &refactoring.Docker{},  //nolint:exhaustruct // not required for test
		Local:              &refactoring.Local{},   //nolint:exhaustruct // not required for test
		Command:            &refactoring.Command{}, //nolint:exhaustruct // not required for test
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

	executionOrder, _ := uut.BuildExecutionOrder(runModel)

	t.Run("should set stages", func(t *testing.T) {
		t.Parallel()
		if len(executionOrder) != 3 {
			t.Errorf("expected execution order to contain 3 stages but was %d", len(executionOrder))
		}
	})

	t.Run("should set stage 1", func(t *testing.T) {
		t.Parallel()
		if len(executionOrder[0]) != 3 {
			t.Errorf("expected execution order to contain 3 runs in stage 1 but was %d", len(executionOrder[0]))
		}
	})

	t.Run("should set stage 2", func(t *testing.T) {
		t.Parallel()
		if len(executionOrder[1]) != 1 {
			t.Errorf("expected execution order to contain 1 run in stage 2 but was %d", len(executionOrder[1]))
		}
	})

	t.Run("should set stage 3", func(t *testing.T) {
		t.Parallel()
		if len(executionOrder[2]) != 1 {
			t.Errorf("expected execution order to contain 1 run in stage 3 but was %d", len(executionOrder[2]))
		}
	})

	t.Run("should set stage 1 step 1", func(t *testing.T) {
		t.Parallel()
		if executionOrder[0][0] != run1 {
			t.Errorf("expected execution order to contain run1 in stage 1 step 1 but was %s", executionOrder[0][0].ID)
		}
	})

	t.Run("should set stage 1 step 2", func(t *testing.T) {
		t.Parallel()
		if executionOrder[0][1] != run2 {
			t.Errorf("expected execution order to contain run2 in stage 1 step 2 but was %s", executionOrder[0][1].ID)
		}
	})

	t.Run("should set stage 1 step 3", func(t *testing.T) {
		t.Parallel()
		if executionOrder[0][2] != run5 {
			t.Errorf("expected execution order to contain run5 in stage 1 step 3 but was %s", executionOrder[0][2].ID)
		}
	})

	t.Run("should set stage 2 step 1", func(t *testing.T) {
		t.Parallel()
		if executionOrder[1][0] != run3deps1 {
			t.Errorf("expected execution order to contain run3deps1 in stage 2 step 1 but was %s", executionOrder[1][0].ID)
		}
	})

	t.Run("should set stage 3 step 1", func(t *testing.T) {
		t.Parallel()
		if executionOrder[2][0] != run4deps1and2and3 {
			t.Errorf("expected execution order to contain run4deps1and2and3 in stage 3 step 1 but was %s", executionOrder[2][0].ID)
		}
	})
}

// endregion

// region Cyclic Dependency Detection

func TestBuildExecutionOrder_CyclicDependencyDetection(t *testing.T) {
	t.Parallel()

	// run1 <- run2 <- run3 <- run4
	//           |_______^
	// run5

	run1 := &refactoring.Run{
		ID:                 "run1",
		Dependencies:       []*refactoring.Run{},
		SupportedLanguages: []string{},
		Docker:             &refactoring.Docker{},  //nolint:exhaustruct // not required for test
		Local:              &refactoring.Local{},   //nolint:exhaustruct // not required for test
		Command:            &refactoring.Command{}, //nolint:exhaustruct // not required for test
	}

	run2 := &refactoring.Run{
		ID:                 "run2",
		Dependencies:       []*refactoring.Run{run1},
		SupportedLanguages: []string{},
		Docker:             &refactoring.Docker{},  //nolint:exhaustruct // not required for test
		Local:              &refactoring.Local{},   //nolint:exhaustruct // not required for test
		Command:            &refactoring.Command{}, //nolint:exhaustruct // not required for test
	}

	run3 := &refactoring.Run{
		ID:                 "run3",
		Dependencies:       []*refactoring.Run{run2},
		SupportedLanguages: []string{},
		Docker:             &refactoring.Docker{},  //nolint:exhaustruct // not required for test
		Local:              &refactoring.Local{},   //nolint:exhaustruct // not required for test
		Command:            &refactoring.Command{}, //nolint:exhaustruct // not required for test
	}

	run4 := &refactoring.Run{
		ID:                 "run4",
		Dependencies:       []*refactoring.Run{run3},
		SupportedLanguages: []string{},
		Docker:             &refactoring.Docker{},  //nolint:exhaustruct // not required for test
		Local:              &refactoring.Local{},   //nolint:exhaustruct // not required for test
		Command:            &refactoring.Command{}, //nolint:exhaustruct // not required for test
	}

	run5 := &refactoring.Run{
		ID:                 "run5",
		Dependencies:       []*refactoring.Run{},
		SupportedLanguages: []string{},
		Docker:             &refactoring.Docker{},  //nolint:exhaustruct // not required for test
		Local:              &refactoring.Local{},   //nolint:exhaustruct // not required for test
		Command:            &refactoring.Command{}, //nolint:exhaustruct // not required for test
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

	_, err := uut.BuildExecutionOrder(runModel)

	if err == nil {
		t.Error("expected error to be returned but was nil")
	}
}

// endregion
