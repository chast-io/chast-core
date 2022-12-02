package pathutils

import (
	"os"
	"path/filepath"
	"strings"

	"chast.io/core/internal/post_processing/merger/pkg/mergeoptions"
	"github.com/joomcode/errorx"
)

func TargetPath(path string, sourceFolder string, targetFolder string) string {
	correctedPath := strings.TrimPrefix(path, sourceFolder)
	targetPath := filepath.Join(targetFolder, correctedPath)

	return targetPath
}

func CleanupPath(path string) error {
	if err := os.RemoveAll(path); err != nil {
		return errorx.ExternalError.Wrap(err, "failed to remove path")
	}

	return nil
}

func IsInMetaFolder(sourcePath string, sourceRootFolder string, options *mergeoptions.MergeOptions) bool {
	trimmedSourcePath := strings.TrimPrefix(sourcePath, sourceRootFolder)

	return strings.HasPrefix(trimmedSourcePath, "/"+options.MetaFilesLocation) ||
		strings.HasPrefix(trimmedSourcePath, options.MetaFilesLocation)
}
