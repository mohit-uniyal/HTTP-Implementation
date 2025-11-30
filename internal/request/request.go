package request

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
)

type RequestState string

const (
	RequestInitialized RequestState = "REQUEST_INITIALIZED"
	RequestDone        RequestState = "REQUEST_DONE"
)

type HttpMethods string

const (
	Get    HttpMethods = "GET"
	Post   HttpMethods = "POST"
	Put    HttpMethods = "PUT"
	Patch  HttpMethods = "Patch"
	Delete HttpMethods = "Delete"
)

type HttpVersion string

const (
	HttpVersion1_1 HttpVersion = "HTTP/1.1"
)

func IsValidMethod(method string) bool {
	switch method {
	case string(Get), string(Post), string(Put), string(Patch), string(Delete):
		return true
	default:
		return false
	}
}

func IsValidRoute(route string) error {

	routeRegex := regexp.MustCompile(`^\/[A-Za-z0-9\/\-\_\.\~]*$`)

	if routeRegex.MatchString(route) {
		return nil
	}

	return fmt.Errorf("invalid route")
}

func IsValidHttpVersion(httpVersion string) bool {
	switch httpVersion {
	case string(HttpVersion1_1):
		return true
	default:
		return false
	}
}

var REQUEST_SEPARATOR = "\r\n"

type Request struct {
	RequestLine RequestLine
	State       RequestState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

// returns number of bytes consumed. Returns 0 and no error if \r\n is not found.
func (r *Request) parseRequestLine(data []byte) (int, error) {

	firstIndex := bytes.Index(data, []byte(REQUEST_SEPARATOR))
	if firstIndex == -1 {
		return 0, nil
	}

	requestLine := data[:firstIndex]

	var parsedRequestLine RequestLine

	requestLineEntities := bytes.Split(requestLine, []byte(" "))
	if len(requestLineEntities) != 3 {
		fmt.Printf("invalid number of arguments in request line")
		return 0, fmt.Errorf("invalid number of arguments in request line")
	}

	//1. Validate the method
	method := string(requestLineEntities[0])

	if !IsValidMethod(method) {
		fmt.Printf("not a valid method")
		return 0, fmt.Errorf("not a valid method")
	}

	parsedRequestLine.Method = method

	//2. Validate the route
	route := string(requestLineEntities[1])

	if err := IsValidRoute(route); err != nil {
		fmt.Println(err)
		return 0, err
	}

	parsedRequestLine.RequestTarget = route

	//3. Validate the HTTP version

	httpVersion := string(requestLineEntities[2])

	if !IsValidHttpVersion(httpVersion) {
		fmt.Println("invalid HTTP version")
		return 0, fmt.Errorf("invalid HTTP version")
	}

	parsedRequestLine.HttpVersion = string(bytes.Split([]byte(httpVersion), []byte("/"))[1])

	r.RequestLine = parsedRequestLine

	return len(requestLine), nil

}

// this function accepts the next slice of bytes that needs to be parsed
func (r *Request) parse(data []byte) (int, error) {

	read := 0

	switch r.State {
	case RequestInitialized:
		numberOfBytesConsumed, err := r.parseRequestLine(data)
		if err != nil {
			return 0, err
		}

		if numberOfBytesConsumed != 0 {
			r.State = RequestDone
		}

		read = numberOfBytesConsumed
	case RequestDone:
		read = 0
	default:
		fmt.Printf("invalid request state: %s", r.State)
		return 0, fmt.Errorf("invalid request state: %s", r.State)
	}

	return read, nil

}

func RequestFromReader(reader io.Reader) (*Request, error) {
	parsedRequest := &Request{
		State: RequestInitialized,
	}

	buf := make([]byte, 1024)
	bufferLength := 0

	for parsedRequest.State != RequestDone {
		n, err := reader.Read(buf[bufferLength:])
		if err != nil {
			fmt.Printf("error reading stream: %v\n", err)
			return nil, fmt.Errorf("error reading the stream: %w", err)
		}

		bufferLength += n
		numberOfBytesParsed, err := parsedRequest.parse(buf[:bufferLength])
		if err != nil {
			fmt.Printf("error parsing the request: %v\n", err)
			return nil, fmt.Errorf("error parsing the request: %w", err)
		}

		if numberOfBytesParsed != 0 {
			copy(buf, buf[numberOfBytesParsed:bufferLength])
			bufferLength -= numberOfBytesParsed
		}

	}

	return parsedRequest, nil

}
