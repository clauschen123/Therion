// Copyright © 2016.6 Claus Chen
//
//

package server

type MsgID uint32

const (
	msg_head_len uint32 = 12
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

//type Message struct {
//	length uint32
//	data   []byte
//}

//!+ New a message
func MakeMessage(pid ProtoID, msgid MsgID, data []byte) *Message {

	var size uint32 = msg_head_len + uint32(len(data))
	return &Message{
		message_head{size, pid, msgid},
		message_body{data}}
}
