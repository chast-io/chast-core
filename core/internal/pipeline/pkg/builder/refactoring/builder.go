package refactoringpipelinebuilder

import (
	"strconv"

	"chast.io/core/internal/internal_util/collection"
	refactoringpipelinemodel "chast.io/core/internal/pipeline/pkg/model/refactoring"
	refactoringRunModelIsolator "chast.io/core/internal/run_model/pkg/isolator/refactoring"
	"chast.io/core/internal/run_model/pkg/model/refactoring"
)

func BuildRunPipeline(runModel *refactoring.RunModel) *refactoringpipelinemodel.Pipeline {
	isolatedExecutionOrder := buildIsolatedExecutionOrder(runModel)

	// TODO make configurable
	pipeline := refactoringpipelinemodel.NewPipeline("/tmp/chast/", "/tmp/chast-changes/", "/")

	for i, runModelsInStage := range isolatedExecutionOrder {
		stage := refactoringpipelinemodel.NewStage(strconv.Itoa(i + 1)) // TODO forward custom stage name

		for i := range runModelsInStage {
			runModel := runModelsInStage[len(runModelsInStage)-1-i]
			step := refactoringpipelinemodel.NewStep(runModel)
			stage.AddStep(step)
		}

		pipeline.AddStage(stage)
	}

	return pipeline
}

func buildExecutionOrder(run []*refactoring.Run) [][]*refactoring.Run {
	usages, modelLookup := collectUsagesAndBuildLookup(run)

	// create a list of stages with tasks that can be executed in parallel
	executionGraphList := make([][]*refactoring.Run, 0)

	for len(usages[0]) > 0 {
		level := make([]*refactoring.Run, 0)
		usagesWithNoDependants := usages[0]
		newUsagesWithNoDependants := make(map[string]refactoring.Run)

		for i := range usagesWithNoDependants {
			run := usagesWithNoDependants[i]
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

		executionGraphList = append([][]*refactoring.Run{level}, executionGraphList...)
	}

	return executionGraphList
}

func collectUsagesAndBuildLookup(run []*refactoring.Run) ([]map[string]refactoring.Run, map[string]int) {
	usages := make([]map[string]refactoring.Run, 0)
	modelLookup := make(map[string]int)

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

func buildIsolatedExecutionOrder(
	runModel *refactoring.RunModel,
) [][]*refactoring.SingleRunModel {
	executionOrder := buildExecutionOrder(runModel.Run)

	return collection.Map(executionOrder, func(run []*refactoring.Run) []*refactoring.SingleRunModel {
		return collection.Map(run, func(run *refactoring.Run) *refactoring.SingleRunModel {
			return refactoringRunModelIsolator.Isolate(runModel, run)
		})
	})
}
