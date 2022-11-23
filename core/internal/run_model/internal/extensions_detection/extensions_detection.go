package extensionsdetection

import (
	chastlog "chast.io/core/internal/logger"
	"os"
	"path/filepath"
	"strings"

	"github.com/joomcode/errorx"
	"github.com/spf13/afero"
)

type Extension struct {
	Name             string
	Count            int
	CommonParentPath string
}

func DetectExtensions(rootPath string) (map[string]*Extension, error) {
	extensions := make(map[string]*Extension)

	chastlog.Log.Tracef("Detecting extensions of files in path: %s", rootPath)

	osFs := afero.NewOsFs()
	if err := afero.Walk(osFs, rootPath, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return errorx.ExternalError.New("Path does not exist")
		}

		if info.IsDir() {
			return nil
		}

		extension := filepath.Ext(path)
		extension = strings.TrimPrefix(extension, ".")

		if _, ok := extensions[extension]; !ok {
			extensions[extension] = &Extension{
				Name:             extension,
				Count:            1,
				CommonParentPath: path,
			}

			return nil
		}

		extensions[extension].Count++
		extensions[extension].CommonParentPath = commonParentPath(extensions[extension].CommonParentPath, path)

		return nil
	}); err != nil {
		return nil, errorx.ExternalError.Wrap(err, "Could not walk through directory")
	}

	return extensions, nil
}

// returns the common parent path of two given paths.
func commonParentPath(path1 string, path2 string) string {
	path1, _ = filepath.Abs(path1)
	path2, _ = filepath.Abs(path2)

	path1 = filepath.ToSlash(path1)
	path2 = filepath.ToSlash(path2)

	path1Parts := strings.Split(path1, "/")
	path2Parts := strings.Split(path2, "/")

	commonParts := []string{"/"}

	for i := 0; i < len(path1Parts) && i < len(path2Parts); i++ {
		if path1Parts[i] != path2Parts[i] {
			break
		}

		commonParts = append(commonParts, path1Parts[i])
	}

	return filepath.Join(commonParts...)
}
