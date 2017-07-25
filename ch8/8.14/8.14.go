/*
	Change the chat server's network protocol so that
	each client provides its name on entering. Use
	that name instead of the network address when
	prefixing each message with its sender's identity.
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
	ch := make(chan string) // outgoing client messages
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()
	ch <- "You are " + who
	messages <- who + " has arrived"
	cli := client{ch, who, new(bool)}
	entering <- cli
	go timeClient(conn, cli.timedOut)

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- who + ": " + input.Text()
		*cli.timedOut = false
	}
	// NOTE: ignoring potential errors from input.Err()

	leaving <- cli
	messages <- who + " has left"
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
