// Package pkglibio provides a concurrency-safe hashMap of system packages.
package packmanio

//!+
import "sync"

var (
	mu         sync.RWMutex
	packageLib packageLibrary
)

// packageLibrary is a hashMap of package names to dependency lists (slices)
type packageLibrary struct {
	pkgDeps map[string][]string
}

func newPackageLibrary() *packageLibrary {
	return &packageLibrary{
		pkgDeps: make(map[string][]string),
	}
}
