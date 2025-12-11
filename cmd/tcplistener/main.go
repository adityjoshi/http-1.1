package main

import (
	"fmt"
	"httpprotocol/internal/request"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Can't accept connection because -> ", err)
			continue
		}
		log.Println("A new client connection has been made:", conn.RemoteAddr())
		line, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatalf("error", err)
		}
		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", line.RequestLine.Method)
		fmt.Printf("- Target: %s\n", line.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", line.RequestLine.HttpVersion)

	}

}
