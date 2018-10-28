package main

import (
	"bufio" // used for taking inputs from user or port
	"fmt"   //libaray use for print function
	"log"
	"net" // used for net
	"os"
	// used to read user inputs
)

func main() {
	fmt.Println("Start")
	l, err := net.Listen("tcp", ":8087") //set a server listening at port 20000
	if err != nil {
		fmt.Println(err)
	}
	for {
		conn, err := l.Accept() // accept incoming client connection
		if err != nil {
			fmt.Println(err)
		}
		go res(conn) //function for reading data from client
	}
}

//function :- This function reads data from CLIENT and display it
func res(conn net.Conn) {
	reader := bufio.NewReader(conn) // set a new reader for tcp
	// for loop is used for continuos reading from port
	file, err := os.Create("test.txt") // Truncates if file already exists, be careful!
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	defer file.Close() // Make sure to close the file when you're done

	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Client DISCONNECTED")
			return
		} else {
			fmt.Println(text) // print the recieved data
			_, err = file.WriteString(text)

			if err != nil {
				log.Fatalf("failed writing to file: %s", err)
			}
		}
	}
}
