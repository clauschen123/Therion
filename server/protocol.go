// Copyright © 2016.6 Claus Chen
//
//

package server

type ProtoID int32

//!+ IProtocol
//目前一个协议组只能被一个连接处理
type IProtocol interface {
	GetProtID() ProtoID

	SetConn(IConnectionHandler)
	GetConn() IConnectionHandler

	RegisterMessage(msgid MsgID, hdlr func(*Message))
	HandleMessage(msg *Message)
}

//!+ Protocol
type Protocol struct {
	Pid         ProtoID
	MsgHandler  map[MsgID]func(*Message)
	ConnHandler IConnectionHandler
}

func (this *Protocol) GetProtID() ProtoID {
	return this.Pid
}

func (this *Protocol) SetConn(conn IConnectionHandler) {
	this.ConnHandler = conn
}

func (this *Protocol) GetConn() IConnectionHandler {
	return this.ConnHandler
}

func (this *Protocol) RegisterMessage(msgid MsgID, hdlr func(*Message)) {
	this.MsgHandler[msgid] = hdlr
}

func (this *Protocol) HandleMessage(msg *Message) {
	this.MsgHandler[msg.msgid](msg)
}
