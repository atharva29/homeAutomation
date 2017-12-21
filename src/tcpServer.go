package main
// this file will connect to tcp client and will act like a tcp server
import (
  "fmt"
  "net"
  "bufio"
  "os"
)

var pass  = make(chan string)

func main(){
fmt.Println("Start")
  l,err := net.Listen("tcp",":2000")
  printErr(err)
    for {
        conn,err := l.Accept()
          if err != nil{
            fmt.Println(err)
            }
            go Reader(conn)
            go res(conn)
            go writer(conn)
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
              if (text == "echo\n"){
                pass <-text
              }
            fmt.Println("text : " + text)
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

func printErr (err error){
  if err != nil{
  fmt.Println(err)
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
        }  else if (values == "close"){
                  fmt.Println("GOTTA")
                  conn.Close()
                  return
            }
          }
       }
     }
   }
