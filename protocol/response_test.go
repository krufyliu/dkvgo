package protocol

import "testing"
import "bufio"
import "bytes"

func Test_Response(t *testing.T) {
	res := NewDkvResponse("OK")
	binContent, err := res.Dumps()
	if err != nil {
		t.Error(err)
	}
	t.Logf("%x", binContent)
	r := bufio.NewReader(bytes.NewReader(binContent))
	res = &DkvResponse{}
	if err = res.PullHeader(r); err != nil {
		t.Error(err)
	}
	t.Logf("%+v", res)
}

func Test_ResonseWithPayload(t *testing.T) {
	res := NewDkvResponseWithPayload("OK", []byte("hello world"))
	binContent, err := res.Dumps()
	if err != nil {
		t.Error(err)
	}
	t.Logf("%x", binContent)
	r := bufio.NewReader(bytes.NewReader(binContent))
	res = &DkvResponse{}
	if err = res.PullHeader(r); err != nil {
		t.Error(err)
	}
	t.Logf("%+v", res)
	if err = res.PullContent(r); err != nil {
		t.Error(err)
	}
	t.Logf("%+v", res)
}
