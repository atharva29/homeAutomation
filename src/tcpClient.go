package main
import (
"fmt"
"net"
"bufio"
"os"   )

var pass = make(chan string)
func main(){

for {
  conn,err:=net.Dial("tcp",":2000")
  if err!= nil{
    fmt.Println("ERR")
    return
  } else {
    go res(conn)
    go Reader(conn)
    go writer(conn)
    }
  }
}


func res(conn net.Conn){
    reader := bufio.NewReader(conn)
    for {
        text,err:=reader.ReadString('\n')
        if err != nil {
          fmt.Println("Client DISCONNECTED")
          pass <- "close"
          return
          } else {
          fmt.Println("text : " + text)
          pass <- text
                  }
          }
  }


  func Reader(conn net.Conn){
    osReader := bufio.NewReader(os.Stdin)
    for {
    msg,err:=osReader.ReadString('\n')
    if err != nil {
      fmt.Println("Client DISCONNECTED")
      return
       }else {
      fmt.Println("recieved from os ***" + msg)
      pass <- msg
        }
      }
  }

  func writer(conn net.Conn){
    writer:= bufio.NewWriter(conn)
    for {
    select{
      case values := <- pass :{
        if (values != "close"){
        writer.WriteString(values)
        writer.Flush()
              } else {
                conn.Close()
                fmt.Println("GOTTA")
                return
              }
            }
          }
        }
      }
