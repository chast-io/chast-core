package refactoringRunModelIsolator

import (
	"chast.io/core/internal/run_model/model/refactoring"
)

//func (s *RunModelIsolator) SplitRunModels(runModel *run_models.RunModel) ([]*run_models.RunModel, error) {
//	switch m := (*runModel).(type) {
//	case *refactoring.RunModel:
//		return s.splitRunModels(m)
//	default:
//		return nil, errors.New("Not a refactoring run model")
//	}
//}

func Isolate(runModel *refactoring.RunModel, run *refactoring.Run) *refactoring.SingleRunModel {
	return &refactoring.SingleRunModel{
		Run: run,
	}
}
