package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	remoteAddr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatal("error: ", err.Error())
	}

	udpConn, err := net.DialUDP("udp", nil, remoteAddr) // network, local address, remote addr(to where upd packet are destined)
	if err != nil {
		log.Fatal("error: ", err.Error())
	}
	defer udpConn.Close()

	read_buffer := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(" >")
		str, err := read_buffer.ReadBytes('\n')
		if err != nil {
			log.Println("error: ", err.Error())
			continue
		}

		udpConn.Write([]byte(str))
	}

}
