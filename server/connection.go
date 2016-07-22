// Copyright © 2016.6 Claus Chen
//
//

package server

import (
	"fmt"
	"net"
)

type SID uint64
type CID uint64
type CIDS []CID
type SIDS []SID

type ConnectionMgr map[uint64](*Connection)

//!+ ConnectionHandler
type IConnectionHandler interface {
	RegisterProtocol(IProtocol)
	HandleMessage(*Message)
	Established(SID, *Connection)

	Send2Client(CID, *Message) bool
	Send2Clients(CIDS, *Message) bool
	Send2Server(SID, *Message) bool
	Send2Servers(SIDS, *Message) bool
}

type ConnectionHandler struct {
	proto    map[ProtoID]IProtocol
	sid2conn map[SID](*Connection)
}

func MakeConnectionHandler() IConnectionHandler {
	return &ConnectionHandler{
		proto:    map[ProtoID]IProtocol{},
		sid2conn: map[SID](*Connection){}}
}

func (this *ConnectionHandler) getConnBySID(sid SID) *Connection {
	for k, v := range this.sid2conn {
		if k == sid {
			return v
		}
	}
	return nil
}

func (this *ConnectionHandler) RegisterProtocol(proto IProtocol) {
	protid := proto.GetProtID()
	this.proto[protid] = proto
}

func (this *ConnectionHandler) HandleMessage(msg *Message) {

	if proto, ok := this.proto[msg.protoid]; ok {
		proto.HandleMessage(msg)
	} else {
		fmt.Println("protoid not find:", msg.protoid)
	}
}

func (this *ConnectionHandler) Established(sid SID, conn *Connection) {
	//TODO
	this.sid2conn[sid] = conn
}

func (this *ConnectionHandler) Send2Client(cid CID, msg *Message) bool {
	//TODO
	return true
}

func (this *ConnectionHandler) Send2Clients(cids CIDS, msg *Message) bool {
	//TODO
	return true
}

func (this *ConnectionHandler) Send2Server(sid SID, msg *Message) bool {
	//TODO
	return true
}

func (this *ConnectionHandler) Send2Servers(sids SIDS, msg *Message) bool {
	//TODO
	if len(sids) == 0 {
		fmt.Println("sids empty !")
		for _, conn := range this.sid2conn {
			conn.Post(msg)
		}
	} else {
		for _, sid := range sids {
			fmt.Println("Find conn:", sid)
			conn := this.getConnBySID(sid)
			conn.Post(msg)
		}
	}
	return true
}

//!+ Connection
type Connection struct {
	id        uint64
	name      string
	connected bool
	ch        chan *Message
	stream    IStream
	handler   IConnectionHandler
}

func MakeConnection(stream net.Conn, hdlr IConnectionHandler) *Connection {

	ch := make(chan *Message)
	name := stream.RemoteAddr().String()

	c := &Connection{
		id:        0,
		name:      name,
		connected: true,
		ch:        ch,
		stream:    MakeStream(stream),
		handler:   hdlr}

	go c.send()
	go c.receive()

	return c
}

func (this *Connection) SetHandler(hdlr IConnectionHandler) {
	this.handler = hdlr
}

func (this *Connection) RawConn() net.Conn {
	return this.stream.RawConn()
}

func (this *Connection) Post(msg *Message) {
	this.ch <- msg
}

func (this *Connection) receive() {

	for {
		msg, err := this.stream.Read()

		if err != nil {
			fmt.Println(this.name, " disconnected!")
			break
		}
		server.Post(this.handler, msg)
	}

	leaving <- this
	//TODO	server.Post(this.handler, MakeMessage(e_protoid_system, e_msgid_disconnect, []byte(this.name+" has left")))

	this.RawConn().Close()
}

func (this *Connection) send() {

	var ch <-chan *Message = this.ch
	for msg := range ch {

		if err := this.stream.Write(msg); err != nil {
			fmt.Println("stream write msg fail:", err)
			break
		}

	}
}
