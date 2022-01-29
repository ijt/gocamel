package main

import "testing"

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
