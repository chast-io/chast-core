package diff

type ChangeDiff struct {
	BaseFolder   string
	ChangedFiles []string
	Diffs        map[string]FsDiff
}

type FsDiff struct {
	FileStatus FileStatus
	Diffs      []FileDiff
}

type FileDiff struct {
	Type Operation
	Text string
}

type Operation int8
type FileStatus int8

const (
	Delete Operation = -1
	Insert Operation = 1
	Equal  Operation = 0
)

const (
	Added    FileStatus = 1
	Modified FileStatus = 2
	Deleted  FileStatus = 3
)
