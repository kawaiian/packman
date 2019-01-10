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
		response = handleIdx(pkgReq, pkgTree)
	case "QUERY":
		log.Printf("Received QUERY request for %s\n", pkgReq.Pkg)
		response = handleQry(pkgReq, pkgTree)
	case "REMOVE":
		log.Printf("Received REMOVE request for %s\n", pkgReq.Pkg)
		response = handleRemove(pkgReq, pkgTree)
	default:
		log.Printf("Received invalid request: %s\n", pkgReq.Cmd)
		response = "ERROR"
	}
	return response
}

func handleQry(pkgReq PkgRequest, pkgTree map[string][]string) string {
	mu.RLock()
	defer mu.RUnlock()

	if pkgTree != nil {
		_, pkgIndexed := pkgTree[pkgReq.Pkg]
		if pkgIndexed {
			log.Printf("Found %s in index", pkgReq.Pkg)
			return "OK"
		}
	}
	log.Printf("Did not find %s in index.", pkgReq.Pkg)
	return "FAIL"
}

func handleIdx(pkgReq PkgRequest, pkgTree map[string][]string) string {

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

func handleRemove(pkgReq PkgRequest, pkgTree map[string][]string) string {
	pkgName := pkgReq.Pkg
	// Check if the package exists in the pkgTree, if not we're done
	// TODO: Need to check that the pkgTree exists, or else return "OK"
	mu.RLock()
	_, pkgIndexed := pkgTree[pkgName]
	mu.RUnlock()
	if !pkgIndexed {
		log.Printf("Package %s is not indexed, so doesn't need to be removed.", pkgName)
		return "OK"
	}

	// Otherwise, iterate the pkgTree, and use a binary search for the package on each dependency list
	// This is O (nlogn) in the worst case, where the dependency does not exist in any other package
	log.Printf("Checking dependencies for %s...", pkgName)
	for _, depList := range pkgTree {
		if len(depList) > 0 {
			idx := sort.Search(len(depList), func(i int) bool { return depList[i] >= pkgName })
			if idx < len(depList) && depList[idx] == pkgName {
				log.Printf("Found dependency for %s; can't remove it.", pkgName)
				return "FAIL"
			}
		}
	}

	delete(pkgTree, pkgName)
	log.Printf("No dependencies for %s detected. Successfully removed.", pkgName)

	return "OK"
}
