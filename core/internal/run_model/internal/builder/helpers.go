package builder

import (
	"path/filepath"
	"strings"

	recipemodel "chast.io/core/internal/recipe/pkg/model"
	"github.com/joomcode/errorx"
)

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

func verifyPathExtension(value string, extensions []string) error {
	if len(extensions) == 0 {
		return nil
	}

	for _, extension := range extensions {
		if strings.HasSuffix(value, extension) {
			return nil
		}
	}

	return errorx.IllegalFormat.New("Path %s does not have a valid extension. Valid extensions: %v", value, extensions)
}
