package recipe

type RefactoringRecipe struct {
	BaseRecipe `yaml:",inline"`
	Run        []Run    `yaml:"run"`
	Tests      []string `yaml:"tests"` // TODO placeholder for tests
}

func (recipe *RefactoringRecipe) GetRecipeType() ChastOperationType {
	return Refactoring
}

type Run struct {
	Id                 string   `yaml:"id,omitempty"`
	Dependencies       []string `yaml:"dependencies,omitempty"`
	SupportedLanguages []string `yaml:"supportedLanguages,omitempty"`
	Docker             Docker   `yaml:"docker"`
	Local              Local    `yaml:"local"`
	Script             []string `yaml:"script"`
	ChangeLocations    []string `yaml:"changeLocations"` // TODO check concrete definition
}

type Docker struct {
	DockerImage string `yaml:"dockerImage"`
}

type Local struct {
	RequiredTools []RequiredTool `yaml:"requiredTools"`
}

type RequiredTool struct {
	Description string `yaml:"description"`
	CheckCmd    string `yaml:"checkCmd"`
}

type ChangeLocation struct {
	Location          string   `yaml:"location"`
	AllowedOperations []string `yaml:"allowedOperations"` // modify, delete, insert
}
