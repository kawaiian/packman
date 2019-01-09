// Just trying to see how quickly I can stand up a golang package server.

package main

import (
	"bufio"
	"errors"
	"log"
	"net"
	"os"
	"regexp"
	"strings"
	"sync"
)

const (
	addr     = host + ":" + port
	connType = "tcp"
	host     = "0.0.0.0"
	port     = "8080"
)

var (
	mu      sync.RWMutex
	pkgTree map[string][]string
)

type pkgRequest struct {
	req         string
	pkg         string
	depndencies []string
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

	// Accept connections, concurrently
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
		msg, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			log.Printf("Error reading request: %v\n", err)
			return
		}

		// validate the request; if valid, return a pkgRequest object
		pkgReq, err := parseRequest(msg)

		// here we handle each case (INDEX, QUERY, REMOVE)
		// using their own functions, each working off the mutex for locking
		// with pkgTree for storage.

		log.Printf("Received %s request!", msg[0])
	}
}

func parseRequest(msg string) (pkgRequest, error) {
	var pkgReq pkgRequest
	validReqs := map[string]struct{}{"INDEX": {}, "QUERY": {}, "REMOVE": {}}

	msgParts := strings.Split(msg, "|")
	req := msgParts[0]
	_, reqOk := validReqs[req]

	if len(msgParts) < 3 || !reqOk {
		return pkgReq, errors.New("Invalid request")
	}

	validator, _ := regexp.Compile(`^[a-zA-Z0-9]+[a-zA-Z0-9\-\_\+]*[a-zA-Z0-9\+]*$`)

	// Otherwise, return a request struct
}
