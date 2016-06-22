package main

import (
	"fmt"
)

var (
	entering = make(chan *Session)
	leaving  = make(chan *Session)
	messages = make(chan *Message)
)

func run() {

	sessionMgr := make(SessionMgr) // all connected sessionMgr
	for {
		select {
		case msg := <-messages:
			fmt.Println("Get 1 message, proto:", msg.protoid,
				"msgid:", msg.msgid,
				"content:", string(msg.data))

			for _, v := range sessionMgr {
				v.ch <- msg
			}

		case session := <-entering:
			sessionMgr[session.id] = session

		case session := <-leaving:
			close(session.ch)
			delete(sessionMgr, session.id)
		}
	}
}
