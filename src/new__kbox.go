package main
import (
  "fmt"
  "net"
  "bufio"
  "strings"
  "net/http"
  "github.com/gorilla/websocket"
  "runtime"
  _ "github.com/mattn/go-sqlite3"
  "database/sql"
  "bytes"
  "strconv"
  "os"
  //"reflect"
)
var CloudStep = 0

var writesensor,writecloud,deviceID,NewClient,query_to_db,res_to_query,data_to_db = make(chan string),make(chan string),make(chan string),make(chan net.Conn),make(chan string),make(chan string), make(chan string)


func CloudWriter(conn net.Conn)  {
	writer := bufio.NewWriter(conn)
	for CloudWriteData := range writecloud {
		if CloudWriteData == "CLOUD DISCONNECTED" {
			return
		}
		fmt.Println("writing on cloud")
		writer.WriteString(CloudWriteData + "\n")
		writer.Flush()
	}
}

func CloudReader(conn net.Conn)  {
	CloudReaderStep := 0
	reader := bufio.NewReader(conn)
	for {
		fmt.Println("CloudReaderStep:" , CloudReaderStep)
		CloudReadData,err := reader.ReadString('\n')
		if err != nil {
		  fmt.Println(err)
		  fmt.Println("No Cloud Connection")
		  conn.Close()
		  CloudStep = 0
		  writecloud <- "CLOUD DISCONNECTED"
		  return
		} else {
			if CloudReaderStep == 0 {
				//fmt.Println("in if CloudReaderStep==0")
				DeviceID := strings.TrimSuffix(CloudReadData, "\n")
				fmt.Println(DeviceID)
				writesensor <- DeviceID
				//fmt.Println("Device name sent via channel writesensor")
				CloudReaderStep = 1
			} else {
			   writesensor <- CloudReadData
			   fmt.Println(CloudReadData)
			   //fmt.Println("CloudReaderStep:" , CloudReaderStep)
			   CloudReaderStep = 0
			}
		}
	}
}



func sensorReader (client net.Conn) {
	sensorReadStep := 0
	readclient := bufio.NewReader(client)
	for {
		data,err := readclient.ReadString('\n')
		if err != nil {
			fmt.Println("CLIENT DISCONNECTED  #######################")
			deviceID <- "DELETE"
			NewClient <- client
			//fmt.Println(err)
			//client.Close()
			return
		} else {
			if sensorReadStep == 0 {
				s:=strings.Split(data,",")
				DevID:= s[0]
				deviceID <- DevID
				NewClient <- client
				sensorReadStep = 1
			} else {
				if CloudStep == 0 {
					data_to_db <- data
				} else {
					data_to_db <- data
					writecloud <- data
					fmt.Println(data)
				}
			}
		}
	}
}

func websocket_server(){
	http.HandleFunc("/",func(w http.ResponseWriter, r *http.Request){
		fmt.Println(runtime.NumGoroutine())
		upgrader := websocket.Upgrader{}
		conn, err := upgrader.Upgrade(w,r,nil)
		//fmt.Println("type of conn:",reflect.TypeOf(conn))
		if err != nil {
			fmt.Println(err)
			return
		}
		defer conn.Close()
		//var res_to_query = make(chan string)
		go socket_writer(conn)
		for{
			_, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("readmessage",err)
				res_to_query <- "CLIENT DISCONNECTED"
				break
			}
			n:= bytes.Index(message,[]byte("\n"))
			//fmt.Println(n)
			s := string(message[:n])
			db_query:=strings.TrimSuffix(s,"\n")
			fmt.Println("db_query :",db_query)
			query_to_db <- db_query
		}
		//query_to_db <- "select * from people"
		
		/*
		default :{
			_, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("readmessage",err)
				return
			}
			n:= bytes.Index(message,[]byte("\n"))
			//fmt.Println(n)
			s := string(message[:n])
			db_query:=strings.TrimSuffix(s,"\n")
			fmt.Println("db_query :",db_query)
			query_to_db <- db_query
		}
		*/
	})	
	http.ListenAndServe(":8080",nil)		
}

func socket_writer(conn *websocket.Conn){
	for resp := range res_to_query{
		if resp == "CLIENT DISCONNECTED" {
			return
		}
		fmt.Println("received")
		conn.WriteMessage(1,[]byte(resp + "\n"))
	}
}

/*
func socket_reader(conn *websocket.Conn){
	fmt.Println("conn reader created")
	for{
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("readmessage",err)
			conn.Close()
			return
		}
		n:= bytes.Index(message,[]byte("\n"))
		//fmt.Println(n)
		s := string(message[:n])
		db_query:=strings.TrimSuffix(s,"\n")
		fmt.Println("db_query :",db_query)
		query_to_db <- db_query
	}
}
*/

	
func handle_db(){
	os.Remove("./foo.db")
	database, _ := sql.Open("sqlite3","./foo.db")
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, dev TEXT NOT NULL, x integer NOT NULL, y integer NOT NULL, z integer NOT NULL, date_time DATETIME DEFAULT CURRENT_TIMESTAMP )")
	statement.Exec()
	fmt.Println("no error till now")
	statement, _ = database.Prepare("INSERT INTO people (id,dev,x,y,z) VALUES (?,?,?,?,?)")
	//statement1,_ := database.Prepare("SELECT * FROM people WHERE id >= (?)")
	id := 0
	point := 0
	for{
		select{
			case values := <- data_to_db :{
				id++
				value := strings.TrimSuffix(values, "\n")
				fmt.Println("value :",value)
				temp :=strings.Split(value, ",")
				fmt.Println(temp)
				dev:=temp[0]
				x,_:=strconv.Atoi(temp[1])
				y,_ :=strconv.Atoi(temp[2])
				z,_:=strconv.Atoi(temp[3])
				fmt.Println("id,dev,x , y ,z=",id,dev,x,y,z)
				statement.Exec(id,dev,x,y,z)
				if (CloudStep == 0) && (point == 0){
					fmt.Println("point=id")
					point = id
				}
				if (CloudStep != 0) && (point !=0){
					//fmt.Println("hi i m here")
					rows ,err:= database.Query("SELECT * FROM people WHERE id >= (?) AND id < (?)",point,id)
					point = 0
					if err != nil{
						fmt.Println("ERROR ROWS",err)
					}
					var id string
					var dev string
					var x string
					var y string
					var z string
					var date_time string
					var total bytes.Buffer
					for rows.Next() {
						fmt.Println("in rows.next for loop")
						rows.Scan(&id,&dev, &x, &y ,&z, &date_time)
						total.WriteString(fmt.Sprintf(id + ": " + dev + " ," + x +","+ y +","+ z +"," + date_time + "\n"))
						//fmt.Println(total)
						//fmt.Println("type of total:",reflect.TypeOf(total))
					}
					writecloud <- (total.String())
					fmt.Println("sent over writecloud")
				}
			}
			case db_query := <- query_to_db :{
				rows, err:= database.Query(db_query)
				if err != nil{
					fmt.Println("ERROR ROWS",err)

				}
				var id string
				var dev string
				var x string
				var y string
				var z string
				var date_time string
				for rows.Next() {
					fmt.Println("in rows.next for loop")
					rows.Scan(&id,&dev, &x, &y ,&z, &date_time)
					total:=  fmt.Sprintf(id + ": " + dev + " ," + x +","+ y +","+ z +"," + date_time)
					fmt.Println(total)
					//fmt.Println("type of total:",reflect.TypeOf(total))
					res_to_query <- total
					fmt.Println("sent over res_to_query")
				}
			}
		}
	}
}
	
	
func Mapping(){
	var Mapper = make(map[string]net.Conn)
	for{
		select{
			case dev:= <- deviceID:{
				c := <- NewClient
				if dev=="DELETE"{
					for k, v := range Mapper {
						if v == c {
							delete(Mapper, k)
							fmt.Println("newmap:" , Mapper)
						}
					}
				} else{
					//c := <- NewClient
					Mapper[dev] = c
					fmt.Println("new connection mapped")
					fmt.Println("map:",Mapper)
				}
			}

			case dev_to_command := <- writesensor:{
				fmt.Println("dev_to_command :" ,dev_to_command)
				fmt.Println("finding device")
				conn, prs := Mapper[dev_to_command]
				if prs == true {
					fmt.Println("Device found")
					writeclient := bufio.NewWriter(conn)
					j:= <- writesensor
					temp :=strings.Split(j, ",")
					for _ , k := range temp{
					  writeclient.WriteString(k + "\n")
					  fmt.Println("k" ,k)
					 writeclient.Flush()
					}
				} else{
				  fmt.Println(prs)
				  writecloud <- "No such device"
				}
			}
		}
	}
}

func Server() {
/*  cert, err := tls.LoadX509KeyPair("public.crt", "private.key")

 	if err != nil {
 		fmt.Println(err)
 	}

 	config := tls.Config{Certificates: []tls.Certificate{cert} }
 config.Rand = rand.Reader

 	listener, err := tls.Listen("tcp", ":6600", &config)
 	if err != nil {
 		fmt.Println(err)
 	}

 	fmt.Println("Server(TLS) up and listening on port 6600")
*/
	listener, err:= net.Listen("tcp", ":6600")
	if err != nil{
		fmt.Println(err)
	}
	for {
		connection,err:=listener.Accept()
		if err != nil{
			fmt.Println(err)
    	} else {
			go sensorReader(connection)
		}
	}
}

func main(){
	//var res_to_query
	go Server()
	go Mapping()
	go websocket_server()
	go handle_db()
	for {
		if CloudStep == 0 {
		  conn,err := net.Dial("tcp",":4000")
		  if err != nil{
			fmt.Println(err)
			//fmt.Println("in main error")
			} else {
				go CloudReader(conn)
				go CloudWriter(conn)
				CloudStep = 1
			}
		}
	}
}


///////////////////////////////////////
/*
type Person struct {
    Id int
    Name string
}
rows, err := db.Query("SELECT id, name FROM person;")
if err != nil {
    // handle err
    log.Fatal(err)
}
columnNames, err := rows.Columns() // []string{"id", "name"}
if err != nil {
    // handle err
    log.Fatal(err)
}
people = make([]Person, 0, 2)
for rows.Next() {
    person := Person{}
    // person == Person{0, ""}
    pointers := make([]interface{}, len(columnNames))
    // pointers == `[]interface{}{nil, nil}`
    structVal := reflect.ValueOf(person)
    for i, colName := range columnNames {
        fieldVal := structVal.FieldByName(strings.Title(colName))
        if !fieldVal.IsValid() {
            log.Fatal("field not valid")
        }
        pointers[i] = fieldVal.Addr().Interface()
    }
    // pointers == `[]interface{}{&int, &string}`
    err := rows.Scan(pointers...)
    if err != nil {
        // handle err
        log.Fatal(err)
    }
    // person == Person{1, "John Doe"}
    people = append(people, person)
}
*/
