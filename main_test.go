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
		// Leave shouting snakes alone.
		{"SNAKE_SHOUT_CASE", "SNAKE_SHOUT_CASE"},
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
		{
			name: "function calls",
			input: `package x

func f() {
	a_b(c_d())
	return e_f()
}
`,
			want: `package x

func f() {
	aB(cD())
	return eF()
}
`,
		},
		{
			name: "method declaration",
			input: `package x

func (a_b *c_d) e_f() (g_h i_j) {
}
`,
			want: `package x

func (aB *cD) eF() (gH iJ) {
}
`,
		},
		{
			name: "method call",
			input: `package x

func f() {
	a_b.c_d().e_f()
}
`,
			want: `package x

func f() {
	aB.cD().eF()
}
`,
		},
		{
			name: "field access",
			input: `package x

func f() int {
	return a.b_c.d_e
}
`,
			want: `package x

func f() int {
	return a.bC.dE
}
`,
		},
		{
			name: "variable declaration",
			input: `package x

var a_b c_d = e_f.g_h().i_j
`,
			want: `package x

var aB cD = eF.gH().iJ
`,
		},
		{
			name: "test names get left alone",
			input: `package x

func TestFoo_Bar(t *testing.T) {
}
`,
			want: `package x

func TestFoo_Bar(t *testing.T) {
}
`,
		},
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
