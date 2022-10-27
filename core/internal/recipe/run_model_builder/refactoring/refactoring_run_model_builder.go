package refactoring

import (
	"chast.io/core/internal/model/recipe"
	"chast.io/core/internal/model/run_models"
	"chast.io/core/internal/model/run_models/refactoring"
	"chast.io/core/pkg/util/collection"
	"github.com/pkg/errors"
	"path/filepath"
	"strings"
)

type RunModelBuilder struct {
}

func NewRunModelBuilder() *RunModelBuilder {
	return &RunModelBuilder{}
}

func (parser *RunModelBuilder) BuildRunModel(
	recipeModel *recipe.Recipe,
	arguments *run_models.ParsedArguments) (*run_models.RunModel, error) {

	switch m := (*recipeModel).(type) {
	case *recipe.RefactoringRecipe:
		return parser.buildRunModel(m, arguments)
	default:
		return nil, errors.New("Not a refactoring recipe")
	}
}

func (parser *RunModelBuilder) buildRunModel(
	recipeModel *recipe.RefactoringRecipe,
	arguments *run_models.ParsedArguments) (*run_models.RunModel, error) {
	// TODO hande additional arguments
	var runModel run_models.RunModel
	mappedRun := collection.Map(recipeModel.Run,
		func(run recipe.Run) refactoring.Run { return convertRun(run, arguments) },
	)
	runModel = refactoring.RunModel{
		SupportedLanguages: recipeModel.SupportedLanguages,
		Run:                mappedRun,
	}

	return &runModel, nil
}

func convertRun(run recipe.Run, arguments *run_models.ParsedArguments) refactoring.Run {
	return refactoring.Run{
		Command: convertCommand(run.Script, arguments),
		Docker:  convertDocker(run.Docker),
		Local:   convertLocal(run.Local),
	}
}

func convertCommand(commands []string, arguments *run_models.ParsedArguments) refactoring.Command {
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

func convertDocker(docker recipe.Docker) refactoring.Docker {
	return refactoring.Docker{
		DockerImage: docker.DockerImage,
	}
}

func convertLocal(local recipe.Local) refactoring.Local {
	return refactoring.Local{
		RequiredTools: collection.Map(local.RequiredTools, convertRequiredTool),
	}
}

func convertRequiredTool(requiredTool recipe.RequiredTool) refactoring.RequiredTool {
	return refactoring.RequiredTool{
		Description: requiredTool.Description,
		CheckCmd:    requiredTool.CheckCmd,
	}
}
