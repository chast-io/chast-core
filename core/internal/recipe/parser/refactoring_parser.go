package parser

import (
	"chast.io/core/internal/model/run_models"
	"chast.io/core/internal/model/run_models/refactoring"
	"chast.io/core/pkg/util/collection"
	"errors"
	"gopkg.in/yaml.v3"
	"log"
	"strings"
)

import (
	. "chast.io/core/internal/recipe/recipe_model"
)

type RefactoringParser struct {
	recipeModel *RefactoringRecipe
}

func (parser *RefactoringParser) ParseRecipe(data *[]byte) {
	var refactoringRecipe RefactoringRecipe

	decoder := yaml.NewDecoder(strings.NewReader(string(*data)))
	decoder.KnownFields(true)

	err := decoder.Decode(&refactoringRecipe)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	parser.recipeModel = &refactoringRecipe
}

func (parser *RefactoringParser) VerifyRecipeAndBuildModel() (*run_models.RunModel, error) {
	recipeModel := parser.recipeModel

	if recipeModel == nil {
		return nil, errors.New("recipe model needs to initialized")
	}

	var runModel run_models.RunModel
	mappedRun := collection.Map(recipeModel.Run, convertRun)
	runModel = refactoring.RunModel{
		SupportedLanguages: recipeModel.SupportedLanguages,
		Run:                mappedRun,
	}

	return &runModel, nil
}

func convertRun(run Run) refactoring.Run {
	return refactoring.Run{
		Command: convertCommand(run.Cmd),
		Docker:  convertDocker(run.Docker),
		Local:   convertLocal(run.Local),
	}
}

func convertCommand(command string) refactoring.Command {
	return refactoring.Command{
		Cmd:              strings.Fields(command),
		WorkingDirectory: "/shared/home/rjenni/Projects/mse-repos/master-thesis/chast/chast-refactoring-antlr/refactorings/rearrange_class_members/chast/run",
	}
}

func convertDocker(docker Docker) refactoring.Docker {
	return refactoring.Docker{
		DockerImage: docker.DockerImage,
	}
}

func convertLocal(local Local) refactoring.Local {
	return refactoring.Local{
		RequiredTools: collection.Map(local.RequiredTools, convertRequiredTool),
	}
}

func convertRequiredTool(requiredTool RequiredTool) refactoring.RequiredTool {
	return refactoring.RequiredTool{
		Description: requiredTool.Description,
		Cmd:         requiredTool.Cmd,
	}
}
