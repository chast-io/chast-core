package recipemodel

type Recipe interface {
	GetRecipeType() ChastOperationType
}

type RecipeInfo struct {
	Version string `yaml:"version"`
	Type    string `yaml:"type"`
}

type BaseRecipe struct {
	Recipe
	RecipeInfo           `yaml:",inline"`
	Name                 string      `yaml:"name"`
	Maintainer           string      `yaml:"maintainer,omitempty"`
	Repository           string      `yaml:"repository,omitempty"`
	PositionalParameters []Parameter `yaml:"positionalParameters,omitempty"`
	Flags                []Flag      `yaml:"flags,omitempty"`
	Documentation        string      `yaml:"documentation"` // placeholder for documentation
}

func (recipe *BaseRecipe) GetRecipeType() ChastOperationType {
	return Unknown
}

func (recipe *BaseRecipe) GetPositionalParameters() []Parameter {
	return recipe.PositionalParameters
}

func (recipe *BaseRecipe) GetFlags() []Flag {
	return recipe.Flags
}

func (recipe *BaseRecipe) GetFlagsMap() map[string]*Flag {
	return flagsToMap(recipe.Flags)
}

func flagsToMap(flags []Flag) map[string]*Flag {
	flagsMap := make(map[string]*Flag)

	for i, flag := range flags {
		flagsMap[flag.Name] = &flags[i]
		if flag.ShortName != "" {
			flagsMap[flag.ShortName] = &flags[i]
		}
	}

	return flagsMap
}

type Parameter struct {
	ID                   string `yaml:"id"`
	RequiredExtension    `yaml:",inline"`
	TypeExtension        `yaml:",inline"`
	DescriptionExtension `yaml:",inline"`
}

type Flag struct {
	Name                 string `yaml:"name"`
	ShortName            string `yaml:"shortName,omitempty"`
	RequiredExtension    `yaml:",inline"`
	TypeExtension        `yaml:",inline"`
	DescriptionExtension `yaml:",inline"`
}

type RequiredExtension struct {
	Required     bool   `yaml:"required,omitempty"`
	DefaultValue string `yaml:"defaultValue,omitempty"`
}

type TypeExtension struct {
	Type       string   `yaml:"type,omitempty"`
	Extensions []string `yaml:"extensions,omitempty"`
}

type DescriptionExtension struct {
	Description     string `yaml:"description,omitempty"`
	LongDescription string `yaml:"longDescription,omitempty"`
}
