package worker

import (
	"bufio"
	"bytes"
	"io"
	"testing"
)

func Test_MatchState(t *testing.T) {
	var testStr = `I0104 19:34:59.578550  6430 test_3d_facebook.cpp:697] [ Input dir: /dkvision/data/jia_yu_guan_20170104/gopro02/item3 ]
I0104 19:34:59.578553  6430 test_3d_facebook.cpp:698] [ Output dir: /dkvision/output/jia_yu_guan_20170104/gopro02/item3/FACEBOOK_3D ]
I0104 19:34:59.578555  6430 test_3d_facebook.cpp:699] [ Process frame: 3600 ]
I0104 19:35:02.013476  6430 test_3d_facebook.cpp:771] [ Prepare images time: 2.43491 second ]
I0104 19:35:05.486690  6430 test_3d_facebook.cpp:261] [ Compute flow time: 3.47317 second ]
I0104 19:35:07.415637  6430 test_3d_facebook.cpp:342] [ Compute novel view time: 1.92889 second ]
I0104 19:35:09.705060  6430 test_3d_facebook.cpp:968] [ Total time: 10.1265 second ]`
	r := bytes.NewReader([]byte(testStr))
	state, err := matchState(bufio.NewReader(r))
	if err != nil && err != io.EOF {
		t.Error(err)
	}
	t.Logf("%+v", state)
}

func Test_MatchStateEOF(t *testing.T) {
	var testStr = `I0104 19:34:59.578550  6430 test_3d_facebook.cpp:697] [ Input dir: /dkvision/data/jia_yu_guan_20170104/gopro02/item3 ]
I0104 19:34:59.578553  6430 test_3d_facebook.cpp:698] [ Output dir: /dkvision/output/jia_yu_guan_20170104/gopro02/item3/FACEBOOK_3D ]
I0104 19:34:59.578555  6430 test_3d_facebook.cpp:699] [ Process frame: 3600 ]
I0104 19:35:02.013476  6430 test_3d_facebook.cpp:771] [ Prepare images time: 2.43491 second ]
I0104 19:35:05.486690  6430 test_3d_facebook.cpp:261] [ Compute flow time: 3.47317 second ]
I0104 19:35:07.415637  6430 test_3d_facebook.cpp:342] [ Compute novel view time: 1.92889 second ]`
	r := bytes.NewReader([]byte(testStr))
	_, err := matchState(bufio.NewReader(r))
	if err != io.EOF {
		t.Error("should return EOF error")
	}
}
