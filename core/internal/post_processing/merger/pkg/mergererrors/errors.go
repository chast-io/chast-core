package mergererrors

import "github.com/joomcode/errorx"

var ErrMergeOverwriteBlock = errorx.InternalError.New(
	"Error due to attempting to merge a file over an existing file in blockOverwrite mode",
)
