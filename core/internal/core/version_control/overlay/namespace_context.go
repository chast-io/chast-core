package overlay

type ToArgsConverter interface {
	toStringArgs() []string
}

type FromArgsConverter interface {
	convertFromStringArgs(args []string)
}

type ArgsConverter interface {
	ToArgsConverter
	FromArgsConverter
}
