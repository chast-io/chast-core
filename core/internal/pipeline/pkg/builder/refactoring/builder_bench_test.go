package refactoringpipelinebuilder

import (
	"strconv"
	"testing"

	"chast.io/core/internal/run_model/pkg/model/refactoring"
	log "github.com/sirupsen/logrus"
)

func BenchmarkBuildExecutionOrder(b *testing.B) {
	runModels := make([]*refactoring.Run, 0)

	for runNumber := 0; runNumber < 10000; runNumber++ {
		dependencies := make([]*refactoring.Run, 0)
		for j := 2; j < runNumber/10; j++ {
			dependencies = append(dependencies, runModels[runNumber/j+j])
		}

		runModel := &refactoring.Run{
			ID:                 "run" + strconv.Itoa(runNumber+1),
			Dependencies:       dependencies,
			SupportedLanguages: []string{},
			Docker:             &refactoring.Docker{},  //nolint:exhaustruct // not required for test
			Local:              &refactoring.Local{},   //nolint:exhaustruct // not required for test
			Command:            &refactoring.Command{}, //nolint:exhaustruct // not required for test
		}

		runModels = append(runModels, runModel)
	}

	runModel := &refactoring.RunModel{
		Run: runModels,
	}

	b.ResetTimer()

	logLevel := log.GetLevel()
	log.SetLevel(log.FatalLevel)

	_, _ = BuildRunPipeline(runModel)

	log.SetLevel(logLevel)
}
