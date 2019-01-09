// Just trying to see how quickly I can stand up a golang package server.

package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

const (
	host     = "0.0.0.0"
	port     = "8080"
	addr     = host + ":" + port
	connType = "tcp"
)

var (
	mu      sync.RWMutex
	pkgTree map[string][]string
)

type Request struct {
	request      string
	pkg          string
	dependencies []string
}

//--- main ---
func main() {
	lstnr, err := net.Listen(connType, addr)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer lstnr.Close()
	log.Print("Packman server listening on " + addr + "...")

	for {
		conn, err := lstnr.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v\n", err)
			continue
		}
		go handleRequest(conn)
	}
}

func handleRequest(c net.Conn) {
	defer c.Close()
	log.Printf("Serving %s\n", c.RemoteAddr().String())

	for {
		req, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			log.Printf("Error reading request: %v\n", err)
			return
		}

		// validate the request; if valid, return a Request object

		// here we handle each case (INDEX, QUERY, REMOVE)
		// using their own functions, each working off the mutex for locking
		// with pkgTree for storage.
		pkgReq := strings.Split(req, "|")
		log.Printf("Received %s request!", req[0])
	}
}

func parseRequest(req string) (Request, error) {
	// Check that the request is one of INDEX, QUERY, REMOVE
	//otherwise return an error

	// For all the package names (both the package and potential dependencies)
	// Validate with a regular expression
}
