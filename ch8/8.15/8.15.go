/*
	Failure of any client program to read data in a timely
	manner ultimately causes all clients to get stuck. Modify
	the broadcaster to skip a message rather than wait if a
	client writer is not ready to accept it. Alternatively, add
	buffering to each client's outgoing message channel so that
	most messages are not dropped; the broadcaster should use a
	non-blocking send to this channel
*/
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

//!+broadcaster
type client struct {
	MsgChan  chan<- string // an outgoing message channel
	Name     string
	timedOut *bool
}

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string) // all incoming client messages
)

func timeClient(conn net.Conn, timedOut *bool) {
	for !*timedOut {
		*timedOut = true
		time.Sleep(30 * time.Second)
	}
	conn.Close()
}

func broadcaster() {
	clients := make(map[client]bool) // all connected clients
	for {
		select {
		case msg := <-messages:
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			for cli := range clients {
				cli.MsgChan <- msg
			}

		case cli := <-entering:
			cli.MsgChan <- "Current users connected: "
			for client := range clients {
				cli.MsgChan <- client.Name
			}
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
			close(cli.MsgChan)
		}
	}
}

//!-broadcaster

//!+handleConn
func handleConn(conn net.Conn) {
	ch := make(chan string, 100) // outgoing client messages
	go clientWriter(conn, ch)

	cli := client{MsgChan: ch, timedOut: new(bool)}
	input := bufio.NewScanner(conn)
	ch <- "Enter your username: "
	input.Scan()
	cli.Name = input.Text()
	messages <- cli.Name + " has arrived"
	entering <- cli
	go timeClient(conn, cli.timedOut)
	for input.Scan() {
		messages <- cli.Name + ": " + input.Text()
		*cli.timedOut = false
	}
	// NOTE: ignoring potential errors from input.Err()

	leaving <- cli
	messages <- cli.Name + " has left"
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}

//!-handleConn

//!+main
func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}
