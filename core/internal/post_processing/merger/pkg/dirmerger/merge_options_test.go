package dirmerger

import (
	"testing"
)

func TestMergeOptions_ShouldSkip(t *testing.T) {
	t.Parallel()

	type fields struct {
		Exclusions []*WildcardString
		Inclusions []*WildcardString
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
				Exclusions: []*WildcardString{
					NewWildcardString("/folderToExclude/*"),
				},
				Inclusions: make([]*WildcardString, 0),
			},
			args: args{
				location: "/folderToExclude/file.txt",
			},
			want: true,
		},
		{
			name: "should not skip when not excluded",
			fields: fields{
				Exclusions: []*WildcardString{
					NewWildcardString("/folderToExclude/*"),
				},
				Inclusions: make([]*WildcardString, 0),
			},
			args: args{
				location: "/folder/file.txt",
			},
			want: false,
		},
		{
			name: "should skip when excluded - implied star at end",
			fields: fields{
				Exclusions: []*WildcardString{
					NewWildcardString("/folderToExclude/"),
				},
				Inclusions: make([]*WildcardString, 0),
			},
			args: args{
				location: "/folderToExclude/file.txt",
			},
			want: true,
		},
		{
			name: "should not skip when not excluded - implied star at end",
			fields: fields{
				Exclusions: []*WildcardString{
					NewWildcardString("/folderToExclude/"),
				},
				Inclusions: make([]*WildcardString, 0),
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
				Inclusions: []*WildcardString{
					NewWildcardString("/folderToInclude/*"),
				},
				Exclusions: make([]*WildcardString, 0),
			},
			args: args{
				location: "/folderToInclude/file.txt",
			},
			want: false,
		},
		{
			name: "should skip when not included",
			fields: fields{
				Exclusions: []*WildcardString{
					NewWildcardString("/folderToInclude/*"),
				},
				Inclusions: make([]*WildcardString, 0),
			},
			args: args{
				location: "/folder/file.txt",
			},
			want: false,
		},
		{
			name: "should not skip when included - implied star at end",
			fields: fields{
				Exclusions: []*WildcardString{
					NewWildcardString("/folderToInclude/"),
				},
				Inclusions: make([]*WildcardString, 0),
			},
			args: args{
				location: "/folderToInclude/file.txt",
			},
			want: true,
		},
		{
			name: "should skip when not included - implied star at end",
			fields: fields{
				Exclusions: []*WildcardString{
					NewWildcardString("/folderToInclude/"),
				},
				Inclusions: make([]*WildcardString, 0),
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
				Inclusions: []*WildcardString{
					NewWildcardString("/folderToInclude/*"),
				},
				Exclusions: []*WildcardString{
					NewWildcardString("/folderToInclude/exclude/*"),
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
				Inclusions: []*WildcardString{
					NewWildcardString("/folderToInclude/*"),
				},
				Exclusions: []*WildcardString{
					NewWildcardString("/folderToInclude/exclude/*"),
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
				Inclusions: []*WildcardString{
					NewWildcardString("/folderToInclude/*"),
				},
				Exclusions: []*WildcardString{
					NewWildcardString("/folderToInclude/exclude/*"),
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

			o := NewMergeOptions()
			o.Exclusions = testCase.fields.Exclusions
			o.Inclusions = testCase.fields.Inclusions

			if got := o.ShouldSkip(testCase.args.location); got != testCase.want {
				t.Errorf("ShouldSkip() = %v, want %v", got, testCase.want)
			}
		})
	}
}
