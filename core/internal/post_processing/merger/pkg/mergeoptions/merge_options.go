package mergeoptions

import (
	wildcardstring "chast.io/core/internal/internal_util/wildcard_string"
	"os"
	"strings"
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
	CopyMode                   bool
	MetaFilesLocation          string
	MetaFilesDeletedExtension  string
	FolderPermission           os.FileMode
	Exclusions                 []*wildcardstring.WildcardString
	Inclusions                 []*wildcardstring.WildcardString
}

func NewMergeOptions() *MergeOptions {
	return &MergeOptions{
		DryRun: false,

		BlockOverwrite:             false,
		MergeMetaFilesFolder:       false,
		DeleteEmptyFolders:         false,
		DeleteMarkedAsDeletedPaths: false,
		CopyMode:                   false,

		MetaFilesLocation:         unionFsMetaFolder,
		MetaFilesDeletedExtension: unionFsHiddenPathSuffix,
		FolderPermission:          defaultFolderPermission,

		Exclusions: []*wildcardstring.WildcardString{},
		Inclusions: []*wildcardstring.WildcardString{},
	}
}

func (o *MergeOptions) ShouldSkip(location string) bool {
	// TODO add test
	cleanedLocation := strings.ReplaceAll(location, o.MetaFilesDeletedExtension, "")

	if len(o.Inclusions) > 0 {
		hasMatch := false

		for _, includeLocation := range o.Inclusions {
			if includeLocation.Matches(cleanedLocation) {
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
		if skipLocation.MatchesPath(cleanedLocation) {
			return true
		}
	}

	return false
}
