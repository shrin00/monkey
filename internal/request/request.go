// this package implements http request message and provides parsing functions
package request

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"
)

type StateParse string

const (
	StateInit StateParse = "init"
	StateDone StateParse = "done"
)

// define a request line
type RequestLine struct {
	Method        string
	RequestTarget string
	HttpVersion   string
}

// define header item
type headerItem struct {
	key   string
	value string
}

// define header
type Header []headerItem

// define a request
type Request struct {
	RequestLine *RequestLine
	Header      Header
	Body        []byte
	State       StateParse // indicates the state of the request parsing
}

// create a new Request
func newRequest() *Request {
	return &Request{
		State: StateInit,
	}
}

// Constant variables
var SEPARATOR = []byte("\r\n")
var HTTP_VERSION = "1.1"
var ErrorMalformedRequestLine = fmt.Errorf("malformed request line")
var ErrorUnsupportedHttpVersion = fmt.Errorf("unsupported http version, only http 1.1 is supported")

// we need to create a function to which a reader can be passed,
// which consist of request data in the bytes format,
// function will return the Request read and the errors if any
// our goal is to Request variable, from the input bytes.
// which is parsing the raw bytes of data into the structred Request.

// creating a Request structure from the raw data of the io.Reader
func RequestFromReader(reader io.Reader) (*Request, error) {
	rq := newRequest()

	buf := make([]byte, 1024) // 1024 size of the buffer, which means maximum of 1024
	// bytes of data can be read at once, if the request message size is greater than 1k, then we need to increase the size
	bufLen := 0 // length of the buffer, number of bytes in the buf, we will set it to 0, since len(buf) is 1024
	// var requestLine *RequestLine
	for !rq.done() {
		fmt.Println("Length of the buf just after initialization", bufLen)
		// readN - how much bytes from raw request is read into buf
		readN, err := reader.Read(buf[bufLen:]) // io.Reader has a Read function, which should read the data into the slice of type []byte
		if err != nil {
			err := fmt.Errorf("failed to read the request message: %v", err.Error())
			log.Println(err)
			return nil, err
		}

		// now that we have n number of bytes read into the buf slice, what should I do with it?
		// now, we have the raw bytes in slice of buf, we want them to be in Request format
		// so, we have to parse each section, Request-line, Header, body
		// 1. Request-line: see #RequestLineParser
		bufLen += readN
		// n how much is consumed by the parse
		n, err := rq.parse(buf[:bufLen])
		if err != nil {
			log.Println(err)
			return nil, err
		}

		copy(buf, buf[n:bufLen]) // copying left over bytes into buf, which bring the rest for the bytes to front
		bufLen -= n              // reseting the bufLen to the, unparsed items of bytes
	}

	return rq, nil
}

// This fucntion will return a RequestLine, by parsing the p into RequestLine
// this function will expect the request-line in []byte format,
// it will extract the method, request target, http version into RequestLine and return
// it will also validate the data assigned, if the value doesn't in the standard principle,
// we will reaturn with error and 0 number bytes read
func parseRequestLine(p []byte) (*RequestLine, int, error) {
	parsedRl := &RequestLine{}
	log.Println("p in parseRequestLine: ", string(p))

	// a request-line ends with \r\n (CRLF), we will check the index of the SEPARTOR
	// existence of the SEPAROTOR, indicates complete request-line bytes
	// if it isn't exists then, data is not complete
	idxSep := bytes.Index(p, SEPARATOR)
	if idxSep == -1 {
		return nil, 0, nil
	}

	// request-line is consist of the 3 parts, method, request-target and http version(see RequestLine struct)
	// at this point we have complete request-line, we extract rl bytes from the []bytes
	// which is bytes untill the index of SEPARATOR not included
	rl := string(p[:idxSep])
	n := len(rl) + len(SEPARATOR)

	// divide the rl into slice at " "(space character) and verify it consist of the 3 items or not
	rl_parts := strings.Split(rl, " ")
	if len(rl_parts) != 3 {
		err := fmt.Errorf("%w: parts of request line, expected 3, got %d", ErrorMalformedRequestLine, len(rl_parts))
		log.Println(err.Error())
		return nil, 0, err
	}
	parsedRl.Method = rl_parts[0]
	parsedRl.RequestTarget = rl_parts[1]
	parsedRl.HttpVersion = strings.Split(rl_parts[2], "/")[1]

	// verify if the METHOD is in capital letter or not
	if strings.ToUpper(parsedRl.Method) != parsedRl.Method {
		log.Println(ErrorMalformedRequestLine.Error() + fmt.Sprintf(": method: %v", parsedRl.Method))
		return nil, 0, ErrorMalformedRequestLine
	}

	// verfiy the version
	if parsedRl.HttpVersion != HTTP_VERSION {
		log.Println(ErrorUnsupportedHttpVersion)
		return nil, 0, ErrorUnsupportedHttpVersion
	}

	return parsedRl, n, nil
}

func (r *Request) parse(p []byte) (int, error) {
	read := 0
outer:
	for {
		switch r.State {
		case StateInit:
			rl, n, err := parseRequestLine(p[read:])
			if err != nil {
				return 0, err
			}

			if n == 0 {
				break outer
			}

			r.RequestLine = rl
			read += n
			r.State = StateDone
		case StateDone:
			break outer
		}
	}

	return read, nil
}

// done checks if the state of the request is done or not
func (r *Request) done() bool {
	return r.State == StateDone
}
