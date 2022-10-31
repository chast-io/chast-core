package overlay

type NamespaceContext struct {
	RootFolder          string
	MergeFolders        []string
	ChangeCaptureFolder string
	OperationDirectory  string
	WorkingDirectory    string
	Commands            [][]string
}

func NewNamespaceContext(
	rootFolder string,
	mergeFolders []string,
	changeCaptureFolder string,
	operationDirectory string,
	workingDirectory string,
	command [][]string) *NamespaceContext {
	return &NamespaceContext{
		rootFolder,
		mergeFolders,
		changeCaptureFolder,
		operationDirectory,
		workingDirectory,
		command,
	}
}
