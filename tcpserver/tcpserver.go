package tcpserver

import (
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"time"
)

func StartANewTCPServer(port int, mon *ResourceMonitor) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
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

	go func() {
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
	}()

	for {
		select {
		case <-interrupt:
			log.Println("tcpserver received interrupts and terminated")
			return
		}
	}
}
