package recipemodel

type RefactoringRecipe struct {
	BaseRecipe       `yaml:",inline"`
	PrimaryParameter *Parameter `yaml:"primaryParameter"`
	Runs             []Run      `yaml:"run"`
	Tests            []string   `yaml:"tests"` // TODO placeholder for tests
}

func (recipe *RefactoringRecipe) GetRecipeType() ChastOperationType {
	return Refactoring
}

type Run struct {
	ID                  string   `yaml:"id,omitempty"`
	Dependencies        []string `yaml:"dependencies,omitempty"`
	SupportedExtensions []string `yaml:"supportedExtensions,omitempty"`
	Flags               []Flag   `yaml:"flags,omitempty"`
	Docker              *Docker  `yaml:"docker"`
	Local               *Local   `yaml:"local"`
	Script              []string `yaml:"script"`
	ChangeLocations     []string `yaml:"changeLocations"` // TODO check concrete definition
}

func (run *Run) GetFlags() []Flag {
	return run.Flags
}

func (run *Run) GetFlagsMap() map[string]*Flag {
	return flagsToMap(run.Flags)
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
