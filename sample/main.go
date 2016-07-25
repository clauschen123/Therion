package main

import (
	"fmt"
	"log"

	"github.com/clauschen123/Therion/server"
)

func StartServer() {

	var game server.Server

	if err := game.Init(server.E_host_game); err != nil {
		log.Fatal(err)
		return
	}

	if err := server.Accept(&game, "0.0.0.0:8000"); err != nil {
		log.Fatal(err)
	}

	var str string
	fmt.Scan(&str)
}

func StartClient() {

	var gate server.Server

	if err := gate.Init(server.E_host_gate); err != nil {
		log.Fatal(err)
		return
	}

	if err := server.Connect(&gate, "127.0.0.1:8000"); err != nil {
		log.Fatal(err)
	}

	var str string
	fmt.Scan(&str)
}

func main() {
	StartServer()
	//StartClient()
}
