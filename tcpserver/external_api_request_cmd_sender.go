package tcpserver

import (
	"time"
)

func RunNoop() error {
	return nil
}

func RunDelayingNoop() error {
	time.Sleep(time.Duration(8) * time.Second)
	return nil
}
