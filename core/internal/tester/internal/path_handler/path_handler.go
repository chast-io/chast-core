package pathhandler

import (
	"path/filepath"
	"strings"

	recipemodel "chast.io/core/internal/recipe/pkg/model"
	refactoringservice "chast.io/core/internal/service/pkg/refactoring"
	"github.com/joomcode/errorx"
)

func AbsolutizePathArgs(recipe *recipemodel.RefactoringRecipe, args []string, workingDir string) []string {
	if len(args) == 0 {
		return args
	}

	convertedArgs := make([]string, len(args))

	pathsRootFolder := filepath.Join(workingDir, "input")

	primaryPath, absPrimPathErr := absolutizePath(args[0], recipe.PrimaryParameter.TypeExtension, pathsRootFolder)
	if absPrimPathErr != nil {
		panic(absPrimPathErr)
	}

	convertedArgs[0] = primaryPath

	index := 1
	for index < len(args) {
		arg := args[index]
		param := recipe.PositionalParameters[index-1].TypeExtension

		path, absPathErr := absolutizePath(arg, param, pathsRootFolder)
		if absPathErr != nil {
			panic(absPathErr)
		}

		convertedArgs[index] = path

		index++
	}

	return convertedArgs
}

func AbsolutizePathFlags(
	recipe *recipemodel.RefactoringRecipe,
	flags []refactoringservice.FlagParameter,
	workingDir string,
) []refactoringservice.FlagParameter {
	if len(flags) == 0 {
		return flags
	}

	convertedFlags := make([]refactoringservice.FlagParameter, len(flags))

	pathsRootFolder := filepath.Join(workingDir, "input")

	index := 0
	for index < len(flags) {
		path, absPathErr := absolutizePath(flags[index].Value, recipe.Flags[index].TypeExtension, pathsRootFolder)
		if absPathErr != nil {
			panic(absPathErr)
		}

		convertedFlags[index] = refactoringservice.FlagParameter{
			Name:  flags[index].Name,
			Value: path,
		}

		index++
	}

	return convertedFlags
}

func absolutizePath(path string, typeExtension recipemodel.TypeExtension, wordingDir string) (string, error) {
	if strings.HasSuffix(typeExtension.Type, "Path") && !strings.HasPrefix(path, "/") {
		abs, err := filepath.Abs(filepath.Join(wordingDir, path))

		if err != nil {
			return path, errorx.ExternalError.Wrap(err, "Could not absolutize path")
		}

		return abs, nil
	}

	return path, nil
}
