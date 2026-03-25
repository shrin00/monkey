package main

import (
	"bytes"
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

	// reading the data
	str := ""
	for {

		// create a slice of byte with size of 8 and read into it
		data := make([]byte, 8)
		n, err := f.Read(data)
		if err != nil {
			if err == io.EOF {
				fmt.Println("DONE!")
				break
			}
			log.Fatal("an error occured while reading the file", err.Error())
		}

		data = data[:n]
		if i := bytes.IndexByte(data, '\n'); i != -1 { // reading 8 bytes at a time, if those 8 bytes contain '\n' new line character
			str += string(data[:i])	// in them, that is the end of the line, bytes upto the index of '\n' character is appeneded to the str
			data = data[i+1: ]     // variable and rest of the bytes are assigned back to data, which ara the bytes from the next line
			fmt.Printf("read: %s\n", str)
			str = ""
		}

		str += string(data)  // continue to append the left over data, if the if conditioned is executed.
	}

}
