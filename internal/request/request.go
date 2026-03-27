package request

import (
	"fmt"
	"io"
	"log"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	Method        string
	RequestTarget string
	HttpVersion   string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	// declare a new request
	msg, err := io.ReadAll(reader)
	if err != nil {
		log.Println("error: ", err.Error())
		return nil, err
	}

	http_message := strings.Split(string(msg), "\r\n")
	fmt.Println(http_message)
	request_line, err := parseRequestLine(http_message[0])
	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: *request_line,
	}, nil
}

func parseRequestLine(requestLine string) (*RequestLine, error) {
	reqLineSections := strings.Split(requestLine, " ")
	// strings.Fields() - a good way to split strings around white space
	if len(reqLineSections) != 3 {
		err := fmt.Errorf("error: malformed request line")
		log.Println(err)
		return nil, err
	}
	method := reqLineSections[0]
	reqTarget := reqLineSections[1]
	httpVersion := strings.Split(reqLineSections[2], "/")[1]

	if strings.ToUpper(method) != method {
		err := fmt.Errorf("error: malformed request method")
		log.Println(err)
		return nil, err
	}

	if httpVersion != "1.1" {
		err := fmt.Errorf("error: unsupported http version")
		return nil, err
	}

	return &RequestLine{
		Method:        method,
		RequestTarget: reqTarget,
		HttpVersion:   httpVersion,
	}, nil
}
