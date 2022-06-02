package commit

import (
	"bytes"
	"fmt"
)

const gitCommitObjectType = "commit"
const newlineByte = 0x0a // \n

type Header struct {
	Value string
}

type Object struct {
	Type    string
	Message []byte
	headers []*Header
}

func (g *Object) Headers() []*Header {
	return g.headers
}

func (g *Object) SetHeaders(headers []*Header) {
	g.headers = headers
}

//ParseGitCommitObject parses a git commit object without its prefix (i.e. 'commit <len>\0')
func ParseGitCommitObject(objectPayload []byte) (*Object, error) {

	buf := bytes.NewBuffer(objectPayload)

	headers, err := parseHeaders(buf)
	if err != nil {
		return nil, fmt.Errorf("invalid commit object prefix: %w", err)
	}

	msg, err := parseCommitMessage(buf)
	if err != nil {
		return nil, fmt.Errorf("invalid commit message: %w", err)
	}

	obj := &Object{
		Type:    gitCommitObjectType,
		Message: msg,
		headers: headers,
	}

	return obj, nil
}

func parseHeaders(buf *bytes.Buffer) ([]*Header, error) {

	var headers []*Header
	for {
		header, err := parseNextHeader(buf)
		if err != nil {
			return nil, fmt.Errorf("failed reading headers: %w", err)
		}
		if header == nil {
			break
		}
		headers = append(headers, header)
	}
	return headers, nil
}

func parseNextHeader(buf *bytes.Buffer) (*Header, error) {
	headerBytes, err := buf.ReadBytes(newlineByte)
	if err != nil {
		return nil, fmt.Errorf("cannot parse commit header: %w", err)
	}
	headerBytes = headerBytes[:len(headerBytes)-1]
	if len(headerBytes) == 0 {
		// we have no more headers, let's stop here
		return nil, nil
	}

	return &Header{Value: string(headerBytes)}, nil
}

func parseCommitMessage(buf *bytes.Buffer) ([]byte, error) {
	msg := buf.Bytes()
	if len(msg) == 0 {
		return nil, fmt.Errorf("empty commit message")
	}
	return msg, nil
}
