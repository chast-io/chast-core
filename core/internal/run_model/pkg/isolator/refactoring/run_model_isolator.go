package refactoringrunmodelisolator

import (
	"chast.io/core/internal/run_model/pkg/model/refactoring"
)

func Isolate(_ *refactoring.RunModel, run *refactoring.Run) *refactoring.SingleRunModel {
	return &refactoring.SingleRunModel{
		Run: run,
	}
}
