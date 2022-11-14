package refactoringrunmodelbuilder

import (
	"path/filepath"
	"strings"

	"chast.io/core/internal/internal_util/collection"
	recipemodel "chast.io/core/internal/recipe/pkg/model"
	"chast.io/core/internal/run_model/internal/builder"
	extensionsdetection "chast.io/core/internal/run_model/internal/extensions_detection"
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
	variables *runmodel.Variables,
	unparsedArguments []string,
	unparsedFlags []runmodel.UnparsedFlag,
) (*runmodel.RunModel, error) {
	switch m := (*recipeModel).(type) {
	case *recipemodel.RefactoringRecipe:
		return parser.buildRunModel(m, variables, unparsedArguments, unparsedFlags)
	default:
		return nil, errors.New("Not a refactoring recipe")
	}
}

func (parser *RunModelBuilder) buildRunModel(
	recipeModel *recipemodel.RefactoringRecipe,
	variables *runmodel.Variables,
	unparsedArguments []string,
	unparsedFlags []runmodel.UnparsedFlag,
) (*runmodel.RunModel, error) {
	if err := builder.HandlePrimaryArgument(recipeModel.PrimaryParameter, variables, unparsedArguments[0]); err != nil {
		return nil, errors.Wrap(err, "Failed to handle primary argument")
	}

	if err := builder.HandlePositionalArguments(recipeModel, variables, unparsedArguments[1:]); err != nil {
		return nil, errors.Wrap(err, "Failed to handle positional arguments")
	}

	if err := builder.HandleFlags(recipeModel, variables, unparsedFlags); err != nil {
		return nil, errors.Wrap(err, "Failed to handle flags")
	}

	var runModel runmodel.RunModel

	filteredRuns, runsFilterError := filterRuns(recipeModel.Runs, variables)
	if runsFilterError != nil {
		return nil, errors.Wrap(runsFilterError, "Failed to filter runs")
	}

	namedRuns := make(map[string]*refactoring.Run)
	mappedRuns := collection.Map(filteredRuns,
		func(run recipemodel.Run) *refactoring.Run { return convertRun(run, variables, namedRuns) },
	)

	runModel = refactoring.RunModel{
		Run: mappedRuns,
	}

	return &runModel, nil
}

func filterRuns(runs []recipemodel.Run, variables *runmodel.Variables) ([]recipemodel.Run, error) {
	extensions, extensionDetectionError := extensionsdetection.DetectExtensions(variables.TypeDetectionPath)
	if extensionDetectionError != nil {
		return nil, errors.Wrap(extensionDetectionError, "Failed to detect extensions")
	}

	filteredRuns := make([]recipemodel.Run, 0)

	for _, run := range runs {
		if run.SupportedExtensions == nil || len(run.SupportedExtensions) == 0 {
			filteredRuns = append(filteredRuns, run)

			continue
		}

		for _, supportedExtension := range run.SupportedExtensions {
			if extensions[supportedExtension] != nil {
				filteredRuns = append(filteredRuns, run)
			}
		}
	}

	return filteredRuns, nil
}

func convertRun(
	run recipemodel.Run,
	variables *runmodel.Variables,
	namedRuns map[string]*refactoring.Run,
) *refactoring.Run {
	dependencies := convertDependencies(run.Dependencies, namedRuns)
	newRun := getOrComputeRunFromNamedRuns(run.ID, namedRuns)

	newRun.ID = run.ID
	newRun.Dependencies = dependencies
	newRun.SupportedLanguages = run.SupportedExtensions
	newRun.Command = convertCommand(run.Script, variables)
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

func convertCommand(commands []string, variables *runmodel.Variables) *refactoring.Command {
	cmds := collection.Map(commands, strings.Fields)
	replaceVariablesWithValuesInCommands(cmds, variables.Map)

	return &refactoring.Command{
		Cmds:             cmds,
		WorkingDirectory: filepath.Join(variables.WorkingDirectory, "run"),
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

func convertDocker(docker *recipemodel.Docker) *refactoring.Docker {
	if docker == nil {
		return nil
	}

	return &refactoring.Docker{
		DockerImage: docker.DockerImage,
	}
}

func convertLocal(local *recipemodel.Local) *refactoring.Local {
	if local == nil {
		return nil
	}

	return &refactoring.Local{
		RequiredTools: collection.Map(local.RequiredTools, convertRequiredTool),
	}
}

func convertRequiredTool(requiredTool recipemodel.RequiredTool) refactoring.RequiredTool {
	return refactoring.RequiredTool{
		Description: requiredTool.Description,
		CheckCmd:    requiredTool.CheckCmd,
	}
}
