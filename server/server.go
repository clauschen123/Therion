// Copyright © 2016.6 Claus Chen
//
//

package server

import (
	"fmt"
	"log"
	"net"
)

type EHostType int8

const (
	E_host_none EHostType = iota
	E_host_client
	E_host_gate
	E_host_game
	E_host_center
	E_host_db
	E_host_logger
	E_host_console
)

var (
	server        IServer = nil
	entering              = make(chan *Connection, 10)
	leaving               = make(chan *Connection, 10)
	queue                 = make(chan func(), 100)
	connectionMgr         = make(ConnectionMgr)
)

func GetServer() IServer {
	return server
}

//!+连接模式
func Connect(svr IServer, addr string) error {

	server = svr

	go server.Run()

	socket, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
		return err
	}
	svr.OnConnected(0, socket) //TODO use 0 temperaly

	return nil
}

//!+侦听模式
func Accept(svr IServer, addr string) error {

	server = svr

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
		return err
	}

	go server.Run()

	go func() {
		for {
			socket, err := listener.Accept()
			if err != nil {
				log.Print(err)
				continue
			}
			fmt.Println("Accpet a connect")
			svr.OnAccepted(0, socket) //TODO
		}
	}()

	return err
}

type IServer interface {
	Init(EHostType) error
	RegisterProtocol(IProtocol)

	OnConnected(SID, net.Conn)
	OnAccepted(SID, net.Conn)

	Run()

	Enter(*Connection)
	Exit(*Connection)
	Post(IConnectionHandler, *Message)
}

//!+ Sample server
type Server struct {
	curConnection *Connection
	svrHandler    IConnectionHandler
}

func (this *Server) Run() {

	for {
		select {
		case handler := <-queue:
			handler()

		case c := <-entering:
			connectionMgr[c.id] = c

		case c := <-leaving:
			close(c.ch)
			delete(connectionMgr, c.id)
		}
	}
}

func (this *Server) Init(host EHostType) error {

	this.svrHandler = MakeConnectionHandler()
	this.RegisterProtocol(ProtoSystem)

	return nil
}

func (this *Server) RegisterProtocol(proto IProtocol) {
	this.svrHandler.RegisterProtocol(proto)
}

func (this *Server) Enter(conn *Connection) {
	entering <- conn
}

func (this *Server) Exit(conn *Connection) {
	leaving <- conn
}

func (this *Server) Post(hdlr IConnectionHandler, msg *Message) {
	queue <- func() { hdlr.HandleMessage(msg) }
}

func (this *Server) OnConnected(sid SID, socket net.Conn) {
	conn := MakeConnection(socket, this.svrHandler)
	this.svrHandler.Established(sid, conn)

	//TODO
	//	conn.Post(MakeMessage(e_protoid_system, e_msgid_auth, []byte("Hello server!")))
	ProtoSystem.(*Protocol).Send2Server_MsgAuthServer(sid)

	this.Enter(conn)
}

func (this *Server) OnAccepted(sid SID, socket net.Conn) {

	conn := MakeConnection(socket, this.svrHandler)
	this.svrHandler.Established(sid, conn)

	//TODO
	//	conn.Post(MakeMessage(e_protoid_system, e_msgid_info, []byte("You are "+conn.name)))

	this.Enter(conn)
}
