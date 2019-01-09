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
	req     string
	pkg     string
	depList []string
}

// --- Main ---
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

		// validate the request
		pkgReq, err := parsePkgRequest(msg)
		if err != nil {
			log.Printf("Invalid request: %v\n", err)
		}

		// here we handle each case (INDEX, QUERY, REMOVE)
		// using their own functions, each working off the mutex for locking
		// with pkgTree for storage.

		log.Printf("Received %s request!", msg[0])
	}
}

func parsePkgRequest(msg string) (pkgRequest, error) {
	var pkgReq pkgRequest

	msgParts := strings.Split(msg, "|")
	if len(msgParts) < 3 {
		return pkgReq, errors.New("invalid request")
	}

	err := parseRequest(msgParts[0])
	if err != nil {
		return pkgReq, err
	}

	err = parsePkgNames(msgParts[1], msgParts[2])
	if err != nil {
		return pkgReq, err
	}

	// Otherwise, return a request struct
	pkgReq = pkgRequest{
		req:     msgParts[0],
		pkg:     msgParts[1],
		depList: strings.Split(msgParts[2], ","),
	}

	return pkgReq, nil
}

func parseRequest(req string) error {
	validReqs := map[string]struct{}{"INDEX": {}, "QUERY": {}, "REMOVE": {}}

	_, reqOk := validReqs[req]

	if !reqOk {
		return errors.New("invalid request")
	}

	return nil
}

func parsePkgNames(pkgName string, pkgDepList string) error {
	// Compose list of package and all dependencies from msg
	fullPkgList := make([]string, len(pkgDepList)+1)
	fullPkgList = append(fullPkgList, pkgName)
	fullPkgList = append(fullPkgList, strings.Split(pkgDepList, ",")...)

	// Check to make sure both the main package and its dependencies are valid
	for _, name := range fullPkgList {
		validName, _ := regexp.MatchString(`^[a-zA-Z0-9]+[a-zA-Z0-9\-\_\+]*[a-zA-Z0-9\+]*$`, name)
		if !validName {
			return errors.New("invalid package name")
		}
	}

	return nil
}
