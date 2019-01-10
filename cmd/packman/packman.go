// -------
// Packman provides a concurrency-safe Linux package index.
// Please see the README for more information!
//
// Author: Kawai Washburn <kawaiwashburn@gmail.com>
// -------

package main

import (
	"bufio"
	"log"
	"net"
	"os"

	"github.com/kawaiian/packman"
)

const (
	addr     = host + ":" + port
	connType = "tcp"
	host     = "0.0.0.0"
	port     = "8080"
)

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
	var response string

	defer c.Close()
	log.Printf("Serving %s\n", c.RemoteAddr().String())

	for {
		msg, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			log.Printf("Error reading message from client: %v\n", err)
			return
		}

		// validate the request
		pkgReq, err := packman.ParsePkgRequest(msg)
		if err != nil {
			log.Printf("Invalid request: %v\n", err)
			response = "ERROR"
		} else {
			response = packman.HandlePkgRequest(pkgReq)
		}
		c.Write([]byte(response + "\n"))
	}
}
