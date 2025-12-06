package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)

	go func() {
		defer close(ch)
		defer f.Close()

		buf := make([]byte, 8)
		line := []byte{}

		for {
			n, err := f.Read(buf)
			if n > 0 {
				for _, b := range buf[:n] {
					if b == '\n' {
						ch <- string(line)
						line = []byte{}
					} else {
						line = append(line, b)
					}
				}
			}

			if err == io.EOF {
				if len(line) > 0 {
					ch <- string(line)
				}
				return
			}
			if err != nil {
				fmt.Fprintln(os.Stderr, "read error:", err)
				return
			}
		}
	}()

	return ch
}

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
		for line := range getLinesChannel(conn) {
			fmt.Println(line)
		}
		fmt.Println("Connection has been closed")
	}

}
