package parser_test

import (
	recipemodel "chast.io/core/internal/recipe/pkg/model"
	"reflect"
	"testing"
)

func testParameter(t *testing.T,
	parameter *recipemodel.Parameter,
	expectedParameter recipemodel.Parameter) {
	t.Helper()

	if parameter == nil {
		t.Error("Expected parameter to be set, but was nil")
	}

	t.Run("ID", func(t *testing.T) {
		t.Parallel()

		if parameter.ID != expectedParameter.ID {
			t.Errorf("Expected parameter ID to be '%v', but was '%v'", expectedParameter.ID, parameter.ID)
		}
	})

	t.Run("Required", func(t *testing.T) {
		t.Parallel()

		if parameter.Required != expectedParameter.Required {
			t.Errorf("Expected parameter required to be '%v', but was '%v'", expectedParameter.Required, parameter.Required)
		}
	})

	t.Run("DefaultValue", func(t *testing.T) {
		t.Parallel()

		if parameter.DefaultValue != expectedParameter.DefaultValue {
			t.Errorf("Expected parameter default value to be '%v', but was '%v'", expectedParameter.DefaultValue, parameter.DefaultValue)
		}
	})

	t.Run("Type", func(t *testing.T) {
		t.Parallel()

		if parameter.Type != expectedParameter.Type {
			t.Errorf("Expected parameter type to be '%v', but was '%v'", expectedParameter.Type, parameter.Type)
		}
	})

	t.Run("TypeExtension", func(t *testing.T) {
		t.Parallel()

		if reflect.DeepEqual(parameter.Extensions, expectedParameter.Extensions) {
			t.Errorf("Expected parameter extensions to be '%v', but was '%v'", expectedParameter.Extensions, parameter.Extensions)
		}
	})

	t.Run("Description", func(t *testing.T) {
		t.Parallel()

		if parameter.Description != expectedParameter.Description {
			t.Errorf("Expected parameter description to be '%v', but was '%v'", expectedParameter.Description, parameter.Description)
		}
	})

	t.Run("LongDescription", func(t *testing.T) {
		t.Parallel()

		if parameter.LongDescription != expectedParameter.LongDescription {
			t.Errorf("Expected parameter long description to be '%v', but was '%v'", expectedParameter.LongDescription, parameter.LongDescription)
		}
	})
}
