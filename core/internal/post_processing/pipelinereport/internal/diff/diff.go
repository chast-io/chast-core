package diff

import (
	"path/filepath"
	"strings"

	refactoringpipelinemodel "chast.io/core/internal/pipeline/pkg/model/refactoring"
	gitDiff "github.com/go-git/go-git/v5/utils/diff"
	"github.com/joomcode/errorx"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/spf13/afero"
	"github.com/ttacon/chalk"
)

const unionFsHiddenPathSuffix = "_HIDDEN~"

func BuildDiff(pipeline *refactoringpipelinemodel.Pipeline, changedFiles []string) (*ChangeDiff, error) {
	changeDiff := ChangeDiff{
		BaseFolder:   pipeline.ChangeCaptureLocation,
		ChangedFiles: changedFiles,
		Diffs:        make(map[string]FsDiff),
	}

	osFileSystem := afero.NewOsFs()

	for _, originalFilePath := range changedFiles {
		newFilePath := filepath.Join(pipeline.ChangeCaptureLocation, originalFilePath)

		if strings.HasSuffix(originalFilePath, unionFsHiddenPathSuffix) {
			changeDiff.Diffs[newFilePath] = FsDiff{FileStatus: Deleted, Diffs: nil}

			continue
		}

		originalExists, originalExistsError := afero.Exists(osFileSystem, originalFilePath)
		if originalExistsError != nil {
			return nil, errorx.ExternalError.Wrap(originalExistsError, "failed to check if original file exists")
		}

		if !originalExists {
			changeDiff.Diffs[newFilePath] = FsDiff{FileStatus: Added, Diffs: nil}

			continue
		}

		isDir, originalIsDirCheckError := afero.IsDir(osFileSystem, originalFilePath)
		if originalIsDirCheckError != nil {
			return nil, errorx.ExternalError.Wrap(originalIsDirCheckError, "failed to check if original file is a directory")
		}

		if isDir {
			changeDiff.Diffs[newFilePath] = FsDiff{FileStatus: Modified, Diffs: nil}

			continue
		}

		diffs := make([]FileDiff, 0)

		originalFileContent, originalReadError := afero.ReadFile(osFileSystem, originalFilePath)
		if originalReadError != nil {
			return nil, errorx.ExternalError.Wrap(originalReadError, "failed to read original file")
		}

		newFileContent, newReadError := afero.ReadFile(osFileSystem, newFilePath)
		if newReadError != nil {
			return nil, errorx.ExternalError.Wrap(newReadError, "failed to read new file")
		}

		fileDiff := gitDiff.Do(string(originalFileContent), string(newFileContent))

		for _, diff := range fileDiff {
			var convertedDiff FileDiff

			switch diff.Type {
			case diffmatchpatch.DiffEqual:
				convertedDiff = FileDiff{Equal, diff.Text}
			case diffmatchpatch.DiffInsert:
				convertedDiff = FileDiff{Insert, diff.Text}
			case diffmatchpatch.DiffDelete:
				convertedDiff = FileDiff{Delete, diff.Text}
			}

			diffs = append(diffs, convertedDiff)
		}

		changeDiff.Diffs[newFilePath] = FsDiff{FileStatus: Modified, Diffs: diffs}
	}

	return &changeDiff, nil
}

func (d *ChangeDiff) ToString(colorize bool) string {
	var stringBuilder strings.Builder

	for _, file := range d.ChangedFiles {
		fileDiff := d.Diffs[d.BaseFolder+file]
		if fileDiff.FileStatus == Modified {
			stringBuilder.WriteString("\n\n")
			stringBuilder.WriteString(file)
			stringBuilder.WriteString("\n\n")
		}

		for _, diff := range fileDiff.Diffs {
			switch diff.Type {
			case Equal:
				stringBuilder.WriteString(prefixEachLine("=", diff.Text))
			case Insert:
				if colorize {
					stringBuilder.WriteString(chalk.Green.Color(prefixEachLine("+", diff.Text)))
				} else {
					stringBuilder.WriteString(prefixEachLine("+", diff.Text))
				}
			case Delete:
				if colorize {
					stringBuilder.WriteString(chalk.Red.Color(prefixEachLine("-", diff.Text)))
				} else {
					stringBuilder.WriteString(prefixEachLine("-", diff.Text))
				}
			}

			stringBuilder.WriteString("\n")
		}
	}

	return stringBuilder.String()
}

func prefixEachLine(prefix string, text string) string {
	return prefix + strings.ReplaceAll(text, "\n", "\n"+prefix)
}
