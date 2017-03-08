package protocol

import (
	"bytes"
	"encoding/binary"
	"io"
)

// Package define binary protocol
type Package struct {
	Version   byte
	Directive int32
	WithPack  byte
	Payload   []byte
}

func NewPackage(directive int32) *Package {
	return &Package{
		Version:   0x01,
		Directive: directive,
		WithPack:  0,
	}
}

func NewPackageWithPayload(directive int32, payload []byte) *Package {
	return &Package{
		Version:   0x01,
		Directive: directive,
		WithPack:  1,
		Payload:   payload,
	}
}

func (p *Package) Marshal() ([]byte, error) {
	buffer := bytes.NewBuffer(make([]byte, 0, 10))
	if err := buffer.WriteByte(p.Version); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, binary.LittleEndian, p.Directive); err != nil {
		return nil, err
	}
	if err := buffer.WriteByte(p.WithPack); err != nil {
		return nil, err
	}
	if p.WithPack == 0 {
		return buffer.Bytes(), nil
	}
	var length = len(p.Payload)
	if err := binary.Write(buffer, binary.LittleEndian, int32(length)); err != nil {
		return nil, err
	}
	if _, err := buffer.Write(p.Payload); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (p *Package) Unmarshal(r io.Reader) error {
	buffer := make([]byte, 6)
	if _, err := io.ReadFull(r, buffer); err != nil {
		return err
	}
	var directive = int32(binary.LittleEndian.Uint32(buffer[1:5]))
	p.Version = buffer[0]
	p.Directive = directive
	p.WithPack = buffer[5]
	if p.WithPack == 0 {
		return nil
	}
	if _, err := io.ReadFull(r, buffer[0:4]); err != nil {
		return err
	}
	var length = int32(binary.LittleEndian.Uint32(buffer[0:4]))
	buffer = make([]byte, length)
	if _, err := io.ReadFull(r, buffer); err != nil {
		return err
	}
	p.Payload = buffer
	return nil
}
