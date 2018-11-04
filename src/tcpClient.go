package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

var indicator = 0
var pass = make(chan string)

func main() {
	fmt.Println("Start")
	for {

		if indicator == 0 {
			fmt.Println("Connecting .....")
			conn, err := net.Dial("tcp", ":8090")
			if err != nil {
			} else {
				indicator = 1
				fmt.Println("Connected")
				go res(conn)
				go Reader(conn)
				go writer(conn)
				go ticker()
			}
		}
	}

}

func res(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Client DISCONNECTED")
			indicator = 0
			pass <- "close"
			return
		} else {
			fmt.Println("text : " + text)
			//  pass <- text
		}
	}
}

func Reader(conn net.Conn) {
	osReader := bufio.NewReader(os.Stdin)
	for {
		msg, err := osReader.ReadString('\n')
		if err != nil {
			fmt.Println("Client DISCONNECTED")
			return
		} else {
			//fmt.Println("recieved from os ***" + msg)
			pass <- msg
		}
	}
}

func ticker() {
	ticker := time.NewTicker(4 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			{
				for i := 0; i < 5; i++ {
					pass <- "2,ACC,0.98\n"
					fmt.Println("2,ACC,0.98")
					time.Sleep(100 * time.Millisecond)
				}
				time.Sleep(2 * time.Second)
				fmt.Println("2,LASER,3.068")
				time.Sleep(1 * time.Millisecond)
				pass <- "2,LASER,2.068\n"
				fmt.Println("2,HUMIDITY,86")
				fmt.Println("2,TEMP,30.00")
				pass <- "2,HUMIDITY,86\n"
				pass <- "2,TEMP,30.00\n"
			}
		}
	}
}

func writer(conn net.Conn) {
	writer := bufio.NewWriter(conn)
	for {
		select {
		case values := <-pass:
			{
				if values != "close" {
					writer.WriteString(values)
					writer.Flush()
				} else {
					conn.Close()
					indicator = 0
					fmt.Println("GOTTA")
					return
				}
			}
		}
	}
}
