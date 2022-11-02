package refactoringrunmodelbuilder

import (
	"path/filepath"
	"strings"

	"chast.io/core/internal/internal_util/collection"
	recipemodel "chast.io/core/internal/recipe/pkg/model"
	runmodel "chast.io/core/internal/run_model/pkg/model"
	"chast.io/core/internal/run_model/pkg/model/refactoring"
	"github.com/pkg/errors"
)

type RunModelBuilder struct{}

func NewRunModelBuilder() *RunModelBuilder {
	return &RunModelBuilder{}
}

func (parser *RunModelBuilder) BuildRunModel(
	recipeModel *recipemodel.Recipe,
	arguments *runmodel.ParsedArguments,
) (*runmodel.RunModel, error) {
	switch m := (*recipeModel).(type) {
	case *recipemodel.RefactoringRecipe:
		return parser.buildRunModel(m, arguments)
	default:
		return nil, errors.New("Not a refactoring recipe")
	}
}

func (parser *RunModelBuilder) buildRunModel(
	recipeModel *recipemodel.RefactoringRecipe,
	arguments *runmodel.ParsedArguments,
) (*runmodel.RunModel, error) {
	// TODO hande additional arguments
	var runModel runmodel.RunModel

	namedRuns := make(map[string]*refactoring.Run)
	mappedRuns := collection.Map(recipeModel.Runs,
		func(run recipemodel.Run) *refactoring.Run { return convertRun(run, arguments, namedRuns) },
	)

	runModel = refactoring.RunModel{
		Run:    mappedRuns,
		Stages: nil, // TODO implement
	}

	return &runModel, nil
}

func convertRun(
	run recipemodel.Run,
	arguments *runmodel.ParsedArguments,
	namedRuns map[string]*refactoring.Run,
) *refactoring.Run {
	dependencies := convertDependencies(run.Dependencies, namedRuns)
	newRun := getOrComputeRunFromNamedRuns(run.ID, namedRuns)

	newRun.ID = run.ID
	newRun.Dependencies = dependencies
	newRun.SupportedLanguages = run.SupportedLanguages
	newRun.Command = convertCommand(run.Script, arguments)
	newRun.Docker = convertDocker(run.Docker)
	newRun.Local = convertLocal(run.Local)

	return newRun
}

func convertDependencies(dependencies []string, namedRuns map[string]*refactoring.Run) []*refactoring.Run {
	convertDependencies := make([]*refactoring.Run, len(dependencies))

	for i, dependency := range dependencies {
		if _, ok := namedRuns[dependency]; !ok {
			namedRuns[dependency] = &refactoring.Run{} //nolint:exhaustruct // NPE prevention
		}

		convertDependencies[i] = namedRuns[dependency]
	}

	return convertDependencies
}

func getOrComputeRunFromNamedRuns(runID string, namedRuns map[string]*refactoring.Run) *refactoring.Run {
	var newRun *refactoring.Run

	if runID != "" {
		if _, ok := namedRuns[runID]; !ok {
			newRun = &refactoring.Run{} //nolint:exhaustruct // NPE prevention
			namedRuns[runID] = newRun
		} else {
			newRun = namedRuns[runID]
		}
	} else {
		newRun = &refactoring.Run{} //nolint:exhaustruct // NPE prevention
	}

	return newRun
}

func convertCommand(commands []string, arguments *runmodel.ParsedArguments) refactoring.Command {
	cmds := collection.Map(commands, strings.Fields)
	replaceVariablesWithValuesInCommands(cmds, arguments.Arguments)

	return refactoring.Command{
		Cmds:             cmds,
		WorkingDirectory: filepath.Join(arguments.WorkingDirectory, "run"),
	}
}

func replaceVariablesWithValuesInCommands(commands [][]string, arguments map[string]string) {
	for i, cmd := range commands {
		for j, cmdPart := range cmd {
			commands[i][j] = replaceVariablesWithValues(cmdPart, arguments)
		}
	}
}

func replaceVariablesWithValues(value string, arguments map[string]string) string {
	// TODO optimize
	for key, val := range arguments {
		value = strings.ReplaceAll(value, "$"+key, val)
		value = strings.ReplaceAll(value, "${"+key+"}", val)
	}

	return value
}

func convertDocker(docker recipemodel.Docker) refactoring.Docker {
	return refactoring.Docker{
		DockerImage: docker.DockerImage,
	}
}

func convertLocal(local recipemodel.Local) refactoring.Local {
	return refactoring.Local{
		RequiredTools: collection.Map(local.RequiredTools, convertRequiredTool),
	}
}

func convertRequiredTool(requiredTool recipemodel.RequiredTool) refactoring.RequiredTool {
	return refactoring.RequiredTool{
		Description: requiredTool.Description,
		CheckCmd:    requiredTool.CheckCmd,
	}
}
