package packman

import (
	"sync"
	"testing"
)

var mu = &sync.Mutex{}

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
		{PkgRequest{Cmd: "QUERY", Pkg: "flask", DepList: []string{""}}, "OK"}, // pkg indexed
		{PkgRequest{Cmd: "QUERY", Pkg: "git", DepList: []string{""}}, "FAIL"}, // pkg not indexed
	}
	for _, test := range tests {
		if got := handleQry(test.input, pkgTree, mu); got != test.want {
			t.Errorf("handleQry(%q) = %v", test.input, got)
		}
	}
}

func TestHandleIdxSucceeds(t *testing.T) {
	pkgTree := map[string][]string{
		"git":       []string{"gcc+"},
		"ilo-tools": []string{""},
	}

	testReq := PkgRequest{Cmd: "INDEX", Pkg: "ab", DepList: []string{""}}
	response := handleIdx(testReq, pkgTree, mu)

	if response == "OK" {
		_, indexed := pkgTree[testReq.Pkg]
		if !indexed {
			t.Errorf("handleIdx(%q) did not properly index package", testReq.Pkg)
		}
	} else {
		t.Errorf("handleIdx(%q) did not successfully complete index operation,", testReq.Pkg)
	}
}

func TestHandleIdxForExisingPkg(t *testing.T) {
	pkgTree := map[string][]string{
		"ab":        []string{""},
		"git":       []string{"gcc+"},
		"ilo-tools": []string{""},
		"make":      []string{""},
	}
	testReq := PkgRequest{Cmd: "INDEX", Pkg: "git", DepList: []string{"ab", "make"}}
	expected := []string{"ab", "make"}

	response := handleIdx(testReq, pkgTree, mu)

	if response == "OK" {
		newDepList := pkgTree["git"]
		for idx, pkgName := range newDepList {
			if pkgName != expected[idx] {
				t.Errorf("handleIdx(%q) did not update deplist properly", testReq.Pkg)
			}
		}
	} else {
		t.Errorf("handleIdx(%q) failed to index properly", testReq.Pkg)
	}

}

func TestHandleRemoveSuceeds(t *testing.T) {
	pkgTree := map[string][]string{
		"ab":        []string{"git", "ilo-tools"},
		"gcc+":      []string{""},
		"git":       []string{"gcc+"},
		"ilo-tools": []string{""},
		"make":      []string{""},
	}

	testReq := PkgRequest{Cmd: "REMOVE", Pkg: "make", DepList: []string{""}}
	response := handleRemove(testReq, pkgTree, mu)
	if response == "OK" {
		_, stillIndexed := pkgTree[testReq.Pkg]
		if stillIndexed {
			t.Errorf("handleRemove(%q) failed to remove package", testReq.Pkg)
		}
	} else {
		t.Errorf("handleRemove(%q) failed to finish remove command,", testReq.Pkg)
	}
}

func TestHandleRemoveFails(t *testing.T) {
	pkgTree := map[string][]string{
		"ab":   []string{"git", "ilo-tools"},
		"gcc+": []string{""},
		"git":  []string{"gcc+"},
	}

	testReq := PkgRequest{Cmd: "REMOVE", Pkg: "git", DepList: []string{""}}
	response := handleRemove(testReq, pkgTree, mu)
	if response == "FAIL" {
		_, stillIndexed := pkgTree[testReq.Pkg]
		if !stillIndexed {
			t.Errorf("handleRemove(%q) removed package when it shouldn't", testReq.Pkg)
		}
	} else {
		t.Errorf("handleRemove(%q) failed to finish remove command,", testReq.Pkg)
	}
}
