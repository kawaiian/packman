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
	indexError := "invalid command: not INDEX, QUERY, or REMOVE"
	var tests = []struct {
		input string
		want  string
	}{
		{"", indexError},
		{"&^@%$#", indexError},
		{"INDE", indexError},
		{"REMVE", indexError},
		{"1QUERY", indexError},
	}
	for _, test := range tests {
		if got := parseCmd(test.input); got.Error() != test.want {
			t.Errorf("parseCmd(%q) = %v", test.input, got)
		}
	}
}
