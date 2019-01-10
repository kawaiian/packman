// unit tests fo handler
package packman

import (
	"testing"
)

func TestHandleQry(t *testing.T) {
	pkgTree := map[string][]string{
		"flask": []string{"mysql"},
		"arc":   []string{"flask"},
		"mysql": []string{""},
	}
	var tests = []struct {
		input PkgRequest
		want  string
	}{
		{PkgRequest{Cmd: "QUERY", Pkg: "flask", DepList: []string{""}}, "OK"},
		{PkgRequest{Cmd: "QUERY", Pkg: "git", DepList: []string{""}}, "FAIL"},
	}
	for _, test := range tests {
		if got := handleQry(test.input, pkgTree); got != test.want {
			t.Errorf("handleQry(%q) = %v", test.input, got)
		}
	}
}
