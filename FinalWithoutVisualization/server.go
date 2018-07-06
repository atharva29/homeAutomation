package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	// "os"
	"runtime"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// Data struct holds the slice(dynamic Array) of info struct
type Data struct {
	// IsNew bool   `json:""`
	Slice []Info `json:""`
}

// Info Holds the Message that to be passed to the website
type Info struct {
	Id          int
	UpdateValue int
	Lat         float32 //`json:""`
	Lnd         float32 //`json:""`
	Message     string  //`json:""`
	Amount      int
}

// type websocketConnection struct {
// 	connection *websocket.Conn
// }

var upgrader = websocket.Upgrader{}
var giveTtnData = make(chan Info)
var addConnection = make(chan *websocket.Conn)
var deleteConnection = make(chan *websocket.Conn)

func readConnection(conn *websocket.Conn) {
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error in reading")
			deleteConnection <- conn
			return
		}

	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("inside websocket Establisher")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("Upgrading To Websocket Error : ", err)
		return
	}
	addConnection <- conn
	go readConnection(conn)
	// write with connection
}

func simpleHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("inside serving map.html")
	http.ServeFile(w, r, "Map.html")
}

// Mapping takes is a gorouting that listen from Every Channel and takes appropiate Decision
func Mapping() {
	temp := 0
	fmt.Println("Inside Mapper")
	fmt.Println(runtime.NumGoroutine())
	var something []Info
	something = append(something, Info{
		Id:          0,
		Lat:         19.022380,
		Lnd:         72.855598,
		Message:     "VJTI Main Gate Dustbin",
		Amount:      20,
		UpdateValue: 0,
	}, Info{
		Id:          1,
		Lat:         19.023362,
		Lnd:         72.856953,
		Message:     "VJTI Back Gate DUSTBIN",
		Amount:      40,
		UpdateValue: 0,
	}, Info{
		Id:          2,
		Lat:         19.024508,
		Lnd:         72.856466,
		Message:     "VJTI Chemical Department Gate DUSTBIN",
		Amount:      60,
		UpdateValue: 0,
	}, Info{
		Id:          3,
		Lat:         19.019235,
		Lnd:         72.855112,
		Message:     "VJTI Hostel Gate DUSTBIN",
		Amount:      80,
		UpdateValue: 0,
	})
	DataHolder := Data{
		Slice: something,
	}
	var Mapper = make(map[int]*websocket.Conn)
	for {
		select {
		case conn := <-addConnection:
			{
				fmt.Println("Message Refresh : Recieved a New Connection")
				Mapper[temp] = conn
				fmt.Println("Mapper ", Mapper)
				temp = temp + 1
				fmt.Println("temp ", temp)
				conn.WriteJSON(DataHolder)
				fmt.Println("Message Refresh : New Data Sent Successfully for Refresh")
			}
		case deleteConn := <-deleteConnection:
			{
				for k, v := range Mapper {
					if v == deleteConn {
						delete(Mapper, k)
						fmt.Println("newmap:", Mapper)
					}
				}
			}
		case t := <-giveTtnData:
			{
				fmt.Println("Message Update : Recieved Data From Channel For Update")
				if t.Id >= 0 {
					count := 0
					for index, value := range DataHolder.Slice {
						count++
						if value.Id == t.Id {
							t.UpdateValue = index
							DataHolder.Slice[index] = t
							fmt.Println("Message Update : Data Updated at Index Value = ", index)
						}
					}
					if t.UpdateValue == 0 {
						t.UpdateValue = count
						DataHolder.Slice = append(DataHolder.Slice, t)
					}
				} else {
					fmt.Println("Message Update : Wrong Device ID")
					break
				}
				for k := range Mapper {
					err := Mapper[k].WriteJSON(DataHolder)
					if err != nil {
						fmt.Println("Message Update : Error Writing Json To Website")
					}
					fmt.Println(DataHolder.Slice)
					DataHolder.Slice[len(DataHolder.Slice)-1].UpdateValue = 0
				}
				fmt.Println("Message Update : Data Sent Successfully")
			}
		}
	}
}

func parseGhPost(rw http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var t Info
	err := decoder.Decode(&t)
	fmt.Println(t)
	if err != nil {
		panic(err)
	}
	giveTtnData <- t
}

func main() {
	go Mapping()
	router := mux.NewRouter()
	router.HandleFunc("/", simpleHandler)
	router.HandleFunc("/coehttpintegration", parseGhPost)
	router.HandleFunc("/ws", handler)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("C:\\Users\\parthesh\\go\\src\\FinalWithoutVisualization")))
	http.ListenAndServe(":8050", router)
	// err := http.ListenAndServe(":"+os.Getenv("PORT"), router)
	// if err != nil {
	// 	panic(err)
	// }
}
