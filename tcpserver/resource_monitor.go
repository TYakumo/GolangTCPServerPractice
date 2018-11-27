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
	IncConn        = 4
	DecConn        = 5
)

type ResourceMonitor struct {
	cmdInQueue    int
	cmdExecuted   int
	cmdThroughput int
	activeConn    int
	statusSignal  chan int
}

func (r *ResourceMonitor) StatusString() string {
	var ret string
	ret += "Cmd In Queue: " + strconv.Itoa(r.cmdInQueue) + "<br>"
	ret += "Total Cmd Executed: " + strconv.Itoa(r.cmdExecuted) + "<br>"
	ret += "Active Conn: " + strconv.Itoa(r.activeConn) + "<br>"
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
		case IncConn:
			r.activeConn++
		case DecConn:
			r.activeConn--
		}
		// r.reportStatus()
	}
}

func (r *ResourceMonitor) stopMonitor() {
	r.statusSignal <- DaemonKilled
}

func StartANewResourceMonitor() (*ResourceMonitor, error) {
	r := ResourceMonitor{0, 0, 0, 0, make(chan int)}
	go r.statusDaemon()
	return &r, nil
}
