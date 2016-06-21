package main

import (
	"fmt"
	"log"
	//"server"
)

func main() {

	err := Start("localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	var str string
	fmt.Scan(&str)
}
