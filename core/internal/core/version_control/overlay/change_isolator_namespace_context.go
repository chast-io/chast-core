package overlay

type NamespaceContext struct {
	RootFolder          string
	ChangeCaptureFolder string
	OperationDirectory  string
	WorkingDirectory    string
	Command             []string
}

func NewNamespaceContext(rootFolder string,
	changeCaptureFolder string,
	operationDirectory string,
	workingDirectory string,
	command ...string) *NamespaceContext {
	return &NamespaceContext{
		rootFolder,
		changeCaptureFolder,
		operationDirectory,
		workingDirectory,
		command,
	}
}

func (ctx *NamespaceContext) convertFromStringArgs(args []string) {
	ctx.RootFolder = args[0]
	ctx.ChangeCaptureFolder = args[1]
	ctx.OperationDirectory = args[2]
	ctx.WorkingDirectory = args[3]
	ctx.Command = args[4:]
}

func newNamespaceContextFromStringArgs(args []string) *NamespaceContext {
	return &NamespaceContext{
		RootFolder:          args[0],
		ChangeCaptureFolder: args[1],
		OperationDirectory:  args[2],
		WorkingDirectory:    args[3],
		Command:             args[4:],
	}
}

func (ctx *NamespaceContext) toStringArgs() []string {
	return append([]string{ctx.RootFolder, ctx.ChangeCaptureFolder, ctx.OperationDirectory, ctx.WorkingDirectory}, ctx.Command...)
}
