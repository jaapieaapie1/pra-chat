package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"os"
)

var Connections []*websocket.Conn
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
func main() {


	http.HandleFunc("/socket", ServeWs)

	err := http.ListenAndServe(":" + os.Getenv("PORT"), nil)
	if err != nil {
		panic(err)
	}
}

func ServeWs(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}

	conn.SetCloseHandler(func(code int, text string) error {
		deleteConn(conn)
		return nil
	})

	Connections = append(Connections, conn)

	for {
		_, message, err := conn.ReadMessage();
		if err != nil {
			_ = conn.Close()
			deleteConn(conn)
			return
		}

		for _, connection := range Connections {
			_ = connection.WriteMessage(websocket.TextMessage, message)
		}

	}
}

func deleteConn(conn *websocket.Conn) {
	index := -1
	for i, connection := range Connections {
		if connection == conn {
			index = i
			break
		}
	}

	if index == -1 {
		return
	}

	Connections = append(Connections[:index], Connections[index+1:]...)


}