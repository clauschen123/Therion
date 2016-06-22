package main

import (
	"fmt"
	"log"
	"net"
)

func Start(addr string) error {

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
		return err
	}

	go run()

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Print(err)
				fmt.Println("Accpet fail:", err)
				continue
			}
			fmt.Println("Accpet a connect")
			MakeSession(conn)
		}
	}()

	return err
}
