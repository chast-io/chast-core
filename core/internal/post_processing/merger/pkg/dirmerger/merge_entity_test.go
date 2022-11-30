package dirmerger_test

import (
	"reflect"
	"testing"

	uut "chast.io/core/internal/post_processing/merger/pkg/dirmerger"
	"chast.io/core/internal/run_model/pkg/model/refactoring"
)

func TestNewMergeEntity(t *testing.T) {
	t.Parallel()

	type args struct {
		sourcePath      string
		changeLocations *refactoring.ChangeLocations
	}

	tests := []struct {
		name string
		args args
		want uut.MergeEntity
	}{
		{
			name: "TestNewMergeEntity Normal",
			args: args{
				sourcePath:      "sourcePath",
				changeLocations: &refactoring.ChangeLocations{}, //nolint:exhaustruct // test data
			},
			want: uut.MergeEntity{
				SourcePath:      "sourcePath",
				ChangeLocations: &refactoring.ChangeLocations{}, //nolint:exhaustruct // test data
			},
		},
		{
			name: "TestNewMergeEntity ChangeFilteringLocations is nil",
			args: args{
				sourcePath:      "sourcePath",
				changeLocations: nil,
			},
			want: uut.MergeEntity{
				SourcePath: "sourcePath",
				ChangeLocations: &refactoring.ChangeLocations{
					Include: make([]string, 0),
					Exclude: make([]string, 0),
				},
			},
		},
	}

	for i := range tests {
		testCase := tests[i]
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			if got := uut.NewMergeEntity(testCase.args.sourcePath, testCase.args.changeLocations); !reflect.DeepEqual(got, testCase.want) {
				t.Errorf("NewMergeEntity() = %v, want %v", got, testCase.want)
			}
		})
	}
}
