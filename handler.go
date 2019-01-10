package packman

import (
	"log"
	"sort"
	"sync"
)

var (
	mu      sync.RWMutex
	pkgTree map[string][]string
)

// HandlePkgRequest reads and writes to the pkgLib in-memory store, based on the associated command
func HandlePkgRequest(pkgReq PkgRequest) string {
	var response string
	switch pkgReq.Cmd {
	case "INDEX":
		log.Printf("Received INDEX request for %s with dependencies %s\n", pkgReq.Pkg, pkgReq.DepList)
		response = handleIdx(pkgReq)
	case "QUERY":
		log.Printf("Received QUERY request for %s\n", pkgReq.Pkg)
		response = handleQry(pkgReq)
	case "REMOVE":
		log.Printf("Received REMOVE request for %s\n", pkgReq.Pkg)
	}
	return response
}

func handleQry(pkgReq PkgRequest) string {
	mu.RLock()
	defer mu.RUnlock()

	_, pkgIndexed := pkgTree[pkgReq.Pkg]
	if pkgIndexed {
		return "OK"
	}
	return "FAIL"
}

func handleIdx(pkgReq PkgRequest) string {

	pkgName := pkgReq.Pkg

	// if there are dependencies, check to see if all of them are indexed, FAIL if not
	// this is O(n) where n is the size of the dependency list
	if len(pkgReq.DepList[0]) > 0 {
		mu.RLock()
		for _, depName := range pkgReq.DepList {
			_, depIndexed := pkgTree[depName]
			if !depIndexed {
				log.Printf("Unable to index %s because dependency %s not indexed", pkgName, depName)
				mu.RUnlock()
				return "FAIL"
			}
		}
		mu.RUnlock()
	}

	// sort the dependency list in ascending oreder before indexing
	// this allows us to use a binary search when checking dependencies in a REMOVE request
	sort.SliceStable(pkgReq.DepList, func(a, b int) bool { return pkgReq.DepList[a] < pkgReq.DepList[b] })

	// index the package and dependency list, note we create the datastore if it doesn't exist
	mu.Lock()
	if pkgTree == nil {
		pkgTree = make(map[string][]string)
	}
	pkgTree[pkgName] = pkgReq.DepList
	mu.Unlock()
	log.Printf("Successfully indexed %s with dependencies %s", pkgName, pkgReq.DepList)
	return "OK"
}

func handleRemove(pkgreq PkgRequest) string {
	// Check if the package exists in the pkgTree
	// If not, return "OK"
	// Otherwise, iterate the pkgTree, then use sort.SearchStrings on each depList to see if the package exists
	// this is O (nlogn) in the worst case, where the dependency does not exist in any other package
	return "OK"
}
