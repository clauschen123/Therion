package main

import (
	"fmt"
	"net"
)

type SessionMgr map[uint](*Session)
type Session struct {
	id        uint
	name      string
	connected bool
	ch        chan *Message
	stream    IStream
}

func (session *Session) GetConn() net.Conn {
	return session.stream.GetConn()
}

func MakeSession(conn net.Conn) *Session {

	ch := make(chan *Message)
	name := conn.RemoteAddr().String()

	session := &Session{0, name, true, ch, MakeStream(conn)}

	go session.session_send()
	session.ch <- MakeMessage(1, 2, []byte("You are "+name))

	go session.session_recv()
	messages <- MakeMessage(1, 2, []byte(name+" has arrived"))
	entering <- session

	return session
}

func (session *Session) session_recv() {

	for {

		msg, err := session.stream.Read()

		if err != nil {
			fmt.Println(session.name, " disconnected!")
			break
		}
		fmt.Println("Succ read 1 message ...")
		messages <- msg

	}

	leaving <- session
	messages <- MakeMessage(1, 2, []byte(session.name+" has left"))
	session.GetConn().Close()
}

func (session *Session) session_send() {
	for msg := range session.ch {
		//fmt.Fprintln(session.GetConn(), msg) // NOTE: ignoring network errors
		if err := session.stream.Write(msg); err != nil {
			fmt.Println("stream write msg fail:", err)
			break
		}
		fmt.Println(" ------Send a msg")
	}
}
