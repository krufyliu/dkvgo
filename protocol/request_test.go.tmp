package protocol

import "testing"
import "bufio"
import "bytes"

func requestEqual(r1, r2 *DkvRequst) bool {
	return r1.Version == r2.Version && r1.Method == r2.Method
}

func Test_Reqeust(t *testing.T) {
	req := NewDkvRequest("RUN", []string{"test", "hello"})
	binContent, err := req.Dumps()
	if err != nil {
		t.Error(err)
	}
	t.Logf("%x", binContent)
	r := bufio.NewReader(bytes.NewReader(binContent))
	req = &DkvRequst{}
	if err = req.PullHeader(r); err != nil {
		t.Error(err)
	}
	t.Logf("%+v", req)
}

func Test_ReqeustWithPayload(t *testing.T) {
	req := NewDkvRequestWithPayload("RUN", []string{"test", "hello"}, []byte("this is just a test"))
	binContent, err := req.Dumps()
	if err != nil {
		t.Error(err)
	}
	t.Logf("%x", binContent)
	r := bufio.NewReader(bytes.NewReader(binContent))
	req = &DkvRequst{}
	if err = req.PullHeader(r); err != nil {
		t.Error(err)
	}
	t.Logf("%+v", req)
	if err = req.PullContent(r); err != nil {
		t.Error(err)
	}
	t.Logf("%+v", req)
}
