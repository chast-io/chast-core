package refactoringRunModelSplitter

import (
	"chast.io/core/internal/model/run_models/refactoring"
)

type RunModelSplitter struct{}

//func (s *RunModelSplitter) SplitRunModels(runModel *run_models.RunModel) ([]*run_models.RunModel, error) {
//	switch m := (*runModel).(type) {
//	case *refactoring.RunModel:
//		return s.splitRunModels(m)
//	default:
//		return nil, errors.New("Not a refactoring run model")
//	}
//}

func (s *RunModelSplitter) SplitRunModels(runModel *refactoring.RunModel) ([]*refactoring.SingleRunModel, error) {
	var runModels []*refactoring.SingleRunModel
	for _, run := range runModel.Run {
		subRunModel := &refactoring.SingleRunModel{
			Run: run,
		}

		runModels = append(runModels, subRunModel)
	}
	return runModels, nil
}
