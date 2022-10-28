package refactoring_pipeline_builder

import (
	"chast.io/core/internal/model/run_models/refactoring"
	"log"
)

func BuildRunPipeline(runModel *refactoring.RunModel) {
	log.Printf("Refactoring BuildRunPipeline")
	log.Printf("Run Command: %s", runModel.Run[0].Command.Cmds)

	buildLocalRunPipeline(runModel)

	//var nsContext = overlay.NewNamespaceContext(
	//	"/",                            // This will be defined by the versioning system
	//	"/tmp/overlay-auto-test/upper", // This will be defined by the versioning system
	//	"/tmp/overlay-auto-test/operationDirectory", // This will be defined by the versioning system
	//	runModel.Run[0].Command.WorkingDirectory,
	//	runModel.Run[0].Command.Cmds[0]..., // TODO support multiple commands
	//)

	//if err := overlay.RunCommandInIsolatedEnvironment(nsContext); err != nil {
	//	log.Fatalf("Error running command in isolated environment - %s", err)
	//}
}

func getExecutionGraph(run []*refactoring.Run) ([][]refactoring.SingleRunModel, error) {
	var usages = make([]map[string]refactoring.Run, 0)
	var modelLookup = make(map[string]int)

	for _, runModel := range run {
		setUsageOfRun(runModel, &modelLookup, &usages)

		for _, dependency := range runModel.Dependencies {
			setUsageOfRun(dependency, &modelLookup, &usages)
		}
	}

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

	return nil, nil
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

func buildLocalRunPipeline(runModel *refactoring.RunModel) {
	getExecutionGraph(runModel.Run)

	// TODO convert execution graph to a pipeline

	//runModels, err := (&refactoringRunModelSplitter.RunModelSplitter{}).SplitRunModels(runModel)
	//if err != nil {
	//	return nil, err
	//}

	//pipeline := refactoringPipelineModel.NewPipeline("", "", "")
	//
	//for _, runModel := range runModels {
	//	step := pipeline.AddStep(runModel)
	//	println(fmt.Sprintf("%#v", step))
	//}
}

func buildDockerRunPipeline(runModel *refactoring.Run) {

}
