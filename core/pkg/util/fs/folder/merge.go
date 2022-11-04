package folder

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

func MergeFolders(sourceFolders []string, targetFolder string, blockOverwrite bool) error {
	if err := os.MkdirAll(targetFolder, 0777); err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to create target folder \"%s\"", targetFolder))
	}

	for _, sourceFolder := range sourceFolders {
		if err := moveFolderContents(sourceFolder, targetFolder, blockOverwrite); err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to merge folder \"%s\" with \"%s\"", sourceFolder, targetFolder))
		}
	}

	return nil
}

var errMergeOverwriteBlock = errors.New(
	"Error due to attempting to merge a file over an existing file in blockOverwrite mode",
)

func moveFolderContents(sourceFolder string, targetFolder string, blockOverwrite bool) error {
	osFileSystem := afero.NewOsFs()
	if walkError := afero.Walk(osFileSystem, sourceFolder, func(path string, info fs.FileInfo, _ error) error {
		if !info.IsDir() {
			done, err := moveFile(path, sourceFolder, targetFolder, osFileSystem, blockOverwrite)
			if done {
				return err
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
) (bool, error) {
	correctedPath := strings.TrimPrefix(path, sourceFolder)
	targetFilePath := filepath.Join(targetFolder, correctedPath)

	exists, existsError := afero.Exists(osFileSystem, targetFilePath)
	if existsError != nil {
		return true, errors.Wrap(existsError, "Failed to check if file exists")
	}

	if exists {
		if blockOverwrite {
			return true, errMergeOverwriteBlock
		}

		log.Debugf("File overwritten during merge of folders. Affected File: %s", targetFilePath)

		if err := osFileSystem.Remove(targetFilePath); err != nil {
			return true, errors.Wrap(err, "Failed to remove file")
		}
	}

	if err := os.MkdirAll(filepath.Dir(targetFilePath), 0755); err != nil {
		return true, errors.Wrap(err, "Failed to create target directory")
	}

	if err := osFileSystem.Rename(path, targetFilePath); err != nil {
		return true, errors.Wrap(err, "Failed to move file")
	}

	return false, nil
}
