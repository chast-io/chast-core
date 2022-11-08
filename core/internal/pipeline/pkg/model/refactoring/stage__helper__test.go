package refactoringpipelinemodel

import (
	"chast.io/core/internal/run_model/pkg/model/refactoring"
)

func dummyPipeline() *Pipeline {
	return NewPipeline("/operationLocation", "/changeCaptureLocation", "/rootFileSystemLocation")
}

func dummyStage() *Stage {
	return NewStage("test-name")
}

func dummyStageWithPipeline() (*Stage, *Pipeline) {
	stage := dummyStage()
	pipeline := dummyPipeline()
	pipeline.AddStage(stage)
	return stage, pipeline
}

func dummyStep() *Step {
	runModel := &refactoring.SingleRunModel{
		Run: &refactoring.Run{
			ID:                 "runId",
			Dependencies:       make([]*refactoring.Run, 0),
			SupportedLanguages: []string{"java"},
			Docker:             refactoring.Docker{},
			Local:              refactoring.Local{},
			Command:            refactoring.Command{},
		},
		Stage: "stage",
	}

	return NewStep(runModel)
}

// TODO add test for AddStep, GetPrevChangeCaptureLocations
