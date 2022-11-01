package refactoring_pipeline_builder

import (
	refactoringPipelineModel "chast.io/core/internal/model/pipeline/refactoring"
	"chast.io/core/internal/model/run_models/refactoring"
	refactoringRunModelIsolator "chast.io/core/internal/recipe/run_model_isolator/refactoring"
	"chast.io/core/pkg/util/collection"
	"log"
)

func BuildRunPipeline(runModel *refactoring.RunModel) *refactoringPipelineModel.Pipeline {
	executionOrder, err := getExecutionOrder(runModel.Run)
	if err != nil {
		log.Fatalf("Error getting execution order - %s", err)
	}

	isolatedExecutionOrder := buildIsolatedExecutionOrder(executionOrder, runModel)

	// TODO make configurable
	pipeline := refactoringPipelineModel.NewPipeline("/tmp/chast/", "/tmp/chast-changes/", "/")

	for _, runModelsInStage := range isolatedExecutionOrder {
		stage := refactoringPipelineModel.NewStage("") // TODO forward custom stage name
		for _, runModel := range runModelsInStage {
			step := refactoringPipelineModel.NewStep(runModel)
			stage.AddStep(step)
		}
		pipeline.AddStageAtStart(stage)
	}

	return pipeline

}

func getExecutionOrder(run []*refactoring.Run) ([][]*refactoring.Run, error) {
	usages, modelLookup := collectUsagesAndBuildLookup(run)

	// create a list of stages with tasks that can be executed in parallel
	executionGraphList := make([][]*refactoring.Run, 0)
	for len(usages[0]) > 0 {
		level := make([]*refactoring.Run, 0)
		usagesWithNoDependants := usages[0]
		newUsagesWithNoDependants := make(map[string]refactoring.Run)
		for _, run := range usagesWithNoDependants {
			level = append(level, &run)
			delete(usages[0], run.GetUUID())

			for _, dependency := range run.Dependencies {
				uuid := dependency.GetUUID()
				lookupIndex := modelLookup[uuid]

				if lookupIndex == 1 {
					newUsagesWithNoDependants[uuid] = *dependency
				} else {
					usages[lookupIndex-1][uuid] = usages[lookupIndex][uuid]
					delete(usages[lookupIndex], uuid)
				}
				modelLookup[uuid]--
			}

		}
		usages[0] = newUsagesWithNoDependants
		executionGraphList = append(executionGraphList, level)
	}

	return executionGraphList, nil
}

func collectUsagesAndBuildLookup(run []*refactoring.Run) ([]map[string]refactoring.Run, map[string]int) {
	var usages = make([]map[string]refactoring.Run, 0)
	var modelLookup = make(map[string]int)

	for _, runModel := range run {
		setUsageOfRun(runModel, &modelLookup, &usages)

		for _, dependency := range runModel.Dependencies {
			setUsageOfRun(dependency, &modelLookup, &usages)
		}
	}
	return usages, modelLookup
}

func setUsageOfRun(
	run *refactoring.Run,
	modelLookup *map[string]int,
	executionGraph *[]map[string]refactoring.Run,
) {

	uuid := run.GetUUID()
	if index, ok := (*modelLookup)[uuid]; ok {
		getMapInArray(executionGraph, index+1)[uuid] = *run
		delete((*executionGraph)[index], uuid)
		(*modelLookup)[uuid] = index + 1
	} else {
		getMapInArray(executionGraph, 0)[uuid] = *run
		(*modelLookup)[uuid] = 0
	}
}

func getMapInArray(maps *[]map[string]refactoring.Run, key int) map[string]refactoring.Run {
	if len(*maps) <= key {
		*maps = append(*maps, make(map[string]refactoring.Run))
	}
	if (*maps)[key] == nil {
		(*maps)[key] = make(map[string]refactoring.Run)
	}
	return (*maps)[key]
}

func buildIsolatedExecutionOrder(executionOrder [][]*refactoring.Run, runModel *refactoring.RunModel) [][]*refactoring.SingleRunModel {
	return collection.Map(executionOrder, func(run []*refactoring.Run) []*refactoring.SingleRunModel {
		return collection.Map(run, func(run *refactoring.Run) *refactoring.SingleRunModel {
			return refactoringRunModelIsolator.Isolate(runModel, run)
		})
	})
}
