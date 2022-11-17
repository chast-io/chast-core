package refactoringrunmodelisolator

import (
	"chast.io/core/internal/run_model/pkg/model/refactoring"
)

func Isolate(runModel *refactoring.RunModel, run *refactoring.Run) *refactoring.SingleRunModel {
	return &refactoring.SingleRunModel{
		Run: run,
	}
}
