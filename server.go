package main

import (
	"fmt"
	"log"
	"net"
)

var (
	entering = make(chan *Session)
	leaving  = make(chan *Session)
	queue    = make(chan *Message)

	sessionMgr = make(SessionMgr)

	protoHandler [e_protoid_num]*Protocol = [e_protoid_num]*Protocol{nil}
)

type IServer interface {
	Init() error
	RegisterProtocol(*Protocol)

	OnConnected(uint32, net.Conn)
	OnAccepted(net.Conn)

	Run()

	Enter(*Session)
	Exit(*Session)
	Post(*Message)
}

//!+a sample server
type Server struct {
}

func (server *Server) Init() error {
	return nil
}

func (server *Server) Enter(session *Session) {
	entering <- session
}

func (server *Server) Exit(session *Session) {
	leaving <- session
}

func (server *Server) Post(msg *Message) {
	queue <- msg
}

func (server *Server) OnConnected(sid uint32, conn net.Conn) {
	session := MakeSession(conn)

	server.Enter(session)
}

func (server *Server) OnAccepted(conn net.Conn) {
	session := MakeSession(conn)

	session.Post(MakeMessage(e_protoid_base, e_msgid_info, []byte("You are "+session.name)))
	server.Post(MakeMessage(e_protoid_base, e_msgid_login, []byte(session.name+" has arrived")))
	server.Enter(session)

}

func (server *Server) Run() {

	for {
		select {
		case msg := <-queue:
			server.HandleMessage(msg)

		case session := <-entering:
			sessionMgr[session.id] = session

		case session := <-leaving:
			close(session.ch)
			delete(sessionMgr, session.id)
		}
	}
}

func (server *Server) RegisterProtocol(proto *Protocol) {
	protoid := proto.protoid
	protoHandler[protoid] = proto
}

func (server *Server) HandleMessage(msg *Message) {
	if msg.protoid < e_protoid_num {
		protoHandler[msg.protoid].HandleMessage(msg)
	}
}

//!+客户端模式
func Connect(server IServer, addr string) error {

	go server.Run()

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
		return err
	}
	server.OnConnected(0, conn)

	return nil
}

//!+服务端模式
func Start(server IServer, addr string) error {

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
		return err
	}

	go server.Run()

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Print(err)
				fmt.Println("Accpet fail:", err)
				continue
			}
			fmt.Println("Accpet a connect")
			server.OnAccepted(conn)
		}
	}()

	return err
}
