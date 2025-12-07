package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {

	addr, err := net.ResolveUDPAddr("UDP", "localhost:42069")
	if err != nil {
		log.Fatal("error is ->", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal("Error", conn)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println(">")
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("read error: %v", err)
			continue
		}
		_, err = conn.Write([]byte(line))
		if err != nil {
			log.Printf("write error: %v", err)
		}
	}

}
