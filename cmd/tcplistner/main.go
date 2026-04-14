package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/shrin00/moneky/internal/request"
)

// func getLinesChannel(f io.ReadCloser) <-chan string {
// 	chanLine := make(chan string)

// 	go func() {
// 		defer close(chanLine)
// 		defer f.Close()

// 		// reading the data
// 		str := ""
// 		for {

// 			// create a slice of byte with size of 8 and read into it
// 			data := make([]byte, 8)
// 			n, err := f.Read(data)
// 			if err != nil {
// 				if err == io.EOF {
// 					break
// 				}
// 				log.Fatal("an error occured while reading the file", err.Error())
// 			}

// 			data = data[:n]
// 			if i := bytes.IndexByte(data, '\n'); i != -1 { // reading 8 bytes at a time, if those 8 bytes contain '\n' new line character
// 				str += string(data[:i]) // in them, that is the end of the line, bytes upto the index of '\n' character is appeneded to the str
// 				data = data[i+1:]       // variable and rest of the bytes are assigned back to data, which ara the bytes from the next line
// 				chanLine <- str
// 				time.Sleep(3 * time.Second)
// 				str = ""
// 			}

// 			str += string(data) // continue to append the left over data, if the if conditioned is executed.
// 		}
// 		if len(str) != 0 {
// 			chanLine <- str
// 		}
// 	}()

// 	return chanLine
// }

func printRequest(rq *request.Request) {
	fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n", rq.RequestLine.Method, rq.RequestLine.RequestTarget, rq.RequestLine.HttpVersion)
	fmt.Println("Headers: ")
	for key, val := range rq.Headers {
		fmt.Printf("- %s: %s\n", key, val)
	}
	fmt.Printf("Body:\n%s\n", rq.Body)
}

func main() {
	// opening a file
	// f, err := os.Open("message.txt")
	// if err != nil {
	// 	log.Fatal("Unable to open the file")
	// }
	// defer f.Close()

	// open a tcp connection and listen on port 42069

	conn, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal("error while opening a tcp connection", err.Error())
		os.Exit(-1)
	}
	log.Println("listening on address: ", conn.Addr().String())

	// udp sender to listner on the localhost:42070
	// udpAddr, _ := net.ResolveUDPAddr("udp", "localhost:42070")
	// udpConn, _ := net.DialUDP("udp", nil, udpAddr)
	// defer udpConn.Close()

	for {
		// accept the inbound data
		f, err := conn.Accept()
		if err != nil {
			log.Fatal("failed to get next connection")
		}

		// lines := getLinesChannel(f)
		// for l := range lines {
		// 	fmt.Printf("read: %s\n", l)
		// 	// udpConn.Write([]byte(l + "\n"))
		// }

		rq, err := request.RequestFromReader(f)
		if err != nil {
			err = fmt.Errorf("error: %w", err)
			log.Println(err)
			continue
		}

		printRequest(rq)
	}

}
