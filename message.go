package main

import (
	"fmt"
)

type MsgID uint16
type ProtoID uint16

const (
	e_protoid_system ProtoID = iota
	e_protoid_num
)

const (
	e_msgid_connected MsgID = iota
	e_msgid_disconnect
	e_msgid_auth
	e_msgid_echo
	e_msgid_info
	e_msgid_num
)

const (
	msg_head_len uint32 = 8
	msg_max_len  uint32 = 10 * 1024
)

//!+ Message
type message_head struct {
	length  uint32
	protoid ProtoID
	msgid   MsgID
}

type message_body struct {
	data []byte
}

type Message struct {
	message_head
	message_body
}

//!+ Protocol
type IProtocol interface {
	GetProtID() ProtoID
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

//!+ New a message
func MakeMessage(pid ProtoID, msgid MsgID, data []byte) *Message {

	var size uint32 = msg_head_len + uint32(len(data))
	return &Message{
		message_head{size, pid, msgid},
		message_body{data}}
}

//!+ New a protocol
func MakeProtocol(protid ProtoID) IProtocol {
	proto := &Protocol{
		pid: protid,
		handler: map[MsgID]func(*Message){

			e_msgid_connected: func(msg *Message) {
				fmt.Println("Get login: ", msg.protoid, msg.msgid, string(msg.data))
			},

			e_msgid_disconnect: func(msg *Message) {
				fmt.Println("Get logout: ", msg.protoid, msg.msgid, string(msg.data))
			},

			e_msgid_echo: func(msg *Message) {
				fmt.Println("Get echo: ", string(msg.data))
			},

			e_msgid_info: func(msg *Message) {
				fmt.Println("Get info: ", msg.protoid, msg.msgid, string(msg.data))
			},

			e_msgid_auth: func(msg *Message) {
				fmt.Println("Get auth: ", msg.protoid, msg.msgid, string(msg.data))
				server.Post(MakeMessage(e_protoid_system, e_msgid_info, []byte("Hello client!")))
			},
		},
	}
	return proto
}
