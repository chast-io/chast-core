package pipelinereport

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	refactoringpipelinemodel "chast.io/core/internal/pipeline/pkg/model/refactoring"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

type Report struct {
	ChangedFiles []string
}

func BuildReport(pipeline *refactoringpipelinemodel.Pipeline) (*Report, error) {
	changedFiles := make([]string, 0)

	osFileSystem := afero.NewOsFs()
	if walkError := afero.Walk(osFileSystem, pipeline.ChangeCaptureFolder,
		func(path string, info fs.FileInfo, _ error) error {
			if !info.IsDir() {
				correctedPath := strings.TrimPrefix(path, pipeline.ChangeCaptureFolder)
				changedFiles = append(changedFiles, correctedPath)
			}

			return nil
		},
	); walkError != nil {
		return nil, errors.Wrap(walkError, "Failed to walk change capture folder")
	}

	return &Report{
		ChangedFiles: changedFiles,
	}, nil
}

func (report *Report) ChangedFilesRelative() ([]string, error) {
	changedFilesRelative := make([]string, 0)

	workingDirPath, getWDError := os.Getwd()
	if getWDError != nil {
		return nil, errors.Wrap(getWDError, "failed get working directory")
	}

	for _, filePath := range report.ChangedFiles {
		rel, relError := filepath.Rel(workingDirPath, filePath)
		if relError != nil {
			return nil, errors.Wrap(relError, "failed getting relative path")
		}

		changedFilesRelative = append(changedFilesRelative, rel)
	}

	return changedFilesRelative, nil
}
