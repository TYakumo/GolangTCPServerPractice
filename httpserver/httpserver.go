package httpserver

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"../tcpserver"
	"github.com/gorilla/websocket"
)

type DashboardViewHandler struct {
	mon *tcpserver.ResourceMonitor
}

func StartANewHTTPServer(port int, mon *tcpserver.ResourceMonitor) {
	dhandler := DashboardViewHandler{mon}
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	http.HandleFunc("/dashboard", dhandler.dashboardHandler)
	http.HandleFunc("/", dhandler.homeHandler)
	ServerAddr := "localhost:" + strconv.Itoa(port)
	h := http.Server{Addr: ServerAddr}

	go func() {
		err := h.ListenAndServe()
		if err != nil {
			log.Printf("Listen and serve error: %v", err)
			return
		}
	}()

	for {
		select {
		case <-interrupt:
			log.Println("Httpserver Received Interrupts")
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			//shutdown the server
			err := h.Shutdown(ctx)
			if err != nil {
				log.Printf("Shutdown request error: %v", err)
			}
			return
		}
	}
}

var upgrader = websocket.Upgrader{} // use default options

func (d *DashboardViewHandler) dashboardHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	timeoutSec := 10

	ticker := time.NewTicker(time.Duration(200) * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			nextTime := time.Now().Add(time.Duration(timeoutSec) * time.Second)
			c.SetWriteDeadline(nextTime)
			message := d.mon.StatusString()
			err = c.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				log.Println("write:", err)
				return
			}
		}
	}
}

func (d *DashboardViewHandler) homeHandler(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/dashboard")
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>
window.addEventListener("load", function(evt) {
    var printBox = document.getElementById("serverMsg");
    var ws;

	if (ws) {
		return false;
	}
	ws = new WebSocket("{{.}}");
	ws.onopen = function(evt) {
		printBox.innerHTML = "Listening to Server";
	}
	ws.onclose = function(evt) {
		ws = null;
	}
	ws.onmessage = function(evt) {
		printBox.innerHTML = evt.data;
	}
	ws.onerror = function(evt) {
		printBox.innerHTML = "ERROR: " + evt.data;
	}
});
</script>
</head>
<body>
<p id="serverMsg">No data yet</p>
</body>
</html>
`))
