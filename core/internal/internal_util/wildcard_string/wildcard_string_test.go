package wildcardstring_test

import (
	"reflect"
	"testing"

	uut "chast.io/core/internal/internal_util/wildcard_string"
)

func TestNewWildcardString(t *testing.T) {
	t.Parallel()

	type args struct {
		value string
	}

	tests := []struct {
		name string
		args args
		want *uut.WildcardString
	}{
		{
			name: "String with no wildcard",
			args: args{
				value: "test",
			},
			want: &uut.WildcardString{
				Pattern: "test",
			},
		},
		{
			name: "String with wildcard",
			args: args{
				value: "test*",
			},
			want: &uut.WildcardString{
				Pattern: "test*",
			},
		},
	}

	for i := range tests {
		testCase := tests[i]

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			if got := uut.NewWildcardString(testCase.args.value); !reflect.DeepEqual(got, testCase.want) {
				t.Errorf("NewWildcardString() = %v, want %v", got, testCase.want)
			}
		})
	}
}

func TestWildcardString_Matches(t *testing.T) { //nolint:maintidx // nested test case
	t.Parallel()

	type fields struct {
		pattern string
	}

	type args struct {
		value string
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "String with no wildcard [equal]",
			fields: fields{
				pattern: "test",
			},
			args: args{
				value: "test",
			},
			want: true,
		},
		{
			name: "String with no wildcard 1 [not equal]",
			fields: fields{
				pattern: "test",
			},
			args: args{
				value: "1test",
			},
			want: false,
		},
		{
			name: "String with no wildcard 2 [not equal]",
			fields: fields{
				pattern: "test",
			},
			args: args{
				value: "test1",
			},
			want: false,
		},
		{
			name: "String with no wildcard 3 [not equal]",
			fields: fields{
				pattern: "test",
			},
			args: args{
				value: "not-same",
			},
			want: false,
		},
		// wildcard at end
		{
			name: "String with wildcard at end 1 [equal]",
			fields: fields{
				pattern: "test*",
			},
			args: args{
				value: "test",
			},
			want: true,
		},
		{
			name: "String with wildcard at end 2 [equal]",
			fields: fields{
				pattern: "test*",
			},
			args: args{
				value: "test1",
			},
			want: true,
		},
		{
			name: "String with wildcard at end 3 [equal]",
			fields: fields{
				pattern: "test*",
			},
			args: args{
				value: "test123",
			},
			want: true,
		},
		{
			name: "String with wildcard at end 1 [not equal]",
			fields: fields{
				pattern: "test*",
			},
			args: args{
				value: "tes1t",
			},
			want: false,
		},
		{
			name: "String with wildcard at end 2 [not equal]",
			fields: fields{
				pattern: "test*",
			},
			args: args{
				value: "tes123t",
			},
			want: false,
		},
		{
			name: "String with wildcard at end 3 [not equal]",
			fields: fields{
				pattern: "test*",
			},
			args: args{
				value: "",
			},
			want: false,
		},
		{
			name: "String with wildcard at end 4 [not equal]",
			fields: fields{
				pattern: "test*",
			},
			args: args{
				value: "1test",
			},
			want: false,
		},
		{
			name: "String with wildcard at end 5 [not equal]",
			fields: fields{
				pattern: "test*",
			},
			args: args{
				value: "1test1",
			},
			want: false,
		},
		// wildcard in middle
		{
			name: "String with wildcard in middle 1 [equal]",
			fields: fields{
				pattern: "te*st",
			},
			args: args{
				value: "test",
			},
			want: true,
		},
		{
			name: "String with wildcard in middle 2 [equal]",
			fields: fields{
				pattern: "te*st",
			},
			args: args{
				value: "te1st",
			},
			want: true,
		},
		{
			name: "String with wildcard in middle 3 [equal]",
			fields: fields{
				pattern: "te*st",
			},
			args: args{
				value: "te123st",
			},
			want: true,
		},
		{
			name: "String with wildcard in middle 1 [not equal]",
			fields: fields{
				pattern: "te*st",
			},
			args: args{
				value: "tst",
			},
			want: false,
		},
		{
			name: "String with wildcard in middle 2 [not equal]",
			fields: fields{
				pattern: "te*st",
			},
			args: args{
				value: "tet",
			},
			want: false,
		},
		{
			name: "String with wildcard in middle 3 [not equal]",
			fields: fields{
				pattern: "te*st",
			},
			args: args{
				value: "tt",
			},
			want: false,
		},
		// wildcard at start
		{
			name: "String with wildcard at start 1 [equal]",
			fields: fields{
				pattern: "*test",
			},
			args: args{
				value: "test",
			},
			want: true,
		},
		{
			name: "String with wildcard at start 2 [equal]",
			fields: fields{
				pattern: "*test",
			},
			args: args{
				value: "1test",
			},
			want: true,
		},
		{
			name: "String with wildcard at start 3 [equal]",
			fields: fields{
				pattern: "*test",
			},
			args: args{
				value: "123test",
			},
			want: true,
		},
		{
			name: "String with wildcard at start 1 [not equal]",
			fields: fields{
				pattern: "*test",
			},
			args: args{
				value: "est",
			},
			want: false,
		},
		{
			name: "String with wildcard at start 2 [not equal]",
			fields: fields{
				pattern: "*test",
			},
			args: args{
				value: "te1st",
			},
			want: false,
		},
		{
			name: "String with wildcard at start 3 [not equal]",
			fields: fields{
				pattern: "*test",
			},
			args: args{
				value: "test1",
			},
			want: false,
		},
		{
			name: "String with wildcard at start 4 [not equal]",
			fields: fields{
				pattern: "*test",
			},
			args: args{
				value: "1te1st",
			},
			want: false,
		},
		{
			name: "String with wildcard at start 5 [not equal]",
			fields: fields{
				pattern: "*test",
			},
			args: args{
				value: "1test1",
			},
			want: false,
		},
		{
			name: "No wildcard but regex special characters 1 [not equal]",
			fields: fields{
				pattern: "test",
			},
			args: args{
				value: "test.",
			},
			want: false,
		},
		{
			name: "No wildcard but regex special characters 2 [not equal]",
			fields: fields{
				pattern: "test",
			},
			args: args{
				value: "test.out",
			},
			want: false,
		},
	}

	for i := range tests {
		testCase := tests[i]

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			w := &uut.WildcardString{
				Pattern: testCase.fields.pattern,
			}
			if got := w.Matches(testCase.args.value); got != testCase.want {
				t.Errorf("Matches(%v) with pattern %v resulted in: %v, want %v",
					testCase.args.value, testCase.fields.pattern, got, testCase.want)
			}
		})
	}
}

func TestWildcardString_MatchesPath(t *testing.T) {
	t.Parallel()

	type fields struct {
		pattern string
	}

	type args struct {
		value string
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "File with no wildcard [equal]",
			fields: fields{
				pattern: "/folder1/file1",
			},
			args: args{
				value: "/folder1/file1",
			},
			want: true,
		},
		{
			name: "File with no wildcard - no slash [not equal]",
			fields: fields{
				pattern: "/folder1/file1",
			},
			args: args{
				value: "/folder1/file1suffix",
			},
			want: false,
		},
		{
			name: "Folder with no wildcard - slash [equal]",
			fields: fields{
				pattern: "/folder1/subfolder1/",
			},
			args: args{
				value: "/folder1/subfolder1/",
			},
			want: true,
		},
		{
			name: "Folder with no wildcard, but check sub folder[equal]",
			fields: fields{
				pattern: "/folder1/",
			},
			args: args{
				value: "/folder1/subfolder1/",
			},
			want: true,
		},
		{
			name: "Folder with no wildcard, but check sub file [equal]",
			fields: fields{
				pattern: "/folder1/",
			},
			args: args{
				value: "/folder1/file1",
			},
			want: true,
		},
		{
			name: "File with no wildcard [not equal]",
			fields: fields{
				pattern: "/folder1/file1",
			},
			args: args{
				value: "/folder1/file2",
			},
			want: false,
		},
		{
			name: "Folder with no wildcard [not equal]",
			fields: fields{
				pattern: "/folder1/subfolder1/",
			},
			args: args{
				value: "/folder1/subfolder2/",
			},
			want: false,
		},
		{
			name: "Folder with wildcard at end - slash [equal]",
			fields: fields{
				pattern: "/folder1/*",
			},
			args: args{
				value: "/folder1/",
			},
			want: true,
		},
		{
			name: "Folder with wildcard at end - no slash [equal]",
			fields: fields{
				pattern: "/folder1/*",
			},
			args: args{
				value: "/folder1",
			},
			want: true,
		},
		{
			name: "Folder with wildcard at end - any [equal]",
			fields: fields{
				pattern: "/folder1/*",
			},
			args: args{
				value: "/folder1/subfolder1/subfolder2",
			},
			want: true,
		},
		{
			name: "Folder with wildcard at end - different start [not equal]",
			fields: fields{
				pattern: "/folder1/*",
			},
			args: args{
				value: "/folder2/subfolder1",
			},
			want: false,
		},
		{
			name: "Folder with wildcard in middle [equal]",
			fields: fields{
				pattern: "/folder1/*/subfolder2/",
			},
			args: args{
				value: "/folder1/subfolder1/subfolder2/",
			},
			want: true,
		},
		{
			name: "Folder with wildcard in middle [not equal]",
			fields: fields{
				pattern: "/folder1/*/subfolder2/",
			},
			args: args{
				value: "/folder1/subfolder2/",
			},
			want: false,
		},
		{
			name: "Folder with wildcard at start - folders [equal]",
			fields: fields{
				pattern: "*/subfolder2/",
			},
			args: args{
				value: "/folder1/subfolder1/subfolder2/",
			},
			want: true,
		},
		{
			name: "Folder with wildcard at start - no folders [equal]",
			fields: fields{
				pattern: "*/folder1/",
			},
			args: args{
				value: "/folder1/",
			},
			want: true,
		},
		{
			name: "Folder with wildcard at start - folders [not equal]",
			fields: fields{
				pattern: "*/subfolder2/",
			},
			args: args{
				value: "/folder1/subfolder1/subsubfolder1/",
			},
			want: false,
		},
		{
			name: "Folder with wildcard at start - no folders [not equal]",
			fields: fields{
				pattern: "*/folder1/",
			},
			args: args{
				value: "/folder2/",
			},
			want: false,
		},
		{
			name: "File with wildcard at end [equal]",
			fields: fields{
				pattern: "/file1*",
			},
			args: args{
				value: "/file1",
			},
			want: true,
		},
		{
			name: "File with wildcard at end [equal]",
			fields: fields{
				pattern: "/file1*",
			},
			args: args{
				value: "/file1suffix",
			},
			want: true,
		},
		{
			name: "Multiple wildcards 1 [equal]",
			fields: fields{
				pattern: "*/folder1/*",
			},
			args: args{
				value: "/folder1/subfolder1/file1",
			},
			want: true,
		},
		{
			name: "Multiple wildcards 2 [equal]",
			fields: fields{
				pattern: "*/subfolder1/*",
			},
			args: args{
				value: "/folder1/subfolder1/file1",
			},
			want: true,
		},
		{
			name: "Multiple wildcards 3 [equal]",
			fields: fields{
				pattern: "*/subfolder1/*",
			},
			args: args{
				value: "/folder1/subfolder1",
			},
			want: true,
		},
		{
			name: "Multiple wildcards 4 [equal]",
			fields: fields{
				pattern: "/folder1/*/*/file1",
			},
			args: args{
				value: "/folder1/subfolder1/subsubfolder1/file1",
			},
			want: true,
		},
		{
			name: "Multiple wildcards 4 [not equal]",
			fields: fields{
				pattern: "/folder1/*/*/file1",
			},
			args: args{
				value: "/folder1/subfolder1/file1",
			},
			want: false,
		},
		{
			name: "No wildcard but regex special characters 1 [not equal]",
			fields: fields{
				pattern: "/folder/test",
			},
			args: args{
				value: "/folder/test.",
			},
			want: false,
		},
		{
			name: "No wildcard but regex special characters 2 [not equal]",
			fields: fields{
				pattern: "/folder/test",
			},
			args: args{
				value: "/folder/test.out",
			},
			want: false,
		},
	}

	for i := range tests {
		testCase := tests[i]

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			w := &uut.WildcardString{
				Pattern: testCase.fields.pattern,
			}
			if got := w.MatchesPath(testCase.args.value); got != testCase.want {
				t.Errorf("MatchesPath(%v) with pattern %v resulted in: %v, want %v",
					testCase.args.value, testCase.fields.pattern, got, testCase.want)
			}
		})
	}
}
