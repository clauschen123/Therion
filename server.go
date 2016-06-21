package main

import (
	"fmt"
	"log"
	"net"
)

var (
	entering = make(chan *Session)
	leaving  = make(chan *Session)
	messages = make(chan string) // all incoming client messages
)

//主逻辑: 遍历每个session，处理消息
func server_run() {
	clients := make(SessionMgr) // all connected clients
	for {
		select {
		case msg := <-messages:
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			for _, v := range clients {
				v.ch <- msg
			}

		case cli := <-entering:
			clients[cli.id] = cli

		case cli := <-leaving:
			close(cli.ch)
			delete(clients, cli.id)
		}
	}
}

// Server start with addr
func Start(addr string) error {

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
		return err
	}

	go server_run()

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Print(err)
				continue
			}
			fmt.Println("Accpet a connect")
			MakeSession(conn)
		}
	}()

	return err
}
