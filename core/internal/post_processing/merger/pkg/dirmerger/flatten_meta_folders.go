package dirmerger

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	chastlog "chast.io/core/internal/logger"
	"github.com/joomcode/errorx"
	"github.com/spf13/afero"
)

func flattenMetaFolder(sourceFolder string, options *MergeOptions) error {
	metaSourceFolder := filepath.Join(sourceFolder, options.MetaFilesLocation)

	metaFolderOptions := *options
	metaFolderOptions.BlockOverwrite = true
	metaFolderOptions.CopyMode = false

	osFileSystem := afero.NewOsFs()

	if walkError := afero.Walk(osFileSystem, metaSourceFolder, func(path string, info fs.FileInfo, _ error) error {
		if info == nil {
			return nil // exit if the file does not exist - this can happen due to a removal in a previous iteration
		}

		if path == sourceFolder {
			return nil // skip the root folder
		}

		sanitizedPath := strings.TrimPrefix(path, metaSourceFolder)
		if metaFolderOptions.ShouldSkip(sanitizedPath) {
			chastlog.Log.Tracef("Skipping path \"%s\" due to being excluded or not included", path)

			return nil // exit if the path is to be skipped
		}

		if info.IsDir() {
			if err := moveMetaFolder(path, metaSourceFolder, sourceFolder, osFileSystem, &metaFolderOptions); err != nil {
				return errorx.InternalError.Wrap(err, "Failed to move folder")
			}
		} else {
			if err := moveMetaFile(path, metaSourceFolder, sourceFolder, osFileSystem, &metaFolderOptions); err != nil {
				return errorx.InternalError.Wrap(err, "Failed to move file")
			}
		}

		return nil
	}); walkError != nil {
		return errorx.ExternalError.Wrap(walkError, "Failed to walk through source folder")
	}

	if !options.DryRun {
		if err := cleanupPath(metaSourceFolder, osFileSystem, options); err != nil {
			return err
		}
	}

	if err := sanitizeMovedMetaPaths(sourceFolder, options); err != nil {
		return errorx.InternalError.Wrap(err, "Failed to sanitize moved meta paths")
	}

	return nil
}

func moveMetaFolder(
	sourcePath string,
	sourceRootFolder string,
	targetRootFolder string,
	osFileSystem afero.Fs,
	options *MergeOptions) error {
	if exists, err := afero.Exists(osFileSystem, sourcePath); err != nil || !exists {
		return nil //nolint:nilerr // If the folder does not exist, ignore it and continue
	}

	isEmpty, isEmptyCheckError := afero.IsEmpty(osFileSystem, sourcePath)
	if isEmptyCheckError != nil {
		return errorx.ExternalError.Wrap(isEmptyCheckError, "Failed to check if folder is empty")
	}

	if !isEmpty {
		chastlog.Log.Tracef("Folder \"%s\" is not empty, skipping -> will be handled later", sourcePath)

		return nil
	}

	targetPath := targetPath(sourcePath, sourceRootFolder, targetRootFolder)

	if !options.DryRun {
		if err := osFileSystem.MkdirAll(targetPath, options.FolderPermission); err != nil {
			return errorx.ExternalError.Wrap(err, fmt.Sprintf("Failed to create folder \"%s\"", targetPath))
		}

		if !options.CopyMode {
			if err := cleanupPath(sourcePath, osFileSystem, options); err != nil {
				return err
			}
		}
	}

	return nil
}

func moveMetaFile(
	sourcePath string,
	sourceRootFolder string,
	targetRootFolder string,
	osFileSystem afero.Fs,
	options *MergeOptions,
) error {
	if exists, err := afero.Exists(osFileSystem, sourcePath); err != nil || !exists {
		return nil //nolint:nilerr // If the folder does not exist, ignore it and continue
	}

	targetPath := targetPath(sourcePath, sourceRootFolder, targetRootFolder)

	if options.DryRun {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(targetPath), options.FolderPermission); err != nil {
		return errorx.ExternalError.Wrap(err, "Failed to create target directory")
	}

	if err := osFileSystem.Rename(sourcePath, targetPath); err != nil {
		return errorx.ExternalError.Wrap(err, "Failed to move file")
	}

	return nil
}

func sanitizeMovedMetaPaths(sourcePath string, options *MergeOptions) error {
	osFileSystem := afero.NewOsFs()

	if walkError := afero.Walk(osFileSystem, sourcePath, func(path string, info fs.FileInfo, _ error) error {
		if info == nil {
			return nil // exit if the file does not exist - this can happen due to a removal in a previous iteration
		}

		if info.IsDir() {
			if strings.HasSuffix(path, options.MetaFilesDeletedExtension) {
				correspondingFolder := strings.TrimSuffix(path, options.MetaFilesDeletedExtension)
				if exists, err := afero.Exists(osFileSystem, correspondingFolder); err != nil || !exists {
					return nil //nolint:nilerr // If the folder does not exist, ignore it and continue
				}

				isEmpty, isEmptyCheckError := afero.IsEmpty(osFileSystem, path)
				if isEmptyCheckError != nil {
					return errorx.ExternalError.Wrap(isEmptyCheckError, "Failed to check if folder is empty")
				}

				if !isEmpty {
					return errorx.InternalError.New("Meta folder is not empty. This should not happen!")
				}

				if err := os.RemoveAll(path); err != nil {
					return errorx.ExternalError.Wrap(err, "Failed to remove folder")
				}

				if err := os.Rename(correspondingFolder, path); err != nil {
					return errorx.ExternalError.Wrap(err, "Failed to rename folder")
				}
			}
		}

		return nil
	}); walkError != nil {
		return errorx.ExternalError.Wrap(walkError, "Failed to walk through source folder")
	}

	return nil
}
