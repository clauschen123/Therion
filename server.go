package main

import (
	"fmt"
	"log"
	"net"
)

type EHostType int8

const (
	e_host_none EHostType = iota
	e_host_client
	e_host_gate
	e_host_game
	e_host_center
	e_host_db
	e_host_logger
	e_host_console
)

var (
	server        IServer = nil
	entering              = make(chan *Connection, 10)
	leaving               = make(chan *Connection, 10)
	queue                 = make(chan func(), 100)
	connectionMgr         = make(ConnectionMgr)
)

//!+连接模式
func Connect(svr IServer, addr string) error {

	server = svr

	go svr.Run()

	socket, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
		return err
	}
	svr.OnConnected(0, socket)

	return nil
}

//!+侦听模式
func Start(svr IServer, addr string) error {

	server = svr

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
		return err
	}

	go svr.Run()

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Print(err)
				continue
			}
			fmt.Println("Accpet a connect")
			svr.OnAccepted(conn)
		}
	}()

	return err
}

type IServer interface {
	Init(EHostType) error

	OnConnected(uint32, net.Conn)
	OnAccepted(net.Conn)

	Run()

	Enter(*Connection)
	Exit(*Connection)
	Post(*Connection, *Message)
}

//!+ Sample server
type Server struct {
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

	return nil
}

func (this *Server) Enter(conn *Connection) {
	entering <- conn
}

func (this *Server) Exit(conn *Connection) {
	leaving <- conn
}

func (this *Server) Post(conn *Connection, msg *Message) {
	queue <- func() { conn.HandleMessage(msg) }
}

func (this *Server) OnConnected(sid uint32, socket net.Conn) {
	conn := MakeConnection(socket)

	conn.RegisterProtocol(MakeProtocol(e_protoid_system))

	conn.Post(MakeMessage(e_protoid_system, e_msgid_auth, []byte("Hello server!")))
	this.Enter(conn)
}

func (this *Server) OnAccepted(socket net.Conn) {
	conn := MakeConnection(socket)

	conn.RegisterProtocol(MakeProtocol(e_protoid_system))

	conn.Post(MakeMessage(e_protoid_system, e_msgid_info, []byte("You are "+conn.name)))

	this.Enter(conn)
	this.Post(conn, MakeMessage(e_protoid_system, e_msgid_connected, []byte(conn.name+" has arrived")))
}
