package request

import (
	"fmt"
	"io"
	"strings"
)

type ParseState string

const (
	Stateinit ParseState = "init"
	StateDone ParseState = "done"
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
var Error_Reading_From_Done_State = fmt.Errorf("error reading from the done state")

func parseRequestLine(b string) (*RequestLine, int, string, error) {

	idx := strings.Index(b, "\r\n")
	if idx == -1 {

		return nil, 0, b, nil
	}

	line := b[:idx]
	rest := b[idx+2:]

	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return nil, 0, rest, Error_Bad_Request_Line
	}

	httpParts := strings.Split(parts[2], "/")
	if len(httpParts) != 2 || httpParts[0] != "HTTP" || httpParts[1] != "1.1" {
		return nil, 0, rest, Error_Bad_Request_Line
	}

	rl := &RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   httpParts[1],
	}

	consumed := idx + 2
	return rl, consumed, rest, nil
}
func newRequest() *Request {
	return &Request{
		State: Stateinit,
	}
}

func (r *Request) Parse(data []byte) (int, error) {

	if r.State == StateDone {
		return 0, Error_Reading_From_Done_State
	}
	if r.State != Stateinit {
		return 0, fmt.Errorf("unknown State")
	}

	rl, noOfBytes, _, err := parseRequestLine(string(data))
	if err != nil {
		return 0, err
	}
	if noOfBytes == 0 {
		return 0, nil
	}
	r.RequestLine = *rl
	r.State = StateDone
	return noOfBytes, nil
}

func (r *Request) done() bool {
	return r.State == StateDone
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
