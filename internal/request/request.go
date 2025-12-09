package request

import (
	"fmt"
	"io"
	"strings"
)

type ParseState string

const (
	Stateinit ParseState = "init"
	StatDone  ParseState = "done"
)

type Request struct {
	RequestLine RequestLine
	State       ParseState
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

	httpParts := strings.Split(parts[2], "/")

	if len(httpParts) != 2 || httpParts[0] != "HTTP" || httpParts[1] != "1.1" {
		return nil, restOfMsg, Error_Bad_Request_Line
	}

	rl := &RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   httpParts[1],
	}

	return rl, restOfMsg, nil
}

func newRequest() *Request {
	return &Request{
		State: Stateinit,
	}
}

func (r *Request) parse(data []byte) (int, error) {

}

func (r *Request) done() bool {
	return r.State == StatDone
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()

	buf := make([]byte, 1024)
	bufLen := 0
	for !request.done() {
		n, err := reader.Read(buf[bufLen:])
		if err != nil {
			return nil, err
		}
		bufLen += n
		readN, err := request.Parse(buf[:bufLen])
		if err != nil {
			return nil, err
		}
		copy(buf, buf[readN:bufLen])
		bufLen -= readN

	}
	return request, nil
}
