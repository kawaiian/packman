// unit tests for validator
package packman

import (
	"testing"
)

func TestParseCmdSucceeds(t *testing.T) {
	var tests = []string{"INDEX", "QUERY", "REMOVE"}

	for _, test := range tests {
		err := parseCmd(test)
		if err != nil {
			t.Errorf("parseCmd(%q) resulted in an error when it should not", test)
		}
	}
}
func TestParseCmdErrors(t *testing.T) {
	parseError := "invalid command: not INDEX, QUERY, or REMOVE"
	var tests = []struct {
		input string
		want  string
	}{
		{"", parseError},
		{"&^@%$#", parseError},
		{"INDE", parseError},
		{"REMVE", parseError},
		{"1QUERY", parseError},
	}
	for _, test := range tests {
		if got := parseCmd(test.input); got.Error() != test.want {
			t.Errorf("parseCmd(%q) = %v", test.input, got)
		}
	}
}

func TestParsePkgNamesErrors(t *testing.T) {
	parseError := "invalid package name"

	var tests = []struct {
		pkgName string
		depList []string
		want    string
	}{
		{"", []string{""}, parseError},                  // both empty
		{"%^&*", []string{"mysql", "git"}, parseError},  // invalid pkgName, valid depList
		{"mysql", []string{"git", "afd$y"}, parseError}, // valid pkgName, partially invalid depList
	}
	for _, test := range tests {
		if got := parsePkgNames(test.pkgName, test.depList); got.Error() != test.want {
			t.Errorf("parsePkgNames(%q, %q) = %v", test.pkgName, test.depList, got)
		}
	}
}

func TestParsePkgNamesSucceeds(t *testing.T) {
	var tests = []struct {
		pkgName string
		depList []string
	}{
		{"mysql", []string{""}},                    // valid pkgName, no dependencies
		{"flask", []string{"mysql"}},               // valid pkgName, one valid dependency
		{"flask", []string{"mysql", "sqlalchemy"}}, // valid pkgName, multiple valid dependencies
	}
	for _, test := range tests {
		err := parsePkgNames(test.pkgName, test.depList)
		if err != nil {
			t.Errorf("parsePkgNames(%q, %q) = %s", test.pkgName, test.depList, err.Error())
		}
	}
}
