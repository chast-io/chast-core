package model

type RefactoringRecipe struct {
	BaseRecipe         `yaml:",inline"`
	SupportedLanguages []string `yaml:"supportedLanguages"`
	Run                []Run    `yaml:"run"`
	Tests              []string `yaml:"tests"` // TODO placeholder for tests
}

type Run struct {
	Docker          Docker   `yaml:"docker"`
	Local           Local    `yaml:"local"`
	Script          []string `yaml:"script"`
	ChangeLocations []string `yaml:"changeLocations"` // TODO check concrete definition
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
