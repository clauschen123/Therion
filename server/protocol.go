// Copyright © 2016.6 Claus Chen
//
//

package server

type ProtoID int32

//!+ Protocol
type IProtocol interface {
	GetProtID() ProtoID

	RegisterMessage(msgid MsgID, hdlr func(*Message))
	HandleMessage(msg *Message)
}

type Protocol struct {
	Pid     ProtoID
	Handler map[MsgID]func(*Message)
}

func (this *Protocol) GetProtID() ProtoID {
	return this.Pid
}

func (this *Protocol) RegisterMessage(msgid MsgID, hdlr func(*Message)) {
	this.Handler[msgid] = hdlr
}

func (this *Protocol) HandleMessage(msg *Message) {
	this.Handler[msg.msgid](msg)
}
