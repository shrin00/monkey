package server

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/shrin00/moneky/internal/request"
	"github.com/shrin00/moneky/internal/response"
)

type Server struct {
	listner     net.Listener
	closed      atomic.Bool
	handlerFunc Handler
}

func Serve(port int, handlerFunc Handler) (*Server, error) {
	listner, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("failed to create a listner: %v", err)
	}

	newServer := &Server{
		listner:     listner,
		handlerFunc: handlerFunc,
	}

	go newServer.listen()

	return newServer, nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listner.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("failed to get or accept next connection: %v", err)
			continue
		}

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	req, err := request.RequestFromReader(conn)
	if err != nil {
		log.Println("error while parsing the request: ", err.Error())
		return
	}

	var b bytes.Buffer
	handleErr := s.handlerFunc(&b, req)
	if handleErr != nil {
		err := WriteHandlerError(conn, handleErr)
		if err != nil {
			log.Println("error while writing error: ", err)
			return
		}
		return // Don't write success response if there was an error
	}

	// Write success response
	body := b.Bytes()
	headers := response.GetDefaultHeaders(len(body))
	if err := response.WriteStatusLine(conn, response.StatusOK); err != nil {
		log.Println("error while writing the response line: ", err.Error())
		return
	}
	if err := response.WriteHeaders(conn, headers); err != nil {
		log.Println("error while writing headers: ", err.Error())
		return
	}
	if _, err := conn.Write(body); err != nil {
		log.Println("error while writing response body: ", err)
		return
	}
}

func (s *Server) Close() error {
	s.closed.Store(true)

	return s.listner.Close()
}
