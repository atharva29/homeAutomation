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



var passFlag = make(chan int)
var flag = 0 // flag
var indicator = 0   // indicator = 0 , not connected to server , if its '1' means connected to server
var num = 1   // auto increment of database
var num1 = 1
var num2 = 1
var pass  = make(chan string)   // channel for Cient data
var data_to_db = make(chan string) // channel for database
//var res_to_query = make(chan string)
var passServer = make(chan string) // channel for RESPONSE OF  Server query database

var query_to_db = make(chan string) //passess query passed by user to database
var deviceID = make(chan string) // it carries device id to map or unmap deviceID on Mapper
var NewClient = make(chan net.Conn)// it carries connection of deviceID,going to mount on Mapper


func main(){
  V1 := make([]float64 , 8)
  V2 := make([]float64 , 8)
  I :=  make([]float64 , 2)

//fmt.Println(V1)
go database()

fmt.Println("Start")
l,err := net.Listen("tcp",":2000")  //set a server listening at port 20000
if err!=nil{
  printErr(err)
  }
go tcpClient() // function for communication with Server
go Mapper()
  for {
            conn,err := l.Accept() // accept incoming client connection
            if err != nil{
                fmt.Println(err)
                }
             go Reader()      //function for reading data from user input
             go res(conn,V1,V2,I)        //function for reading data from client
             go writer(conn)     // writer writes data to port
        }
  }




//function :- This function reads data from CLIENT and display it
func res(conn net.Conn,V1 []float64 ,V2 []float64 ,I []float64){
    var dev string
    readerStep := 0
    reader := bufio.NewReader(conn)     // set a new reader for tcp
                  // for loop is used for continuos reading from port
    for {
        text,err:=reader.ReadString('\n')
        if err != nil {
          fmt.Println("Client DISCONNECTED")
        //  pass <- "close"   // this will close the writer
          deviceID <-"DELETE"
          NewClient <-conn
          conn.Close()
          return
          } else {
            if readerStep == 0{
              temp :=strings.Split(text, ",")
              dev =temp[1]
              dev1 :=strings.TrimSuffix(dev,"\n")
              deviceID <- dev1
              NewClient <- conn
              fmt.Println("deviceID :=" , dev)
              readerStep = 1
            }


           fmt.Println("text : " + text)   // print the recieved data
           data_to_db <- text
            }
          }
      }


// this function maintains a map of active connections with their device ID
func Mapper(){
  var Mapper = make(map[string]net.Conn)
  for{
      select{
      case dev:=<-deviceID:{
        c := <- NewClient

        if dev=="DELETE"{
          for k, v := range Mapper {
            if v == c {
              delete(Mapper, k)
            //  fmt.Println("newmap:" , Mapper)
                }
                  }
              } else{
              Mapper[dev] = c
          //    fmt.Println("new connection mapped")
          //    fmt.Println("map:",Mapper)
          //    fmt.Println(len(Mapper))
                }
        }
     }
     missingConnection(Mapper)

  }
}

func missingConnection(Mapper map[string]net.Conn){
  var key =make([]string,0)
  fmt.Println(len(Mapper))

  if len(Mapper) == 3 {
    fmt.Println("ALL DEVICES Connected .......................")
    //return
  }else if len(Mapper) == 0{
    fmt.Println("No Devices Connected....................")
    flag = 0
    return
  }
  for k,_:= range Mapper{
  key = append(key,k)
  }
  //fmt.Println("K:=",key)

  i := 0

    if len(key)==1{
        if key[i]=="ard1"{
          fmt.Println("Missing ARD2 and ARD3 ..............")
          flag = 1

          } else if key[i] == "ard2"{
          fmt.Println("Missing ARD1 and ARD3 .............")
          flag = 2

          } else{
            fmt.Println("Missing ARD1 and ARD2 .............")
            flag = 3
          }

        } else if len(key)==2{
              if key[i] =="ard1" || key[i] == "ard2"{
                if key[i+1] == "ard2" || key[i+1] =="ard1"{
                      fmt.Println("Missing ARD3 ................")
                      flag = 4
                      }
                }
                if key[i] == "ard2" || key[i] == "ard3"{
                      if key[i+1] =="ard3" || key[i+1] =="ard2" {
                        fmt.Println("Missing ARD1 ...............")
                        flag = 5
                      }
                }
                if key[i] == "ard1" || key[i] == "ard3"{
                      if key[i+1] =="ard3" || key[i+1] =="ard1" {
                        fmt.Println("Missing ARD2 ...............")
                        flag = 6
                      }
                 }
           }
}


func makeZero(temp []float64){
  for i :=0 ; i <len(temp) ; i++{
    temp[i] = 0
  }
}

//
// func intToString(values []int){
//     valuesText := []string{}
//     for i := range values {
//         number := values[i]
//         text := strconv.Itoa(number)
//         valuesText = append(valuesText, text)
//       }
//       result := strings.Join(valuesText, "+")
//       return result
// }


//function:- this function read input from user
func Reader(){
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
        //if (values != "close"){
        fmt.Println("Data Sent")
        writer.WriteString(values)    //write string to port
        writer.Flush()                // clear the buffer
            }
         }
      }
    }



   func database(){
     os.Remove("./attendance.db")
   	 database, _ := sql.Open("sqlite3","./attendance.db")       // creates a new db file
     statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS student (num INTEGER PRIMARY KEY,ID INTEGER NULL, NAME TEXT NULL,date_time DATETIME DEFAULT CURRENT_TIMESTAMP )")  //create table
     statement.Exec()  // execute create table statement
     statement, _ = database.Prepare("INSERT INTO student (num,ID,NAME) VALUES (?,?,?)") // make statement for entering values afterwards


          statement1, _ := database.Prepare("CREATE TABLE IF NOT EXISTS student1 (num INTEGER PRIMARY KEY,ID INTEGER NULL, NAME STRING NULL,date_time DATETIME DEFAULT CURRENT_TIMESTAMP )")  //create table
          statement1.Exec()  // execute create table statement
          statement1, _ = database.Prepare("INSERT INTO student1 (num,ID,NAME) VALUES (?,?,?)") // make statement for entering values afterwards

          statement2, _ := database.Prepare("CREATE TABLE IF NOT EXISTS student2 (num INTEGER PRIMARY KEY,ID INTEGER NULL, NAME STRING NULL,date_time DATETIME DEFAULT CURRENT_TIMESTAMP )")  //create table
          statement2.Exec()  // execute create table statement
          statement2, _ = database.Prepare("INSERT INTO student2 (num,ID,NAME) VALUES (?,?,?)") // make statement for entering values afterwards






                        for{
                          select {
                            // reads data from sensors and put them in database
                            case data:=<-data_to_db:{
                              data1:=strings.TrimSuffix(data,"\n") // remove '\n'
                              temp:=strings.Split(data1,",") // split comma separated data into array
                              ID,_:=strconv.Atoi(temp[0])
                              NAME:=temp[1]
                              if NAME =="ard1"{
                                statement.Exec(num,ID,NAME) // put data to database
                                fmt.Println(num,ID,NAME)
                                num = num+1
                              }
                                 if NAME =="ard2"{
                                  statement1.Exec(num1,ID,NAME) // put data to database
                                  fmt.Println(num1,ID,NAME)
                                  num1 = num1+1
                                }
                                 if  NAME =="ard3"{
                                  statement2.Exec(num2,ID,NAME) // put data to database
                                  fmt.Println(num2,ID,NAME)
                                  num2 = num2+1
                                }
    }


                            case db_query := <- query_to_db :{
                              ///query is passed to database and then response is sent to passServer channel
                              rows, err:= database.Query(db_query)
                      				if err != nil{
                      					fmt.Println("ERROR ROWS",err)

                        break

                                }

              				var num string
                      				var ID string
                      				var NAME string
                      				var date_time string
                      				for rows.Next() {
                      					rows.Scan(&num,&ID, &NAME, &date_time)
                      					total:=  fmt.Sprintf(date_time + "::::" + num + ": " + ID + " ," + NAME )
                      					fmt.Println(total)
                                passServer <- total
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
         }
      }
    }
}
