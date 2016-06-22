package main

import (
	"fmt"
	"log"
	//"server"
)

func server() {
	err := Start("localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	var str string
	fmt.Scan(&str)
}

func client() {
	err := Connect("localhost:8000")
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
