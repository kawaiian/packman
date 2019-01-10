// Just trying to see how quickly I can stand up a golang package server.

package main

import (
	"bufio"
	"errors"
	"fmt"
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
			log.Printf("Error reading message from client: %v\n", err)
			return
		}

		// validate the request
		pkgReq, err := parsePkgRequest(msg)
		if err != nil {
			log.Printf("Invalid request: %v\n", err)
		} else {
			switch pkgReq.req {
			case "INDEX":
				fmt.Printf("Received INDEX request for %s with dependencies %s\n", pkgReq.pkg, pkgReq.depList)
			case "QUERY":
				fmt.Printf("Received QUERY request for %s\n", pkgReq.pkg)
			case "REMOVE":
				fmt.Printf("Received REMOVE request for %s\n", pkgReq.pkg)
			}
		}
	}
}

func parsePkgRequest(msg string) (pkgRequest, error) {
	var pkgReq pkgRequest

	// Requests are expected to be "<command>|<pkg>|<dependency>\n"
	msgParts := strings.Split(msg, "|")
	if len(msgParts) < 3 {
		return pkgReq, errors.New("invalid request format")
	}

	err := parseRequest(msgParts[0])
	if err != nil {
		return pkgReq, err
	}

	dListString := strings.TrimSuffix(msgParts[2], "\n")
	pkgDepList := strings.Split(dListString, ",")
	err = parsePkgNames(msgParts[1], pkgDepList)
	if err != nil {
		return pkgReq, err
	}

	// Otherwise, return a request struct
	pkgReq = pkgRequest{
		req:     msgParts[0],
		pkg:     msgParts[1],
		depList: pkgDepList,
	}

	return pkgReq, nil
}

// parseRequest Validates request type is one of INDEX, QUERY, or REMOVE
func parseRequest(req string) error {
	validReqs := map[string]struct{}{"INDEX": {}, "QUERY": {}, "REMOVE": {}}

	_, reqOk := validReqs[req]
	if !reqOk {
		return errors.New("invalid request type: not INDEX, QUERY, or REMOVE")
	}

	return nil
}

// parsePkgNames validates the format of package and dependency names
func parsePkgNames(pkgName string, pkgDepList []string) error {

	fullPkgList := make([]string, 1)
	fullPkgList[0] = pkgName
	if pkgDepList[0] != "" {
		fullPkgList = append(fullPkgList, pkgDepList...)
	}

	validName := false
	validator := regexp.MustCompile(`^[a-zA-Z0-9]+[a-zA-Z0-9\-\_\+]*[a-zA-Z0-9\+]*$`)
	for _, name := range fullPkgList {
		validName = validator.MatchString(name)
		if !validName {
			return errors.New("invalid package name")
		}
	}

	return nil
}
