package main

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{}

func simpleHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func tickTick(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)
	go func(conn *websocket.Conn) {
		ch := time.Tick(5 * time.Second)
		for range ch {
			conn.WriteJSON(myStruct{
				// Lat:    19.0274,
				// Lng:    72.856689,
				Lat : 18.9322 ,
				Lng : 72.8264 ,
				Colour: 0,
			})
		}
	}(conn)
}

func main() {
	router := mux.NewRouter()
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("/home/vjti/website")))
	router.HandleFunc("/", simpleHandler)
	router.HandleFunc("/ws1", tickTick)

	http.ListenAndServe(":3000", nil)
}

type myStruct struct {
	Lat    float32 `json:""`
	Lng    float32 `json:""`
	Colour int     `json:""`
}
