package refactoring

import (
	refactoringService "chast.io/core/internal/service/pkg/refactoring"
	util "chast.io/core/pkg/util/fs/file"
	log "github.com/sirupsen/logrus"
)

func Run(recipe *util.File, args ...string) {
	err := refactoringService.Run(recipe, args...)
	if err != nil {
		log.Fatalln(err)
	}
}
