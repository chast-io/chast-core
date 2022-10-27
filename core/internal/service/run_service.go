package service

import (
	refactoringPipelineBuilder "chast.io/core/internal/core/pipeline_builder/refactoring"
	"chast.io/core/internal/model/run_models/refactoring"
	"chast.io/core/internal/recipe/run_model_builder"
	util "chast.io/core/pkg/util"
	"github.com/pkg/errors"
	"log"
)
import (
	"chast.io/core/internal/recipe/parser"
)

func Run(recipeFile *util.File, args ...string) (*Pipeline, error) {
	//var model generalRunModel.RunModel
	parsedRecipe, err := parser.ParseRecipe(recipeFile)
	if err != nil {
		panic(err)
	}
	log.Printf("Parsed recipe: %v", parsedRecipe)

	runModel, err := run_model_builder.BuildRunModel(parsedRecipe, args, recipeFile.ParentDirectory)
	if err != nil {
		return nil, err
	}

	switch m := (*runModel).(type) {
	case refactoring.RunModel:
		refactoringPipelineBuilder.BuildRunPipeline(&m)
	default:
		return nil, errors.Errorf("No pipline builder for provided run model")
	}

	println(runModel)

	return &Pipeline{}, nil
}

// TODO
type Pipeline struct {
}
