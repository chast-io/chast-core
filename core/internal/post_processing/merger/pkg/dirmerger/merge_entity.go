package dirmerger

import "chast.io/core/internal/run_model/pkg/model/refactoring"

type MergeEntity struct {
	SourcePath      string
	ChangeLocations *refactoring.ChangeLocations
}

func NewMergeEntity(sourcePath string, changeLocations *refactoring.ChangeLocations) MergeEntity {
	if changeLocations == nil {
		changeLocations = &refactoring.ChangeLocations{
			Include: make([]string, 0),
			Exclude: make([]string, 0),
		}
	}

	return MergeEntity{
		SourcePath:      sourcePath,
		ChangeLocations: changeLocations,
	}
}
