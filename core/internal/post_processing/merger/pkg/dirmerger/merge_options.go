package dirmerger

import (
	"os"
)

const unionFsMetaFolder = ".unionfs-fuse"
const unionFsHiddenPathSuffix = "_HIDDEN~"
const defaultFolderPermission = 0755

type MergeOptions struct {
	DryRun                     bool
	BlockOverwrite             bool
	MergeMetaFilesFolder       bool
	DeleteEmptyFolders         bool
	DeleteMarkedAsDeletedPaths bool
	MetaFilesLocation          string
	MetaFilesDeletedExtension  string
	FolderPermission           os.FileMode
	Exclusions                 []*WildcardString
	Inclusions                 []*WildcardString
}

func NewMergeOptions() *MergeOptions {
	return &MergeOptions{
		DryRun: false,

		BlockOverwrite:             false,
		MergeMetaFilesFolder:       true,
		DeleteEmptyFolders:         false,
		DeleteMarkedAsDeletedPaths: false,

		MetaFilesLocation:         unionFsMetaFolder,
		MetaFilesDeletedExtension: unionFsHiddenPathSuffix,
		FolderPermission:          defaultFolderPermission,

		Exclusions: []*WildcardString{},
		Inclusions: []*WildcardString{},
	}
}

func (o *MergeOptions) ShouldSkip(location string) bool {
	if len(o.Inclusions) > 0 {
		hasMatch := false

		for _, includeLocation := range o.Inclusions {
			if includeLocation.Matches(location) {
				hasMatch = true

				break
			}
		}

		if !hasMatch {
			return true
		}
	}

	// continue to check if it is excluded

	for _, skipLocation := range o.Exclusions {
		if skipLocation.MatchesPath(location) {
			return true
		}
	}

	return false
}
