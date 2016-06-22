package main

const (
	msg_head_len uint32 = 8
	msg_max_len  uint32 = 10 * 1024
)

var (
	total_msgid uint16 = 0
)

type message_head struct {
	length  uint32
	protoid uint16
	msgid   uint16
}

type message_body struct {
	data []byte
}

type Message struct {
	message_head
	message_body
}

func MakeMessage(protoid uint16, msgid uint16, data []byte) *Message {
	msgid = total_msgid + 1
	total_msgid = msgid

	var size uint32 = msg_head_len + uint32(len(data))
	return &Message{
		message_head{size, protoid, msgid},
		message_body{data}}
}
