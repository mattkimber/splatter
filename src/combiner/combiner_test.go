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
		wantResult map[string][]string
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
			wantResult: map[string][]string{
				"foo_a": {"foo_1_a.png", "foo_2_a.png"},
				"foo_b": {"foo_1_b.png", "foo_2_b.png"},
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
			wantResult: map[string][]string{
				"foo_a": {"foo_1_a.png", "foo_2_a.png"},
				"foo_b": {"foo_1_b.png", "foo_2_b.png"},
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
