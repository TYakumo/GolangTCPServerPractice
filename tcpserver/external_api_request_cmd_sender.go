package tcpserver

import (
	"time"
)

func RunNoop(errChan chan error) {
	errChan <- nil
}

func RunDelayingNoop(errChan chan error) {
	time.Sleep(time.Duration(8) * time.Second)
	errChan <- nil
}
