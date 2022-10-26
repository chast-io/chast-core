package run_models

type RunModel interface {
}

type RunModelArgsMerger interface {
	MergeArgsIntoRunModel(runModel *RunModel, args ...string) (*RunModel, error)
}
