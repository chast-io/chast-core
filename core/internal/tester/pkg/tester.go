package tester

import (
	"chast.io/core/internal/tester/internal/comparer"
	pathhandler "chast.io/core/internal/tester/internal/path_handler"
	"path/filepath"
	"strings"

	"chast.io/core/internal/internal_util/collection"
	chastlog "chast.io/core/internal/logger"
	recipemodel "chast.io/core/internal/recipe/pkg/model"
	"chast.io/core/internal/recipe/pkg/parser"
	refactoringservice "chast.io/core/internal/service/pkg/refactoring"
	util "chast.io/core/pkg/util/fs/file"
	"github.com/joomcode/errorx"
)

func Test(recipeFile *util.File) {
	parsedRecipe, recipeParseError := parser.ParseRecipe(recipeFile)
	if recipeParseError != nil {
		panic(recipeParseError)
	}

	switch concreteRecipe := (*parsedRecipe).(type) {
	case *recipemodel.RefactoringRecipe:
		workingDir := recipeFile.ParentDirectory

		if len(concreteRecipe.Tests) == 0 {
			chastlog.Log.Infof("No tests found for recipe %s", recipeFile.AbsolutePath)

			return
		}

		for index, test := range concreteRecipe.Tests {
			testWorkingDir := filepath.Join(workingDir, "tests", test.ID)

			args := pathhandler.AbsolutizePathArgs(concreteRecipe, test.Args, testWorkingDir)
			flags := pathhandler.AbsolutizePathFlags(concreteRecipe, convertFlags(test.Flags), testWorkingDir)

			pipeline, recipeRunError := refactoringservice.Run(
				recipeFile,
				args,
				flags,
			)

			if recipeRunError != nil {
				panic(recipeRunError)
			}

			comparer.CompareResults(&concreteRecipe.Tests[index], pipeline, workingDir)
		}
	default:
		panic(errorx.UnsupportedOperation.New("No run model builder for recipe of type %T", concreteRecipe.GetRecipeType()))
	}
}

func convertFlags(flags []string) []refactoringservice.FlagParameter {
	return collection.Map(flags, func(flag string) refactoringservice.FlagParameter {
		split := strings.Split(flag, "=")

		return refactoringservice.FlagParameter{
			Name:  split[0],
			Value: split[1],
		}
	})
}
