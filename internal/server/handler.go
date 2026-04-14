package server

import (
	"io"

	"github.com/shrin00/moneky/internal/request"
	"github.com/shrin00/moneky/internal/response"
)

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w io.Writer, r *request.Request) *HandlerError

func WriteHandlerError(w io.Writer, hErr *HandlerError) error {
	if err := response.WriteStatusLine(w, hErr.StatusCode); err != nil {
		return err
	}
	if err := response.WriteHeaders(w, response.GetDefaultHeaders(len([]byte(hErr.Message)))); err != nil {
		return err
	}
	if _, err := w.Write([]byte(hErr.Message)); err != nil {
		return err
	}
	return nil
}
