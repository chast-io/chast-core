package refactoringpipelinebuilder //nolint:testpackage // access to private members required
import (
	log "github.com/sirupsen/logrus"
	"strconv"
	"testing"

	"chast.io/core/internal/run_model/pkg/model/refactoring"
)

func BenchmarkBuildExecutionOrder(b *testing.B) {
	runModels := make([]*refactoring.Run, 0)

	for i := 0; i < 10000; i++ {
		dependencies := make([]*refactoring.Run, 0)
		for j := 2; j < i/10; j++ {
			dependencies = append(dependencies, runModels[i/j+j])
		}

		runModel := &refactoring.Run{
			ID:                 "run" + strconv.Itoa(i+1),
			Dependencies:       dependencies,
			SupportedLanguages: []string{},
			Docker:             refactoring.Docker{},  //nolint:exhaustruct // not required for test
			Local:              refactoring.Local{},   //nolint:exhaustruct // not required for test
			Command:            refactoring.Command{}, //nolint:exhaustruct // not required for test
		}

		runModels = append(runModels, runModel)
	}

	runModel := &refactoring.RunModel{
		Run: runModels,
	}

	b.ResetTimer()

	logLevel := log.GetLevel()
	log.SetLevel(log.FatalLevel)

	BuildRunPipeline(runModel)

	log.SetLevel(logLevel)
}
