package refactoringpipelinebuilder

import (
	"chast.io/core/internal/run_model/pkg/model/refactoring"
	"strings"
	"testing"
)

func TestBuildRunPipeline_SingleRun(t *testing.T) {
	runModel1 := &refactoring.RunModel{
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
		Stages: []string{"stage1"},
	}

	//pipeline1 := refactoringpipelinemodel.Pipeline{
	//	OperationLocation:      "/tmp/chast",
	//	ChangeCaptureLocation:    "/tmp/chast-changes",
	//	RootFileSystemLocation: "/",
	//	UUID: "",
	//	Stages:
	//}

	actualPipeline := BuildRunPipeline(runModel1)

	t.Run("should set UUID prefix", func(t *testing.T) {
		t.Parallel()
		if strings.HasPrefix(actualPipeline.UUID, "PIPELINE-") == false {
			t.Errorf("Expected pipeline UUID to start with 'PIPELINE-', but was '%s'", actualPipeline.UUID)
		}
	})

	t.Run("should set correct UUID ", func(t *testing.T) {
		t.Parallel()
		if len(actualPipeline.UUID) != len("PIPELINE-")+len("00000000-0000-0000-0000-000000000000") {
			t.Errorf("Expected pipeline UUID to be 36 characters long, but was %d", len(actualPipeline.UUID))
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
