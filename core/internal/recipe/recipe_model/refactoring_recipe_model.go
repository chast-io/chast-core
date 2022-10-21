package recipe_model

type RefactoringRecipe struct {
	BaseRecipe         `yaml:",inline"`
	SupportedLanguages []string `yaml:"supportedLanguages"`
	Run                []Run    `yaml:"run"`
	Tests              []string `yaml:"tests"` // placeholder for tests
}

type Run struct {
	Docker Docker `yaml:"docker"`
	Local  Local  `yaml:"local"`
	Cmd    string `yaml:"cmd"`
}

type Docker struct {
	DockerImage string `yaml:"dockerImage"`
}

type Local struct {
	RequiredTools []RequiredTool `yaml:"requiredTools"`
}

type RequiredTool struct {
	Description string `yaml:"description"`
	Cmd         string `yaml:"cmd"`
}
