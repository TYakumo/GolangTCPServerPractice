package tcpserver

import (
	"log"
	"strconv"
)

var (
	DaemonKilled   = -1
	IncCmdInQue    = 1
	DecCmdInQue    = 2
	IncCmdExecuted = 3
)

type ResourceMonitor struct {
	cmdInQueue    int
	cmdExecuted   int
	cmdThroughput int
	statusSignal  chan int
}

func (r *ResourceMonitor) StatusString() string {
	var ret string
	ret += "Total Cmd Executed: " + strconv.Itoa(r.cmdExecuted) + "\n"
	return ret
}

func (r *ResourceMonitor) reportStatus() {
	log.Println("cmdInQueue: ", r.cmdInQueue, " cmdExecuted: ", r.cmdExecuted)
}

func (r *ResourceMonitor) statusDaemon() {
	for {
		signalReceived := <-r.statusSignal
		switch signalReceived {
		case DaemonKilled:
			return
		case IncCmdInQue:
			r.cmdInQueue++
		case DecCmdInQue:
			r.cmdInQueue--
		case IncCmdExecuted:
			r.cmdExecuted++
		}
		r.reportStatus()
	}
}

func (r *ResourceMonitor) stopMonitor() {
	r.statusSignal <- DaemonKilled
}

func StartANewResourceMonitor() (*ResourceMonitor, error) {
	r := ResourceMonitor{0, 0, 0, make(chan int)}
	go r.statusDaemon()
	return &r, nil
}
