package pipelinereport

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	refactoringpipelinemodel "chast.io/core/internal/pipeline/pkg/model/refactoring"
	"chast.io/core/internal/post_processing/pipelinereport/internal/diff"
	filetree "chast.io/core/internal/post_processing/pipelinereport/internal/tree"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

type Report struct {
	ChangedPaths []string
	ChangeDiff   *diff.ChangeDiff
	Pipeline     *refactoringpipelinemodel.Pipeline
}

func BuildReport(pipeline *refactoringpipelinemodel.Pipeline) (*Report, error) {
	changedPaths := make([]string, 0)

	osFileSystem := afero.NewOsFs()
	if walkError := afero.Walk(osFileSystem, pipeline.ChangeCaptureFolder,
		func(path string, info fs.FileInfo, _ error) error {
			if info.IsDir() {
				folderIsEmpty, isEmptyCheckError := afero.IsEmpty(osFileSystem, path)
				if isEmptyCheckError != nil {
					return errors.Wrap(isEmptyCheckError, "failed to check if folder is empty")
				}
				if folderIsEmpty {
					correctedPath := strings.TrimPrefix(path, pipeline.ChangeCaptureFolder)
					changedPaths = append(changedPaths, correctedPath)
				}
			} else {
				correctedPath := strings.TrimPrefix(path, pipeline.ChangeCaptureFolder)
				changedPaths = append(changedPaths, correctedPath)
			}

			return nil
		},
	); walkError != nil {
		return nil, errors.Wrap(walkError, "Failed to walk change capture folder")
	}

	changeDiff, diffBuildError := diff.BuildDiff(pipeline, changedPaths)
	if diffBuildError != nil {
		return nil, errors.Wrap(diffBuildError, "failed to build diffs")
	}

	return &Report{
		ChangedPaths: changedPaths,
		ChangeDiff:   changeDiff,
		Pipeline:     pipeline,
	}, nil

}

func (report *Report) ChangedFilesRelative() ([]string, error) {
	changedFilesRelative := make([]string, 0)

	workingDirPath, getWDError := os.Getwd()
	if getWDError != nil {
		return nil, errors.Wrap(getWDError, "failed get working directory")
	}

	for _, filePath := range report.ChangedPaths {
		rel, relError := filepath.Rel(workingDirPath, filePath)
		if relError != nil {
			return nil, errors.Wrap(relError, "failed getting relative path")
		}

		changedFilesRelative = append(changedFilesRelative, rel)
	}

	return changedFilesRelative, nil
}

func (report *Report) PrintFileTree(colorize bool) {
	println(filetree.ToString(report.Pipeline.ChangeCaptureFolder, report.ChangeDiff, false, colorize))
}

func (report *Report) PrintChanges(colorize bool) {
	println(report.ChangeDiff.ToString(colorize))
}
