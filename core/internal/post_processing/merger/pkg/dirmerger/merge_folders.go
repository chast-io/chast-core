package dirmerger

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

const unionFsHiddenPathSuffix = "_HIDDEN~"
const defaultFolderPermission = 0777

var errMergeOverwriteBlock = errors.New(
	"Error due to attempting to merge a file over an existing file in blockOverwrite mode",
)

func MergeFolders(sourceFolders []string, targetFolder string, blockOverwrite bool) error {
	if err := os.MkdirAll(targetFolder, defaultFolderPermission); err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to create target folder \"%s\"", targetFolder))
	}

	for _, sourceFolder := range sourceFolders {
		if err := moveFolderContents(sourceFolder, targetFolder, blockOverwrite); err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to merge folder \"%s\" with \"%s\"", sourceFolder, targetFolder))
		}
	}

	return nil
}

func moveFolderContents(sourceFolder string, targetFolder string, blockOverwrite bool) error {
	osFileSystem := afero.NewOsFs()

	if exists, err := afero.Exists(osFileSystem, sourceFolder); err != nil || !exists {
		return nil //nolint:nilerr // If the folder does not exist, ignore it and continue
	}

	if walkError := afero.Walk(osFileSystem, sourceFolder, func(path string, info fs.FileInfo, _ error) error {
		if info == nil {
			return nil
		}

		if info.IsDir() {
			if err := createFolder(path, sourceFolder, targetFolder, osFileSystem, blockOverwrite); err != nil {
				return errors.Wrap(err, "Failed to move folder")
			}
		} else {
			if err := moveFile(path, sourceFolder, targetFolder, osFileSystem, blockOverwrite); err != nil {
				return errors.Wrap(err, "Failed to move file")
			}
		}

		return nil
	}); walkError != nil {
		return errors.Wrap(walkError, "Failed to walk through source folder")
	}

	if err := os.RemoveAll(sourceFolder); err != nil {
		return errors.Wrap(err, "failed to remove merge source directory")
	}

	return nil
}

func moveFile(
	path string,
	sourceFolder string,
	targetFolder string,
	osFileSystem afero.Fs,
	blockOverwrite bool,
) error {
	targetPath := targetPath(path, sourceFolder, targetFolder)

	if err := handlePossibleMarkedAsDeletedPath(targetPath, osFileSystem, blockOverwrite); err != nil {
		return errors.Wrap(err, "Failed to handle possible marked as deleted path")
	}

	exists, existsError := afero.Exists(osFileSystem, targetPath)
	if existsError != nil {
		return errors.Wrap(existsError, "Failed to check if path exists")
	}

	if exists {
		if blockOverwrite {
			return errMergeOverwriteBlock
		}

		log.Debugf("File overwritten during merge of folders. Affected File: %s", targetPath)

		if err := osFileSystem.Remove(targetPath); err != nil {
			return errors.Wrap(err, "Failed to remove file")
		}
	}

	if err := os.MkdirAll(filepath.Dir(targetPath), defaultFolderPermission); err != nil {
		return errors.Wrap(err, "Failed to create target directory")
	}

	if err := osFileSystem.Rename(path, targetPath); err != nil {
		return errors.Wrap(err, "Failed to move file")
	}

	return nil
}

func createFolder(
	path string,
	sourceFolder string,
	targetFolder string,
	osFileSystem afero.Fs,
	blockOverwrite bool,
) error {
	targetPath := targetPath(path, sourceFolder, targetFolder)

	if err := handlePossibleMarkedAsDeletedPath(targetPath, osFileSystem, blockOverwrite); err != nil {
		return errors.Wrap(err, "Failed to handle possible marked as deleted path")
	}

	if err := osFileSystem.MkdirAll(targetPath, defaultFolderPermission); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to create folder \"%s\"", targetPath))
	}

	return nil
}

func targetPath(path string, sourceFolder string, targetFolder string) string {
	correctedPath := strings.TrimPrefix(path, sourceFolder)
	targetPath := filepath.Join(targetFolder, correctedPath)

	return targetPath
}

func handlePossibleMarkedAsDeletedPath(targetPath string, osFileSystem afero.Fs, blockOverwrite bool) error {
	deletedTargetPath := targetPath + unionFsHiddenPathSuffix

	exists, existsError := afero.Exists(osFileSystem, deletedTargetPath)
	if existsError != nil {
		return errors.Wrap(existsError, "Failed to check if path exists")
	}

	if exists {
		if blockOverwrite {
			return errMergeOverwriteBlock
		}

		if err := osFileSystem.RemoveAll(deletedTargetPath); err != nil {
			return errors.Wrap(err, "Failed to remove marked-as-deleted flag file")
		}
	}

	return nil
}
