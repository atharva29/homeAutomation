package main

import (
  "fmt"
  "net"
  "bufio"
  "os"
)

func main(){

  conn,_:=net.Dial("tcp",":2000")
  //writer:=bufio.NewWriter(conn)
  readerConn :=bufio.NewReader(conn)
    //    reader := bufio.NewReader(os.Stdin)
        text,_:=reader.ReadString('\n')
      //  writer.WriteString(text)
fmt.Println("msg",text)

}
