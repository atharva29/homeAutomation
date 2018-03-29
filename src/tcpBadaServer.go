package main

import (
  "fmt"     //libaray use for print function
  "net"     // used for net
  "bufio"   // used for taking inputs from user or port
  "os"      // used to read user inputs
)


var pass  = make(chan string)   // channel for Cient data

func main(){
fmt.Println("Start")
l,err := net.Listen("tcp",":6600")  //set a server listening at port 6600
if err!=nil{
  printErr(err)
  }
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
                      conn.Close()      // closes the connection
                      return
                }
              }
           }
        }
      }
