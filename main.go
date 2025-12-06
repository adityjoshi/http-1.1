package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	filePath := "message.txt"

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to open the file")
	}
	defer file.Close()

	scan := bufio.NewScanner(file)
	for scan.Scan() {
		line := scan.Text()
		fmt.Println("read", line)
	}
}
