package main

import (
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})


	http.HandleFunc("/ws1", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil)
		go func(conn *websocket.Conn) {
			ch := time.Tick(5 * time.Second)

			for range ch {
				conn.WriteJSON(myStruct{
					Lat:  19.0274,
					Lng: 72.856689,
				})
			}
		}(conn)
	})

	http.ListenAndServe(":3000", nil)
}

type myStruct struct {
	Lat  float32 `json:""`
	Lng float32 `json:""`
}
