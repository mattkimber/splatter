package combiner

import (
	"reflect"
	"testing"
)

func TestGetImageMap(t *testing.T) {
	type args struct {
		definitions []SheetDefinition
		files       []string
	}

	tests := []struct {
		name       string
		args       args
		wantResult ImageSpecMap
	}{
		{
			name: "Ordered file list",
			args: args{
				definitions: []SheetDefinition{{
					Prefix:   "foo",
					Suffixes: []string{"a", "b"},
				}},
				files: []string{"foo_1_a.png", "foo_1_b.png", "foo_2_a.png", "foo_2_b.png"},
			},
			wantResult: ImageSpecMap{
				"foo_a": ImageSpec{Files: []string{"foo_1_a.png", "foo_2_a.png"}},
				"foo_b": ImageSpec{Files: []string{"foo_1_b.png", "foo_2_b.png"}},
			},
		},
		{
			name: "Unordered file list",
			args: args{
				definitions: []SheetDefinition{{
					Prefix:   "foo",
					Suffixes: []string{"a", "b"},
				}},
				files: []string{"foo_2_a.png", "foo_1_b.png", "foo_1_a.png", "foo_2_b.png"},
			},
			wantResult: ImageSpecMap{
				"foo_a": ImageSpec{Files: []string{"foo_1_a.png", "foo_2_a.png"}},
				"foo_b": ImageSpec{Files: []string{"foo_1_b.png", "foo_2_b.png"}},
			},
		},
		{
			name: "Prefix array",
			args: args{
				definitions: []SheetDefinition{{
					Prefixes: []string{"foo_1", "foo_2"},
					Suffixes: []string{"a", "b"},
				}},
				files: []string{"foo_2_a.png", "foo_1_b.png", "foo_1_a.png", "foo_2_b.png"},
			},
			wantResult: ImageSpecMap{
				"foo_1_a": ImageSpec{Files: []string{"foo_1_a.png"}},
				"foo_2_a": ImageSpec{Files: []string{"foo_2_a.png"}},
				"foo_1_b": ImageSpec{Files: []string{"foo_1_b.png"}},
				"foo_2_b": ImageSpec{Files: []string{"foo_2_b.png"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := GetImageMap(tt.args.definitions, tt.args.files); !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("GetImageMap() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
