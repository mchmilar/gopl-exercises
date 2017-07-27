// Clock is a TCP server that periodically writes the time.
/*
	Modify clock2 to accept a port number, and write a program,
	clockwall, that acts as a client of several clock servers
	at once, reading the times from each one and displating the
	the results in a tale, akin to the wall of clocks seen in
	some business offices. If you have access to geopgraphically
	distributed computers, run instances remotely; otherwise run
	local instances on different ports with fake timezones
*/
package main

import (
	"flag"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

var port = flag.Int("port", 8000, "define the port to listen on")

func handleConn(c net.Conn) {
	defer c.Close()
	for {
		_, err := io.WriteString(c, time.Now().Format("15:04:05\n"))
		if err != nil {
			return // e.g., client disconnected
		}
		time.Sleep(1 * time.Second)
	}
}

func main() {
	flag.Parse()
	laddr := "localhost:" + strconv.Itoa(*port)
	listener, err := net.Listen("tcp", laddr)
	if err != nil {
		log.Fatal(err)
	}
	//!+
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		go handleConn(conn) // handle connections concurrently
	}
	//!-
}
