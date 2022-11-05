package diff

import (
	"os"
	"path/filepath"

	refactoringpipelinemodel "chast.io/core/internal/pipeline/pkg/model/refactoring"
	"chast.io/core/internal/post_processing/pipelinereport"
	"github.com/go-git/go-git/v5/utils/diff"
	"github.com/pkg/errors"
)

func BuildDiff(pipeline *refactoringpipelinemodel.Pipeline, report *pipelinereport.Report) error {
	files := report.ChangedFiles

	for _, originalFilePath := range files {
		newFilePath := filepath.Join(pipeline.ChangeCaptureFolder, originalFilePath)

		originalFileContent, originalReadError := os.ReadFile(originalFilePath)
		if originalReadError != nil {
			return errors.Wrap(originalReadError, "failed to read original file")
		}

		newFileContent, newReadError := os.ReadFile(newFilePath)
		if newReadError != nil {
			return errors.Wrap(newReadError, "failed to read new file")
		}

		fileDiff := diff.Do(string(originalFileContent), string(newFileContent))

		println(fileDiff)
	}

	return nil
}
