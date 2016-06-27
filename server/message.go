// Copyright © 2016.6 Claus Chen
//
//

package main

type MsgID uint16

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

//!+ New a message
func MakeMessage(pid ProtoID, msgid MsgID, data []byte) *Message {

	var size uint32 = msg_head_len + uint32(len(data))
	return &Message{
		message_head{size, pid, msgid},
		message_body{data}}
}
