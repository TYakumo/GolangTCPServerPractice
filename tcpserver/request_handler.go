package tcpserver

import (
	"log"
	"net"
	"time"
)

func HandleRequest(conn net.Conn, cmdHandler *CmdHandler) {
	log.Println("A new connection accepted.")
	defer conn.Close()
	defer log.Println("Connection closed.")
	bufSize := 1024
	buf := make([]byte, bufSize)

	timeoutSec := 30

	for {
		nextTime := time.Now().Add(time.Duration(timeoutSec) * time.Second)
		conn.SetReadDeadline(nextTime)
		size, err := conn.Read(buf)

		if err != nil {
			log.Println(err)
			return
		}

		data := buf[:size]
		log.Println("Received API Request", string(data))
		opcode, err := cmdHandler.ExecuteCommand(string(data))
		if err != nil {
			log.Println("Executing API Request failed", err)
			conn.Write([]byte(err.Error() + "\n"))
		}

		if opcode == QuitOpcode {
			return
		}

		if opcode < 0 {
			conn.Write([]byte("Received unknown command\n"))
		}
	}
}
