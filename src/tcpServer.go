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
  "strings"
  "database/sql"
  "strconv"
)


var indicator = 1

var pass  = make(chan string)   // channel for data
var query_to_db = make(chan string) // channel for database
var res_to_query = make(chan string)
var passServer = make (chan string)


func main(){
fmt.Println("Start")
go tcpClient()
l,err := net.Listen("tcp",":2000")  //set a server listening at port 20000
printErr(err)

  for {
            conn,err := l.Accept()
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
              if (text == "echo\n"){    // this condition is used for checking connection
                pass <-text
              }
            fmt.Println("text : " + text)   // print the recieved data
            //dbChannel <- text
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
                    fmt.Println("GOTTA")
                    conn.Close()      // closes the connection
                    return
              }
            }
         }
      }
    }



   func database(conn net.Conn){
     database, _ := sql.Open("sqlite3","./attendance.db")       // creates a new db file
     statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS student (ID INTEGER NULL, NAME TEXT NULL)")  //create table
     statement.Exec()  // execute create table statement
     statement, _ = database.Prepare("INSERT INTO student (ID,NAME) VALUES (?,?)") // make statement for entering values afterwards

                      Writer := bufio.NewWriter(conn)// create a new writer
                      go func(){
                        for{
                          select {
                            case data:=<-query_to_db:{
                              data1:=strings.TrimSuffix(data,"\n")
                              temp:=strings.Split(data1,",")
                              ID,_:=strconv.Atoi(temp[0])
                              NAME:=temp[1]
                              fmt.Println(temp)
                              statement.Exec(ID,NAME)
                              Writer.WriteString("Executed in DB")
                              fmt.Println(ID,NAME)
                                  }
                                }
                            }
                          }()
                }



func TcpServerWriter(conn net.Conn){
  writer:= bufio.NewWriter(conn) // makes a new writer for port
  for {
      select{
      case values := <- passServer :{
        if (values != "close"){
        writer.WriteString(values)    //write string to port
        writer.Flush()                // clear the buffer
          }  else if (values == "close"){ // if reading error is found then close the connection
                    fmt.Println("GOTTA")
                    conn.Close()      // closes the connection
                    return
              }
            }
         }
      }
}


func TcpServerReader(conn net.Conn){
  reader := bufio.NewReader(conn)     // set a new reader for tcp
                // for loop is used for continuos reading from port
  for {
      text,err:=reader.ReadString('\n')
      if err != nil {
        fmt.Println("Server DISCONNECTED")
        passServer <- "close"   // this will close the writer
        return
        } else {
            if (text == "echo\n"){    // this condition is used for checking connection
              passServer <-text
            }
          fmt.Println("text : " + text)   // print the recieved data
          //dbChannel <- text
          }
        }
  }

  func tcpClient(){

      //  if indicator ==1 {
          conn,err:=net.Dial("tcp",":6600")
          if err != nil{
            fmt.Println(err)
            return
          }
        //  } else {
            go TcpServerReader(conn)
            go TcpServerWriter(conn)
            indicator = 1
        for{}
        //  }
        }
