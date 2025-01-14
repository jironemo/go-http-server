package main

import (
	"fmt"
	"net"
	"reflect"
	"strings"
)

var route_map = map[string]reflect.Value{}

func mapper(path string, v net.Conn, data string, headers ...interface{}) error {
	handler, exists := route_map[path]
	if !exists {
		return fmt.Errorf("no handler for this route")
	}
	params_arr := make([]reflect.Value, len(headers))
	params_arr[0] = reflect.ValueOf(conn)
	params_arr[1] = reflect.ValueOf(data)
	for i, param := range headers {
		params_arr[i+2] = reflect.ValueOf(param)
	}
	handler.Call(params_arr)

	return nil
}

func registerRoute(route string, handler interface{}) {
	route_map[route] = reflect.ValueOf(handler)
}

func test(conn net.Conn, data string, headers []string) {
	writeResponse(conn, "Echo: "+data, int(200), headers)
}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")

	if err != nil {
		fmt.Println("Failed to bind to port 4221")
	}

	for {
		conn, err := l.Accept()

		registerRoute("/test", test)

		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	req := make([]byte, 1024)

	_, err := conn.Read(req)

	if err != nil {
		fmt.Println("Failed to read data")
	}

	lines := strings.Split(string(req), "\r\n")

	first := strings.Split(lines[0], " ")
	data := lines[len(lines)-1]
	headers := lines[1 : len(lines)-2]
	//Extract method (GET,POST,etc) and path "/test","/user" from the first line
	protocol := first[0]
	path := first[1]

	switch protocol {
	case "GET":
		err = mapper(path, conn, "", headers)
	case "POST":
		err = mapper(path, conn, data, headers)
	}

	if err != nil {
		writeResponse(conn, err.Error(), int(404), nil)
	}
	conn.Close()
}

func writeResponse(conn net.Conn, data string, status int, headers []string) {

	response := "HTTP/1.1 "
	switch status {
	case 200:
		response += "200 OK"
	case 404:
		response += "404 Not Found"
	case 500:
		response += "500 Internal Server Error"
	}

	//TODO: handle headers
	print(headers)

	response += "\r\n\r\n" + data + "\r\n"
	fmt.Println(response)
	conn.Write([]byte(response))
}
