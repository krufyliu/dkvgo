package dkvgo

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
)

// Payload present action payload
type Payload map[string]string

// DkvProtocol describe the communication protocal of worker and scheduler
type DkvProtocol struct {
	Action  string  `json:"action"`
	Payload Payload `json:"payload"`
}

// NewDkvProtocol make a new DkvProtocal
func NewDkvProtocol() *DkvProtocol {
	return &DkvProtocol{Payload: make(Payload)}
}

// AddPayload add (key, value) pair to payload field
func (p *DkvProtocol) AddPayload(key, value string) {
	p.Payload[key] = value
}

// Encode encode DkvProtocal to length and json based binary
func (p DkvProtocol) Encode() ([]byte, error) {
	content, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	length := int32(len(content))
	buffer := bytes.NewBuffer(make([]byte, 0, int(length)+binary.Size(length)))
	err = binary.Write(buffer, binary.BigEndian, length)
	if err != nil {
		return nil, err
	}
	_, err = buffer.Write(content)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// DecodeFrom decode binary which red from io.Reader to current DkvProtocol object
func (p *DkvProtocol) DecodeFrom(r io.Reader) error {
	var length int32
	err := binary.Read(r, binary.BigEndian, &length)
	if err != nil {
		return err
	}
	buffer := make([]byte, length)
	_, err = r.Read(buffer)
	if err != nil {
		return err
	}
	err = json.Unmarshal(buffer, p)
	if err != nil {
		return err
	}
	return nil
}

// DecodeFrom decode binary which read from io.Reader to DkvProtocol object
func DecodeFrom(r io.Reader) (*DkvProtocol, error) {
	p := NewDkvProtocol()
	err := p.DecodeFrom(r)
	if err != nil {
		return nil, err
	}
	return p, nil
}
