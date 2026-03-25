package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	// opening a file
	f, err := os.Open("message.txt")
	if err != nil {
		log.Fatal("Unable to open the file")
	}
	defer f.Close()

	for {
		data := make([]byte, 8)
		n, err := f.Read(data)
		if err != nil {
			if err == io.EOF {
				fmt.Println("DONE!")
				break
			}
			log.Fatal("an error occured while reading the file", err.Error())
		}

		fmt.Printf("read: %s\n", data[:n])
	}

}
