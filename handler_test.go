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

func TestHandleIdxForExisingPkg(t *testing.T) {
	pkgTree := map[string][]string{
		"ab":        []string{""},
		"git":       []string{"gcc+"},
		"ilo-tools": []string{""},
		"make":      []string{""},
	}
	testPkg := PkgRequest{Cmd: "INDEX", Pkg: "git", DepList: []string{"ab", "make"}}
	expected := []string{"ab", "make"}

	response := handleIdx(testPkg, pkgTree)

	if response != "FAIL" {
		newDepList := pkgTree["git"]
		for idx, pkgName := range newDepList {
			if pkgName != expected[idx] {
				t.Errorf("handleIdx(%q) did not update deplist properly", testPkg.Pkg)
			}
		}
	} else {
		t.Errorf("handleIdx(%q) failed to index properly", testPkg.Pkg)
	}

}
