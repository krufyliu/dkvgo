package protocol

import (
	"bytes"
	"testing"
)

func Test_PackageEncodeAndDecode(t *testing.T) {
	p := NewPackage(0x11)
	message, err := p.Marshal()
	if err != nil {
		t.Error(err)
	}
	t.Logf("%x", message)
	if len(message) != 6 {
		t.Error("bad marshal")
	}

	q := &Package{}
	if err := q.Unmarshal(bytes.NewReader(message)); err != nil {
		t.Error(err)
	}
	t.Logf("%+v", q)
}

func Test_PackageEncodeAndDecodeWithPayload(t *testing.T) {
	payload := []byte("hello world")
	p := NewPackageWithPayload(0x11, payload)
	message, err := p.Marshal()
	if err != nil {
		t.Error(err)
	}
	t.Logf("%x", message)
	if len(message) != 10+len(payload) {
		t.Error("bad marshal")
	}

	q := &Package{}
	if err := q.Unmarshal(bytes.NewReader(message)); err != nil {
		t.Error(err)
	}
	t.Logf("%+v", q)
}
