package protocol

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"strings"
)

// DkvRequst describe request protocol
type DkvRequst struct {
	Version string
	Method  string
	Args    []string
	Payload []byte
}

// NewDkvRequest create a DkvRequst
func NewDkvRequest(method string, args []string) *DkvRequst {
	return &DkvRequst{
		Version: "V1.0",
		Method:  method,
		Args:    args,
	}
}

// NewDkvRequestWithPayload create a DkvRequst with payload
func NewDkvRequestWithPayload(method string, args []string, payload []byte) *DkvRequst {
	return &DkvRequst{
		Version: "V1.0",
		Method:  method,
		Args:    args,
		Payload: payload,
	}
}

// Dumps serialize DkvRequst to binary
func (req *DkvRequst) Dumps() ([]byte, error) {
	var arrBuff = make([]string, 0, len(req.Args)+2)
	arrBuff = append(arrBuff, req.Version, req.Method)
	arrBuff = append(arrBuff, req.Args...)
	headerLine := strings.Join(arrBuff, " ")
	if len(req.Payload) == 0 {
		return []byte(headerLine), nil
	}
	var payloadLen = len(req.Payload)
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
	if _, err := buffer.Write(req.Payload); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// PullHeader parse DkvRequst header
func (req *DkvRequst) PullHeader(r *bufio.Reader) error {
	line, _, err := r.ReadLine()
	if err != nil {
		return nil
	}
	arr := strings.Split(string(line), " ")
	req.Version = arr[0]
	req.Method = arr[1]
	req.Args = arr[2:]
	return nil
}

// PullContent parse request content
func (req *DkvRequst) PullContent(r *bufio.Reader) error {
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
	req.Payload = buffer
	return nil
}
