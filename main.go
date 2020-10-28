package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/mrod502/ddb/ddb"
)

func main() {
	err := ddb.OpenDB()
	if err != nil {
		panic(err)
	}
	defer ddb.CloseDB()
	go ddb.ServeTLS()
	handleExit()
}

func handleExit() {
	exitChan := make(chan os.Signal)

	signal.Notify(exitChan, os.Interrupt)
	<-exitChan
	fmt.Println("Ctrl+C pressed, exiting")
}
