package refactoringpipelinemodel

type ExecutionGroup struct {
	Steps []*Step

	pipeline *Pipeline
}

func NewExecutionGroup() *ExecutionGroup {
	return &ExecutionGroup{ //nolint:exhaustruct // rest initialized in withPipeline
		Steps: make([]*Step, 0),
	}
}

func (s *ExecutionGroup) AddStep(step *Step) {
	if step == nil {
		return
	}

	s.Steps = append(s.Steps, step)
	if s.pipeline != nil {
		s.setStepDetails(step)
	}
}

func (s *ExecutionGroup) withPipeline(targetPipeline *Pipeline) {
	s.pipeline = targetPipeline

	for _, step := range s.Steps {
		s.setStepDetails(step)
	}
}

func (s *ExecutionGroup) setStepDetails(step *Step) {
	step.withPipeline(s.pipeline)
}
