package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/shrin00/moneky/internal/headers"
)

type StatusCode int

const (
	StatusOK                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	var err error
	switch statusCode {
	case StatusOK:
		_, err = w.Write(startLine(StatusOK, "OK"))
	case StatusBadRequest:
		_, err = w.Write(startLine(StatusBadRequest, "Bad Request"))
	case StatusInternalServerError:
		_, err = w.Write(startLine(StatusInternalServerError, "Internal Server Error"))
	default:
		_, err = fmt.Fprintf(w, "HTTP/1.1 %d ", statusCode)
	}

	if err != nil {
		return err
	}
	return nil
}

func startLine(code StatusCode, msg string) []byte {
	return fmt.Appendf(nil, "HTTP/1.1 %d %s\r\n", code, msg)
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	return headers.Headers{"Content-Length": strconv.Itoa(contentLen), "Connection": "close", "Content-Type": "text/plain"}
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, val := range headers {
		if _, err := fmt.Fprintf(w, "%s: %s\r\n", key, val); err != nil {
			return err
		}
	}

	if _, err := io.WriteString(w, "\r\n"); err != nil {
		return err
	}

	return nil
}
