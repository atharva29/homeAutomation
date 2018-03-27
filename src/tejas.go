package main
import(
  "fmt"
  "net"
  "bufio"
  "strings"
  "runtime"
  "os"
  "database/sql"
  "bytes"
  "strconv"
  "net/http"
  _ "github.com/mattn/go-sqlite3"
  "github.com/gorilla/websocket"
)
//var CloudStep = 0
var conn net.Conn


func main(){
/////////////////////////////////////////////////DATABASE
fmt.Println("main")
os.Remove("./tejas.db")       // removes database file
database, _ := sql.Open("sqlite3","./tejas.db")       // creates a new db file
statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS student (PKEY INTEGER NULL,ID INTEGER NULL, NAME TEXT NULL)")  //create table
statement.Exec()  // execute create table statement
statement, _ = database.Prepare("INSERT INTO student (PKEY,ID,NAME) VALUES (?,?,?)") // make statement for entering values afterwards
////////////////////////////////////////////////END OF DATABASE



  //Connect to tcp connection
go func(){    listener,err :=net.Listen("tcp",":6000")// create  a listener
                  if err != nil {
                    fmt.Println(err)
                  }
                  conn,err := listener.Accept() // accept connection from CLIENT
                  if err != nil {
                    fmt.Println(err)
                  }
                  Reader :=bufio.NewReader(conn)  // create reader for reading from client

               go func(){
                ///////////////////////////////READER/////////////
                for{
                data,err:=Reader.ReadString('\n')
                if err != nil {
                    fmt.Println("CLIENT DISCONNECTED  #######################")
                    fmt.Println(err)
                    conn.Close()
                    return
                    }
                  data1:=strings.TrimSuffix(data,"\n")
                  temp:=strings.Split(data1,",")
                  PKEY,_:=strconv.Atoi(temp[0])
                  ID,_:=strconv.Atoi(temp[1])
                  NAME:=temp[2]
                  statement.Exec(PKEY,ID,NAME)
                  }
                  }()
                  ///////////////////////END OF Reader

          }()


            http.HandleFunc("/",func(w http.ResponseWriter, r *http.Request){
            fmt.Println(runtime.NumGoroutine())
            upgrader := websocket.Upgrader{}
            connect, err := upgrader.Upgrade(w,r,nil)
            if err != nil {
              fmt.Println(err)
            }
            defer connect.Close()
            for {
            _, message, err := connect.ReadMessage()
              if err != nil {
                fmt.Println("echo",err)
                return
              }
            n:= bytes.Index(message,[]byte("\n"))
            s := string(message[:n])
            j:=strings.TrimSuffix(s,"\n")
            fmt.Println("Message :",j)
            rows, err:= database.Query(j)
            if err != nil{
              fmt.Println("ERROR ROWS",err)
            }
            fmt.Println("Inside GIVEDB")
            var PKEY string
            var ID string
            var NAME string
            for rows.Next() {
                rows.Scan(&PKEY,&ID,&NAME)//,&InsertDateTime)
                total:=  fmt.Sprintf(PKEY + ": " + ID+ " ," + NAME )// +", "+ data)//+" "+ InsertDateTime)
                fmt.Println(total)
                connect.WriteMessage(1,[]byte(total + "\n"))
              }
            }
          })
          http.ListenAndServe(":8080",nil)

}
