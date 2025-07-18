package request

import (
    "bytes"
	"errors"
    "fmt"
	"io"
	"strings"
)

const (
    crlf = "\r\n"
    bufferSize = 8
)

const (
    requestInitialized requestState = iota
    requestDone
)

type Request struct {
    RequestLine RequestLine
    state requestState
}

type RequestLine struct {
    HttpVersion string
    RequestTarget string
    Method string
}

type requestState int

func RequestFromReader(reader io.Reader) (*Request, error) {
    buf := make([]byte, bufferSize, bufferSize)
    readToIndex := 0
    request := &Request{
        state: requestInitialized,
    }

    for request.state != requestDone {
        if readToIndex >= len(buf) {
            newBuf := make([]byte, len(buf)*2)
            copy(newBuf, buf)
            buf = newBuf
        }

        bytesRead, err := reader.Read(buf[readToIndex:])
        if err != nil {
            if errors.Is(err, io.EOF) {
                request.state = requestDone
                break
            }
            return nil, err
        }
        readToIndex += bytesRead
        bytesParsed, err := request.parse(buf[:readToIndex])
        if err != nil {
            return nil, err
        }

        copy(buf, buf[bytesParsed:])
        readToIndex -= bytesParsed
    }

    return request, nil
}

func (r *Request) parse(data []byte) (int, error) {
    switch r.state {
    case requestInitialized:
        requestLine, numBytes, err := parseRequestLine(data)
        if err != nil {
            return 0, err
        }

        if numBytes == 0 {
            return 0, nil
        }
        r.RequestLine = *requestLine
        r.state = requestDone
        return numBytes, nil
    case requestDone:
        return 0, fmt.Errorf("error: trying to read data in a done state")
    default:
        return 0, fmt.Errorf("error: unknown state")
    }
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
    idx := bytes.Index(data, []byte(crlf))
    if idx == -1 {
        return nil, 0, nil
    }

    requestLineText := string(data[:idx])
    requestLine, err := requestLineFromString(requestLineText)
    if err != nil {
        return nil, 0, err
    }
    return requestLine, idx+2, nil
}

func requestLineFromString(str string) (*RequestLine, error) {
    parts := strings.Split(str, " ")
    if len(parts) != 3 {
        return nil, fmt.Errorf("Poorly formatted request-line: %s", str)
    }

    method := parts[0]
    switch method {
    case "GET", "POST", "PUT", "DELETE", "UPDATE":
    default:
        return nil, fmt.Errorf("invalid method: %s", method)
    }

    requestTarget := parts[1]

    versionParts := strings.Split(parts[2], "/")
    if len(versionParts) != 2 {
        return nil, fmt.Errorf("malformed start-line: %s", str)
    }

    httpPart := versionParts[0]
    if httpPart != "HTTP" {
        return nil, fmt.Errorf("Unrecognized HTTP-version: %s", httpPart)
    }

    version := versionParts[1]
    if version != "1.1" {
        return nil, fmt.Errorf("Unrecognized HTTP-version: %s", version)
    }

    return &RequestLine {
        HttpVersion: version,
        RequestTarget: requestTarget,
        Method: method,
    }, nil
}
