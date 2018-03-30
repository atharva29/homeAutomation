package main
// this file will connect to tcp client and will act like a tcp server
/*  NEVER SEND DATA TO CHANNEL , WHICH IS NOT GETTING READER
ALWAYS READ DATA FROM INCOMING CHANNEL
*/


import (
  "fmt"     //libaray use for print function
  "net"     // used for net
  "bufio"   // used for taking inputs from user or port
  "os"      // used to read user inputs
  _ "github.com/mattn/go-sqlite3" // database
  "strings"     //
  "database/sql"
  "strconv"
)


var indicator = 0   // indicator = 0 , not connected to server , if its '1' means connected to server
var num = 0
var pass  = make(chan string)   // channel for Cient data
var data_to_db = make(chan string) // channel for database
//var res_to_query = make(chan string)
var passServer = make (chan string) // channel for RESPONSE OF  Server query database

var query_to_db = make (chan string)

func main(){
fmt.Println("Start")
l,err := net.Listen("tcp",":2000")  //set a server listening at port 20000
if err!=nil{
  printErr(err)
  }
go tcpClient() // function for communication with Server

  for {
            conn,err := l.Accept() // accept incoming client connection
            if err != nil{
                fmt.Println(err)
                }
             go Reader(conn)      //function for reading data from user input
             go res(conn)        //function for reading data from client
             go writer(conn)     // writer writes data to port
        }
  }




//function :- This function reads data from CLIENT and display it
func res(conn net.Conn){
    reader := bufio.NewReader(conn)     // set a new reader for tcp
                  // for loop is used for continuos reading from port
    for {
        text,err:=reader.ReadString('\n')
        if err != nil {
          fmt.Println("Client DISCONNECTED")
          pass <- "close"   // this will close the writer
          return
          } else {
            fmt.Println("text : " + text)   // print the recieved data
           data_to_db <- text
            }
          }
      }


//function:- this function read input from user
func Reader(conn net.Conn){
  osReader := bufio.NewReader(os.Stdin)   // make a READER
  for {
  msg,err:=osReader.ReadString('\n')
  if err != nil {
    fmt.Println("Client DISCONNECTED")
    return
    }else {
    pass <- msg         // pass the user input to channel

      }
    }
  }


//function :- this function is used for printing error
func printErr (err error){    // this function is used for printing error
  if err != nil{
  fmt.Println(err)
  }
}



// function :- this function writes data to tcp port
func writer(conn net.Conn){
  writer:= bufio.NewWriter(conn) // makes a new writer for port
  for {
      select{
      case values := <- pass :{
        if (values != "close"){
        writer.WriteString(values)    //write string to port
        writer.Flush()                // clear the buffer
          }  else if (values == "close"){ // if reading error is found then close the connection
                    conn.Close()      // closes the connection
                    return
              }
            }
         }
      }
    }



   func database(conn net.Conn){
     os.Remove("./attendance.db")
   	 database, _ := sql.Open("sqlite3","./attendance.db")       // creates a new db file
     statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS student (num INTEGER PRIMARY KEY,ID INTEGER NULL, NAME TEXT NULL,date_time DATETIME DEFAULT CURRENT_TIMESTAMP )")  //create table
     statement.Exec()  // execute create table statement
     statement, _ = database.Prepare("INSERT INTO student (num,ID,NAME) VALUES (?,?,?)") // make statement for entering values afterwards

     //writer:= bufio.NewWriter(conn) // makes a new writer for port

                        for{
                          select {
                            case data:=<-data_to_db:{
                              num = num+1
                              data1:=strings.TrimSuffix(data,"\n")
                              temp:=strings.Split(data1,",")
                              ID,_:=strconv.Atoi(temp[0])
                              NAME:=temp[1]
                              statement.Exec(num,ID,NAME)
                              fmt.Println(num,ID,NAME)
                                  }


                            case db_query := <- query_to_db :{
                              //query := strings.TrimSuffix(db_query,"\n")
                              rows, err:= database.Query(db_query)
                      				if err != nil{
                      					fmt.Println("ERROR ROWS",err)
                                return
                                }
                      				var num string
                      				var ID string
                      				var NAME string
                      				var date_time string
                      				for rows.Next() {
                      					//fmt.Println("in rows.next for loop")
                      					rows.Scan(&num,&ID, &NAME, &date_time)
                      					total:=  fmt.Sprintf(num + ": " + ID + " ," + NAME )
                      					fmt.Println(total)
                                //fmt.Println("total -----")
                                //fmt.Println("type of total:",reflect.TypeOf(total))
                      				  passServer <- total
                      					//fmt.Println("sent over res_to_query")
                                //    writer.WriteString(total)
                              }
                      			}
                          }
                        }
                }


// function for writing data to server , used to pass response for query to server
func TcpServerWriter(conn net.Conn){
  writer:= bufio.NewWriter(conn) // makes a new writer for port

  for {
      select{
      case values := <- passServer :{
        if (values != "close"){
        writer.WriteString(values +"\n")    //write string to port
        writer.Flush()                // clear the buffer
          }  else if (values == "close"){ // if reading error is found then close the connection
                    conn.Close()      // closes the connection
                    indicator = 0     // server DISCONNECTED
                    return
              }
            }
         }
      }
}

func serverWriter(conn net.Conn){
  writer:= bufio.NewWriter(conn) // makes a new writer for port
  for values := range passServer{
      if (values != "close"){
      writer.WriteString(values +"\n")    //write string to port
      writer.Flush()                // clear the buffer
        }  else if (values == "close"){ // if reading error is found then close the connection
                  conn.Close()      // closes the connection
                  indicator = 0     // server DISCONNECTED
                  return
          }
    }
}


// function for reading commands from server
func TcpServerReader(conn net.Conn){
  reader := bufio.NewReader(conn)     // set a new reader for tcp
                // for loop is used for continuos reading from port
  for {
      text,err:=reader.ReadString('\n') // read string from server
      if err != nil {
        fmt.Println("Server DISCONNECTED")
        passServer <- "close"   // this will close the writer
        return
        } else {

          fmt.Println("Server COMMAND : " + text)   // print the recieved data
          query_to_db <-text
          //passServer <-text
          }
        }
  }


// function for communication with server
func tcpClient(){
    for {
    // tries to connect with server
    if indicator == 0 {
      fmt.Println("Connecting .....")
      conn,err := net.Dial("tcp",":6600")
        if err != nil{

        } else {
          // connected to server
                indicator = 1
                fmt.Println("Connected")
                go TcpServerReader(conn)// read incoming data
                go TcpServerWriter(conn) // write data to server
                go database(conn)
         }
      }
    }
  }
