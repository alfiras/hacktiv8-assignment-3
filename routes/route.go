package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"hacktiv8-assignment-3/configs"
	"hacktiv8-assignment-3/data"
	"hacktiv8-assignment-3/socket"

	"github.com/gorilla/websocket"
)

func RunRouting() {
	http.HandleFunc("/", handleWebsocket)
	http.HandleFunc("/polling", handlePolling)
	http.HandleFunc("/long-polling", handleLongPolling)

	http.HandleFunc("/data", handleGetData)

	http.HandleFunc("/ws", handleWs)

	fmt.Printf("Server starting in port %s \n", configs.PORT)
	http.ListenAndServe(configs.PORT, nil)
}

func handleWebsocket(w http.ResponseWriter, r *http.Request) {
	html, err := ioutil.ReadFile("pages/websocket.html")
	if err != nil {
		http.Error(w, "Cant open file", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%s", html)
}

func handleLongPolling(w http.ResponseWriter, r *http.Request) {
	html, err := ioutil.ReadFile("pages/long-polling.html")
	if err != nil {
		http.Error(w, "Cant open file", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%s", html)
}

func handlePolling(w http.ResponseWriter, r *http.Request) {
	html, err := ioutil.ReadFile("pages/polling.html")
	if err != nil {
		http.Error(w, "Cant open file", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%s", html)
}

func handleWs(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrade(w, r, w.Header(), configs.READBUFFSIZE, configs.WRITEBUFFSIZE)
	if err != nil {
		http.Error(w, "Cant open websocket connection", http.StatusBadRequest)
	}
	socket.AddConnection(conn)
	go handleIOWebsocket(conn)
}

func handleGetData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	dataStatus := data.Data{}
	if r.Method == "GET" {
		dataStatus = data.ReadFromJson()
		json.NewEncoder(w).Encode(dataStatus)
		return
	}
	http.Error(w, "Method not allowed", http.StatusBadRequest)
}

func handleIOWebsocket(conn *websocket.Conn) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Err", fmt.Sprintf("%v", r))
		}
	}()

	conn.WriteJSON(data.ReadFromJson())
	fmt.Println("initial read on new connection")

	for {

		payload := socket.SocketPayload{}
		err := conn.ReadJSON(&payload)
		if err != nil {
			if strings.Contains(err.Error(), "websocket: close") {
				socket.RemoveConnection(conn)
				return
			}

			log.Println("Err", err.Error())
			continue
		}

	}
}
