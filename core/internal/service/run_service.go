package service

import (
	refactoringPipelineBuilder "chast.io/core/internal/core/pipeline_builder/refactoring"
	refactoringRunModel "chast.io/core/internal/model/run_models/refactoring"
	util "chast.io/core/pkg/util"
	"log"
)
import (
	"chast.io/core/internal/recipe/parser"
)

func BuildPipeline(recipe util.FileReader) Pipeline {
	//var model generalRunModel.RunModel
	model, err := parser.ParseRecipe(recipe)
	if err != nil {
		panic(err)
	}
	log.Printf("Parsed recipe_model: %v", model)

	switch m := (*model).(type) {
	case refactoringRunModel.RunModel:
		refactoringPipelineBuilder.BuildRunPipeline(&m)
	default:
		panic("unknown run model")
	}

	//	Build Pipeline

	return Pipeline{}
}

// TODO
type Pipeline struct {
}
