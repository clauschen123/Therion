// Copyright © 2016.6 Claus Chen
//
//
package server

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
)

var (
	MessageDataSizeInvalid = errors.New("ReadPacket: package crack, invalid size")
	MessageTooBig          = errors.New("ReadPacket: package too big")
)

type IStream interface {
	Read() (*Message, error)
	Write(*Message) error
	RawConn() net.Conn
}

type Stream struct {
	conn net.Conn
}

func MakeStream(sock net.Conn) IStream {
	return &Stream{sock}
}

func (stream *Stream) RawConn() net.Conn {
	return stream.conn
}

func (stream *Stream) Read() (msg *Message, err error) {

	head := make([]byte, msg_head_len)

	if _, err = io.ReadFull(stream.conn, head); err != nil {
		return nil, err
	}

	msg = &Message{}

	headbuf := bytes.NewReader(head)

	if err = binary.Read(headbuf, binary.LittleEndian, &msg.length); err != nil {
		fmt.Println("read len err:", err)
		return nil, err
	}

	if err = binary.Read(headbuf, binary.LittleEndian, &msg.protoid); err != nil {
		fmt.Println("read protoid err:", err)
		return nil, err
	}

	if err = binary.Read(headbuf, binary.LittleEndian, &msg.msgid); err != nil {
		fmt.Println("read msgid err:", err)
		return nil, err
	}

	if msg.length > msg_max_len {
		return nil, MessageTooBig
	}

	size := msg.length - msg_head_len
	if size < 0 {
		return nil, MessageDataSizeInvalid
	}

	msg.data = make([]byte, size)
	if _, err = io.ReadFull(stream.conn, msg.data); err != nil {
		return nil, err
	}

	return
}

func (stream *Stream) Write(msg *Message) (err error) {

	buff := bytes.NewBuffer([]byte{})

	if err = binary.Write(buff, binary.LittleEndian, msg.length); err != nil {
		return
	}

	if err = binary.Write(buff, binary.LittleEndian, msg.protoid); err != nil {
		return
	}

	if err = binary.Write(buff, binary.LittleEndian, msg.msgid); err != nil {
		return
	}

	if _, err = stream.conn.Write(buff.Bytes()); err != nil {
		return err
	}

	if _, err = stream.conn.Write(msg.data); err != nil {
		return err
	}

	return
}
