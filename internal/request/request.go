package request

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"
)

type StateParse string

type RequestLine struct {
	Method        string
	RequestTarget string
	HttpVersion   string
}

type Request struct {
	RequestLine RequestLine
	State       StateParse
}

const (
	StateInit StateParse = "initialized"
	StateDone StateParse = "done"
)

var ErrorMalformedRequestLine = fmt.Errorf("malformed request line")
var ErrorUnsupportedHttpVersion = fmt.Errorf("unsupported http version")
var SEPARATOR = []byte("\r\n")

func newRequest() *Request {
	return &Request{
		State: StateInit,
	}
}

func parseRequestLine(message []byte) (*RequestLine, int, error) {

	// check if the message contains SEPARATOR
	sepIndex := bytes.Index(message, SEPARATOR)
	if sepIndex == -1 {
		return nil, 0, nil
		// if we don't find the SEPARATOR in the sent message,
		// then it means we don't have the complete request line and can't parse the request-line
	}

	// if SEPARATOR exists, then slice the byte till the sepIndex and we get the request-line
	rl := string(message[:sepIndex])
	numOfBytes := len(rl) + len(SEPARATOR)
	fmt.Println("request line: ", rl)

	parts := strings.Split(rl, " ")
	// strings.Fields() - a good way to split strings around white space
	if len(parts) != 3 {
		err := fmt.Errorf("error: malformed request line")
		log.Println(err)
		return nil, 0, err
	}
	method := parts[0]
	reqTarget := parts[1]
	httpVersion := strings.Split(parts[2], "/")[1]

	if strings.ToUpper(method) != method {
		log.Println("error: ", ErrorMalformedRequestLine)
		return nil, 0, fmt.Errorf("to upper")
	}

	if httpVersion != "1.1" {
		log.Println("error: ", ErrorUnsupportedHttpVersion)
		return nil, 0, ErrorUnsupportedHttpVersion
	}

	return &RequestLine{
		Method:        method,
		RequestTarget: reqTarget,
		HttpVersion:   httpVersion,
	}, numOfBytes, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	rq := newRequest()

	// create a slice of byte to read the data into memory, which we call buf
	buf := make([]byte, 1024)
	bufLen := 0

	for !rq.done() {
		n, err := reader.Read(buf[bufLen:])
		if err != nil {
			log.Println("error: ", err.Error())
			return nil, err
		}

		bufLen += n
		readN, err := rq.parse(buf[:bufLen])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[readN:bufLen])
		bufLen -= readN
	}

	// http_message := strings.Split(string(msg), "\r\n")
	// // fmt.Println(http_message)
	// request_line, n, err := parseRequestLine(buf)
	// if err != nil {
	// 	return nil, err
	// }

	// return &Request{
	// 	RequestLine: *request_line,
	// }, nil

	return rq, nil
}

func (r *Request) done() bool {
	return r.State == StateDone
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0

outer:
	for {
		switch r.State {
		case StateInit:
			rl, n, err := parseRequestLine(data[read:])
			if err != nil {
				return 0, err
			}
			if n == 0 {
				break outer
			}

			r.RequestLine = *rl
			read += n

			r.State = StateDone

		case StateDone:
			break outer
		}
	}
	return read, nil
}
