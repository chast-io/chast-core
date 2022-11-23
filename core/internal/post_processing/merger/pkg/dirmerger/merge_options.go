package dirmerger

import "os"

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
	SkipLocations              []WildcardString
	SkipExtensions             []WildcardString
	IncludeLocations           []WildcardString
	IncludeExtensions          []WildcardString
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

		SkipLocations:  []WildcardString{},
		SkipExtensions: []WildcardString{},

		IncludeLocations:  []WildcardString{},
		IncludeExtensions: []WildcardString{},
	}
}
