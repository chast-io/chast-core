package refactoring

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	chastlog "chast.io/core/internal/logger"
	refactoringService "chast.io/core/internal/service/pkg/refactoring"
	util "chast.io/core/pkg/util/fs/file"
	"github.com/joomcode/errorx"
)

func Run(recipe *util.File, args ...string) {
	pipeline, runError := refactoringService.Run(recipe, args, nil)
	if runError != nil {
		chastlog.Log.Fatalf("%+v", errorx.EnsureStackTrace(runError))
	}

	if err := refactoringService.ShowReport(pipeline); err != nil {
		chastlog.Log.Fatalf("%+v", errorx.EnsureStackTrace(err))
	}

	result := StringPrompt("Do you want to apply the refactoring? (y/N)")

	if result == "y" || result == "Y" {
		if err := refactoringService.ApplyChanges(pipeline); err != nil {
			chastlog.Log.Fatalf("%+v", errorx.EnsureStackTrace(err))
		}
	}
}

func StringPrompt(label string) string {
	var str string

	reader := bufio.NewReader(os.Stdin)

	for {
		_, _ = fmt.Fprint(os.Stderr, label+" ")

		str, _ = reader.ReadString('\n')

		if str != "" {
			break
		}
	}

	return strings.TrimSpace(str)
}
