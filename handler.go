package packman

import (
	"fmt"
	"sync"
)

var (
	mu      sync.RWMutex
	pkgTree map[string][]string
)

// HandlePkgRequest reads and writes to the pkgLib in-memory store, based on the associated command
func HandlePkgRequest(pkgReq PkgRequest) string {
	switch pkgReq.cmd {
	case "INDEX":
		fmt.Printf("Received INDEX request for %s with dependencies %s\n", pkgReq.pkg, pkgReq.depList)
	case "QUERY":
		fmt.Printf("Received QUERY request for %s\n", pkgReq.pkg)
	case "REMOVE":
		fmt.Printf("Received REMOVE request for %s\n", pkgReq.pkg)
	}
}
