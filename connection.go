// Copyright © 2016.6 claus chen
// Connection: the management of the socket connectivity

package main

import (
	"fmt"
	"net"
)

type ConnectionMgr map[uint64](*Connection)

type Connection struct {
	id        uint64
	name      string
	connected bool
	ch        chan *Message
	stream    IStream
	handler   map[ProtoID]IProtocol //每个连接所对应的处理器
}

func MakeConnection(stream net.Conn) *Connection {

	ch := make(chan *Message)
	name := stream.RemoteAddr().String()

	c := &Connection{0, name, true, ch, MakeStream(stream), make(map[ProtoID]IProtocol)}

	go c.send()
	go c.receive()

	return c
}

func (this *Connection) RawConn() net.Conn {
	return this.stream.RawConn()
}

func (this *Connection) RegisterProtocol(proto IProtocol) {
	protid := proto.GetProtID()
	this.handler[protid] = proto
}

func (this *Connection) Post(msg *Message) {
	this.ch <- msg
}

func (this *Connection) HandleMessage(msg *Message) {

	if proto, ok := this.handler[msg.protoid]; ok {
		proto.HandleMessage(msg)
	} else {
		fmt.Println("protoid not find:", msg.protoid)
	}
}

func (this *Connection) receive() {

	for {
		msg, err := this.stream.Read()

		if err != nil {
			fmt.Println(this.name, " disconnected!")
			break
		}
		server.Post(this, msg)
	}

	leaving <- this
	server.Post(this, MakeMessage(e_protoid_system, e_msgid_disconnect, []byte(this.name+" has left")))

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
