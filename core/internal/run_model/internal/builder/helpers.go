package builder

import (
	recipemodel "chast.io/core/internal/recipe/pkg/model"
	"github.com/pkg/errors"
	"path/filepath"
	"strings"
)

func absolutizePath(flagValue string, typeExtension recipemodel.TypeExtension, wordingDir string) (string, error) {
	if strings.HasSuffix(typeExtension.Type, "Path") && !strings.HasPrefix(flagValue, "/") {
		return filepath.Abs(filepath.Join(wordingDir, flagValue))
	}

	return flagValue, nil
}

func verifyPathExtension(value string, extensions []string) error {
	if len(extensions) == 0 {
		return nil
	}

	for _, extension := range extensions {
		if strings.HasSuffix(value, extension) {
			return nil
		}
	}

	return errors.Errorf("Path %s does not have a valid extension. Valid extensions: %v", value, extensions)
}
