package tcpserver

import (
	"log"
	"net"
	"strconv"
	"time"
)

func StartANewServer(port int) {
	mon, err := StartANewResourceMonitor()

	if err != nil {
		log.Panicln(err)
	}

	rateCntr, err := NewRateLimitController(30, time.Second/30)

	if err != nil {
		log.Panicln(err)
	}

	l, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Panicln(err)
	}
	log.Println("Listening to connections on port", strconv.Itoa(port))
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Panicln(err)
		}

		cmdHandler, err := StartANewCommandHandler(mon.statusSignal, rateCntr)

		if err != nil {
			log.Panicln(err)
		}

		go HandleRequest(conn, cmdHandler)
	}
}
