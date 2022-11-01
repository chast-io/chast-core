package model

type Recipe interface {
	GetRecipeType() ChastOperationType
}

type RecipeInfo struct {
	Version string `yaml:"version"`
	Type    string `yaml:"type"`
}

type BaseRecipe struct {
	Recipe
	RecipeInfo    `yaml:",inline"`
	Name          string     `yaml:"name"`
	Maintainer    string     `yaml:"maintainer,omitempty"`
	Repository    string     `yaml:"repository,omitempty"`
	Arguments     []Argument `yaml:"args"`
	Documentation string     `yaml:"documentation"` // placeholder for documentation
}

func (recipe *BaseRecipe) GetRecipeType() ChastOperationType {
	return Unknown
}

type Argument struct {
	ID               string `yaml:"id"`
	Type             string `yaml:"type,omitempty"`
	ShortDescription string `yaml:"shortDescription"`
	Required         bool   `yaml:"required,omitempty"`
	Description      string `yaml:"description,omitempty"`
}
