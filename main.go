package main

import (
	"fmt"
	"io"
	"log"
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
				log.Fatal(err)
			}
		}
	}()

	return ch
}

func main() {
	f, err := os.Open("message.txt")
	if err != nil {
		log.Fatal(err)
	}

	for line := range getLinesChannel(f) {
		fmt.Printf("read: %s\n", line)
	}
}
