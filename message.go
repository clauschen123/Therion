package main

import (
	"fmt"
)

type MsgID uint16
type ProtoID uint16

const (
	e_protoid_base ProtoID = iota
	e_protoid_num
)

const (
	e_msgid_base MsgID = iota
	e_msgid_login
	e_msgid_logout
	e_msgid_echo
	e_msgid_info
	e_msgid_num
)

const (
	msg_head_len uint32 = 8
	msg_max_len  uint32 = 10 * 1024
)

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

type Protocol struct {
	protoid ProtoID
	handler map[MsgID]func(*Message)
}

func (protocol *Protocol) HandleMessage(msg *Message) {
	protocol.handler[msg.msgid](msg)
}

func MakeMessage(protoid ProtoID, msgid MsgID, data []byte) *Message {

	var size uint32 = msg_head_len + uint32(len(data))
	return &Message{
		message_head{size, protoid, msgid},
		message_body{data}}
}

func MakeProtocol(protoid ProtoID) *Protocol {
	proto := &Protocol{
		protoid: protoid,
		handler: map[MsgID]func(*Message){

			e_msgid_login:  func(msg *Message) { fmt.Println("Get login: ", msg.protoid, msg.msgid, string(msg.data)) },
			e_msgid_logout: func(msg *Message) { fmt.Println("Get logout: ", msg.protoid, msg.msgid, string(msg.data)) },
			e_msgid_echo:   func(msg *Message) { fmt.Println("Get echo: ", msg.protoid, msg.msgid, string(msg.data)) },
			e_msgid_info:   func(msg *Message) { fmt.Println("Get info: ", msg.protoid, msg.msgid, string(msg.data)) },
		},
	}
	return proto
}
