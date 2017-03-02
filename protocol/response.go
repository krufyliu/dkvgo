package protocol

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"strings"
)

// DkvResponse define a response structure
type DkvResponse struct {
	Version string
	Status  string
	Payload []byte
}

// NewDkvResponse create a simple resonse
func NewDkvResponse(status string) *DkvResponse {
	return &DkvResponse{
		Version: "V1.0",
		Status:  status,
	}
}

// NewDkvResponseWithPayload create a reponse with payload
func NewDkvResponseWithPayload(status string, payload []byte) *DkvResponse {
	return &DkvResponse{
		Version: "V1.0",
		Status:  status,
		Payload: payload,
	}
}

// Dumps serialize response to binary
func (res *DkvResponse) Dumps() ([]byte, error) {
	headerLine := res.Version + " " + res.Status
	if len(res.Payload) == 0 {
		return []byte(headerLine), nil
	}
	var payloadLen = len(res.Payload)
	var buffer = bytes.NewBuffer(make([]byte, 0, len(headerLine)+payloadLen+1))
	if _, err := buffer.Write([]byte(headerLine)); err != nil {
		return nil, err
	}
	if err := buffer.WriteByte('\n'); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, binary.LittleEndian, int32(payloadLen)); err != nil {
		return nil, err
	}
	if _, err := buffer.Write(res.Payload); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// PullHeader parse response header
func (res *DkvResponse) PullHeader(r *bufio.Reader) error {
	line, _, err := r.ReadLine()
	if err != nil {
		return nil
	}
	arr := strings.Split(string(line), " ")
	res.Version = arr[0]
	res.Status = arr[1]
	return nil
}

// PullContent parse response payload
func (res *DkvResponse) PullContent(r *bufio.Reader) error {
	var contentLen int32
	if err := binary.Read(r, binary.LittleEndian, &contentLen); err != nil {
		return err
	}
	buffer := make([]byte, int(contentLen))
	for n := 0; n < int(contentLen); {
		rn, err := r.Read(buffer[n:])
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		n += rn
	}
	res.Payload = buffer
	return nil
}
