package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	//	"fmt"
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
	GetConn() net.Conn
}

type Stream struct {
	conn net.Conn
}

func MakeStream(conn net.Conn) IStream {
	return &Stream{conn}
}

func (stream *Stream) GetConn() net.Conn {
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
		return nil, err
	}
	//fmt.Println("**** length:", msg.length)

	if err = binary.Read(headbuf, binary.LittleEndian, &msg.protoid); err != nil {
		return nil, err
	}
	//fmt.Println("**** protoid:", msg.protoid)

	if err = binary.Read(headbuf, binary.LittleEndian, &msg.msgid); err != nil {
		return nil, err
	}
	//fmt.Println("**** msgid:", msg.msgid)

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

	// 防止将Send放在go内造成的多线程冲突问题
	//	self.sendtagGuard.Lock()
	//	defer self.sendtagGuard.Unlock()

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
