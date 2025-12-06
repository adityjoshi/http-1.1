package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	file, err := os.Open("message.txt")
	if err != nil {
		log.Fatal("unable to open the file")
	}
	for line := range getLineChannel(file) {
		fmt.Printf("read %s\n", line)
	}

}

func getLineChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)

	go func() {
		defer close(ch)
		defer f.Close()

		buf := make([]byte, 8)
		line := []byte{}

	}()
}
