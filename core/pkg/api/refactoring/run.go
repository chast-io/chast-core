package refactoring

import util "chast.io/core/pkg/util"
import "chast.io/core/internal/service"

func Run(recipe util.FileReader, args ...string) {
	service.Run(recipe, args...)
}
