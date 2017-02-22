package dkvgo

import (
	"bytes"
	"testing"
)

func dkvProtocolEqual(left, right *DkvProtocol) bool {
	if left.Action != right.Action {
		return false
	}
	if len(left.Payload) != len(left.Payload) {
		return false
	}
	for key, value := range left.Payload {
		if right.Payload[key] != value {
			return false
		}
	}
	return true
}

func Test_Encode_And_Decode(t *testing.T) {
	p := NewDkvProtocol()
	p.Action = "ping"
	p.AddPayload("ts", "141123123123")
	buffer, err := p.Encode()
	if err != nil {
		t.Error(err.Error())
	}
	t.Logf("%x", buffer)
	r := bytes.NewReader(buffer)
	dp, err := DecodeFrom(r)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(dp)
	if !dkvProtocolEqual(p, dp) {
		t.Error("decode result is not equal with original data")
	}
}
