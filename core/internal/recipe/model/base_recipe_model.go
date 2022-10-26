package model

type RecipeInfo struct {
	Version string `yaml:"version"`
	Type    string `yaml:"type"`
}

type BaseRecipe struct {
	RecipeInfo    `yaml:",inline"`
	Maintainer    string `yaml:"maintainer"`
	Name          string `yaml:"name"`
	Args          []Args `yaml:"args"`
	Documentation string `yaml:"documentation"` // placeholder for documentation
}

type Args struct {
	ID               string `yaml:"id"`
	Type             string `yaml:"type,omitempty"`
	ShortDescription string `yaml:"shortDescription"`
	Required         bool   `yaml:"required,omitempty"`
	Description      string `yaml:"description,omitempty"`
}
