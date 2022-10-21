package api

type Runner interface {
	Run(command Command) (string, error)
}

type Command struct {
	cmd string
}
