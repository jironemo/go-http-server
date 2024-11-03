package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "net" and "os" imports above (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
	}
	req := make([]byte, 1024)
	_, err = conn.Read(req)
	if err != nil {
		fmt.Println("Failed to read data")
	}
	requests := strings.Split(string(req), "\r\n")
	uri := strings.Split(requests[0], " ")[1]
	print(uri)
	if uri != "/" {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	} else {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	}
	conn.Close()
}
