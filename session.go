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

func (session *Session) Post(msg *Message) {
	session.ch <- msg
}

func MakeSession(conn net.Conn) *Session {

	ch := make(chan *Message)
	name := conn.RemoteAddr().String()

	session := &Session{0, name, true, ch, MakeStream(conn)}

	go session.send()
	go session.receive()

	return session
}

func (session *Session) receive() {

	for {
		msg, err := session.stream.Read()

		if err != nil {
			fmt.Println(session.name, " disconnected!")
			break
		}
		fmt.Println("Succ read 1 message ...")
		queue <- msg
	}

	leaving <- session
	queue <- MakeMessage(e_protoid_base, e_msgid_logout, []byte(session.name+" has left"))
	session.GetConn().Close()
}

func (session *Session) send() {

	var ch <-chan *Message = session.ch
	for msg := range ch {
		//fmt.Fprintln(session.GetConn(), msg) // NOTE: ignoring network errors
		if err := session.stream.Write(msg); err != nil {
			fmt.Println("stream write msg fail:", err)
			break
		}
		fmt.Println(" ------Send a msg")
	}
}
