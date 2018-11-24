package main

import (
	"log"

	"./httpserver"
	"./tcpserver"
)

func main() {

	mon, err := tcpserver.StartANewResourceMonitor()

	if err != nil {
		log.Panicln(err)
	}

	go httpserver.StartANewHTTPServer(8080, mon)
	tcpserver.StartANewTCPServer(3333, mon)
}
