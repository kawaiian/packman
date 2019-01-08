// Just trying to see how quickly I can stand up a golang package server.

package main

import (
	"log"
	"net"
	"os"
)

const (
	host     = "0.0.0.0"
	port     = "8080"
	address  = host + ":" + port
	connType = "tcp"
)

var (
	mu		sync.Mutex	// guards package library
	pkgLib	
)

// Request holds the message from a client
type Request struct {
	request      string
	pkg          string
	dependencies []string
}

// packageLibrary is a hashMap of package names to dependency lists
type packageLibrary struct {
	pkgDeps		map[string] []string
}

func newPackageLibrary() *packageLibrary {
	return &packageLibrary {
		pkgDeps: make(map[string] []string)
	}
}


//--- main ---
func main() {
	listener, err := net.Listen(connType, address)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer listener.Close()
	log.Print("Packman server listening on " + address + "...")

	// need to create some sort of datastore object here?

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v\n", err)
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(c net.Conn) {
	defer c.Close()
	for {

	}
}
