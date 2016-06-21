package main

import (
	"bufio"
	"fmt"
	"net"
)

type Session struct {
	id        uint
	name      string
	connected bool
	ch        chan string
	conn      net.Conn
}

type SessionMgr map[uint](*Session)

func MakeSession(conn net.Conn) {

	ch := make(chan string)
	//go session.session_send(conn, ch) //这里如果不立即go一下，ch <- "You are " + name这句会死

	name := conn.RemoteAddr().String()
	//ch <- "You are " + name

	session := &Session{0, name, true, ch, conn}
	fmt.Println(name, " has arrived 1")

	messages <- name + " has arrived"
	entering <- session

	go session.session_recv()
	go session.session_send() //这里如果不立即go一下，ch <- "You are " + name这句会死

	fmt.Println(name, " has arrived 2")
	session.ch <- "You are " + name
	fmt.Println(name, " has arrived 3")
}

func (session *Session) session_recv() {
	input := bufio.NewScanner(session.conn)
	for input.Scan() {
		fmt.Println(input.Text())
		messages <- session.name + ": " + input.Text()
	}
	// NOTE: ignoring potential errors from input.Err()
	fmt.Println(session.name + " has left")

	leaving <- session
	messages <- session.name + " has left"
	session.conn.Close()
}

func (session *Session) session_send() {
	for msg := range session.ch {
		fmt.Fprintln(session.conn, msg) // NOTE: ignoring network errors
	}
}
