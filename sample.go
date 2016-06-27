package main

import (
	"fmt"
	"log"

	"github.com/clauschen123/Therion/"
)

func StartServer() {

	var game Server

	if err := game.Init(e_host_game); err != nil {
		log.Fatal(err)
		return
	}

	if err := Accept(&game, "0.0.0.0:8000"); err != nil {
		log.Fatal(err)
	}

	var str string
	fmt.Scan(&str)
}

func StartClient() {

	var gate Server

	if err := gate.Init(e_host_gate); err != nil {
		log.Fatal(err)
		return
	}

	if err := Connect(&gate, "127.0.0.1:8000"); err != nil {
		log.Fatal(err)
	}

	var str string
	fmt.Scan(&str)
}

func main() {
	//StartServer()
	StartClient()
}
