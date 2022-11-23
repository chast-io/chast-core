package refactoring

import (
	chastlog "chast.io/core/internal/logger"
	refactoringService "chast.io/core/internal/service/pkg/refactoring"
	util "chast.io/core/pkg/util/fs/file"
	"github.com/joomcode/errorx"
)

func Run(recipe *util.File, args ...string) {
	err := refactoringService.Run(recipe, args, nil)
	if err != nil {
		chastlog.Log.Fatalf("%+v", errorx.EnsureStackTrace(err))
	}
}
