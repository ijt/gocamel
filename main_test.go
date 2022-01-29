package main

import (
	"reflect"
	"testing"
)

func Test_snakeToCamel(t *testing.T) {
	tests := []struct {
		ident string
		want  string
	}{
		{"", ""},
		{"a", "a"},
		{"ab", "ab"},
		{"a_b", "aB"},
		{"a_b_c", "aBC"},
		{"a_bc", "aBc"},
		{"ab_c", "abC"},
	}
	for _, tt := range tests {
		t.Run(tt.ident, func(t *testing.T) {
			if got := snakeToCamel(tt.ident); got != tt.want {
				t.Errorf("snakeToCamel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_snakeCaseToCamelCaseFile(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "empty",
			input:   "",
			want:    "",
			wantErr: true,
		},
		{
			name: "struct",
			input: `package x

type a_b struct {
	some_field type_of_field
}
`,
			want: `package x

type aB struct {
	someField typeOfField
}
`,
		},
		{
			name: "function declaration",
			input: `package x

func a_b(c_d e_f) g_h {
	return i_j
}
`,
			want: `package x

func aB(cD eF) gH {
	return iJ
}
`,
		},
		// { name: "function call" },
		// { name: "method declaration" },
		// { name: "method call" },
		// { name: "variable declaration" },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := snakeCaseToCamelCaseFile("foo.go", []byte(tt.input), false /* willPrintAST */)
			if (err != nil) != tt.wantErr {
				t.Errorf("snakeCaseToCamelCaseFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotStr := string(got)
			if !reflect.DeepEqual(gotStr, tt.want) {
				t.Errorf("snakeCaseToCamelCaseFile() got = '%s', want '%s'", gotStr, tt.want)
			}
		})
	}
}
