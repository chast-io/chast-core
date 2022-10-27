package refactoring

import (
	util "chast.io/core/pkg/util"
	log "github.com/sirupsen/logrus"
)
import "chast.io/core/internal/service"

func Run(recipe *util.File, args ...string) {
	_, err := service.Run(recipe, args...)
	if err != nil {
		log.Fatalln(err)
	}
}
