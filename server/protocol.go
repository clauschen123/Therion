// Copyright © 2016.6 Claus Chen
//
//

package server

import (
	"fmt"
)

type ProtoID uint16

const (
	e_protoid_system ProtoID = iota
	e_protoid_num
)

//TODO
var SystemProto IProtocol = MakeProtocol(e_protoid_system)

//!+ Protocol
type IProtocol interface {
	GetProtID() ProtoID

	RegisterMessage(msgid MsgID, hdlr func(*Message))
	HandleMessage(msg *Message)
}

type Protocol struct {
	pid     ProtoID
	handler map[MsgID]func(*Message)
}

func (this *Protocol) GetProtID() ProtoID {
	return e_protoid_system
}

func (this *Protocol) RegisterMessage(msgid MsgID, hdlr func(*Message)) {
	this.handler[msgid] = hdlr
}

func (this *Protocol) HandleMessage(msg *Message) {
	this.handler[msg.msgid](msg)
}

//!+ New a protocol
func MakeProtocol(protid ProtoID) IProtocol {
	proto := &Protocol{

		pid: protid,

		handler: map[MsgID]func(*Message){

			e_msgid_connected: func(msg *Message) {
				fmt.Println("Get connected: ", msg.protoid, msg.msgid, string(msg.data))
			},

			e_msgid_disconnect: func(msg *Message) {
				fmt.Println("Get disconnect: ", msg.protoid, msg.msgid, string(msg.data))
			},

			e_msgid_echo: func(msg *Message) {
				fmt.Println("Get echo: ", msg.protoid, msg.msgid, string(msg.data))
			},

			e_msgid_info: func(msg *Message) {
				fmt.Println("Get info: ", msg.protoid, msg.msgid, string(msg.data))
			},

			e_msgid_auth: func(msg *Message) {
				fmt.Println("Get auth: ", msg.protoid, msg.msgid, string(msg.data))
				//server.Post(MakeMessage(e_protoid_system, e_msgid_info, []byte("Hello client!")))
			},
		},
	}
	return proto
}
