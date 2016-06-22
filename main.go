package main

import (
	"fmt"
	"log"
)

func server() {

	var server Server

	server.Init()
	server.RegisterProtocol(MakeProtocol(e_protoid_base))

	err := Start(&server, "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	var str string
	fmt.Scan(&str)
}

func client() {

	var server Server

	server.Init()
	server.RegisterProtocol(MakeProtocol(e_protoid_base))

	err := Connect(&server, "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	var str string
	fmt.Scan(&str)
}

func main() {
	//server()
	client()
}
