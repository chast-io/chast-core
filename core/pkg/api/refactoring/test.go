package refactoring

import (
	tester "chast.io/core/internal/tester/pkg"
	util "chast.io/core/pkg/util/fs/file"
)

func Test(recipe *util.File, args ...string) {
	tester.Test(recipe)
}
