package headers

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return Headers{}
}

var SEP = []byte("\r\n")

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	sepIdx := bytes.Index(data, SEP)
	if sepIdx == -1 {
		return 0, false, nil
	}
	data = data[:sepIdx]
	n = len(data) + len(SEP)
	if sepIdx == 0 {
		return n, true, nil
	}

	colonIdx := bytes.Index(data, []byte(":"))
	if colonIdx == -1 {
		err = fmt.Errorf("error: malformed header field line")
		return 0, false, err
	}
	key, value := strings.ToLower(string(data[:colonIdx])), strings.TrimSpace(string(data[colonIdx+1:]))

	// check if key contains space or not
	if len(key) != len(strings.TrimSpace(key)) {
		err := fmt.Errorf("error: malformed header")
		return 0, false, err
	}

	if !isValidFieldNameToken(key) {
		err = fmt.Errorf("error: malformed field name")
		return 0, false, err
	}

	if val, ok := h[key]; ok {
		h[key] = val + ", " + value
	} else {
		h[key] = value
	}
	return n, false, nil
}

func (h Headers) Get(key string) string {
	key = strings.ToLower(key)
	return h[key]
}

func isValidFieldNameToken(s string) bool {
	fnToken := regexp.MustCompile(`^[!#$%&'*+\-.^_\x60|~0-9a-zA-Z]+$`)

	return fnToken.MatchString(s)
}
