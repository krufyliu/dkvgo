package protocol

import (
    "testing"
	"bytes"
)

func Test_ProtocolEncodeAndDecode(t *testing.T) {
    p := NewProtocol(0x11)
    message, err := p.Marshal()
    if err != nil {
        t.Error(err)
    }
    t.Logf("%x", message)
    if len(message) != 6 {
        t.Error("bad marshal")
    }

    q := &Protocol{}
    if err := q.Unmarshal(bytes.NewReader(message)); err != nil {
        t.Error(err)
    }
    t.Logf("%+v", q)
}

func Test_ProtocolEncodeAndDecodeWithPayload(t *testing.T) {
    payload := []byte("hello world")
    p := NewProtocolWithPayload(0x11, payload)
    message, err := p.Marshal()
    if err != nil {
        t.Error(err)
    }
    t.Logf("%x", message)
    if len(message) != 10 + len(payload) {
        t.Error("bad marshal")
    }

    q := &Protocol{}
    if err := q.Unmarshal(bytes.NewReader(message)); err != nil {
        t.Error(err)
    }
    t.Logf("%+v", q)
}