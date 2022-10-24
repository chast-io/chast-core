package version_control

type VersionControl struct {
	versions            []Version
	currentVersionIndex int
	basePath            string
}

func NewVersionControl(basePath string) *VersionControl {
	return &VersionControl{basePath: basePath, versions: []Version{
		{source: basePath, target: basePath},
	}}
}

func (vc *VersionControl) CurrentVersion() Version {
	return vc.versions[vc.currentVersionIndex]
}

func (vc *VersionControl) StartNewVersion(versionTag string) Version {
	version := *NewVersion(
		vc.basePath,
		vc.CurrentVersion().target,
		versionTag)
	vc.versions = append(vc.versions, version)
	return version
}
