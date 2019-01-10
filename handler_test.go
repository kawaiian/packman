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

func TestHandleIdx(t *testing.T) {
	pkgTree := map[string][]string{
		"git":       []string{"gcc+"},
		"ilo-tools": []string{""},
	}
	var tests = []struct {
		input PkgRequest
		want  string
	}{
		{PkgRequest{Cmd: "INDEX", Pkg: "flask", DepList: []string{"mysql", "nginx"}}, "FAIL"}, // dependencies not indexed yet
		{PkgRequest{Cmd: "INDEX", Pkg: "ab", DepList: []string{""}}, "OK"},                    // no dependencies
	}
	for _, test := range tests {
		if got := handleIdx(test.input, pkgTree); got != test.want {
			t.Errorf("handleIdx(%q) = %v", test.input, got)
		}
	}
}
