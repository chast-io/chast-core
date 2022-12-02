package mergeoptions_test

import (
	"testing"

	wildcardstring "chast.io/core/internal/internal_util/wildcard_string"
	"chast.io/core/internal/post_processing/merger/pkg/mergeoptions"
)

func TestMergeOptions_ShouldSkip(t *testing.T) {
	t.Parallel()

	type fields struct {
		Exclusions []*wildcardstring.WildcardString
		Inclusions []*wildcardstring.WildcardString
	}

	type args struct {
		location string
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// exclude
		{
			name: "should skip when excluded",
			fields: fields{
				Exclusions: []*wildcardstring.WildcardString{
					wildcardstring.NewWildcardString("/folderToExclude/*"),
				},
				Inclusions: make([]*wildcardstring.WildcardString, 0),
			},
			args: args{
				location: "/folderToExclude/file.txt",
			},
			want: true,
		},
		{
			name: "should not skip when not excluded",
			fields: fields{
				Exclusions: []*wildcardstring.WildcardString{
					wildcardstring.NewWildcardString("/folderToExclude/*"),
				},
				Inclusions: make([]*wildcardstring.WildcardString, 0),
			},
			args: args{
				location: "/folder/file.txt",
			},
			want: false,
		},
		{
			name: "should skip when excluded - implied star at end",
			fields: fields{
				Exclusions: []*wildcardstring.WildcardString{
					wildcardstring.NewWildcardString("/folderToExclude/"),
				},
				Inclusions: make([]*wildcardstring.WildcardString, 0),
			},
			args: args{
				location: "/folderToExclude/file.txt",
			},
			want: true,
		},
		{
			name: "should not skip when not excluded - implied star at end",
			fields: fields{
				Exclusions: []*wildcardstring.WildcardString{
					wildcardstring.NewWildcardString("/folderToExclude/"),
				},
				Inclusions: make([]*wildcardstring.WildcardString, 0),
			},
			args: args{
				location: "/folder/file.txt",
			},
			want: false,
		},
		// include
		{
			name: "should not skip when included",
			fields: fields{
				Inclusions: []*wildcardstring.WildcardString{
					wildcardstring.NewWildcardString("/folderToInclude/*"),
				},
				Exclusions: make([]*wildcardstring.WildcardString, 0),
			},
			args: args{
				location: "/folderToInclude/file.txt",
			},
			want: false,
		},
		{
			name: "should skip when not included",
			fields: fields{
				Exclusions: []*wildcardstring.WildcardString{
					wildcardstring.NewWildcardString("/folderToInclude/*"),
				},
				Inclusions: make([]*wildcardstring.WildcardString, 0),
			},
			args: args{
				location: "/folder/file.txt",
			},
			want: false,
		},
		{
			name: "should not skip when included - implied star at end",
			fields: fields{
				Exclusions: []*wildcardstring.WildcardString{
					wildcardstring.NewWildcardString("/folderToInclude/"),
				},
				Inclusions: make([]*wildcardstring.WildcardString, 0),
			},
			args: args{
				location: "/folderToInclude/file.txt",
			},
			want: true,
		},
		{
			name: "should skip when not included - implied star at end",
			fields: fields{
				Exclusions: []*wildcardstring.WildcardString{
					wildcardstring.NewWildcardString("/folderToInclude/"),
				},
				Inclusions: make([]*wildcardstring.WildcardString, 0),
			},
			args: args{
				location: "/folder/file.txt",
			},
			want: false,
		},
		// exclude and include
		{
			name: "should not skip when included and not excluded",
			fields: fields{
				Inclusions: []*wildcardstring.WildcardString{
					wildcardstring.NewWildcardString("/folderToInclude/*"),
				},
				Exclusions: []*wildcardstring.WildcardString{
					wildcardstring.NewWildcardString("/folderToInclude/exclude/*"),
				},
			},
			args: args{
				location: "/folderToInclude/file.txt",
			},
			want: false,
		},
		{
			name: "should skip when included and excluded",
			fields: fields{
				Inclusions: []*wildcardstring.WildcardString{
					wildcardstring.NewWildcardString("/folderToInclude/*"),
				},
				Exclusions: []*wildcardstring.WildcardString{
					wildcardstring.NewWildcardString("/folderToInclude/exclude/*"),
				},
			},
			args: args{
				location: "/folderToInclude/exclude/file.txt",
			},
			want: true,
		},
		{
			name: "should skip when not included and not excluded",
			fields: fields{
				Exclusions: []*wildcardstring.WildcardString{
					wildcardstring.NewWildcardString("/folderToExclude/*"),
				},
				Inclusions: make([]*wildcardstring.WildcardString, 0),
			},
			args: args{
				location: "/folderToExclude/" + mergeoptions.NewMergeOptions().MetaFilesDeletedExtension +
					"/file.txt" + mergeoptions.NewMergeOptions().MetaFilesDeletedExtension,
			},
			want: true,
		},
		{
			name: "should also match deleted files",
			fields: fields{
				Inclusions: []*wildcardstring.WildcardString{
					wildcardstring.NewWildcardString("/folderToInclude/*"),
				},
				Exclusions: []*wildcardstring.WildcardString{
					wildcardstring.NewWildcardString("/folderToInclude/exclude/*"),
				},
			},
			args: args{
				location: "/otherFolder/exclude/file.txt",
			},
			want: true,
		},
	}

	for i := range tests {
		testCase := tests[i]
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			o := mergeoptions.NewMergeOptions()
			o.Exclusions = testCase.fields.Exclusions
			o.Inclusions = testCase.fields.Inclusions

			if got := o.ShouldSkip(testCase.args.location); got != testCase.want {
				t.Errorf("ShouldSkip() = %v, want %v", got, testCase.want)
			}
		})
	}
}
