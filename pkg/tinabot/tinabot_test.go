package tinabot

import (
	"testing"
)

func TestSplitSep(t *testing.T) {

	tests := map[string][]string{
		"a&b":              {"a", "b"},
		"a&b&c&d":          {"a", "b", "c", "d"},
		"a\\&b&c\\&d":      {"a&b", "c&d"},
		"a\\&b\\&c\\&d":    {"a&b&c&d"},
		"abcd":             {"abcd"},
		"&ab&cd":           {"", "ab", "cd"},
		"&ab&cd&":          {"", "ab", "cd", ""},
		"&ab&&&cd&":        {"", "ab", "", "", "cd", ""},
		"&ab&\\&&cd&":      {"", "ab", "&", "cd", ""},
		"&ab&\\&\\\\q&cd&": {"", "ab", "&\\q", "cd", ""},
		"a\\\\&b":          {"a\\&b"},
	}

	for i := range tests {
		out := splitEsc(i, "&")
		for j := range out {
			if out[j] != tests[i][j] {
				t.Fatalf("Error, wanted %v, got %v", tests[i], out)
			}
		}

	}
}
