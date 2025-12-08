package request

import (
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

var Error_Bad_Request_Line = fmt.Errorf("malformed http")

func parseRequestLine(b string) (*RequestLine, string, error) {
	idx := strings.Index(b, "\r\n")
	if idx == -1 {
		return nil, b, nil
	}
	start := b[:idx]
	restOfMsg := b[idx+len("\r\n"):]

	parts := strings.Split(start, " ")
	if len(parts) != 3 {
		return nil, restOfMsg, Error_Bad_Request_Line
	}

	return &RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   parts[2],
	}, restOfMsg, nil

}

func RequestFromReader(reader io.Reader) (*Request, error) {
	line, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	str := string(line)

	rl, _, err := parseRequestLine(str)
}
