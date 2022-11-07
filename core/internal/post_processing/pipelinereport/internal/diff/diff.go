package diff

import (
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/spf13/afero"
	"github.com/ttacon/chalk"
	"path/filepath"
	"strings"

	refactoringpipelinemodel "chast.io/core/internal/pipeline/pkg/model/refactoring"
	gitDiff "github.com/go-git/go-git/v5/utils/diff"
	"github.com/pkg/errors"
)

const unionFsHiddenPathSuffix = "_HIDDEN~"

func BuildDiff(pipeline *refactoringpipelinemodel.Pipeline, changedFiles []string) (*ChangeDiff, error) {
	changeDiff := ChangeDiff{BaseFolder: pipeline.ChangeCaptureFolder, ChangedFiles: changedFiles, Diffs: make(map[string]FsDiff)}
	osFileSystem := afero.NewOsFs()

	for _, originalFilePath := range changedFiles {
		newFilePath := filepath.Join(pipeline.ChangeCaptureFolder, originalFilePath)

		if strings.HasSuffix(originalFilePath, unionFsHiddenPathSuffix) {
			changeDiff.Diffs[newFilePath] = FsDiff{FileStatus: Deleted}
			continue
		}

		originalExists, originalExistsError := afero.Exists(osFileSystem, originalFilePath)
		if originalExistsError != nil {
			return nil, errors.Wrap(originalExistsError, "failed to check if original file exists")
		}

		if !originalExists {
			changeDiff.Diffs[newFilePath] = FsDiff{FileStatus: Added}
			continue
		}

		diffs := make([]FileDiff, 0)

		originalFileContent, originalReadError := afero.ReadFile(osFileSystem, originalFilePath)
		if originalReadError != nil {
			return nil, errors.Wrap(originalReadError, "failed to read original file")
		}

		newFileContent, newReadError := afero.ReadFile(osFileSystem, newFilePath)
		if newReadError != nil {
			return nil, errors.Wrap(newReadError, "failed to read new file")
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
	var sb strings.Builder

	for _, file := range d.ChangedFiles {
		fileDiff := d.Diffs[d.BaseFolder+file]
		if fileDiff.FileStatus == Modified {
			sb.WriteString("\n\n")
			sb.WriteString(file)
			sb.WriteString("\n\n")
		}
		for _, diff := range fileDiff.Diffs {
			switch diff.Type {
			case Equal:
				sb.WriteString(prefixEachLine("=", diff.Text))
			case Insert:
				if colorize {
					sb.WriteString(chalk.Green.Color(prefixEachLine("+", diff.Text)))
				} else {
					sb.WriteString(prefixEachLine("+", diff.Text))
				}
			case Delete:
				if colorize {
					sb.WriteString(chalk.Red.Color(prefixEachLine("-", diff.Text)))
				} else {
					sb.WriteString(prefixEachLine("-", diff.Text))
				}
			}
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

func prefixEachLine(prefix string, text string) string {
	return prefix + strings.Replace(text, "\n", "\n"+prefix, -1)
}
