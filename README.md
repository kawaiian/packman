# Packman - A Linux package indexer

Packman is a simple golang application for in-memory Linux package indexing. The index is a simple golang Map of package name to a string slice of associated package dependencies. It uses Goroutines and the golang sync library's Mutex construct to achieve concurrency-safe request processing.

## What it does

The packman binary listens for TCP connections on all interfaces on port 8080 of the host. It responds to client requests of the form `<command>|<package>|<dependency_list>?\n`, where the dependency list is optional and only used for INDEX commands.

Packages are expected to be alphanumeric, with potential underscores `_`, dashes `-` or plus `+` symbols.

Commands include:
1. QUERY - check to see if a given package has been indexed. If server returns "OK\n" if the package exists in the index, "FAIL\n" if not.
1. INDEX - attempt to index the package. An index command can only succeed if all of the package's dependencies have already been indexed. If the package was previously indexed, and the new index command succeeds, the previous dependency list is updated (overwritten). The server returns "OK\n" if the package is indexed, "FAIL\n" if not.
1. REMOVE - attempt to remove a package from the index. If the package is a dependency for any other package, the command fails. The server returns "OK\n" if the package is removed, "FAIL\n" if not.

All commands return "ERROR\n" if the command isn't understood or if the command is malformed.

## How to install it

* via Docker: from the repo root, run `docker build -t packman:0.1 .`. The docker image can then be run by mapping the application's listening port (8080) to the host interface, i.e. `docker run -p 8080:8080 -d --name packman packman:0.1`. All of this can be further automated by the user as needed.

* via building the source: I must admit I'm still adapting to Go's get/install/deploy structure, so the directory structure of this repo probably doesn't support a direct install via `go get` or `go install`. Still, you can clone this repo, and then from `packman/cmd/packman` you can run `go build packman.go` to get the `packman` binary for a given platform. The binary will then listen on `0.0.0.0:8080` when run on the host.

## How to test it

* from the repo root, run `go test`. 
* there are unit tests for the `handler` and `validator` sections of the code; current coverage is 64.7% of statements, which should ideally be increased.
* while testing locally, I was using a test harness that auotmated 10-100 concurrent clients, sending a randomized pool of commands (with varying levels of incorrectness). That test harness hasn't been included for licensing reasons, but it enhanced coverage beyond the unit tests here, and essentially provided 'integration'-level testing.

## Design rationale

* mutex (instead of channels) because packman is dealing with simple shared state, not information exchanged between goroutinues.
* map of `[string][]string` for in-memory data storage; one of the first enhancements packman could use is an abstracted storage type/struct/class that allows for looser coupling of read/write operations from stdlib (i.e. give a flexible storage solution that could include in-memory or on-disk).
* the process of checking for dependencies on `INDEX` and `REMOVE` is an expensive operation. A brute-force solution for dependency checking would result in worst-case O(n^2) runtime; I improved on that by sorting dependency lists on insertion, and using a binary search on those dependency lists (slices), which led to worst-case O(n log n) runtime.
* in the interest of readability the code base is broken into three source code files, separated by responsbilities (parsing and command handling). They aren't separate packages, which from a structural perspective isn't the standard Go layout, but sometimes (in my opinion) Go sacrifices readability for simplicity, which I push back on. You'll see the same with variables in the source code; they tend to be more verbose than idiomatic Go, but that is by design.
